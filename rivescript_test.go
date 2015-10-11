package rivescript_test

import (
	"fmt"
	"regexp"
	"testing"
)

////////////////////////////////////////////////////////////////////////////////
// BEGIN Block Tests
////////////////////////////////////////////////////////////////////////////////

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

////////////////////////////////////////////////////////////////////////////////
// Bot Variable Tests
////////////////////////////////////////////////////////////////////////////////

func TestBotVariables(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		! var name = Aiden
		! var age = 5

		+ what is your name
		- My name is <bot name>.

		+ how old are you
		- I am <bot age>.

		+ what are you
		- I'm <bot gender>.

		+ happy birthday
		- <bot age=6>Thanks!
	`)
	bot.reply("What is your name?", "My name is Aiden.")
	bot.reply("How old are you?", "I am 5.")
	bot.reply("What are you?", "I'm undefined.")
	bot.reply("Happy birthday!", "Thanks!")
	bot.reply("How old are you?", "I am 6.")
}

func TestGlobalVariables(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		! global debug = false

		+ debug mode
		- Debug mode is: <env debug>

		+ set debug mode *
		- <env debug=<star>>Switched to <star>.
	`)
	bot.reply("Debug mode.", "Debug mode is: false")
	bot.reply("Set debug mode true", "Switched to true.")
	bot.reply("Debug mode?", "Debug mode is: true")
}

////////////////////////////////////////////////////////////////////////////////
// Substitution Tests
////////////////////////////////////////////////////////////////////////////////

func TestSubstitutions(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ whats up
		- nm.

		+ what is up
		- Not much.
	`)
	bot.reply("whats up", "nm.")
	bot.reply("what's up?", "nm.")
	bot.reply("what is up?", "Not much.")

	bot.extend(`
		! sub whats  = what is
		! sub what's = what is
	`)
	bot.reply("whats up", "Not much.")
	bot.reply("what's up?", "Not much.")
	bot.reply("What is up?", "Not much.")
}

func TestPersonSubstitutions(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ say *
		- <person>
	`)
	bot.reply("say I am cool", "i am cool")
	bot.reply("say You are dumb", "you are dumb")

	bot.extend(`
		! person i am    = you are
		! person you are = I am
	`)
	bot.reply("say I am cool", "you are cool")
	bot.reply("say You are dumb", "I am dumb")
}

////////////////////////////////////////////////////////////////////////////////
// Trigger Tests
////////////////////////////////////////////////////////////////////////////////

func TestAtomicTriggers(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ hello bot
		- Hello human.

		+ what are you
		- I am a RiveScript bot.
	`)
	bot.reply("Hello bot", "Hello human.")
	bot.reply("What are you?", "I am a RiveScript bot.")
}

func TestWildcardTriggers(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ my name is *
		- Nice to meet you, <star>.

		+ * told me to say *
		- Why did <star1> tell you to say <star2>?

		+ i am # years old
		- A lot of people are <star>.

		+ i am _ years old
		- Say that with numbers.

		+ i am * years old
		- Say that with fewer words.
	`)
	bot.reply("my name is Bob", "Nice to meet you, bob.")
	bot.reply("bob told me to say hi", "Why did bob tell you to say hi?")
	bot.reply("i am 5 years old", "A lot of people are 5.")
	bot.reply("i am five years old", "Say that with numbers.")
	bot.reply("i am twenty five years old", "Say that with fewer words.")
}

