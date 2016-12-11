// Package config provides the RiveScript configuration type.
package config

import (
	"github.com/aichaos/rivescript-go/sessions"
	"github.com/aichaos/rivescript-go/sessions/memory"
)

// Type Config configures a RiveScript instance.
type Config struct {
	// Debug enables verbose debug logging to your standard output.
	Debug bool

	// Strict enables strict syntax checking.
	Strict bool

	// Depth sets the recursion depth limit. The zero value will default to
	// 50 levels deep.
	Depth uint

	// UTF8 enables UTF-8 support for user messages and triggers.
	UTF8 bool

	// SessionManager chooses a session manager for user variables.
	SessionManager sessions.SessionManager
}

// Basic creates a default configuration:
//
// - Strict: true
// - Depth: 50
// - UTF8: false
func Basic() *Config {
	return &Config{
		Strict:         true,
		Depth:          50,
		UTF8:           false,
		SessionManager: memory.New(),
	}
}

// UTF8 creates a default configuration with UTF-8 mode enabled.
//
// - Strict: true
// - Depth: 50
// - UTF8: true
func UTF8() *Config {
	return &Config{
		Strict:         true,
		Depth:          50,
		UTF8:           true,
		SessionManager: memory.New(),
	}
}
