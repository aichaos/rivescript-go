package rivescript

import (
	"testing"

	"github.com/aichaos/rivescript-go/sessions/memory"
)

func TestConfigAPI(t *testing.T) {
	// A helper function to handleVars error reporting when we do (or don't) expect errors.
	handleVars := func(fn func(string) (string, error), param string, expected bool) {
		_, err := fn(param)
		if (err == nil && expected) || (err != nil && !expected) {
			t.Errorf(
				"Expected errors for var %s: %v; but we got as the error: [%s]",
				param,
				expected,
				err,
			)
		}
	}

	// A helper to assert a substitution exists.
	assertIn := func(dict map[string]string, key string, expected bool) {
		_, ok := dict[key]
		if ok != expected {
			t.Errorf(
				"Expected key %s to exist: %v; but we got %v",
				key, expected, ok,
			)
		}
	}

	bot := NewTest(t)
	bot.extend(`
		+ hello go
		- <call>hello-go</call>
	`)
	rs := bot.bot

	// Setting a Go object macro.
	rs.SetSubroutine("hello-go", func(rs *RiveScript, args []string) string {
		return "Hello world"
	})
	bot.reply("Hello go", "Hello world")

	// Deleting it.
	rs.DeleteSubroutine("hello-go")
	bot.reply("Hello go", "[ERR: Object Not Found]")

	// Global variables.
	handleVars(rs.GetGlobal, "global test", true)
	rs.SetGlobal("global test", "on")
	handleVars(rs.GetGlobal, "global test", false)
	rs.SetGlobal("global test", "undefined")
	handleVars(rs.GetGlobal, "global test", true)

	// Bot variables.
	handleVars(rs.GetVariable, "var test", true)
	rs.SetVariable("var test", "on")
	handleVars(rs.GetVariable, "var test", false)
	rs.SetVariable("var test", "undefined")
	handleVars(rs.GetVariable, "var test", true)

	// Substitutions.
	assertIn(rs.sub, "what's", false)
	rs.SetSubstitution("what's", "what is")
	assertIn(rs.sub, "what's", true)
	rs.SetSubstitution("what's", "undefined")
	assertIn(rs.sub, "what's", false)

	// Person substitutions.
	assertIn(rs.person, "you", false)
	rs.SetPerson("you", "me")
	assertIn(rs.person, "you", true)
	rs.SetPerson("you", "undefined")
	assertIn(rs.person, "you", false)
}

func TestDebug(t *testing.T) {
	bot := NewTestWithConfig(t, true, false, memory.New())
	bot.bot.say("Debug line.")
	bot.bot.warn("Warning line.")
}
