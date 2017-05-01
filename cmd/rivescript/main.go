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

	"github.com/aichaos/rivescript-go"
	"github.com/aichaos/rivescript-go/lang/javascript"
)

// Build is the git commit hash that the binary was built from.
var Build = "-unknown-"

var (
	// Command line arguments.
	version  bool
	debug    bool
	utf8     bool
	depth    uint
	nostrict bool
	nocolor  bool
)

func init() {
	flag.BoolVar(&version, "version", false, "Show the version number and exit.")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode.")
	flag.BoolVar(&utf8, "utf8", false, "Enable UTF-8 mode.")
	flag.UintVar(&depth, "depth", 50, "Recursion depth limit")
	flag.BoolVar(&nostrict, "nostrict", false, "Disable strict syntax checking")
	flag.BoolVar(&nocolor, "nocolor", false, "Disable ANSI colors")
}

func main() {
	// Collect command line arguments.
	flag.Parse()
	args := flag.Args()

	if version {
		fmt.Printf("RiveScript-Go version %s\n", rivescript.Version)
		os.Exit(0)
	}

	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: rivescript [options] </path/to/documents>")
		os.Exit(1)
	}

	root := args[0]

	// Initialize the bot.
	bot := rivescript.New(&rivescript.Config{
		Debug:  debug,
		Strict: !nostrict,
		Depth:  depth,
		UTF8:   utf8,
	})

	// JavaScript object macro handler.
	bot.SetHandler("javascript", javascript.New(bot))

	// Load the target directory.
	err := bot.LoadDirectory(root)
	if err != nil {
		fmt.Printf("Error loading directory: %s", err)
		os.Exit(1)
	}

	bot.SortReplies()

	fmt.Printf(`
      .   .
     .:...::      RiveScript Interpreter (Go)
    .::   ::.     Library Version: v%s (build %s)
 ..:;;. ' .;;:..
    .  '''  .     Type '/quit' to quit.
     :;,:,;:      Type '/help' for more options.
     :     :

Using the RiveScript bot found in: %s
Type a message to the bot and press Return to send it.
`, rivescript.Version, Build, root)

	// Drop into the interactive command shell.
	reader := bufio.NewReader(os.Stdin)
	for {
		color(yellow, "You>")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) == 0 {
			continue
		}

		if strings.Contains(text, "/help") {
			help()
		} else if strings.Contains(text, "/quit") {
			os.Exit(0)
		} else if strings.Contains(text, "/debug t") {
			bot.SetGlobal("debug", "true")
			color(cyan, "Debug mode enabled.", "\n")
		} else if strings.Contains(text, "/debug f") {
			bot.SetGlobal("debug", "false")
			color(cyan, "Debug mode disabled.", "\n")
		} else if strings.Contains(text, "/debug") {
			debug, _ := bot.GetGlobal("debug")
			color(cyan, "Debug mode is currently:", debug, "\n")
		} else if strings.Contains(text, "/dump t") {
			bot.DumpTopics()
		} else if strings.Contains(text, "/dump s") {
			bot.DumpSorted()
		} else {
			reply, err := bot.Reply("localuser", text)
			if err != nil {
				color(red, "Error>", err.Error(), "\n")
			} else {
				color(green, "RiveScript>", reply, "\n")
			}
		}
	}
}

// Names for pretty ANSI colors.
const (
	red    = `31;1`
	yellow = `33;1`
	green  = `32;1`
	cyan   = `36;1`
)

func color(color string, text ...string) {
	if nocolor {
		fmt.Printf(
			"%s %s",
			text[0],
			strings.Join(text[1:], " "),
		)
	} else {
		fmt.Printf(
			"\x1b[%sm%s\x1b[0m %s",
			color,
			text[0],
			strings.Join(text[1:], " "),
		)
	}
}

func help() {
	fmt.Printf(`Supported commands:
- /help
    Show this text.
- /quit
    Exit the program.
- /debug [true|false]
    Enable or disable debug mode. If no setting is given, it prints
    the current debug mode.
- /dump <topics|sorted>
    For debugging purposes, dump the topic and sorted trigger trees.
`)
}
