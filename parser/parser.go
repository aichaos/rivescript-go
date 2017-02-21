/*
Package parser is a RiveScript language parser.

This package can be used as a stand-alone parser for third party developers
to use, if you want to be able to simply parse (and syntax check!)
RiveScript source code and get an "abstract syntax tree" back from it.
*/
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aichaos/rivescript-go/ast"
)

const RS_VERSION float64 = 2.0

/*
ParserConfig configures the parser.

Configuration Options

	Strict: Enable strict syntax checking. Syntax errors will be considered
		fatal and abandon the parsing process.
	UTF8: Enable UTF-8 mode. When enabled, this allows triggers to contain
		foreign symbols without raising a syntax error.
	OnDebug: A function handler for receiving debug information from this
		package, if you want that information.
	OnWarn: A function handler for receiving warnings (non-fatal errors) from
		this package.

All options have meaningful zero values.
*/
type ParserConfig struct {
	Strict bool // Strict syntax checking enable (true by default)
	UTF8   bool // Enable UTF-8 mode (false by default)

	// Optional handlers for the caller to get debug information out.
	OnDebug func(message string, a ...interface{})
	OnWarn  func(message, filename string, lineno int, a ...interface{})
}

type Parser struct {
	C ParserConfig
}

// New creates and returns a new instance of a RiveScript Parser.
func New(config ParserConfig) *Parser {
	return &Parser{config}
}

// say proxies to the OnDebug handler.
func (self *Parser) say(message string, a ...interface{}) {
	if self.C.OnDebug != nil {
		self.C.OnDebug(message, a...)
	}
}

// warn proxies to the OnWarn handler.
func (self *Parser) warn(message, filename string, lineno int, a ...interface{}) {
	if self.C.OnWarn != nil {
		self.C.OnWarn(message, filename, lineno, a...)
	}
}

