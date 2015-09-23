package rivescript

// Loading and Parsing Methods

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/*
LoadFile loads a single RiveScript source file from disk.

Params:

	path: File path to
*/
func (rs *RiveScript) LoadFile(path string) error {
	rs.say("Load RiveScript file: %s", path)

	fh, err := os.Open(path)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open file %s: %s", path, err))
	}

	defer fh.Close()
	scanner := bufio.NewScanner(fh)
	scanner.Split(bufio.ScanLines)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return rs.parse(path, lines)
}

/*
LoadDirectory loads multiple RiveScript documents from a folder on disk.

Params:

	path: Path to the directory on disk
	extensions...: List of file extensions to filter on, default is '.rive' and '.rs'
*/
func (rs *RiveScript) LoadDirectory(path string, extensions ...string) error {
	if len(extensions) == 0 {
		extensions = []string{".rive", ".rs"}
	}

	files, err := filepath.Glob(fmt.Sprintf("%s/*", path))
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to open folder %s: %s", path, err))
	}

	for _, f := range files {
		// Restrict file extensions.
		validExtension := false
		for _, exten := range extensions {
			if strings.HasSuffix(f, exten) {
				validExtension = true
				break
			}
		}

		if validExtension {
			err := rs.LoadFile(f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

/*
Stream loads RiveScript code from a text buffer.

Params:
	code: Raw source code of a RiveScript document, with line breaks after each line.
*/
func (rs *RiveScript) Stream(code string) error {
	lines := strings.Split(code, "\n")
	return rs.parse("Stream()", lines)
}

/*
parse loads the RiveScript code into the bot's memory.
*/
func (rs *RiveScript) parse(path string, lines []string) error {
	rs.say("Parsing code!")

	// Get the "abstract syntax tree" of this file.
	ast, err := rs.parseSource(path, lines)
	if err != nil {
		return err
	}

	// Get all of the "begin" type variables
	for k, v := range ast.begin.global {
		if v == "<undef>" {
			delete(rs.global, k)
		} else {
			rs.global[k] = v
		}
	}
	for k, v := range ast.begin.var_ {
		if v == "<undef>" {
			delete(rs.var_, k)
		} else {
			rs.var_[k] = v
		}
	}
	for k, v := range ast.begin.sub {
		if v == "<undef>" {
			delete(rs.sub, k)
		} else {
			rs.sub[k] = v
		}
	}
	for k, v := range ast.begin.person {
		if v == "<undef>" {
			delete(rs.person, k)
		} else {
			rs.person[k] = v
		}
	}
	for k, v := range ast.begin.array {
		rs.array[k] = v
	}

	// Consume all the parsed triggers.
	for topic, data := range ast.topics {
		// Keep a map of the topics that are included/inherited under this topic.
		if _, ok := rs.includes[topic]; !ok {
			rs.includes[topic] = map[string]bool{}
		}
		if _, ok := rs.inherits[topic]; !ok {
			rs.inherits[topic] = map[string]bool{}
		}

		// Merge in the topic inclusions/inherits.
		for included, _ := range data.includes {
			rs.includes[topic][included] = true
		}
		for inherited, _ := range data.inherits {
			rs.inherits[topic][inherited] = true
		}

		// Initialize the topic structure.
		if _, ok := rs.topics[topic]; !ok {
			rs.topics[topic] = new(astTopic)
			rs.topics[topic].triggers = []*astTrigger{}
		}

		// Consume the AST triggers into the brain.
		for _, trigger := range data.triggers {
			rs.topics[topic].triggers = append(rs.topics[topic].triggers, trigger)

			// Does this one have a %Previous? If so, make a pointer to this
			// exact trigger in rs.thats
			if trigger.previous != "" {
				// Initialize the structure first.
				if _, ok := rs.thats[topic]; !ok {
					rs.thats[topic] = new(thatTopic)
					rs.say("%q", rs.thats[topic])
					rs.thats[topic].triggers = map[string]*thatTrigger{}
				}
				if _, ok := rs.thats[topic].triggers[trigger.trigger]; !ok {
					rs.say("%q", rs.thats[topic].triggers[trigger.trigger])
					rs.thats[topic].triggers[trigger.trigger] = new(thatTrigger)
					rs.thats[topic].triggers[trigger.trigger].previous = map[string]*astTrigger{}
				}
				rs.thats[topic].triggers[trigger.trigger].previous[trigger.previous] = trigger
			}
		}
	}

	// Load all the parsed objects.
	for _, object := range ast.objects {
		// Have a language handler for this?
		if _, ok := rs.handlers[object.language]; ok {
			rs.say("Loading object macro %s (%s)", object.name, object.language)
			rs.handlers[object.language].Load(object.name, object.code)
			rs.objlangs[object.name] = object.language
		}
	}

	return nil
}

/*
SortReplies sorts the reply structures in memory for optimal matching.

After you have finished loading your RiveScript code, call this method to
populate the various sort buffers. This is absolutely necessary for reply
matching to work efficiently!
*/
func (rs *RiveScript) SortReplies() {
	// (Re)initialize the sort cache.
	rs.sorted.topics = map[string][]sortedTriggerEntry{}
	rs.sorted.thats = map[string][]sortedTriggerEntry{}
	rs.say("Sorting triggers...")

	// Loop through all the topics.
	for topic, _ := range rs.topics {
		rs.say("Analyzing topic %s", topic)

		// Collect a list of all the triggers we're going to worry about. If this
		// topic inherits another topic, we need to recursively add those to the
		// list as well.
		allTriggers := rs.getTopicTriggers(topic, rs.topics, nil)

		// Sort these triggers.
		rs.sorted.topics[topic] = rs.sortTriggerSet(allTriggers, true)

		// Get all of the %Previous triggers for this topic.
		thatTriggers := rs.getTopicTriggers(topic, nil, rs.thats)

		// And sort them, too.
		rs.sorted.thats[topic] = rs.sortTriggerSet(thatTriggers, false)
	}

	// Sort the substitution lists.
	rs.sorted.sub = sortList(rs.sub)
	rs.sorted.person = sortList(rs.person)
}
