package rivescript_test

// This test file contains the unit tests that had to be segregated from the
// others in the `src/` package.
//
// The only one here so far is an object macro test. It needed to use the public
// RiveScript API because the JavaScript handler expects an object of that type,
// and so it couldn't be in the `src/` package or it would create a dependency
// cycle.

import (
	"testing"

	rivescript "github.com/aichaos/rivescript-go"
	"github.com/aichaos/rivescript-go/lang/javascript"
)

// This one has to test the public interface because of the JavaScript handler
// expecting a *RiveScript of the correct color.
func TestJavaScript(t *testing.T) {
	rs := rivescript.New(nil)
	rs.SetHandler("javascript", javascript.New(rs))
	rs.Stream(`
		> object reverse javascript
			var msg = args.join(" ");
			return msg.split("").reverse().join("");
		< object

		> object nolang
			return "No language provided!"
		< object

		+ reverse *
		- <call>reverse <star></call>

		+ no lang
		- <call>nolang</call>
	`)
	rs.SortReplies()

	// Helper function to assert replies via the public interface.
	assert := func(input, expected string) {
		reply, err := rs.Reply("local-user", input)
		if err != nil {
			t.Errorf("Got error when trying to get a reply: %v", err)
		} else if reply != expected {
			t.Errorf("Got unexpected reply. Expected %s, got %s", expected, reply)
		}
	}

	assert("reverse hello world", "dlrow olleh")
	assert("no lang", "[ERR: Object Not Found]")

	// Disable support.
	rs.RemoveHandler("javascript")
	assert("reverse hello world", "[ERR: Object Not Found]")
}
