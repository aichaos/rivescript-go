package redis_test

import (
	"testing"

	rivescript "github.com/aichaos/rivescript-go"
	"github.com/aichaos/rivescript-go/sessions"
	"github.com/aichaos/rivescript-go/sessions/redis"
)

// This script tests the 'integration' of the RiveScript public API with the
// RiveScript-Redis public API.

func TestIntegration(t *testing.T) {
	bot := rivescript.New(&rivescript.Config{
		SessionManager: redis.New(&redis.Config{
			Prefix: "rivescript:integration/",
		}),
	})
	bot.Stream(`
        + hello bot
        - Hello human.

        + my name is *
        - <set name=<formal>>Nice to meet you, <get name>.

        + who am i
        - Your name is <get name>.

        + i am # years old
        - <set age=<star>>I will remember you are <get age> years old.

        + how old am i
        - You are <get age>.

        + today is my birthday
        - <add age=1>Happy birthday!
    `)
	expectVar(t, bot, "alice", "name", "")
	bot.SortReplies()

	// See if Redis can remember things.
	expectReply(t, bot, "alice", "my name is Alice", "Nice to meet you, Alice.")
	expectReply(t, bot, "alice", "I am 5 years old", "I will remember you are 5 years old.")
	expectVar(t, bot, "alice", "name", "Alice")
	expectVar(t, bot, "alice", "age", "5")

	// Freeze variables and restore them.
	bot.FreezeUservars("alice")
	expectReply(t, bot, "alice", "Today is my birthday", "Happy birthday!")
	expectVar(t, bot, "alice", "age", "6")
	bot.ThawUservars("alice", sessions.Thaw)
	expectVar(t, bot, "alice", "age", "5")

	// Clean up.
	bot.ClearAllUservars()
}

func expectReply(t *testing.T, bot *rivescript.RiveScript, username, input, expected string) {
	reply, _ := bot.Reply(username, input)
	if reply != expected {
		t.Errorf(
			"got unexpected reply for: [%s] %s\n"+
				"expected: %s\n"+
				"     got: %s",
			username,
			input,
			expected,
			reply,
		)
	}
}

func expectVar(t *testing.T, bot *rivescript.RiveScript, username, name, expected string) {
	value, _ := bot.GetUservar(username, name)
	if value != expected {
		t.Errorf(
			"got unexpected user variable for user '%s'\n"+
				"expected: '%s'='%s'\n"+
				"     got: '%s'",
			username,
			name,
			expected,
			value,
		)
	}
}
