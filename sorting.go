package rivescript

// Data sorting functions

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

// Sort buffer data, for RiveScript.SortReplies()
type sortBuffer struct {
	topics map[string][]sortedTriggerEntry // Topic name -> array of triggers
	thats  map[string][]sortedTriggerEntry
	sub    []string // Substitutions
	person []string // Person substitutions
}

// Holds a sorted trigger and the pointer to that trigger's data
type sortedTriggerEntry struct {
	trigger string
	pointer *astTrigger
}

// Temporary categorization of triggers while sorting
type sortTrack struct {
	atomic map[int][]sortedTriggerEntry // Sort by number of whole words
	option map[int][]sortedTriggerEntry // Sort optionals by number of words
	alpha  map[int][]sortedTriggerEntry // Sort alpha wildcards by no. of words
	number map[int][]sortedTriggerEntry // Sort numeric wildcards by no. of words
	wild   map[int][]sortedTriggerEntry // Sort wildcards by no. of words
	pound  []sortedTriggerEntry         // Triggers of just '#'
	under  []sortedTriggerEntry         // Triggers of just '_'
	star   []sortedTriggerEntry         // Triggers of just '*'
}

/*
SortReplies sorts the reply structures in memory for optimal matching.

After you have finished loading your RiveScript code, call this method to
populate the various sort buffers. This is absolutely necessary for reply
matching to work efficiently!

If the bot has loaded no topics, or if it ends up with no sorted triggers at
the end, it will return an error saying such. This usually means the bot didn't
load any RiveScript code, for example because it looked in the wrong directory.
*/
func (rs *RiveScript) SortReplies() error {
	// (Re)initialize the sort cache.
	rs.sorted.topics = map[string][]sortedTriggerEntry{}
	rs.sorted.thats = map[string][]sortedTriggerEntry{}
	rs.say("Sorting triggers...")

	// If there are no topics, give an error.
	if len(rs.topics) == 0 {
		return errors.New("SortReplies: no topics were found; did you load any RiveScript code?")
	}

	// Loop through all the topics.
	for topic := range rs.topics {
		rs.say("Analyzing topic %s", topic)

		// Collect a list of all the triggers we're going to worry about. If this
		// topic inherits another topic, we need to recursively add those to the
		// list as well.
		allTriggers := rs.getTopicTriggers(topic, false)

		// Sort these triggers.
		rs.sorted.topics[topic] = rs.sortTriggerSet(allTriggers, true)

		// Get all of the %Previous triggers for this topic.
		thatTriggers := rs.getTopicTriggers(topic, true)

		// And sort them, too.
		rs.sorted.thats[topic] = rs.sortTriggerSet(thatTriggers, false)
	}

	// Sort the substitution lists.
	rs.sorted.sub = sortList(rs.sub)
	rs.sorted.person = sortList(rs.person)

	// Did we sort anything at all?
	if len(rs.sorted.topics) == 0 && len(rs.sorted.thats) == 0 {
		return errors.New("SortReplies: ended up with empty trigger lists; did you load any RiveScript code?")
	}

	return nil
}

