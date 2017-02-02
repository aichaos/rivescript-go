package rivescript

import "testing"

func TestConcat(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		// Default concat mode = none
		+ test concat default
		- Hello
		^ world!

		! local concat = space
		+ test concat space
		- Hello
		^ world!

		! local concat = none
		+ test concat none
		- Hello
		^ world!

		! local concat = newline
		+ test concat newline
		- Hello
		^ world!

		// invalid concat setting is equivalent to 'none'
		! local concat = foobar
		+ test concat foobar
		- Hello
		^ world!

		// the option is file scoped so it can be left at
		// any setting and won't affect subsequent parses
		! local concat = newline
	`)
	bot.extend(`
		// concat mode should be restored to the default in a
		// separate file/stream parse
		+ test concat second file
		- Hello
		^ world!
	`)

	bot.reply("test concat default", "Helloworld!")
	bot.reply("test concat space", "Hello world!")
	bot.reply("test concat none", "Helloworld!")
	bot.reply("test concat newline", "Hello\nworld!")
	bot.reply("test concat foobar", "Helloworld!")
	bot.reply("test concat second file", "Helloworld!")
}
