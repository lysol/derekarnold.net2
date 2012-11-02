<!--
{
    "Title": "Let's go!",
    "Date": "2012-11-01 07:58 PM",
    "Tags": ["meta","golang"]
}
-->

I've made a new rule for my own site: write it in a fad language. This year's
fad language is [Go](http://golang.org). (The term "fad" here is used endearingly.)

"Wikipedia says" is the new "Webster's Dictionary defines", but Wikipedia really
does sum it up best:

> Go aims to provide the efficiency of a statically typed compiled language
> with the ease of programming of a dynamic language.

<img class="left-body" src="http://i.imgur.com/vTAzf.jpg" alt="BLLLLLLRREEEEEE" />

I like writing C, but also hate writing it. The memory management is akin to that 
part of the Mega Man games where you have to jump on disappearing blocks in a specific order.

Go's [http](http://golang.org/pkg/net/http/) library is fantastic. There is no 
need for web micro-frameworks, minus some basic URL parsing.

The [template](http://golang.org/pkg/html/template) library is also a welcome 
addition but is rough around the edges. The control structures within template
tags are fairly primitive, requiring you to embed your own functions to accomplish
basic boolean logic within the template itself. Exceptions within the template
may not always create an exception within the application, instead truncating
output. Otherwise, the scope control within the template itself is pretty
useful and even a net gain over other template engines like [Jinja](http://jinja.pocoo.org/) and [Django](https://www.djangoproject.com/)'s.

Would I use Go for other web projects? Yeah, there's some pain during development
but the highly structured nature of it leads to fairly sturdy code.
