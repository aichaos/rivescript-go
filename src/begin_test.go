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
	bot.undefined("name")
	bot.reply("My name is bob", "Hello, Bob.")
	bot.uservar("name", "Bob")
	bot.reply("Hello Bot", "Bob: Hello human.")
}

func TestDefinitions(t *testing.T) {
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
		! global g1 = one
		! global g2 = two
		! global g2 = <undef>

		! var v1 = one
		! var v2 = two
		! var v2 = <undef>

		! sub what's = what is
		! sub who's = who is
		! sub who's = <undef>

		! person you = me
		! person me = you
		! person your = my
		! person my = your
		! person your = <undef>
		! person my = <undef>
	`)
	rs := bot.bot

	assertIn(rs.global, "g1", true)
	assertIn(rs.global, "g2", false)
	assertIn(rs.vars, "v1", true)
	assertIn(rs.vars, "v2", false)
	assertIn(rs.sub, "what's", true)
	assertIn(rs.sub, "who's", false)
	assertIn(rs.person, "you", true)
	assertIn(rs.person, "your", false)
}
