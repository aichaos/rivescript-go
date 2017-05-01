package rivescript

// Debugging methods

import (
	"fmt"
)

// say prints a debugging message
func (rs *RiveScript) say(message string, a ...interface{}) {
	if rs.Debug {
		fmt.Printf(message+"\n", a...)
	}
}

// warn prints a warning message for non-fatal errors
func (rs *RiveScript) warn(message string, a ...interface{}) {
	if !rs.Quiet {
		fmt.Printf("[WARN] "+message+"\n", a...)
	}
}

// warnSyntax is like warn but takes a filename and line number.
func (rs *RiveScript) warnSyntax(message string, filename string, lineno int, a ...interface{}) {
	message += fmt.Sprintf(" at %s line %d", filename, lineno)
	rs.warn(message, a...)
}

/*
DumpTopics is a debug method which pretty-prints the topic tree structure from
the bot's memory.
*/
func (rs *RiveScript) DumpTopics() {
	for topic, data := range rs.topics {
		fmt.Printf("Topic: %s\n", topic)
		for _, trigger := range data.triggers {
			fmt.Printf("  + %s\n", trigger.trigger)
			if trigger.previous != "" {
				fmt.Printf("    %% %s\n", trigger.previous)
			}
			for _, cond := range trigger.condition {
				fmt.Printf("    * %s\n", cond)
			}
			for _, reply := range trigger.reply {
				fmt.Printf("    - %s\n", reply)
			}
			if trigger.redirect != "" {
				fmt.Printf("    @ %s\n", trigger.redirect)
			}
		}
	}
}

/*
DumpSorted is a debug method which pretty-prints the sort tree of topics from
the bot's memory.
*/
func (rs *RiveScript) DumpSorted() {
	rs._dumpSorted(rs.sorted.topics, "Topics")
	rs._dumpSorted(rs.sorted.thats, "Thats")
	rs._dumpSortedList(rs.sorted.sub, "Substitutions")
	rs._dumpSortedList(rs.sorted.person, "Person Substitutions")
}
func (rs *RiveScript) _dumpSorted(tree map[string][]sortedTriggerEntry, label string) {
	fmt.Printf("Sort Buffer: %s\n", label)
	for topic, data := range tree {
		fmt.Printf("  Topic: %s\n", topic)
		for _, trigger := range data {
			fmt.Printf("    + %s\n", trigger.trigger)
		}
	}
}
func (rs *RiveScript) _dumpSortedList(list []string, label string) {
	fmt.Printf("Sort buffer: %s\n", label)
	for _, item := range list {
		fmt.Printf("  %s\n", item)
	}
}
