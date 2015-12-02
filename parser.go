package rivescript

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func (rs *RiveScript) parseSource(filename string, code []string) (*astRoot, error) {
	rs.say("In parse!")

	// Eventual return structure
	ast := newAST()
	ast.begin.global["hi"] = "true"

	// Track temporary variables
	topic := "random"       // Default topic = random
	lineno := 0             // Line numbers for syntax tracking
	comment := false        // In a multi-line comment
	inobj := false          // In an object macro
	objName := ""           // Name of the object we're in
	objLang := ""           // The programming language of the object
	objBuf := []string{}    // Source code buffer of the object
	isThat := ""            // Is a %Previous trigger
	var curTrig *astTrigger // Pointer to the current trigger
	curTrig = nil

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
					newObject := new(astObject)
					newObject.name = objName
					newObject.language = objLang
					newObject.code = objBuf
					ast.objects = append(ast.objects, newObject)
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
			rs.warnSyntax("Weird single-character line '%s' found", filename, lineno, line)
			continue
		}
		cmd := string(line[0])
		line = strings.TrimSpace(line[1:])

		// Ignore in-line comments if there's a space before and after the "//"
		if strings.Index(line, " // ") > -1 {
			line = strings.Split(line, " // ")[0]
		}

		// TODO: check syntax

		// Reset the %Previous state if this is a new +Trigger.
		if cmd == "+" {
			isThat = ""
		}

		rs.say("Cmd: %s; line: %s", cmd, line)

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

				rs.say("\tLookahead %d: %s %s", li, lookCmd, lookahead)

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
			halves := strings.SplitN(line, "=", 2)
			left := strings.Split(strings.TrimSpace(halves[0]), " ")
			value := ""
			type_ := ""
			name := ""
			if len(halves) == 2 {
				value = strings.TrimSpace(halves[1])
			}
			if len(left) >= 1 {
				type_ = strings.TrimSpace(left[0])
				if len(left) >= 2 {
					left = left[1:]
					name = strings.TrimSpace(strings.Join(left, " "))
				}
			}

			// Remove 'fake' line breaks unless this is an array.
			if type_ != "array" {
				crlfReplacer := strings.NewReplacer("<crlf>", "")
				value = crlfReplacer.Replace(value)
			}

			// Handle version numbers
			if type_ == "version" {
				parsedVersion, _ := strconv.ParseFloat(value, 32)
				if parsedVersion > RS_VERSION {
					return nil, errors.New(
						fmt.Sprintf("Unsupported RiveScript version. We only support %f at %s line %d", RS_VERSION, filename, lineno),
					)
				}
				continue
			}

			// All other types of define's require a value and a variable name.
			if len(name) == 0 {
				rs.warnSyntax("Undefined variable name", filename, lineno)
				continue
			}
			if len(value) == 0 {
				rs.warnSyntax("Undefined variable value", filename, lineno)
				continue
			}

			// Handle the rest of the !Define types.
			switch type_ {
			case "local":
				// Local file-scoped parser options
				rs.say("\tSet local parser option %s = %s", name, value)
				localOptions[name] = value
			case "global":
				// Set a 'global' variable.
				rs.say("\tSet global %s = %s", name, value)
				ast.begin.global[name] = value
			case "var":
				// Set a bot variable.
				rs.say("\tSet bot variable %s = %s", name, value)
				ast.begin.var_[name] = value
			case "array":
				// Set an array
				rs.say("\tSet array %s = %s", name, value)

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
				for i, _ := range fields {
					spaceReplacer := strings.NewReplacer("\\s", " ")
					fields[i] = spaceReplacer.Replace(fields[i])
				}

				ast.begin.array[name] = fields
			case "sub":
				// Substitutions
				rs.say("\tSet substitution %s = %s", name, value)
				ast.begin.sub[name] = value
			case "person":
				// Person substitutions
				rs.say("\tSet person substitution %s = %s", name, value)
				ast.begin.person[name] = value
			default:
				rs.warnSyntax("Unknown definition type '%s'", filename, lineno, type_)
			}
		case ">": // > Label
			temp := strings.Split(strings.TrimSpace(line), " ")
			type_ := temp[0]
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
			if type_ == "begin" {
				rs.say("Found the BEGIN block.")
				type_ = "topic"
				name = "__begin__"
			}
			if type_ == "topic" {
				rs.say("Set topic to %s", name)
				curTrig = nil
				topic = name

				// Initialize the topic tree.
				ast = initTopic(ast, topic)

				// Does this topic include or inherit another one?
				mode := ""
				if len(fields) >= 2 {
					for _, field := range fields {
						if field == "includes" || field == "inherits" {
							mode = field
						} else if mode == "includes" {
							ast.topics[topic].includes[field] = true
						} else if mode == "inherits" {
							ast.topics[topic].inherits[field] = true
						}
					}
				}
			} else if type_ == "object" {
				// If a field was provided, it should be the programming language.
				lang := ""
				if len(fields) > 0 {
					lang = strings.ToLower(fields[0])
				}

				// Missing language?
				if lang == "" {
					rs.warnSyntax("No programming language specified for object '%s'", filename, lineno, name)
					continue
				}

				// Start reading the object code.
				objName = name
				objLang = lang
				objBuf = []string{}
				inobj = true
			} else {
				rs.warnSyntax("Unknown label type '%s'", filename, lineno, type_)
			}
		case "<": // < Label
			type_ := line

			if type_ == "begin" || type_ == "topic" {
				rs.say("\tEnd the topic label.")
				topic = "random" // Go back to default topic
			} else if type_ == "object" {
				rs.say("\tEnd the object label.")
				inobj = false
			}
		case "+": // +Trigger
			rs.say("\tTrigger pattern: %s", line)

			// Initialize the trigger tree.
			curTrig = new(astTrigger)
			curTrig.trigger = line
			curTrig.reply = []string{}
			curTrig.condition = []string{}
			curTrig.redirect = ""
			curTrig.previous = isThat
			ast.topics[topic].triggers = append(ast.topics[topic].triggers, curTrig)
		case "-": // -Response
			if curTrig == nil {
				rs.warnSyntax("Response found before trigger", filename, lineno)
				continue
			}

			rs.say("\tResponse: %s", line)
			curTrig.reply = append(curTrig.reply, line)
		case "*": // *condition
			if curTrig == nil {
				rs.warnSyntax("Condition found before trigger", filename, lineno)
				continue
			}

			rs.say("\tCondition: %s", line)
			curTrig.condition = append(curTrig.condition, line)
		case "%": // %Previous
			continue // This was handled above
		case "^": // ^Continue
			continue // This was handled above
		case "@": // @Redirect
			if curTrig == nil {
				rs.warnSyntax("Redirect found before trigger", filename, lineno)
				continue
			}

			rs.say("\tRedirect response to: %s", line)
			curTrig.redirect = line
		default:
			rs.warnSyntax("Unknown command '%s'", filename, lineno, cmd)
		}
	}

	return ast, nil
}
