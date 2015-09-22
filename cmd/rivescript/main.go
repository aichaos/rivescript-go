package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	rivescript "github.com/aichaos/rivescript-go"
)

func main() {
	bot := rivescript.New()
	bot.Debug = true
	fmt.Printf("RiveScript version v%s\n", bot.Version())
	//bot.Debug = true
	//bot.Stream("+ hello bot\n" +
	// 	"- Hello human.\n" +
	// 	"- How are you?\n" +
	// 	"+ goodbye robot\n" +
	// 	"% hello human\n" +
	// 	"- Test\n" +
	// 	"+ * or something{weight=10}\n" +
	// 	"- Or something. <@>\n" +
	// 	"+ *\n" +
	// 	"- I dunno.")
	// bot.LoadFile("eg/brain/rpg.rive")
	//bot.LoadDirectory("eg/brain")
	bot.LoadFile("testsuite.rive")

	bot.SortReplies()
	// bot.Debug = true
	// bot.DumpTopics()
	// bot.DumpSorted()
	//bot.Debug = true
	// fmt.Printf("--- DONE SORTING ---\n")
	// bot.DumpSorted()
	// bot.LoadFile("eg/brain/begin.rive")
	// bot.LoadFile("eg/brain/admin.rive")
	//bot.LoadFile("test.rive")
	// bot.DumpTopics()
	bot.Debug = true

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}

		reply := bot.Reply("local-user", text)
		fmt.Printf("Bot> %s\n", reply)
	}
}
