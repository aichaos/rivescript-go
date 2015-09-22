/*
Package rivescript implements the RiveScript chatbot scripting language.
*/
package rivescript

import (
	"fmt"
	"regexp"
)

// Constants
const RS_VERSION float64 = 2.0

/******************************************************************************
 * Constructor and Debug Methods                                              *
 ******************************************************************************/

type RiveScript struct {
	// Parameters
	Debug  bool // Debug mode
	Strict bool // Strictly enforce RiveScript syntax
	Depth  int  // Max depth for recursion
	UTF8   bool // Support UTF-8 RiveScript code
	UnicodePunctuation *regexp.Regexp

	// Internal data structures
	global   map[string]string          // 'global' variables
	var_     map[string]string          // 'var' bot variables
	sub      map[string]string          // 'sub' substitutions
	person   map[string]string          // 'person' substitutions
	array    map[string][]string        // 'array'
	users    map[string]*UserData       // user variables
	freeze   map[string]*UserData       // frozen user variables
	includes map[string]map[string]bool // included topics
	inherits map[string]map[string]bool // inherited topics
	objlangs map[string]string          // object macro languages
	handlers map[string]*MacroHandler  // object language handlers
	topics   map[string]*astTopic       // main topic structure
	thats    map[string]*thatTopic      // %Previous mapper
	sorted   *sortBuffer                // Sorted data from SortReplies()

	// State information.
	currentUser string
}

func New() *RiveScript {
	rs := new(RiveScript)
	rs.Debug = false
	rs.Strict = true
	rs.Depth = 50
	rs.UTF8 = false
	rs.UnicodePunctuation = regexp.MustCompile(`[.,!?;:]`)

	// Initialize all the data structures.
	rs.global = map[string]string{}
	rs.var_ = map[string]string{}
	rs.sub = map[string]string{}
	rs.person = map[string]string{}
	rs.array = map[string][]string{}
	rs.users = map[string]*UserData{}
	rs.freeze = map[string]*UserData{}
	rs.includes = map[string]map[string]bool{}
	rs.inherits = map[string]map[string]bool{}
	rs.handlers = map[string]*MacroHandler{}
	rs.topics = map[string]*astTopic{}
	rs.thats = map[string]*thatTopic{}
	rs.sorted = new(sortBuffer)

	// Initialize Golang handler.
	//rs.handlers["go"] = new(GolangHandler)
	return rs
}

func (rs RiveScript) Version() string {
	// TODO: versioning
	return "0.0.1"
}

// say prints a debugging message
func (rs RiveScript) say(message string, a ...interface{}) {
	if rs.Debug {
		fmt.Printf(message+"\n", a...)
	}
}

// warn prints a warning message for non-fatal errors
func (rs RiveScript) warn(message string, a ...interface{}) {
	fmt.Printf("[WARN] "+message+"\n", a...)
}

// syntax is like warn but takes a filename and line number.
func (rs RiveScript) warnSyntax(message string, filename string, lineno int, a ...interface{}) {
	message += fmt.Sprintf(" at %s line %d", filename, lineno)
	rs.warn(message, a...)
}
