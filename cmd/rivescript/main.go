/*
Stand-alone RiveScript Interpreter.

This is an example program included with the RiveScript Go library. It serves as
a way to quickly demo and test a RiveScript bot.

Usage

	rivescript [options] /path/to/rive/files

Options

	--debug     Enable debug mode.
	--utf8      Enable UTF-8 support within RiveScript.
	--depth     Override the recursion depth limit (default 50)
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	rivescript "github.com/aichaos/rivescript-go"
	"github.com/aichaos/rivescript-go/lang/rivescript_js"
)

func main() {
	// Collect command line arguments.
	debug := flag.Bool("debug", false, "Enable debug mode.")
	utf8 := flag.Bool("utf8", false, "Enable UTF-8 mode.")
	depth := flag.Int("depth", 50, "Recursion depth limit (default 50)")
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: rivescript [options] </path/to/documents>")
		os.Exit(1)
	}

	root := args[0]

	// Initialize the bot.
	bot := rivescript.New()
	bot.Debug = *debug
	bot.UTF8 = *utf8
	bot.Depth = *depth

	// JavaScript object macro handler.
	jsHandler := rivescript_js.New(bot)
	bot.SetHandler("javascript", jsHandler)

	// Load the target directory.
	err := bot.LoadDirectory(root)
	if err != nil {
		fmt.Printf("Error loading directory: %s", err)
		os.Exit(1)
	}

	bot.SortReplies()

	fmt.Printf(`RiveScript Interpreter (Golang) -- Interactive Mode
---------------------------------------------------
RiveScript version: %s
        Reply root: %s

You are now chatting with the RiveScript bot. Type a message
and press Return to send it. When finished, type '/quit' to
exit the program. Type '/help' for other options.
`, bot.Version(), root)

	// Drop into the interactive command shell.
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("You> ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}

		if strings.Index(text, "/help") == 0 {
			help()
		} else if strings.Index(text, "/quit") == 0 {
			os.Exit(0)
		} else {
			reply := bot.Reply("localuser", text)
			fmt.Printf("Bot> %s\n", reply)
		}
	}
}

func help() {
	fmt.Printf(`Supported commands:
- /help : Show this text.
- /quit : Exit the program.
`)
}
