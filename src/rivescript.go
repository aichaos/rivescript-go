/*
Package rivescript contains all of the private use functions of RiveScript.

Do Not Use This Package Directly

You should not use any exported symbols from this package directly. They are
not stable and are subject to change at any time without notice.

As an end user of the RiveScript library you should stick purely to the exported
API functions of the base RiveScript package and any other subpackages
(for example: parser and ast) but leave the src package alone!

Documentation for most exported functions is available in the root RiveScript
package, which acts as a wrapper. Go there for documentation. Stop looking at
this package lest you be tempted to use it (don't).

You've been warned. Here be dragons.
*/
package rivescript

import (
	"regexp"
	"sync"

	"github.com/aichaos/rivescript-go/macro"
	"github.com/aichaos/rivescript-go/parser"
	"github.com/aichaos/rivescript-go/sessions"
	"github.com/aichaos/rivescript-go/sessions/memory"
)

// RiveScript is the bot instance.
type RiveScript struct {
	// Parameters
	Debug              bool // Debug mode
	Strict             bool // Strictly enforce RiveScript syntax
	Depth              uint // Max depth for recursion
	UTF8               bool // Support UTF-8 RiveScript code
	Quiet              bool // Suppress all warnings from being emitted
	UnicodePunctuation *regexp.Regexp

	// Internal helpers
	parser *parser.Parser

	// Internal data structures
	cLock       sync.Mutex                      // Lock for config variables.
	global      map[string]string               // 'global' variables
	vars        map[string]string               // 'var' bot variables
	sub         map[string]string               // 'sub' substitutions
	person      map[string]string               // 'person' substitutions
	array       map[string][]string             // 'array'
	sessions    sessions.SessionManager         // user variable session manager
	includes    map[string]map[string]bool      // included topics
	inherits    map[string]map[string]bool      // inherited topics
	objlangs    map[string]string               // object macro languages
	handlers    map[string]macro.MacroInterface // object language handlers
	subroutines map[string]Subroutine           // Golang object handlers
	topics      map[string]*astTopic            // main topic structure
	sorted      *sortBuffer                     // Sorted data from SortReplies()

	// State information.
	inReplyContext bool
	currentUser    string
}

/******************************************************************************
 * Constructor and Debug Methods                                              *
 ******************************************************************************/

// New creates a new RiveScript instance with the default configuration.
func New() *RiveScript {
	rs := &RiveScript{
		// Set the default config objects that don't have good zero-values.
		Strict:   true,
		Depth:    50,
		sessions: memory.New(),

		// Default punctuation that gets removed from messages in UTF-8 mode.
		UnicodePunctuation: regexp.MustCompile(`[.,!?;:]`),

		// Initialize all internal data structures.
		global:      map[string]string{},
		vars:        map[string]string{},
		sub:         map[string]string{},
		person:      map[string]string{},
		array:       map[string][]string{},
		includes:    map[string]map[string]bool{},
		inherits:    map[string]map[string]bool{},
		objlangs:    map[string]string{},
		handlers:    map[string]macro.MacroInterface{},
		subroutines: map[string]Subroutine{},
		topics:      map[string]*astTopic{},
		sorted:      new(sortBuffer),
	}

	// Helpers.
	rs.parser = parser.New(parser.ParserConfig{
		Strict:  true,
		OnDebug: rs.say,
		OnWarn:  rs.warnSyntax,
	})

	return rs
}

// Configure is a convenience function for the public API to set all of its
// settings at once.
func (rs *RiveScript) Configure(debug, strict, utf8 bool, depth uint,
	sessions sessions.SessionManager) {
	rs.Debug = debug
	rs.Strict = strict
	rs.UTF8 = utf8
	rs.Depth = depth
	rs.sessions = sessions

	// Reconfigure the parser.
	rs.parser = parser.New(parser.ParserConfig{
		Strict:  strict,
		UTF8:    utf8,
		OnDebug: rs.say,
		OnWarn:  rs.warnSyntax,
	})
}

// SetUnicodePunctuation allows for overriding the regexp for punctuation.
func (rs *RiveScript) SetUnicodePunctuation(value string) {
	rs.UnicodePunctuation = regexp.MustCompile(value)
}
