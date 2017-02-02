package rivescript

import "testing"

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

		// Verify that wildcards in optionals are not matchable.
		+ my favorite [_] is *
		- Why is it <star1>?

		+ i have [#] questions about *
		- Well I don't have any answers about <star1>.
	`)
	bot.reply("my name is Bob", "Nice to meet you, bob.")
	bot.reply("bob told me to say hi", "Why did bob tell you to say hi?")
	bot.reply("i am 5 years old", "A lot of people are 5.")
	bot.reply("i am five years old", "Say that with numbers.")
	bot.reply("i am twenty five years old", "Say that with fewer words.")

	bot.reply("my favorite color is red", "Why is it red?")
	bot.reply("i have 2 questions about bots", "Well I don't have any answers about bots.")
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
	bot.replyError("aabogus", ErrNoTriggerMatched)
	bot.replyError("bogus", ErrNoTriggerMatched)

	bot.reply("hi Aiden", "Matched.")
	bot.reply("hi bot how are you?", "Matched.")
	bot.reply("yo computer what time is it?", "Matched.")
	bot.replyError("yoghurt is yummy", ErrNoTriggerMatched)
	bot.replyError("hide and seek is fun", ErrNoTriggerMatched)
	bot.replyError("hip hip hurrah", ErrNoTriggerMatched)
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
	bot.replyError("What color is my pink house?", ErrNoTriggerMatched)
	bot.reply("What color is my dark blue jacket?", "Your jacket is dark blue.")
	bot.reply("What color was Napoleoan's white horse?", "It was white.")
	bot.reply("What color was my red shirt?", "It was red.")
	bot.reply("I have a blue car.", "Tell me more about your car.")
	bot.replyError("I have a cyan car.", ErrNoTriggerMatched)
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
