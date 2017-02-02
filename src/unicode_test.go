package rivescript

import "testing"

func TestUnicode(t *testing.T) {
	bot := NewTestWithUTF8(t)
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
	bot := NewTestWithUTF8(t)
	bot.extend(`
		+ hello bot
		- Hello human!
	`)

	bot.reply("Hello bot", "Hello human!")
	bot.reply("Hello, bot!", "Hello human!")
	bot.reply("Hello: Bot", "Hello human!")
	bot.reply("Hello... bot?", "Hello human!")

	bot.bot.SetUnicodePunctuation(`xxx`)
	bot.reply("Hello bot", "Hello human!")
	bot.replyError("Hello, bot!", ErrNoTriggerMatched)
}
