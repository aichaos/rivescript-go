package src_test

// NOTE: while these test files live in the 'src' package, they import the
// public facing API from the root rivescript-go package.

import (
	"fmt"
	"testing"

	"github.com/aichaos/rivescript-go"
	"github.com/aichaos/rivescript-go/config"
)

type RiveScriptTest struct {
	bot      *rivescript.RiveScript
	t        *testing.T
	username string
}

func NewTest(t *testing.T) *RiveScriptTest {
	return &RiveScriptTest{
		bot:      rivescript.New(config.Basic()),
		t:        t,
		username: "soandso",
	}
}

func NewTestWithConfig(t *testing.T, config *config.Config) *RiveScriptTest {
	return &RiveScriptTest{
		bot:      rivescript.New(config),
		t:        t,
		username: "soandso",
	}
}

// RS exposes the underlying RiveScript API.
func (rst *RiveScriptTest) RS() *rivescript.RiveScript {
	return rst.bot
}

// extend updates the RiveScript source code.
func (rst RiveScriptTest) extend(code string) {
	rst.bot.Stream(code)
	rst.bot.SortReplies()
}

// reply asserts that a given input gets the expected reply.
func (rst RiveScriptTest) reply(message string, expected string) {
	reply := rst.bot.Reply(rst.username, message)
	if reply != expected {
		rst.t.Error(fmt.Sprintf("Expected %s, got %s", expected, reply))
	}
}

// uservar asserts a user variable.
func (rst RiveScriptTest) uservar(name string, expected string) {
	value, _ := rst.bot.GetUservar(rst.username, name)
	if value != expected {
		rst.t.Error(fmt.Sprintf("Uservar %s expected %s, got %s", name, expected, value))
	}
}