func TestAlternativesAndOptionals(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ what (are|is) you
		- I am a robot.

		+ what is your (home|office|cell) [phone] number
		- It is 555-1234.

		+ [please|can you] ask me a question
		- Why is the sky blue?

		+ (aa|bb|cc) [bogus]
		- Matched.

		+ (yo|hi) [computer|bot] *
		- Matched.
	`)
	bot.reply("What are you?", "I am a robot.")
	bot.reply("What is you?", "I am a robot.")

	bot.reply("What is your home phone number?", "It is 555-1234.")
	bot.reply("What is your home number?", "It is 555-1234.")
	bot.reply("What is your cell phone number?", "It is 555-1234.")
	bot.reply("What is your office number?", "It is 555-1234.")

	bot.reply("Can you ask me a question?", "Why is the sky blue?")
	bot.reply("Please ask me a question?", "Why is the sky blue?")
	bot.reply("Ask me a question.", "Why is the sky blue?")

	bot.reply("aa", "Matched.")
	bot.reply("bb", "Matched.")
	bot.reply("aa bogus", "Matched.")
	bot.reply("aabogus", "ERR: No Reply Matched")
	bot.reply("bogus", "ERR: No Reply Matched")

	bot.reply("hi Aiden", "Matched.")
	bot.reply("hi bot how are you?", "Matched.")
	bot.reply("yo computer what time is it?", "Matched.")
	bot.reply("yoghurt is yummy", "ERR: No Reply Matched")
	bot.reply("hide and seek is fun", "ERR: No Reply Matched")
	bot.reply("hip hip hurrah", "ERR: No Reply Matched")
}

func TestTriggerArrays(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		! array colors = red blue green yellow white
		  ^ dark blue|light blue

		+ what color is my (@colors) *
		- Your <star2> is <star1>.

		+ what color was * (@colors) *
		- It was <star2>.

		+ i have a @colors *
		- Tell me more about your <star>.
	`)
	bot.reply("What color is my red shirt?", "Your shirt is red.")
	bot.reply("What color is my blue car?", "Your car is blue.")
	bot.reply("What color is my pink house?", "ERR: No Reply Matched")
	bot.reply("What color is my dark blue jacket?", "Your jacket is dark blue.")
	bot.reply("What color was Napoleoan's white horse?", "It was white.")
	bot.reply("What color was my red shirt?", "It was red.")
	bot.reply("I have a blue car.", "Tell me more about your car.")
	bot.reply("I have a cyan car.", "ERR: No Reply Matched")
}

func TestWeightedTriggers(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ * or something{weight=10}
		- Or something. <@>

		+ can you run a google search for *
		- Sure!

		+ hello *{weight=20}
		- Hi there!
	`)
	bot.reply("Hello robot.", "Hi there!")
	bot.reply("Hello or something.", "Hi there!")
	bot.reply("Can you run a Google search for Node", "Sure!")
	bot.reply("Can you run a Google search for Node or something", "Or something. Sure!")
}

////////////////////////////////////////////////////////////////////////////////
// Reply Tests
////////////////////////////////////////////////////////////////////////////////

func TestPrevious(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		! sub who's  = who is
		! sub it's   = it is
		! sub didn't = did not

		+ knock knock
		- Who's there?

		+ *
		% who is there
		- <sentence> who?

		+ *
		% * who
		- Haha! <sentence>!

		+ *
		- I don't know.
	`)
	bot.reply("knock knock", "Who's there?")
	bot.reply("Canoe", "Canoe who?")
	bot.reply("Canoe help me with my homework?", "Haha! Canoe help me with my homework!")
	bot.reply("hello", "I don't know.")
}

func TestContinuations(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ tell me a poem
		- There once was a man named Tim,\s
		^ who never quite learned how to swim.\s
		^ He fell off a dock, and sank like a rock,\s
		^ and that was the end of him.
	`)
	bot.reply("Tell me a poem.", "There once was a man named Tim, who never quite learned how to swim. He fell off a dock, and sank like a rock, and that was the end of him.")
}

func TestRedirects(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ hello
		- Hi there!

		+ hey
		@ hello

		+ hi there
		- {@hello}
	`)
	bot.reply("hello", "Hi there!")
	bot.reply("hey", "Hi there!")
	bot.reply("hi there", "Hi there!")
}

func TestConditionals(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ i am # years old
		- <set age=<star>>OK.

		+ what can i do
		* <get age> == undefined => I don't know.
		* <get age> >  25 => Anything you want.
		* <get age> == 25 => Rent a car for cheap.
		* <get age> >= 21 => Drink.
		* <get age> >= 18 => Vote.
		* <get age> <  18 => Not much of anything.

		+ am i your master
		* <get master> == true => Yes.
		- No.
	`)
	age_q := "What can I do?"
	bot.reply(age_q, "I don't know.")

	ages := map[string]string{
		"16": "Not much of anything.",
		"18": "Vote.",
		"20": "Vote.",
		"22": "Drink.",
		"24": "Drink.",
		"25": "Rent a car for cheap.",
		"27": "Anything you want.",
	}
	for age, expect := range ages {
		bot.reply(fmt.Sprintf("I am %s years old.", age), "OK.")
		bot.reply(age_q, expect)
	}

	bot.reply("Am I your master?", "No.")
	bot.bot.SetUservar(bot.username, "master", "true")
	bot.reply("Am I your master?", "Yes.")
}

func TestEmbeddedTags(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ my name is *
		* <get name> != undefined => <set oldname=<get name>>I thought\s
		  ^ your name was <get oldname>?
		  ^ <set name=<formal>>
		- <set name=<formal>>OK.

		+ what is my name
		- Your name is <get name>, right?

		+ html test
		- <set name=<b>Name</b>>This has some non-RS <em>tags</em> in it.
	`)
	bot.reply("What is my name?", "Your name is undefined, right?")
	bot.reply("My name is Alice.", "OK.")
	bot.reply("My name is Bob.", "I thought your name was Alice?")
	bot.reply("What is my name?", "Your name is Bob, right?")
	bot.reply("HTML Test", "This has some non-RS <em>tags</em> in it.")
}

