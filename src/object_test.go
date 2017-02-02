package rivescript

import (
	"strings"
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

// Mock up an object macro handler for the private API testing, for code
// coverage. The MockHandler just returns its text as a string.
func TestMacroParsing(t *testing.T) {
	bot := NewTest(t)
	bot.bot.SetHandler("text", &MockHandler{
		codes: map[string]string{},
	})
	bot.extend(`
		> object hello text
			Hello world!
		< object

		> object goodbye javascript
			return "Goodbye";
		< object

		+ hello
		- <call>hello</call>

		+ goodbye
		- <call>goodbye</call>
	`)
	bot.reply("hello", "Hello world!")
	bot.reply("goodbye", "[ERR: Object Not Found]")
}

// Mock macro handler.
type MockHandler struct {
	codes map[string]string
}

func (m *MockHandler) Load(name string, code []string) {
	m.codes[name] = strings.Join(code, "\n")
}

func (m *MockHandler) Call(name string, fields []string) string {
	return m.codes[name]
}
