package main

import (
	"fmt"
	rivescript "github.com/aichaos/rivescript-go"
)

func main() {
	bot := rivescript.New()
	fmt.Printf("RiveScript version v%s\n", bot.Version())
	bot.Debug = true
	bot.Stream("+ hello bot\n" +
		"- Hello human.\n" +
		"- How are you?\n" +
		"+ hi\n" +
		"% hello human\n" +
		"- Test")
	// bot.LoadFile("eg/brain/rpg.rive")
	bot.SortReplies()
	//bot.LoadDirectory("eg/brain")
	// bot.LoadFile("eg/brain/begin.rive")
	// bot.LoadFile("eg/brain/admin.rive")
	//bot.LoadFile("test.rive")
	// bot.DumpTopics()
}
