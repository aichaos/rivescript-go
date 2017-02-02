package rivescript

import "testing"

func TestNoBeginBlock(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ hello bot
		- Hello human.
	`)
	bot.reply("Hello bot", "Hello human.")
}

func TestSimpleBeginBlock(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		> begin
			+ request
			- {ok}
		< begin

		+ hello bot
		- Hello human.
	`)
	bot.reply("Hello bot", "Hello human.")
}

func TestConditionalBeginBlock(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		> begin
			+ request
			* <get met> == undefined => <set met=true>{ok}
			* <get name> != undefined => <get name>: {ok}
			- {ok}
		< begin

		+ hello bot
		- Hello human.

		+ my name is *
		- <set name=<formal>>Hello, <get name>.
	`)
	bot.reply("Hello bot", "Hello human.")
	bot.uservar("met", "true")
	bot.uservar("name", "undefined")
	bot.reply("My name is bob", "Hello, Bob.")
	bot.uservar("name", "Bob")
	bot.reply("Hello Bot", "Bob: Hello human.")
}
