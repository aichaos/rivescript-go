package rivescript

/*
	NOTE: This module is a wrapper around the bulk of the actual source code
	under the 'src/' subpackage. This gives multiple benefits such as keeping
	the root of the git repo as tidy as possible (low number of source files)
	and keeping the public facing, official API in one small place in the code.

	Everything exported from the 'src' subpackage should not be used directly
	by third party developers. A lot of the symbols from the src package must
	be exported to get this wrapper program to work (and keep a nice looking
	module import path), but only this public facing API module should be used.
*/

import (
	"math/rand"
	"regexp"
	"sync"
	"time"

	"github.com/aichaos/rivescript-go/macro"
	"github.com/aichaos/rivescript-go/parser"
	"github.com/aichaos/rivescript-go/sessions"
	"github.com/aichaos/rivescript-go/sessions/memory"
)

// Version number for the RiveScript library.
const Version = "0.3.0"

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

	// The random number god.
	random     rand.Source
	rng        *rand.Rand
	randomLock sync.Mutex

	// State information.
	inReplyContext bool
	currentUser    string
}

/*
New creates a new RiveScript instance.

A RiveScript instance represents one chat bot personality; it has its own
replies and its own memory of user data. You could make multiple bots in the
same program, each with its own replies loaded from different sources.
*/
func New(cfg *Config) *RiveScript {
	// If no config was given, default to the BasicConfig.
	if cfg == nil {
		cfg = &Config{
			Strict: true,
			Depth:  50,
		}
	}

	// Sensible default config options.
	if cfg.Depth == 0 {
		cfg.Depth = 50
	}
	if cfg.SessionManager == nil {
		cfg.SessionManager = memory.New()
	}

	// Random number seed.
	var random rand.Source
	if cfg.Seed != 0 {
		random = rand.NewSource(cfg.Seed)
	} else {
		random = rand.NewSource(time.Now().UnixNano())
	}

	rs := &RiveScript{
		// Set the default config objects that don't have good zero-values.
		Debug:    cfg.Debug,
		Strict:   cfg.Strict,
		Depth:    cfg.Depth,
		UTF8:     cfg.UTF8,
		sessions: cfg.SessionManager,

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

		random: random,
		rng:    rand.New(random),
	}

	// Helper modules.
	rs.parser = parser.New(parser.ParserConfig{
		Strict:  cfg.Strict,
		UTF8:    cfg.UTF8,
		OnDebug: rs.say,
		OnWarn:  rs.warnSyntax,
	})

	return rs
}

// Forms of undefined.
const (
	// UNDEFINED is the text "undefined", the default text for variable getters.
	UNDEFINED = "undefined"

	// UNDEFTAG is the "<undef>" tag for unsetting variables in !Definitions.
	UNDEFTAG = "<undef>"
)

// Subroutine is a function prototype for defining custom object macros in Go.
type Subroutine func(*RiveScript, []string) string

// SetUnicodePunctuation allows you to override the text of the unicode
// punctuation regexp. Provide a string literal that will validate in
// `regexp.MustCompile()`
func (rs *RiveScript) SetUnicodePunctuation(value string) {
	rs.UnicodePunctuation = regexp.MustCompile(value)
}
