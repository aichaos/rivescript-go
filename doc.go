/*
Package rivescript implements the RiveScript chatbot scripting language.

About RiveScript

RiveScript is a scripting language for authoring chatbots. It has a very
simple syntax and is designed to be easy to read and fast to write.

A simple example of what RiveScript looks like:

	+ hello bot
	- Hello human.

This matches a user's message of "hello bot" and would reply "Hello human."
Or for a slightly more complicated example:

	+ my name is *
	* <formal> == <bot name> => <set name=<formal>>Wow, we have the same name!
	* <get name> != undefined => <set name=<formal>>Did you change your name?
	- <set name=<formal>>Nice to meet you, <get name>!

The official website for RiveScript is https://www.rivescript.com/

To test drive RiveScript in your web browser, try the
[RiveScript Playground](https://play.rivescript.com/).

Object Macros

A common feature in many RiveScript implementations is the object macro, which
enables you to write dynamic program code (in your favorite programming
language) to add extra capabilities to your bot. For example, your bot could
answer a question of `what is the weather like in _____` by running some
code to look up their answer via a web API.

The Go version of RiveScript has support for object macros written in Go
(at compile time of your application). It also has optional support for
JavaScript object macros using the Otto library.

UTF-8 Support

UTF-8 support in RiveScript is considered an experimental feature. It is
disabled by default. Enable it by setting `RiveScript.SetUTF8(true)`.

By default (without UTF-8 mode on), triggers may only contain basic ASCII
characters (no foreign characters), and the user's message is stripped of all
characters except letters, numbers and spaces. This means that, for example,
you can't capture a user's e-mail address in a RiveScript reply, because of
the @ and . characters.

When UTF-8 mode is enabled, these restrictions are lifted. Triggers are only
limited to not contain certain metacharacters like the backslash, and the
user's message is only stripped of backslashes and HTML angled brackets
(to protect from obvious XSS if you use RiveScript in a web application).
Additionally, common punctuation characters are stripped out, with the default
set being `/[.,!?;:]/g`. This can be overridden by providing a new regexp
string literal to the `RiveScript.SetUnicodePunctuation` function. Example:

	// Make a new bot with UTF-8 mode enabled.
	bot := rivescript.New(config.UTF8())

	// Override the punctuation characters that get stripped from the
	// user's message.
	bot.SetUnicodePunctuation(`[.,!?;:]`);

The `<star>` tags in RiveScript will capture the user's "raw" input, so you can
write replies to get the user's e-mail address or store foreign characters in
their name.

See Also

The official homepage of RiveScript, http://www.rivescript.com/

*/
package rivescript