/*
Parse reads and parses RiveScript source code.

This will return an AST Root object containing all of the relevant
information parsed from the source code.

In case of errors (e.g. a syntax error while Strict Mode is enabled) will
return a nil AST root and an error object.

Parameters

	filename: An arbitrary name for the source code being parsed. It will be
		used when reporting warnings from this package.
	code: An array of lines of RiveScript source code.
*/
func (self *Parser) Parse(filename string, code []string) (*ast.Root, error) {
	self.say("In parse!")

	// Eventual return structure.
	// NOTE: the all caps AST is the instance, and the lowercase ast is the
	//       package that defines the types.
	AST := ast.New()

	// Track temporary variables
	var (
		topic   = "random"   // Default topic = random
		lineno  int          // Line numbers for syntax tracking
		comment bool         // In a multi-line comment
		inobj   bool         // In an object macro
		objName string       // Name of the object we're in
		objLang string       // The programming language of the object
		objBuf  = []string{} // Source code buffer of the object
		isThat  string       // Is a %Previous trigger
		curTrig *ast.Trigger // Pointer to the current trigger
	)

	// Local (file-scoped) parser options.
	localOptions := map[string]string{
		"concat": "none",
	}
	concatModes := map[string]string{
		// Supported concat modes
		"none":    "",
		"newline": "\n",
		"space":   " ",
	}

	// Go through the lines of code.
	for lp, line := range code {
		lineno = lp + 1

		// Strip the line
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue // Skip blank lines!
		}

		//------------------------------//
		// Are we inside an `> object`? //
		//------------------------------//
		if inobj {
			// End of the object?
			if strings.Contains(line, "< object") || strings.Contains(line, "<object") {
				// End the object
				if len(objName) > 0 {
					newObject := new(ast.Object)
					newObject.Name = objName
					newObject.Language = objLang
					newObject.Code = objBuf
					AST.Objects = append(AST.Objects, newObject)
				}
				inobj = false
			} else {
				objBuf = append(objBuf, line)
			}
			continue
		}

		//-------------------//
		// Look for comments //
		//-------------------//
		if strings.Index(line, "//") == 0 {
			continue // Single line comment
		} else if strings.Index(line, "/*") == 0 {
			// Start of a multi-line comment.
			if strings.Index(line, "*/") > -1 {
				// The end comment is on the same line!
				continue
			}

			// We're now inside a multi-line comment.
			comment = true
			continue
		} else if strings.Index(line, "*/") > -1 {
			// End of a multi-line comment.
			comment = false
			continue
		} else if comment {
			continue
		}

		// Separate the command from its data.
		if len(line) < 2 {
			self.warn("Weird single-character line '%s' found", filename, lineno, line)
			continue
		}
		cmd := string(line[0])
		line = line[1:]

		// Ignore in-line comments if there's a space before and after the "//"
		if strings.Index(line, " // ") > -1 {
			line = strings.Split(line, " // ")[0]
		}

		line = strings.TrimSpace(line)

		// TODO: check syntax

		// Reset the %Previous state if this is a new +Trigger.
		if cmd == "+" {
			isThat = ""
		}

		self.say("Cmd: %s; line: %s", cmd, line)

		// Do a look-ahead for ^Continue and %Previous commands.
		if cmd != "^" {
			for li, lookahead := range code[lp+1:] {
				lookahead = strings.TrimSpace(lookahead)
				if len(lookahead) < 2 {
					continue
				}
				lookCmd := string(lookahead[0])
				lookahead = strings.TrimSpace(lookahead[1:])

				// We only care about a couple lookahead command types.
				if lookCmd != "%" && lookCmd != "^" {
					break
				}

				// Only continue if the lookahead has any data
				if len(lookahead) == 0 {
					break
				}

				self.say("\tLookahead %d: %s %s", li, lookCmd, lookahead)

				// If the current command is a +, see if the following is a %
				if cmd == "+" {
					if lookCmd == "%" {
						isThat = lookahead
						break
					} else {
						isThat = ""
					}
				}

				// If the current command is a ! and the next command(s) are ^,
				// we'll tack each extension on as a line break (which is useful
				// information for arrays).
				if cmd == "!" {
					if lookCmd == "^" {
						line += fmt.Sprintf("<crlf>%s", lookahead)
					}
					continue
				}

				// If the current command is not a ^, and the line after is not a %,
				// but the line after IS a ^, then tack it on to the end of the
				// current line.
				if cmd != "^" && lookCmd != "%" {
					if lookCmd == "^" {
						// Which character to concatenate with?
						// TODO: if concatModes[blah] isnt undefined
						line += concatModes[localOptions["concat"]] + lookahead
					}
				}
			}
		}

		// Handle the types of RiveScript commands
		switch cmd {
		case "!": // ! Define
			var (
				halves = strings.SplitN(line, "=", 2)
				left   = strings.Split(strings.TrimSpace(halves[0]), " ")
				value  string
				kind   string // global, var, sub, ...
				name   string
			)

			if len(halves) == 2 {
				value = strings.TrimSpace(halves[1])
			}
			if len(left) >= 1 {
				kind = strings.TrimSpace(left[0])
				if len(left) >= 2 {
					left = left[1:]
					name = strings.TrimSpace(strings.Join(left, " "))
				}
			}

			// Remove 'fake' line breaks unless this is an array.
			if kind != "array" {
				crlfReplacer := strings.NewReplacer("<crlf>", "")
				value = crlfReplacer.Replace(value)
			}

			// Handle version numbers
			if kind == "version" {
				parsedVersion, _ := strconv.ParseFloat(value, 32)
				if parsedVersion > RS_VERSION {
					return nil, fmt.Errorf(
						"Unsupported RiveScript version. We only support %f at %s line %d",
						RS_VERSION, filename, lineno,
					)
				}
				continue
			}

			// All other types of define's require a value and a variable name.
			if len(name) == 0 {
				self.warn("Undefined variable name", filename, lineno)
				continue
			}
			if len(value) == 0 {
				self.warn("Undefined variable value", filename, lineno)
				continue
			}

			// Handle the rest of the !Define types.
			switch kind {
			case "local":
				// Local file-scoped parser options
				self.say("\tSet local parser option %s = %s", name, value)
				localOptions[name] = value
			case "global":
				// Set a 'global' variable.
				self.say("\tSet global %s = %s", name, value)
				AST.Begin.Global[name] = value
			case "var":
				// Set a bot variable.
				self.say("\tSet bot variable %s = %s", name, value)
				AST.Begin.Var[name] = value
			case "array":
				// Set an array
				self.say("\tSet array %s = %s", name, value)

				// Did we have multiple parts?
				parts := strings.Split(value, "<crlf>")

				// Process each line of array data.
				fields := []string{}
				for _, val := range parts {
					if strings.Contains(val, "|") {
						fields = append(fields, strings.Split(val, "|")...)
					} else {
						fields = append(fields, strings.Split(val, " ")...)
					}
				}

				// Convert any remaining \s's over.
				for i := range fields {
					spaceReplacer := strings.NewReplacer("\\s", " ")
					fields[i] = spaceReplacer.Replace(fields[i])
				}

				AST.Begin.Array[name] = fields
			case "sub":
				// Substitutions
				self.say("\tSet substitution %s = %s", name, value)
				AST.Begin.Sub[name] = value
			case "person":
				// Person substitutions
				self.say("\tSet person substitution %s = %s", name, value)
				AST.Begin.Person[name] = value
			default:
				self.warn("Unknown definition type '%s'", filename, lineno, kind)
			}
		case ">": // > Label
			temp := strings.Split(strings.TrimSpace(line), " ")
			kind := temp[0]
			temp = temp[1:]
			name := ""
			fields := []string{}
			if len(temp) > 0 {
				name = temp[0]
				temp = temp[1:]
			}
			if len(temp) > 0 {
				fields = temp
			}

			// Handle the label types.
			if kind == "begin" {
				self.say("Found the BEGIN block.")
				kind = "topic"
				name = "__begin__"
			}
			if kind == "topic" {
				self.say("Set topic to %s", name)
				curTrig = nil
				topic = name

				// Initialize the topic tree.
				AST.AddTopic(topic)

				// Does this topic include or inherit another one?
				mode := ""
				if len(fields) >= 2 {
					for _, field := range fields {
						if field == "includes" || field == "inherits" {
							mode = field
						} else if mode == "includes" {
							AST.Topics[topic].Includes[field] = true
						} else if mode == "inherits" {
							AST.Topics[topic].Inherits[field] = true
						}
					}
				}
			} else if kind == "object" {
				// If a field was provided, it should be the programming language.
				lang := ""
				if len(fields) > 0 {
					lang = strings.ToLower(fields[0])
				}

				// Missing language?
				if lang == "" {
					self.warn("No programming language specified for object '%s'", filename, lineno, name)
					inobj = true
					objName = name
					objLang = "__unknown__"
					continue
				}

				// Start reading the object code.
				objName = name
				objLang = lang
				objBuf = []string{}
				inobj = true
			} else {
				self.warn("Unknown label type '%s'", filename, lineno, kind)
			}
		case "<": // < Label
			kind := line

			if kind == "begin" || kind == "topic" {
				self.say("\tEnd the topic label.")
				topic = "random" // Go back to default topic
			} else if kind == "object" {
				self.say("\tEnd the object label.")
				inobj = false
			}
		case "+": // +Trigger
			self.say("\tTrigger pattern: %s", line)

			// Initialize the trigger tree.
			curTrig = new(ast.Trigger)
			curTrig.Trigger = line
			curTrig.Reply = []string{}
			curTrig.Condition = []string{}
			curTrig.Redirect = ""
			curTrig.Previous = isThat
			AST.Topics[topic].Triggers = append(AST.Topics[topic].Triggers, curTrig)
		case "-": // -Response
			if curTrig == nil {
				self.warn("Response found before trigger", filename, lineno)
				continue
			}

			self.say("\tResponse: %s", line)
			curTrig.Reply = append(curTrig.Reply, line)
		case "*": // *condition
			if curTrig == nil {
				self.warn("Condition found before trigger", filename, lineno)
				continue
			}

			self.say("\tCondition: %s", line)
			curTrig.Condition = append(curTrig.Condition, line)
		case "%": // %Previous
			continue // This was handled above
		case "^": // ^Continue
			continue // This was handled above
		case "@": // @Redirect
			if curTrig == nil {
				self.warn("Redirect found before trigger", filename, lineno)
				continue
			}

			self.say("\tRedirect response to: %s", line)
			curTrig.Redirect = line
		default:
			self.warn("Unknown command '%s'", filename, lineno, cmd)
		}
	}

	return AST, nil
}
