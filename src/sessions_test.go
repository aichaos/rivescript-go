package rivescript

import (
	"testing"

	"github.com/aichaos/rivescript-go/sessions"
	"github.com/aichaos/rivescript-go/sessions/null"
)

var commonSessionTest = `
	+ my name is *
	- <set name=<formal>>Nice to meet you, <get name>.

	+ who am i
	- Aren't you <get name>?

	+ what did i just say
	- You just said: <input1>

	+ what did you just say
	- I just said: <reply1>

	+ i hate you
	- How mean!{topic=apology}

	> topic apology
		+ *
		- Nope, I'm mad at you.
	< topic
`

func TestNullSession(t *testing.T) {
	bot := NewTestWithConfig(t, false, false, null.New())
	bot.bot.Quiet = true // Suppress warnings

	bot.extend(commonSessionTest)
	bot.reply("My name is Aiden", "Nice to meet you, undefined.")
	bot.reply("Who am I?", "Aren't you undefined?")
	bot.reply("What did I just say?", "You just said: undefined")
	bot.reply("What did you just say?", "I just said: undefined")
	bot.reply("I hate you", "How mean!")
	bot.reply("My name is Aiden", "Nice to meet you, undefined.")
}

func TestMemorySession(t *testing.T) {
	bot := NewTest(t)
	bot.extend(commonSessionTest)
	bot.reply("My name is Aiden", "Nice to meet you, Aiden.")
	bot.reply("What did I just say?", "You just said: my name is aiden")
	bot.reply("Who am I?", "Aren't you Aiden?")
	bot.reply("What did you just say?", "I just said: Aren't you Aiden?")
	bot.reply("I hate you!", "How mean!")
	bot.reply("My name is Bob", "Nope, I'm mad at you.")
}

func TestFreezeThaw(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ my name is *
		- <set name=<formal>>Nice to meet you, <get name>.

		+ who am i
		- Aren't you <get name>?
	`)
	bot.reply("My name is Aiden", "Nice to meet you, Aiden.")
	bot.reply("Who am I?", "Aren't you Aiden?")

	bot.bot.FreezeUservars(bot.username)
	bot.reply("My name is Bob", "Nice to meet you, Bob.")
	bot.reply("Who am I?", "Aren't you Bob?")

	bot.bot.ThawUservars(bot.username, sessions.Thaw)
	bot.reply("Who am I?", "Aren't you Aiden?")
	bot.bot.FreezeUservars(bot.username)

	bot.reply("My name is Bob", "Nice to meet you, Bob.")
	bot.reply("Who am I?", "Aren't you Bob?")
	bot.bot.ThawUservars(bot.username, sessions.Discard)
	bot.reply("Who am I?", "Aren't you Bob?")
}
