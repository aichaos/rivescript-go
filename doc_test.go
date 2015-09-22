package rivescript_test

import (
	"fmt"
	rivescript "github.com/aichaos/rivescript-go"
	"github.com/aichaos/rivescript-go/lang/rivescript_js"
)

func ExampleRiveScript() {
	bot := rivescript.New()

	// Load a directory full of RiveScript documents (.rive files)
	bot.LoadDirectory("eg/brain")

	// Load an individual file.
	bot.LoadFile("testsuite.rive")

	// Sort the replies after loading them!
	bot.SortReplies()

	// Get a reply.
	reply := bot.Reply("local-user", "Hello, bot!")
	fmt.Printf("The bot says: %s", reply)
}

func ExampleRiveScript_JavaScript() {
	// Example for configuring the JavaScript object macro handler via Otto.

	bot := rivescript.New()

	// Create the JS handler.
	jsHandler := rivescript_js.New()
	bot.SetHandler("javascript", jsHandler)

	// Now we can use object macros written in JS!
	bot.Stream(`
		> object add javascript
			var a = args[0];
			var b = args[1];
			return parseInt(a) + parseInt(b);
		< object

		+ add # and #
		- <star1> + <star2> = <call>add <star1> <star2></call>
	`)
	bot.SortReplies()

	reply := bot.Reply("local-user", "Add 5 and 7")
	fmt.Printf("Bot: %s", reply)
}
