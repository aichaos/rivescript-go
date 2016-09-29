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
package src

import (
	"github.com/aichaos/rivescript-go/macro"
	"github.com/aichaos/rivescript-go/parser"
	"regexp"
)

type RiveScript struct {
	// Parameters
	Debug              bool // Debug mode
	Strict             bool // Strictly enforce RiveScript syntax
	Depth              int  // Max depth for recursion
	UTF8               bool // Support UTF-8 RiveScript code
	UnicodePunctuation *regexp.Regexp

	// Internal helpers
	parser *parser.Parser

	// Internal data structures
	global      map[string]string               // 'global' variables
	var_        map[string]string               // 'var' bot variables
	sub         map[string]string               // 'sub' substitutions
	person      map[string]string               // 'person' substitutions
	array       map[string][]string             // 'array'
	users       map[string]*userData            // user variables
	freeze      map[string]*userData            // frozen user variables
	includes    map[string]map[string]bool      // included topics
	inherits    map[string]map[string]bool      // inherited topics
	objlangs    map[string]string               // object macro languages
	handlers    map[string]macro.MacroInterface // object language handlers
	subroutines map[string]Subroutine           // Golang object handlers
	topics      map[string]*astTopic            // main topic structure
	thats       map[string]*thatTopic           // %Previous mapper
	sorted      *sortBuffer                     // Sorted data from SortReplies()

	// State information.
	currentUser string
}

/******************************************************************************
 * Constructor and Debug Methods                                              *
 ******************************************************************************/

func New() *RiveScript {
	rs := new(RiveScript)
	rs.Debug = false
	rs.Strict = true
	rs.Depth = 50
	rs.UTF8 = false
	rs.UnicodePunctuation = regexp.MustCompile(`[.,!?;:]`)

	// Initialize helpers.
	rs.parser = parser.New(parser.ParserConfig{
		Strict:  rs.Strict,
		UTF8:    rs.UTF8,
		OnDebug: rs.say,
		OnWarn:  rs.warnSyntax,
	})

	// Initialize all the data structures.
	rs.global = map[string]string{}
	rs.var_ = map[string]string{}
	rs.sub = map[string]string{}
	rs.person = map[string]string{}
	rs.array = map[string][]string{}
	rs.users = map[string]*userData{}
	rs.freeze = map[string]*userData{}
	rs.includes = map[string]map[string]bool{}
	rs.inherits = map[string]map[string]bool{}
	rs.objlangs = map[string]string{}
	rs.handlers = map[string]macro.MacroInterface{}
	rs.subroutines = map[string]Subroutine{}
	rs.topics = map[string]*astTopic{}
	rs.thats = map[string]*thatTopic{}
	rs.sorted = new(sortBuffer)

	// Initialize Golang handler.
	//rs.handlers["go"] = new(golangHandler)
	return rs
}

func (rs *RiveScript) SetDebug(value bool) {
	rs.Debug = value
}

func (rs *RiveScript) SetUTF8(value bool) {
	rs.UTF8 = value
}

func (rs *RiveScript) SetStrict(value bool) {
	rs.Strict = value
}

func (rs *RiveScript) SetDepth(value int) {
	rs.Depth = value
}

func (rs *RiveScript) SetUnicodePunctuation(value string) {
	rs.UnicodePunctuation = regexp.MustCompile(value)
}
