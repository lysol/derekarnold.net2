package main

import (
    "fmt"
    "os"
    "time"
    "strings"
    "strconv"
    "io/ioutil"
    "path/filepath"
    "net/http"
    "html/template"
    "encoding/json"
    "github.com/russross/blackfriday"
)

const (
    PostsPerPage int = 5
    PostDateFormat string = "2006-01-02 03:04 PM"
    DisplayDateFormat string = "January _2, 2006"
)

type DerekError struct {
    Message    string
}

func (e DerekError) Error() string {
    return fmt.Sprintf("%v", e.Message)
}

type Nav struct {
    Token string
    Name string
}

type Post struct {
    FileName string
    Body template.HTML
    Meta Metadata
    Token string
}

func (p Post) NiceDate() string {
    return p.RealDate().Format(DisplayDateFormat)
}

func (p Post) RealDate() time.Time {
    t, err := time.Parse(PostDateFormat, p.Meta.Date)
    if err != nil {
        t2, _ := time.Parse(PostDateFormat, "0000-01-01 00:00 AM")
        return t2
    }
    return t
}

type Metadata struct {
    Title string
    Date string
    Tags []string
}

type PageContext struct {
    Navs []Nav
    CurrentNav *Nav
    Page int
    Posts []*Post
    NextLink string
    PrevLink string
    Tag string
    Post *Post
}

var Navs = []Nav {
    Nav{Token: "", Name: "Blog"},
    Nav{Token: "resume", Name: "Résumé"},
    Nav{Token: "about", Name: "About"},
}

func setTagMap() (map[string][][]*Post) {

    tagMap := map[string][][]*Post{}
    for _, post := range postMap {
        for _, tag := range post.Meta.Tags {
            if tagMap[tag] == nil {
                tagMap[tag] = make([][]*Post, 0)
                tagMap[tag] = append(tagMap[tag], []*Post{})
            }
            curInd := len(tagMap[tag]) - 1
            if len(tagMap[tag][curInd]) >= PostsPerPage {
                tagMap[tag] = append(tagMap[tag], []*Post{})
                curInd++
            }
            tagMap[tag][curInd] = append(tagMap[tag][curInd], post)
        }
    }
    return tagMap
}

var articleNames []string

func getArticleNames() ([]string) {
    if len(articleNames) == 0 {
        walker := func(path string, f os.FileInfo, err error) error {
            fname := f.Name()
            if (len(fname) > 3 && fname[len(fname)-3:] == ".md") ||
                (len(fname) > 9 && fname[len(fname)-9:] == ".markdown") {
                articleNames = append(articleNames, "articles/" + f.Name())
            }
            return nil
        }
        filepath.Walk("articles", walker)
    }
    return articleNames
}

var cachedGetPosts []*Post

func getPosts() ([]*Post) {
    if len(cachedGetPosts) > 0 {
        return cachedGetPosts
    }
    articleNames := getArticleNames()

    posts := []*Post{}
    for _, aname := range articleNames {
        buf, err := ioutil.ReadFile(aname)
        contents := string(buf)
        if err != nil {
            continue
        }
        splat := strings.Split(contents, "-->")
        var m Metadata
        justJson := strings.Replace(splat[0][4:], "\n", "", -1)
        err = json.Unmarshal([]byte(justJson), &m)
        if err != nil {
            fmt.Println(err)
            continue
        }
        token := strings.Split(aname[9:], ".")[0]
        post := Post{
            FileName: aname,
            Body: template.HTML(blackfriday.MarkdownBasic(buf)),
            Meta: m,
            Token: token,
        }
        posts = append(posts, &post)
        postMap[token] = &post
    }
    cachedGetPosts = posts
    return posts
}

var cachedSortedPosts []*Post

func sortedPosts() ([]*Post) {
    if len(cachedSortedPosts) > 0 {
        return cachedSortedPosts
    }
    posts := getPosts()
    // lil bub
    var swaps int
    for {
        swaps = 0
        for aindex, post := range posts[:len(posts)-1] {
            if post.RealDate().Before(posts[aindex+1].RealDate()) {
                posts[aindex] = posts[aindex+1]
                posts[aindex+1] = post
                swaps++
            }
        }
        if swaps == 0 { break }
    }
    cachedSortedPosts = posts
    return posts
}

func setPages() ([][]*Post) {
    posts := sortedPosts()

    articleNames := getArticleNames()
    maxPages := len(articleNames) / PostsPerPage
    if len(articleNames) % 10 > 0 {
        maxPages++
    }

    pages := [][]*Post{}

    for aindex, post := range posts {
        pageNumber := int(aindex / PostsPerPage)
        for len(pages) <= pageNumber {
            pages = append(pages, []*Post{})
        }
        pages[pageNumber] = append(pages[pageNumber], post)
    }

    return pages
}

