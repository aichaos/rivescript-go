package rivescript_test

import (
	"fmt"

	"github.com/aichaos/rivescript-go"
	"github.com/aichaos/rivescript-go/lang/javascript"
)

func Example() {
	// Create a new RiveScript instance, which represents an individual bot
	// with its own brain and memory of users.
	//
	// You can provide a rivescript.Config struct to configure the bot and
	// provide values that differ from the defaults:
	bot := rivescript.New(&rivescript.Config{
		UTF8:  true, // enable UTF-8 mode
		Debug: true, // enable debug mode
	})

	// Or if you want the default configuration, provide a nil config.
	// See the documentation for the rivescript.Config type for information
	// on what the defaults are.
	bot = rivescript.New(nil)

	// Load a directory full of RiveScript documents (.rive files)
	bot.LoadDirectory("eg/brain")

	// Load an individual file.
	bot.LoadFile("testsuite.rive")

	// Stream in more RiveScript code dynamically from a string.
	bot.Stream(`
		+ hello bot
		- Hello, human!
	`)

	// Sort the replies after loading them!
	bot.SortReplies()

	// Get a reply.
	reply, _ := bot.Reply("local-user", "Hello, bot!")
	fmt.Printf("The bot says: %s", reply)
}

func ExampleRiveScript_utf8() {
	// Examples of using UTF-8 mode in RiveScript.
	bot := rivescript.New(rivescript.WithUTF8())

	bot.Stream(`
		// Without UTF-8 mode enabled, this trigger would be a syntax error
		// for containing non-ASCII characters; but in UTF-8 mode you can use it.
		+ comment ça va
		- ça va bien.
	`)

	// Always call SortReplies when you're done loading replies.
	bot.SortReplies()

	// Without UTF-8 mode enabled, the user's message "comment ça va" would
	// have the ç symbol removed; but in UTF-8 mode it's preserved and can
	// match the trigger we defined.
	reply, _ := bot.Reply("local-user", "Comment ça va?")
	fmt.Println(reply) // "ça va bien."
}

func ExampleRiveScript_javascript() {
	// Example for configuring the JavaScript object macro handler via Otto.
	bot := rivescript.New(nil)

	// Create the JS handler.
	bot.SetHandler("javascript", javascript.New(bot))

	// Now we can use object macros written in JS!
	bot.Stream(`
		> object add javascript
			var a = args[0];
			var b = args[1];
			return parseInt(a) + parseInt(b);
		< object

		> object setname javascript
			// Set the user's name via JavaScript
			var uid = rs.CurrentUser();
			rs.SetUservar(uid, args[0], args[1])
		< object

		+ add # and #
		- <star1> + <star2> = <call>add <star1> <star2></call>

		+ my name is *
		- I will remember that.<call>setname <id> <formal></call>

		+ what is my name
		- You are <get name>.
	`)
	bot.SortReplies()

	reply, _ := bot.Reply("local-user", "Add 5 and 7")
	fmt.Printf("Bot: %s\n", reply)
}

func ExampleRiveScript_subroutine() {
	// Example for defining a Go function as an object macro.
	bot := rivescript.New(nil)

	// Define an object macro named `setname`
	bot.SetSubroutine("setname", func(rs *rivescript.RiveScript, args []string) string {
		uid, _ := rs.CurrentUser()
		rs.SetUservar(uid, args[0], args[1])
		return ""
	})

	// Stream in some RiveScript code.
	bot.Stream(`
		+ my name is *
		- I will remember that.<call>setname <id> <formal></call>

		+ what is my name
		- You are <get name>.
	`)
	bot.SortReplies()

	bot.Reply("local-user", "my name is bob")
	reply, _ := bot.Reply("local-user", "What is my name?")
	fmt.Printf("Bot: %s\n", reply)
}