/*
sortTriggerSet sorts a group of triggers in an optimal sorting order.

This function has two use cases:

1. Create a sort buffer for "normal" (matchable) triggers, which are triggers
   that are NOT accompanied by a %Previous tag.
2. Create a sort buffer for triggers that had %Previous tags.

Use the `excludePrevious` parameter to control which one is being done. This
function will return a list of sortedTriggerEntry items, and it's intended to
have no duplicate trigger patterns (unless the source RiveScript code explicitly
uses the same duplicate pattern twice, which is a user error).
*/
func (rs *RiveScript) sortTriggerSet(triggers []sortedTriggerEntry, excludePrevious bool) []sortedTriggerEntry {
	// Create a priority map, of priority numbers -> their triggers.
	prior := map[int][]sortedTriggerEntry{}

	// Go through and bucket each trigger by weight (priority).
	for _, trig := range triggers {
		if excludePrevious && trig.pointer.previous != "" {
			continue
		}

		// Check the trigger text for any {weight} tags, default being 0
		match := reWeight.FindStringSubmatch(trig.trigger)
		weight := 0
		if len(match) > 0 {
			weight, _ = strconv.Atoi(match[1])
		}

		// First trigger of this priority? Initialize the weight map.
		if _, ok := prior[weight]; !ok {
			prior[weight] = []sortedTriggerEntry{}
		}

		prior[weight] = append(prior[weight], trig)
	}

	// Keep a running list of sorted triggers for this topic.
	running := []sortedTriggerEntry{}

	// Sort the priorities with the highest number first.
	var sortedPriorities []int
	for k := range prior {
		sortedPriorities = append(sortedPriorities, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedPriorities)))

	// Go through each priority set.
	for _, p := range sortedPriorities {
		rs.say("Sorting triggers with priority %d", p)

		// So, some of these triggers may include an {inherits} tag, if they
		// came from a topic which inherits another topic. Lower inherits values
		// mean higher priority on the stack. Triggers that have NO inherits
		// value at all (which will default to -1), will be moved to the END of
		// the stack at the end (have the highest number/lowest priority).
		inherits := -1        // -1 means no {inherits} tag
		highestInherits := -1 // Highest number seen so far

		// Loop through and categorize these triggers.
		track := map[int]*sortTrack{}
		track[inherits] = initSortTrack()

		// Loop through all the triggers.
		for _, trig := range prior[p] {
			pattern := trig.trigger
			rs.say("Looking at trigger: %s", pattern)

			// See if the trigger has an {inherits} tag.
			match := reInherits.FindStringSubmatch(pattern)
			if len(match) > 0 {
				inherits, _ = strconv.Atoi(match[1])
				if inherits > highestInherits {
					highestInherits = inherits
				}
				rs.say("Trigger belongs to a topic that inherits other topics. "+
					"Level=%d", inherits)
				pattern = reInherits.ReplaceAllString(pattern, "")
			} else {
				inherits = -1
			}

			// If this is the first time we've seen this inheritance level,
			// initialize its sort track structure.
			if _, ok := track[inherits]; !ok {
				track[inherits] = initSortTrack()
			}

			// Start inspecting the trigger's contents.
			if strings.Index(pattern, "_") > -1 {
				// Alphabetic wildcard included.
				cnt := wordCount(pattern, false)
				rs.say("Has a _ wildcard with %d words", cnt)
				if cnt > 0 {
					if _, ok := track[inherits].alpha[cnt]; !ok {
						track[inherits].alpha[cnt] = []sortedTriggerEntry{}
					}
					track[inherits].alpha[cnt] = append(track[inherits].alpha[cnt], trig)
				} else {
					track[inherits].under = append(track[inherits].under, trig)
				}
			} else if strings.Index(pattern, "#") > -1 {
				// Numeric wildcard included.
				cnt := wordCount(pattern, false)
				rs.say("Has a # wildcard with %d words", cnt)
				if cnt > 0 {
					if _, ok := track[inherits].number[cnt]; !ok {
						track[inherits].number[cnt] = []sortedTriggerEntry{}
					}
					track[inherits].number[cnt] = append(track[inherits].number[cnt], trig)
				} else {
					track[inherits].pound = append(track[inherits].pound, trig)
				}
			} else if strings.Index(pattern, "*") > -1 {
				// Wildcard included.
				cnt := wordCount(pattern, false)
				rs.say("Has a * wildcard with %d words", cnt)
				if cnt > 0 {
					if _, ok := track[inherits].wild[cnt]; !ok {
						track[inherits].wild[cnt] = []sortedTriggerEntry{}
					}
					track[inherits].wild[cnt] = append(track[inherits].wild[cnt], trig)
				} else {
					track[inherits].star = append(track[inherits].star, trig)
				}
			} else if strings.Index(pattern, "[") > -1 {
				// Optionals included.
				cnt := wordCount(pattern, false)
				rs.say("Has optionals with %d words", cnt)
				if _, ok := track[inherits].option[cnt]; !ok {
					track[inherits].option[cnt] = []sortedTriggerEntry{}
				}
				track[inherits].option[cnt] = append(track[inherits].option[cnt], trig)
			} else {
				// Totally atomic.
				cnt := wordCount(pattern, false)
				rs.say("Totally atomic trigger with %d words", cnt)
				if _, ok := track[inherits].atomic[cnt]; !ok {
					track[inherits].atomic[cnt] = []sortedTriggerEntry{}
				}
				track[inherits].atomic[cnt] = append(track[inherits].atomic[cnt], trig)
			}
		}

		// Move the no-{inherits} triggers to the bottom of the stack.
		track[highestInherits+1] = track[-1]
		delete(track, -1)

		// Sort the track from the lowest to the highest.
		var trackSorted []int
		for k := range track {
			trackSorted = append(trackSorted, k)
		}
		sort.Ints(trackSorted)

		// Go through each priority level from greatest to smallest.
		for _, ip := range trackSorted {
			rs.say("ip=%d", ip)

			// Sort each of the main kinds of triggers by their word counts.
			running = sortByWords(running, track[ip].atomic)
			running = sortByWords(running, track[ip].option)
			running = sortByWords(running, track[ip].alpha)
			running = sortByWords(running, track[ip].number)
			running = sortByWords(running, track[ip].wild)

			// Add the single wildcard triggers, sorted by length.
			running = sortByLength(running, track[ip].under)
			running = sortByLength(running, track[ip].pound)
			running = sortByLength(running, track[ip].star)
		}
	}

	return running
}

