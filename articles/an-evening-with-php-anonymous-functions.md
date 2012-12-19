<!--
{
    "Title": "An Evening with PHP Anonymous Functions",
    "Date": "2012-12-18 09:29 PM",
    "Tags": ["php","javascript"]
}
-->

Hey there, sailor. Light a couple of candles, pour yourself some wine (preferably from a box),
and put on something more comfortable. We're going to _write some PHP_.

As of PHP 5.4, anonymous functions are now not only available, but make their scope
available via the $this variable like regular objects. They're also represented by a
Closure class, which provides a couple of methods for manipulating its scope.

PHP gives you a couple of ways to mess with anonymous function scope. The first:

    $foo = 'bar';
  	$anon = function() use ($foo) {
  			print $foo;
  	};
  	$anon(); // 'bar'

This is the most straightforward route to using variables from the enclosing scope
in anonymous functions with PHP. It's not as simple as Javascript, since JS automatically
encloses all the variables from outside of the function inside of it, but it's just one
of those cases where you have to be more explicit that take advantage of the implicit
style of Javascript programming.

Binding
-------

So in Javascript, you can easily bind .

    var scope = {foo: 'bar'};
  	var anon = (function() { console.log(this.foo); }).bind(scope);
  	anon(); // returns 'bar'

PHP 5.4's Closure class gives you the ability to do something similar.

    $scope = array('foo' => 'bar');
  	$anon = function() { print $this->foo; };
  	$anon->bindTo($scope);
  	$anon();

But! The above doesn't work. You can only bind anonymous functions to objects in PHP. So
you have to do something along the lines of:

    $scope = new stdClass;
  	$scope->foo = 'bar';
  	$anon = function() { print $this->foo; };
  	$anon->bindTo($scope);
  	$anon() // 'bar'

This is both powerful and useless, if you're writing effective object-oriented code.
It's functionally equivalent to Javascript's function bind method, but because PHP's
$this support is bolted on and can only be used within instance methods and anonymous functions,
I can't foresee a pattern where this is useful, unless you use nothing but stdClass instances
like associative arrays. Which would be silly.

PHP's anonymous functions would be a bit more interesting if you could subclass or modify the
Closure class, to do things like add a functional equivalent to the apply method in Javascript.
However, the class is not extendable, instantiable, or generally useful in ways other objects are.

You can't even add properties to anonymous functions, so you can't even bolt on an apply() method
after the fact. The Closure class is completely nerfed.

I wanted to do something really interesting and twisted with anonymous functions in PHP but the
non-extensibility once again prevents anything really crafty from happening. Bummer.
