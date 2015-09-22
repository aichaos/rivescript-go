package rivescript

// Miscellaneous structures

// User data, key/value pairs about the user.
type userData struct {
	data         map[string]string
	lastMatch    string
	inputHistory []string
	replyHistory []string
}

// newUser creates a new user profile.
func newUser() *userData {
	user := new(userData)
	user.data = map[string]string{}
	user.data["topic"] = "random"

	user.lastMatch = ""

	user.inputHistory = []string{
		"undefined", "undefined", "undefined", "undefined", "undefined",
		"undefined", "undefined", "undefined", "undefined",
	}
	user.replyHistory = []string{
		"undefined", "undefined", "undefined", "undefined", "undefined",
		"undefined", "undefined", "undefined", "undefined",
	}
	return user
}

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

// This is like astTopic but is just for %Previous mapping
type thatTopic struct {
	triggers map[string]*thatTrigger
}

type thatTrigger struct {
	previous map[string]*astTrigger
}
