package rivescript

import (
	"fmt"
	"testing"
)

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

// Test matching a large number of stars, greater than <star1>-<star9>
func TestManyStars(t *testing.T) {
	bot := NewTest(t)
	bot.extend(`
		// 10 stars.
		+ * * * * * * * * * *
		- That's a long one. 1=<star1>; 5=<star5>; 9=<star9>; 10=<star10>;

		// 16 stars!
		+ * * * * * * * * * * * * * * * *
		- Wow! 1=<star1>; 3=<star3>; 7=<star7>; 14=<star14>; 15=<star15>; 16=<star16>;
	`)
	bot.reply(
		"one two three four five six seven eight nine ten eleven",
		"That's a long one. 1=one; 5=five; 9=nine; 10=ten eleven;",
	)
	bot.reply(
		"0 1 2 3 4 5 6 7 8 9 A B C D E F G H I J K",
		"Wow! 1=0; 3=2; 7=6; 14=d; 15=e; 16=f g h i j k;",
	)
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

		// Infinite recursion between these two.
		+ one
		@ two
		+ two
		@ one

		// Variables can throw off redirects with their capitalizations,
		// so make sure redirects handle this properly.
		! var master = Kirsle

		+ my name is (<bot master>)
		- <set name=<formal>>That's my botmaster's name too.

		+ call me <bot master>
		@ my name is <bot master>
	`)
	bot.reply("hello", "Hi there!")
	bot.reply("hey", "Hi there!")
	bot.reply("hi there", "Hi there!")

	bot.reply("my name is Kirsle", "That's my botmaster's name too.")
	bot.reply("call me kirsle", "That's my botmaster's name too.")

	bot.replyError("one", ErrDeepRecursion)
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
	ageQ := "What can I do?"
	bot.reply(ageQ, "I don't know.")

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
		bot.reply(ageQ, expect)
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
