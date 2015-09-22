/*
Package rivescript implements the RiveScript chatbot scripting language.

UTF-8 Support

UTF-8 support in RiveScript is considered an experimental feature. It is
disabled by default. Enable it by setting `RiveScript.UTF8 = true`.

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
set being `/[.,!?;:]/g`. This can be overridden by providing a new `Regexp`
object as the `RiveScript.UnicodePunctuation` attribute. Example:

	// Make a new bot with UTF-8 mode enabled.
	bot := rivescript.New()
	bot.UTF8 = true

	// Override the punctuation characters that get stripped from the
	// user's message.
	bot.UnicodePunctuation = regexp.MustCompile(`[.,!?;:]`);

The `<star>` tags in RiveScript will capture the user's "raw" input, so you can
write replies to get the user's e-mail address or store foreign characters in
their name.

See Also

The official homepage of RiveScript, http://www.rivescript.com/

*/
package rivescript
