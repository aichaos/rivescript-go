package rivescript

// NOTE: while these test files live in the 'src' package, they import the
// public facing API from the root rivescript-go package.

import (
	"testing"

	"github.com/aichaos/rivescript-go/sessions"
	"github.com/aichaos/rivescript-go/sessions/memory"
)

type RiveScriptTest struct {
	bot      *RiveScript
	t        *testing.T
	username string
}

func NewTest(t *testing.T) *RiveScriptTest {
	return NewTestWithConfig(t, false, false, memory.New())
}

func NewTestWithUTF8(t *testing.T) *RiveScriptTest {
	return NewTestWithConfig(t, false, true, memory.New())
}

func NewTestWithConfig(t *testing.T, debug, utf8 bool, ses sessions.SessionManager) *RiveScriptTest {
	test := &RiveScriptTest{
		bot:      New(),
		t:        t,
		username: "soandso",
	}
	test.bot.Configure(debug, true, utf8, 50, 0, ses)
	return test
}

// extend updates the RiveScript source code.
func (rst RiveScriptTest) extend(code string) {
	rst.bot.Stream(code)
	rst.bot.SortReplies()
}

// reply asserts that a given input gets the expected reply.
func (rst RiveScriptTest) reply(message, expected string) {
	reply, err := rst.bot.Reply(rst.username, message)
	if err != nil {
		rst.t.Errorf("Got an error when checking a reply to '%s': %s", message, err)
	} else if reply != expected {
		rst.t.Errorf("Expected %s, got %s", expected, reply)
	}
}

// replyError asserts that a given input gives an error.
func (rst RiveScriptTest) replyError(message string, expected error) {
	if reply, err := rst.bot.Reply(rst.username, message); err == nil {
		rst.t.Errorf(
			"Reply to '%s' was expected to error; but it returned %s",
			message,
			reply,
		)
	} else if err != expected {
		rst.t.Errorf(
			"Reply to '%s' got different error than expected; wanted %s, got %s",
			message,
			expected,
			err,
		)
	}
}

// assertEqual checks if two strings are equal.
func (rst RiveScriptTest) assertEqual(a, b string) {
	if a != b {
		rst.t.Errorf("assertEqual: %s != %s", a, b)
	}
}

// uservar asserts a user variable is defined and has the expected value.
func (rst RiveScriptTest) uservar(name string, expected string) {
	value, err := rst.bot.GetUservar(rst.username, name)
	if err != nil {
		rst.t.Errorf("Got an error when asserting variable %s: %s", name, err)
	} else if value != expected {
		rst.t.Errorf("Uservar %s expected %s, got %s", name, expected, value)
	}
}

// undefined asserts that a user variable is not set.
func (rst RiveScriptTest) undefined(name string) {
	if value, err := rst.bot.GetUservar(rst.username, name); err == nil {
		rst.t.Errorf("Uservar %s was expected to be undefined; but was %s", name, value)
	}
}