func TestSetUservars(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ what is my name
		- Your name is <get name>.

		+ how old am i
		- You are <get age>.
	`)
	bot.bot.SetUservars(bot.username, map[string]string{
		"name": "Aiden",
		"age":  "5",
	})
	bot.reply("What is my name?", "Your name is Aiden.")
	bot.reply("How old am I?", "You are 5.")
}

func TestQuestionmark(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ google *
		- <a href="https://www.google.com/search?q=<star>">Results are here</a>
	`)
	bot.reply("google golang",
		`<a href="https://www.google.com/search?q=golang">Results are here</a>`,
	)
}

////////////////////////////////////////////////////////////////////////////////
// Object Macro Tests
////////////////////////////////////////////////////////////////////////////////

// TODO

////////////////////////////////////////////////////////////////////////////////
// Topic Tests
////////////////////////////////////////////////////////////////////////////////

func TestPunishmentTopic(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		+ hello
		- Hi there!

		+ swear word
		- How rude! Apologize or I won't talk to you again.{topic=sorry}

		+ *
		- Catch-all.

		> topic sorry
			+ sorry
			- It's ok!{topic=random}

			+ *
			- Say you're sorry!
		< topic
	`)
	bot.reply("hello", "Hi there!")
	bot.reply("How are you?", "Catch-all.")
	bot.reply("Swear word!", "How rude! Apologize or I won't talk to you again.")
	bot.reply("hello", "Say you're sorry!")
	bot.reply("How are you?", "Say you're sorry!")
	bot.reply("Sorry!", "It's ok!")
	bot.reply("hello", "Hi there!")
	bot.reply("How are you?", "Catch-all.")
}

func TestTopicInheritance(t *testing.T) {
	bot := NewTest(t)
	RS_ERR_MATCH := "ERR: No Reply Matched"
	bot.extend(`
		> topic colors
			+ what color is the sky
			- Blue.
			+ what color is the sun
			- Yellow.
		< topic

		> topic linux
			+ name a red hat distro
			- Fedora.
			+ name a debian distro
			- Ubuntu.
		< topic

		> topic stuff includes colors linux
			+ say stuff
			- "Stuff."
		< topic

		> topic override inherits colors
			+ what color is the sun
			- Purple.
		< topic

		> topic morecolors includes colors
			+ what color is grass
			- Green.
		< topic

		> topic evenmore inherits morecolors
			+ what color is grass
			- Blue, sometimes.
		< topic
	`)
	bot.bot.SetUservar(bot.username, "topic", "colors")
	bot.reply("What color is the sky?", "Blue.")
	bot.reply("What color is the sun?", "Yellow.")
	bot.reply("What color is grass?", RS_ERR_MATCH)
	bot.reply("Name a Red Hat distro.", RS_ERR_MATCH)
	bot.reply("Name a Debian distro.", RS_ERR_MATCH)
	bot.reply("Say stuff.", RS_ERR_MATCH)

	bot.bot.SetUservar(bot.username, "topic", "linux")
	bot.reply("What color is the sky?", RS_ERR_MATCH)
	bot.reply("What color is the sun?", RS_ERR_MATCH)
	bot.reply("What color is grass?", RS_ERR_MATCH)
	bot.reply("Name a Red Hat distro.", "Fedora.")
	bot.reply("Name a Debian distro.", "Ubuntu.")
	bot.reply("Say stuff.", RS_ERR_MATCH)

	bot.bot.SetUservar(bot.username, "topic", "stuff")
	bot.reply("What color is the sky?", "Blue.")
	bot.reply("What color is the sun?", "Yellow.")
	bot.reply("What color is grass?", RS_ERR_MATCH)
	bot.reply("Name a Red Hat distro.", "Fedora.")
	bot.reply("Name a Debian distro.", "Ubuntu.")
	bot.reply("Say stuff.", `"Stuff."`)

	bot.bot.SetUservar(bot.username, "topic", "override")
	bot.reply("What color is the sky?", "Blue.")
	bot.reply("What color is the sun?", "Purple.")
	bot.reply("What color is grass?", RS_ERR_MATCH)
	bot.reply("Name a Red Hat distro.", RS_ERR_MATCH)
	bot.reply("Name a Debian distro.", RS_ERR_MATCH)
	bot.reply("Say stuff.", RS_ERR_MATCH)

	bot.bot.SetUservar(bot.username, "topic", "morecolors")
	bot.reply("What color is the sky?", "Blue.")
	bot.reply("What color is the sun?", "Yellow.")
	bot.reply("What color is grass?", "Green.")
	bot.reply("Name a Red Hat distro.", RS_ERR_MATCH)
	bot.reply("Name a Debian distro.", RS_ERR_MATCH)
	bot.reply("Say stuff.", RS_ERR_MATCH)

	bot.bot.SetUservar(bot.username, "topic", "evenmore")
	bot.reply("What color is the sky?", "Blue.")
	bot.reply("What color is the sun?", "Yellow.")
	bot.reply("What color is grass?", "Blue, sometimes.")
	bot.reply("Name a Red Hat distro.", RS_ERR_MATCH)
	bot.reply("Name a Debian distro.", RS_ERR_MATCH)
	bot.reply("Say stuff.", RS_ERR_MATCH)
}

////////////////////////////////////////////////////////////////////////////////
// Parser Option Tests
////////////////////////////////////////////////////////////////////////////////

func TestConcat(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		// Default concat mode = none
		+ test concat default
		- Hello
		^ world!

		! local concat = space
		+ test concat space
		- Hello
		^ world!

		! local concat = none
		+ test concat none
		- Hello
		^ world!

		! local concat = newline
		+ test concat newline
		- Hello
		^ world!

		// invalid concat setting is equivalent to 'none'
		! local concat = foobar
		+ test concat foobar
		- Hello
		^ world!

		// the option is file scoped so it can be left at
		// any setting and won't affect subsequent parses
		! local concat = newline
	`)
	bot.extend(`
		// concat mode should be restored to the default in a
		// separate file/stream parse
		+ test concat second file
		- Hello
		^ world!
	`)

	bot.reply("test concat default", "Helloworld!")
	bot.reply("test concat space", "Hello world!")
	bot.reply("test concat none", "Helloworld!")
	bot.reply("test concat newline", "Hello\nworld!")
	bot.reply("test concat foobar", "Helloworld!")
	bot.reply("test concat second file", "Helloworld!")
}

