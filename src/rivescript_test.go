package rivescript

import "testing"

// Test that you get an error if you didn't call SortReplies().
func TestNoSorting(t *testing.T) {
	bot := NewTest(t)
	bot.bot.Quiet = true // Suppress warnings
	bot.replyError("hello bot", ErrRepliesNotSorted)
}

// Test failing to load replies.
func TestFailedLoading(t *testing.T) {
	bot := New()
	err := bot.LoadFile("/root/notexist345613123098")
	if err == nil {
		t.Errorf("I tried to load an obviously missing file, but I succeeded unexpectedly")
	}

	err = bot.LoadDirectory("/root/notexist412901890281")
	if err == nil {
		t.Errorf("I tried to load an obviously missing directory, but I succeeded unexpectedly")
	}

	err = bot.SortReplies()
	if err == nil {
		t.Errorf("I was expecting an error from SortReplies, but I didn't get one")
	}
}
