/*
Package javascript implements JavaScript object macros for RiveScript.

This is powered by the Otto JavaScript engine[1], which is a JavaScript engine
written in pure Go. It is not the V8 engine used by Node, so expect possible
compatibility issues to arise.

Usage is simple. In your Golang code:

	import (
		rivescript "github.com/aichaos/rivescript-go"
		"github.com/aichaos/rivescript-go/lang/javascript"
	)

	func main() {
		bot := rivescript.New(nil)
		jsHandler := javascript.New(bot)
		bot.SetHandler("javascript", jsHandler)

		// and go on as normal
	}

And in your RiveScript code, you can load and run JavaScript objects:

	> object add javascript
		var a = args[0];
		var b = args[1];
		return parseInt(a) + parseInt(b);
	< object

	> object setname javascript
		// Set the user's name via JavaScript
		var uid = rs.CurrentUser();
		rs.SetUservar(uid, args[0], args[1])
	< object

	+ add # and #
	- <star1> + <star2> = <call>add <star1> <star2></call>

	+ my name is *
	- I will remember that.<call>setname <id> <formal></call>

	+ what is my name
	- You are <get name>.

[1]: https://github.com/robertkrimen/otto
*/
package javascript

import (
	"fmt"
	"strings"

	"github.com/aichaos/rivescript-go"
	"github.com/robertkrimen/otto"
)

type JavaScriptHandler struct {
	vm        *otto.Otto
	bot       *rivescript.RiveScript
	functions map[string]string
}

// New creates an object handler for JavaScript with its own Otto VM.
func New(rs *rivescript.RiveScript) *JavaScriptHandler {
	js := new(JavaScriptHandler)
	js.vm = otto.New()
	js.bot = rs
	js.functions = map[string]string{}

	return js
}

// Load loads a new JavaScript object macro into the VM.
func (js JavaScriptHandler) Load(name string, code []string) {
	// Create a unique function name called the same as the object macro name.
	js.functions[name] = fmt.Sprintf(`
		function object_%s(rs, args) {
			%s
		}
	`, name, strings.Join(code, "\n"))

	// Run this code to load the function into the VM.
	js.vm.Run(js.functions[name])
}

// Call executes a JavaScript macro and returns its results.
func (js JavaScriptHandler) Call(name string, fields []string) string {
	// Make the RiveScript object available to the JS.
	v, err := js.vm.ToValue(js.bot)
	if err != nil {
		fmt.Printf("Error binding RiveScript object to Otto: %s", err)
	}

	// Convert the fields into a JavaScript object.
	jsFields, err := js.vm.ToValue(fields)
	if err != nil {
		fmt.Printf("Error binding fields to Otto: %s", err)
	}

	// Run the JS function call and get the result.
	result, err := js.vm.Call(fmt.Sprintf("object_%s", name), nil, v, jsFields)
	if err != nil {
		fmt.Printf("Error: %s", err)
	}

	reply := ""
	if result.IsDefined() {
		reply, _ = result.ToString()
	}

	// Return it.
	return reply
}