////////////////////////////////////////////////////////////////////////////////
// Unicode Tests
////////////////////////////////////////////////////////////////////////////////

func TestUnicode(t *testing.T) {
	bot := NewTest(t)
	bot.bot.UTF8 = true
	bot.extend(`
		! sub who's = who is
		+ äh
		- What's the matter?

		+ ブラッキー
		- エーフィ

		// Make sure %Previous continues working in UTF-8 mode.
		+ knock knock
		- Who's there?

		+ *
		% who is there
		- <sentence> who?

		+ *
		% * who
		- Haha! <sentence>!

		// And with UTF-8.
		+ tëll më ä pöëm
		- Thërë öncë wäs ä män nämëd Tïm

		+ more
		% thërë öncë wäs ä män nämëd tïm
		- Whö nëvër qüïtë lëärnëd höw tö swïm

		+ more
		% whö nëvër qüïtë lëärnëd höw tö swïm
		- Hë fëll öff ä döck, änd sänk lïkë ä röck

		+ more
		% hë fëll öff ä döck änd sänk lïkë ä röck
		- Änd thät wäs thë ënd öf hïm.
	`)

	bot.reply("äh", "What's the matter?")
	bot.reply("ブラッキー", "エーフィ")
	bot.reply("knock knock", "Who's there?")
	bot.reply("Orange", "Orange who?")
	bot.reply("banana", "Haha! Banana!")
	bot.reply("tëll më ä pöëm", "Thërë öncë wäs ä män nämëd Tïm")
	bot.reply("more", "Whö nëvër qüïtë lëärnëd höw tö swïm")
	bot.reply("more", "Hë fëll öff ä döck, änd sänk lïkë ä röck")
	bot.reply("more", "Änd thät wäs thë ënd öf hïm.")
}

func TestPunctuation(t *testing.T) {
	bot := NewTest(t)
	bot.bot.UTF8 = true
	bot.extend(`
		+ hello bot
		- Hello human!
	`)

	bot.reply("Hello bot", "Hello human!")
	bot.reply("Hello, bot!", "Hello human!")
	bot.reply("Hello: Bot", "Hello human!")
	bot.reply("Hello... bot?", "Hello human!")

	bot.bot.UnicodePunctuation = regexp.MustCompile(`xxx`)
	bot.reply("Hello bot", "Hello human!")
	bot.reply("Hello, bot!", "ERR: No Reply Matched")
}
