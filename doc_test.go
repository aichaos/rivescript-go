package rivescript_test

import (
	"fmt"
	rivescript "github.com/aichaos/rivescript-go"
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
