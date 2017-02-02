package rivescript

import "testing"

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
	bot.replyError("What color is grass?", ErrNoTriggerMatched)
	bot.replyError("Name a Red Hat distro.", ErrNoTriggerMatched)
	bot.replyError("Name a Debian distro.", ErrNoTriggerMatched)
	bot.replyError("Say stuff.", ErrNoTriggerMatched)

	bot.bot.SetUservar(bot.username, "topic", "linux")
	bot.replyError("What color is the sky?", ErrNoTriggerMatched)
	bot.replyError("What color is the sun?", ErrNoTriggerMatched)
	bot.replyError("What color is grass?", ErrNoTriggerMatched)
	bot.reply("Name a Red Hat distro.", "Fedora.")
	bot.reply("Name a Debian distro.", "Ubuntu.")
	bot.replyError("Say stuff.", ErrNoTriggerMatched)

	bot.bot.SetUservar(bot.username, "topic", "stuff")
	bot.reply("What color is the sky?", "Blue.")
	bot.reply("What color is the sun?", "Yellow.")
	bot.replyError("What color is grass?", ErrNoTriggerMatched)
	bot.reply("Name a Red Hat distro.", "Fedora.")
	bot.reply("Name a Debian distro.", "Ubuntu.")
	bot.reply("Say stuff.", `"Stuff."`)

	bot.bot.SetUservar(bot.username, "topic", "override")
	bot.reply("What color is the sky?", "Blue.")
	bot.reply("What color is the sun?", "Purple.")
	bot.replyError("What color is grass?", ErrNoTriggerMatched)
	bot.replyError("Name a Red Hat distro.", ErrNoTriggerMatched)
	bot.replyError("Name a Debian distro.", ErrNoTriggerMatched)
	bot.replyError("Say stuff.", ErrNoTriggerMatched)

	bot.bot.SetUservar(bot.username, "topic", "morecolors")
	bot.reply("What color is the sky?", "Blue.")
	bot.reply("What color is the sun?", "Yellow.")
	bot.reply("What color is grass?", "Green.")
	bot.replyError("Name a Red Hat distro.", ErrNoTriggerMatched)
	bot.replyError("Name a Debian distro.", ErrNoTriggerMatched)
	bot.replyError("Say stuff.", ErrNoTriggerMatched)

	bot.bot.SetUservar(bot.username, "topic", "evenmore")
	bot.reply("What color is the sky?", "Blue.")
	bot.reply("What color is the sun?", "Yellow.")
	bot.reply("What color is grass?", "Blue, sometimes.")
	bot.replyError("Name a Red Hat distro.", ErrNoTriggerMatched)
	bot.replyError("Name a Debian distro.", ErrNoTriggerMatched)
	bot.replyError("Say stuff.", ErrNoTriggerMatched)
}
