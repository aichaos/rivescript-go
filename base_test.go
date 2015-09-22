package rivescript_test

import (
	"fmt"
	rivescript "github.com/aichaos/rivescript-go"
	"testing"
)

type RiveScriptTest struct {
	bot      *rivescript.RiveScript
	t        *testing.T
	username string
}

func NewTest(t *testing.T) *RiveScriptTest {
	tester := new(RiveScriptTest)
	tester.bot = rivescript.New()
	tester.t = t
	tester.username = "soandso"
	return tester
}

func (rst RiveScriptTest) extend(code string) {
	rst.bot.Stream(code)
	rst.bot.SortReplies()
}

func (rst RiveScriptTest) reply(message string, expected string) {
	reply := rst.bot.Reply(rst.username, message)
	if reply != expected {
		rst.t.Error(fmt.Sprintf("Expected %s, got %s", expected, reply))
	}
}

func (rst RiveScriptTest) uservar(name string, expected string) {
	value, _ := rst.bot.GetUservar(rst.username, name)
	if value != expected {
		rst.t.Error(fmt.Sprintf("Uservar %s expected %s, got %s", name, expected, value))
	}
}
