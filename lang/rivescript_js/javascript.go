/*
Package rivescript_js implements JavaScript object macros for RiveScript.

Usage is simple. In your Golang code:

	import (
		rivescript "github.com/aichaos/rivescript-go"
		"github.com/aichaos/rivescript-go/lang/rivescript_js"
	)

	func main() {
		bot := rivescript.New()
		jsHandler := rivescript_js.New()
		bot.SetHandler("javascript", jsHandler)

		// and go on as normal
	}

And in your RiveScript code, you can load and run JavaScript objects:

	> object add javascript
		var a = args[0];
		var b = args[1];
		return parseInt(a) + parseInt(b);
	< object

	+ add # and #
	- <star1> + <star2> = <call>add <star1> <star2></call>
*/
package rivescript_js

import (
	"fmt"
	"strings"
	"github.com/robertkrimen/otto"
)

type JavaScriptHandler struct {
	vm *otto.Otto
	functions map[string]string
}

// New creates an object handler for JavaScript with its own Otto VM.
func New() *JavaScriptHandler {
	js := new(JavaScriptHandler)
	js.vm = otto.New()
	js.functions = map[string]string{}
	return js
}

// Load loads a new JavaScript object macro into the VM.
func (js JavaScriptHandler) Load(name string, code []string) {
	// Create a unique function name called the same as the object macro name.
	js.functions[name] = fmt.Sprintf(`
		function object_%s(args) {
			%s
		}
	`, name, strings.Join(code, "\n"))

	// Run this code to load the function into the VM.
	js.vm.Run(js.functions[name])
}

// Call executes a JavaScript macro and returns its results.
func (js JavaScriptHandler) Call(name string, fields []string) string {
	// Turn the array of arguments into a JavaScript list of quoted strings.
	jsFields := ""
	for _, field := range fields {
		field = strings.Replace(field, `"`, `\"`, -1)
		if len(jsFields) > 0 {
			jsFields += ", "
		}
		jsFields += fmt.Sprintf(`"%s"`, field)
	}

	// Run the JS function call and get the result.
	js.vm.Run(fmt.Sprintf(`var result = object_%s([%s]);`, name, jsFields))

	// Retrieve the result from the VM.
	result, err := js.vm.Get("result")
	if err != nil {
		return fmt.Sprintf("[Error in JavaScript object: %s]", err)
	}

	// Return it.
	reply, _ := result.ToString()
	return reply
}
