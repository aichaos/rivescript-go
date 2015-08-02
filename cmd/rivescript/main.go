package main

import (
	//"fmt"
	rivescript "github.com/aichaos/rivescript-go"
)

func main() {
	bot := rivescript.New()
	bot.Debug = false
	bot.LoadDirectory("eg/brain")
	// bot.LoadFile("eg/brain/begin.rive")
	// bot.LoadFile("eg/brain/admin.rive")
	//bot.LoadFile("test.rive")
	bot.DumpTopics()
}
