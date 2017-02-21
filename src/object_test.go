package rivescript

import (
	"strings"
	"testing"
)

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
