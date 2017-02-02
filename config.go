package rivescript

import "github.com/aichaos/rivescript-go/sessions"

/*
Type Config provides options to configure the RiveScript bot.

Create a pointer to this type and send it to the New() constructor to change
the default settings. You only need to provide settings you want to override;
the zero-values of all the options are handled appropriately by the RiveScript
library.

The default values are documented below.
*/
type Config struct {
	// Debug enables verbose logging to standard output. Default false.
	Debug bool

	// Strict enables strict syntax checking, where a syntax error in RiveScript
	// code is considered fatal at parse time. Default true.
	Strict bool

	// UTF8 enables UTF-8 mode within the bot. Default false.
	//
	// When UTF-8 mode is enabled, triggers in the RiveScript source files are
	// allowed to contain foreign characters. Additionally, the user's incoming
	// messages are left *mostly* intact, so that they send messages with
	// foreign characters to the bot.
	UTF8 bool

	// Depth controls the global limit for recursive functions within
	// RiveScript. Default 50.
	Depth uint

	// SessionManager is an implementation of the same name for managing user
	// variables for the bot. The default is the in-memory session handler.
	SessionManager sessions.SessionManager
}

// WithUTF8 provides a Config object that enables UTF-8 mode.
func WithUTF8() *Config {
	return &Config{
		UTF8: true,
	}
}
