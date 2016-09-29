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