// Function to retrieve all appropriately named templates in the template
// directory.
func setTemplates() (map[string]*template.Template) {
    templateNames := []string {}
    texas := func(path string, f os.FileInfo, err error) error {
        // In the eyes of the ranger
        fname := f.Name()
        if fname[len(fname)-5:] == ".html" && fname != "layout.html" {
            // The unsuspected stranger
            templateNames = append(templateNames, "templates/" + f.Name())
        }
        return nil
    }
    filepath.Walk("templates", texas)

    templates := make(map[string]*template.Template)
    for _, tname := range templateNames {
        templates[tname[10:]] = template.Must(template.New(tname).Funcs(template.FuncMap{
                        "eq": func(a, b string) bool {
                                return a == b
                        },
                        "ne": func(a string, b string) bool {
                                return a != b
                        },
                        "gt": func(a int, b int) bool {
                                return a > b
                        },
                        }).ParseFiles(tname, "templates/layout.html"))
    }
    return templates
}

func getNav(token string) (*Nav, *DerekError) {
    for _, nav := range Navs {
        if nav.Token == token {
            return &nav, nil
        }
    }
    return nil, &DerekError{"No nav found"}
}

// Load all the templates.
var templates = setTemplates()
var postMap = map[string]*Post{}
var pages = setPages()
var tagMap = setTagMap()

func renderTemplate(w http.ResponseWriter, name string, ctx interface{}) {
    templates[name].ExecuteTemplate(w, "base", ctx)
}

func ezRender(w http.ResponseWriter, name string, currentNav *Nav) {
    renderTemplate(w, name, PageContext{Navs: Navs, CurrentNav: currentNav})
}

func main() {

    /*navs := map[string] string {
        "": "Blog",
        "resume": "Résumé",
        "about": "About",
    }*/

    //config := map[string] string {}
    http.Handle("/images/",
        http.StripPrefix("/images/", http.FileServer(http.Dir("images/"))))

    http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
        currentNav, err := getNav("about")
        if err != nil {
            fmt.Fprintf(w, err.Message)
            return
        }
        ezRender(w, "about.html", currentNav)
    })

    http.HandleFunc("/resume", func(w http.ResponseWriter, r *http.Request) {
        currentNav, err := getNav("resume")
        if err != nil {
            fmt.Fprintf(w, err.Message)
            return
        }
        ezRender(w, "resume.html", currentNav)
    })

    http.HandleFunc("/post/", func(w http.ResponseWriter, r *http.Request) {
        currentNav, err := getNav("")
        if err != nil {
            fmt.Fprintf(w, err.Message)
            return
        }
        token := r.URL.Path[6:]
        if postMap[token] == nil {
            http.NotFound(w, r)
            return
        }
        renderTemplate(w, "post.html", PageContext{
            Navs: Navs,
            CurrentNav: currentNav,
            Post: postMap[token],
            Page: 0,
        })
    })

    blogHandler := func(w http.ResponseWriter, r *http.Request) {
        currentNav, err := getNav("")
        tag := ""
        if err != nil {
            fmt.Fprintf(w, err.Message)
            return
        }
        page := 1
        parts := strings.Split(r.URL.Path, "/")
        var posts []*Post
        maxPage := len(pages) - 1
        if parts[1] == "page" {
            newpage, err2 := strconv.Atoi(parts[2])
            if err2 != nil {
                http.NotFound(w, r)
                return
            }
            page = newpage
            if page - 1 > len(pages) {
                http.NotFound(w, r)
                return
            }
            posts = pages[page-1]
        } else if parts[1] == "tags" {
            tag = parts[2]
            if tagMap[tag] == nil {
                http.NotFound(w, r)
                return
            }
            if len(parts) >= 5 && parts[3] == "page" {
                newpage, err2 := strconv.Atoi(parts[4])
                if err2 != nil {
                    http.NotFound(w, r)
                    return
                }
                page = newpage
            }
            if page - 1 > len(tagMap[tag]) {
                http.NotFound(w, r)
                return
            }
            posts = tagMap[tag][page-1]
            maxPage = len(tagMap[tag]) - 1
        } else if r.URL.Path == "/" {
            posts = pages[0]
        } else {
            http.NotFound(w, r)
            return
        }
        page--
        var nextLink string
        var prevLink string

        if page > 0 && tag != "" {
            prevLink = fmt.Sprintf("/tags/%s/page/%d", tag, page)
        } else if page > 0 {
            prevLink = fmt.Sprintf("/page/%d", page)
        }
        if page < maxPage && tag != "" {
            nextLink = fmt.Sprintf("/tags/%s/page/%d", tag, page+2)
        } else if page < maxPage {
            nextLink = fmt.Sprintf("/page/%d", page+2)
        }

        renderTemplate(w, "blog.html", PageContext{
            Navs: Navs,
            CurrentNav: currentNav,
            Page: page + 1,
            Posts: posts,
            NextLink: nextLink,
            PrevLink: prevLink,
            Tag: tag,
        })

    }

    http.HandleFunc("/tags", blogHandler)
    http.HandleFunc("/page", blogHandler)
    http.HandleFunc("/", blogHandler)

    http.ListenAndServe(":8090", nil)

}