// sortList sorts lists (like substitutions) from a string:string map.
func sortList(dict map[string]string) []string {
	output := []string{}

	// Track by number of words.
	track := map[int][]string{}

	// Loop through each item.
	for item := range dict {
		cnt := wordCount(item, true)
		if _, ok := track[cnt]; !ok {
			track[cnt] = []string{}
		}
		track[cnt] = append(track[cnt], item)
	}

	// Sort them by word count, descending.
	sortedCounts := []int{}
	for cnt := range track {
		sortedCounts = append(sortedCounts, cnt)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedCounts)))

	for _, cnt := range sortedCounts {
		// Sort the strings of this word-count by their lengths.
		sortedLengths := track[cnt]
		sort.Sort(sort.Reverse(byLength(sortedLengths)))
		for _, item := range sortedLengths {
			output = append(output, item)
		}
	}

	return output
}

/*
sortByWords sorts a set of triggers by word count and overall length.

This is a helper function for sorting the `atomic`, `option`, `alpha`, `number`
and `wild` attributes of the sortTrack and adding them to the running sort
buffer in that specific order. Since attribute lookup by reflection is expensive
in Go, this function is given the relevant sort buffer directly, and the current
running sort buffer to add the results to.

The `triggers` parameter is a map between word counts and the triggers that
fit that number of words.
*/
func sortByWords(running []sortedTriggerEntry, triggers map[int][]sortedTriggerEntry) []sortedTriggerEntry {
	// Sort the triggers by their word counts from greatest to smallest.
	var sortedWords []int
	for wc := range triggers {
		sortedWords = append(sortedWords, wc)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedWords)))

	for _, wc := range sortedWords {
		// Triggers with equal word lengths should be sorted by overall trigger length.
		var sortedPatterns []string
		patternMap := map[string][]sortedTriggerEntry{}

		for _, trig := range triggers[wc] {
			sortedPatterns = append(sortedPatterns, trig.trigger)
			if _, ok := patternMap[trig.trigger]; !ok {
				patternMap[trig.trigger] = []sortedTriggerEntry{}
			}
			patternMap[trig.trigger] = append(patternMap[trig.trigger], trig)
		}
		sort.Sort(sort.Reverse(byLength(sortedPatterns)))

		// Add the triggers to the running triggers bucket.
		for _, pattern := range sortedPatterns {
			running = append(running, patternMap[pattern]...)
		}
	}

	return running
}

/*
sortByLength sorts a set of triggers purely by character length.

This is like `sortByWords`, but it's intended for triggers that consist solely
of wildcard-like symbols with no real words. For example a trigger of `* * *`
qualifies for this, and it has no words, so we sort by length so it gets a
higher priority than simply `*`.
*/
func sortByLength(running []sortedTriggerEntry, triggers []sortedTriggerEntry) []sortedTriggerEntry {
	var sortedPatterns []string
	patternMap := map[string][]sortedTriggerEntry{}
	for _, trig := range triggers {
		sortedPatterns = append(sortedPatterns, trig.trigger)
		if _, ok := patternMap[trig.trigger]; !ok {
			patternMap[trig.trigger] = []sortedTriggerEntry{}
		}
		patternMap[trig.trigger] = append(patternMap[trig.trigger], trig)
	}
	sort.Sort(sort.Reverse(byLength(sortedPatterns)))

	// Only loop through unique patterns.
	patternSet := map[string]bool{}

	// Add them to the running triggers bucket.
	for _, pattern := range sortedPatterns {
		if _, ok := patternSet[pattern]; ok {
			continue
		}
		patternSet[pattern] = true
		running = append(running, patternMap[pattern]...)
	}

	return running
}

// initSortTrack initializes a new, empty sortTrack object.
func initSortTrack() *sortTrack {
	return &sortTrack{
		atomic: map[int][]sortedTriggerEntry{},
		option: map[int][]sortedTriggerEntry{},
		alpha:  map[int][]sortedTriggerEntry{},
		number: map[int][]sortedTriggerEntry{},
		wild:   map[int][]sortedTriggerEntry{},
		pound:  []sortedTriggerEntry{},
		under:  []sortedTriggerEntry{},
		star:   []sortedTriggerEntry{},
	}
}
