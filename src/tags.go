package rivescript

// Tag processing functions.

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/aichaos/rivescript-go/sessions"
)

// formatMessage formats a user's message for safe processing.
func (rs *RiveScript) formatMessage(msg string, botReply bool) string {
	// Lowercase it.
	msg = strings.ToLower(msg)

	// Run substitutions and sanitize what's left.
	msg = rs.substitute(msg, rs.sub, rs.sorted.sub)

	// In UTF-8 mode, only strip metacharacters and HTML brackets (to protect
	// against obvious XSS attacks).
	if rs.UTF8 {
		msg = reMeta.ReplaceAllString(msg, "")
		msg = rs.UnicodePunctuation.ReplaceAllString(msg, "")

		// For the bot's reply, also strip common punctuation.
		if botReply {
			msg = reSymbols.ReplaceAllString(msg, "")
		}
	} else {
		// For everything else, strip all non-alphanumerics.
		msg = stripNasties(msg)
	}

	return msg
}

// triggerRegexp prepares a trigger pattern for the regular expression engine.
func (rs *RiveScript) triggerRegexp(username string, pattern string) string {
	// If the trigger is simply '*' then the * needs to become (.*?)
	// to match the blank string too.
	pattern = reZerowidthstar.ReplaceAllString(pattern, "<zerowidthstar>")

	// Simple replacements.
	pattern = strings.Replace(pattern, "*", `(.+?)`, -1)
	pattern = strings.Replace(pattern, "#", `(\d+?)`, -1)
	pattern = strings.Replace(pattern, "_", `(\w+?)`, -1)
	pattern = reWeight.ReplaceAllString(pattern, "")   // Remove {weight} tags
	pattern = reInherits.ReplaceAllString(pattern, "") // Remove {inherits} tags
	pattern = strings.Replace(pattern, "<zerowidthstar>", `(.*?)`, -1)

	// UTF-8 mode special characters.
	if rs.UTF8 {
		// Literal @ symbols (like in an e-mail address) conflict with arrays.
		pattern = strings.Replace(pattern, `\@`, `\u0040`, -1)
	}

	// Optionals.
	match := reOptional.FindStringSubmatch(pattern)
	var giveup uint
	for len(match) > 0 {
		giveup++
		if giveup > rs.Depth {
			rs.warn("Infinite loop when trying to process optionals in a trigger!")
			return ""
		}

		parts := strings.Split(match[1], "|")
		opts := []string{}
		for _, p := range parts {
			opts = append(opts, fmt.Sprintf(`(?:\s|\b)+%s(?:\s|\b)+`, p))
		}

		// If this optional had a star or anything in it, make it non-matching.
		pipes := strings.Join(opts, "|")
		pipes = strings.Replace(pipes, `(.+?)`, `(?:.+?)`, -1)
		pipes = strings.Replace(pipes, `(\d+?)`, `(?:\d+?)`, -1)
		pipes = strings.Replace(pipes, `(\w+?)`, `(?:\w+?)`, -1)

		pattern = regReplace(pattern,
			fmt.Sprintf(`\s*\[%s\]\s*`, quotemeta(match[1])),
			fmt.Sprintf(`(?:%s|(?:\s|\b)+)`, pipes))
		match = reOptional.FindStringSubmatch(pattern)
	}

	// _ wildcards can't match numbers! Quick note on why I did it this way:
	// the initial replacement above (_ => (\w+?)) needs to be \w because the
	// square brackets in [\s\d] will confuse the optionals logic just above.
	// So then we switch it back down here. Also, we don't just use \w+ because
	// that matches digits, and similarly [A-Za-z] doesn't work with Unicode.
	pattern = strings.Replace(pattern, `\w`, `[^\s\d]`, -1)

	// Filter in arrays.
	giveup = 0
	for strings.Index(pattern, "@") > -1 {
		giveup++
		if giveup > rs.Depth {
			break
		}

		match := reArray.FindStringSubmatch(pattern)
		if len(match) > 0 {
			name := match[1]
			rep := ""
			if _, ok := rs.array[name]; ok {
				rep = fmt.Sprintf(`(?:%s)`, strings.Join(rs.array[name], "|"))
			}
			pattern = strings.Replace(pattern, fmt.Sprintf(`@%s`, name), rep, -1)
		}
	}

	// Filter in bot variables.
	giveup = 0
	for strings.Index(pattern, "<bot ") > -1 {
		giveup++
		if giveup > rs.Depth {
			break
		}

		match := reBotvars.FindStringSubmatch(pattern)
		if len(match) > 0 {
			name := match[1]
			rep := ""
			if _, ok := rs.vars[name]; ok {
				rep = stripNasties(rs.vars[name])
			}
			pattern = strings.Replace(pattern, fmt.Sprintf(`<bot %s>`, name), strings.ToLower(rep), -1)
		}
	}

	// Filter in user variables.
	giveup = 0
	for strings.Index(pattern, "<get ") > -1 {
		giveup++
		if giveup > rs.Depth {
			break
		}

		match := reUservars.FindStringSubmatch(pattern)
		if len(match) > 0 {
			name := match[1]

			value, err := rs.sessions.Get(username, name)
			if err != nil {
				value = UNDEFINED
			}

			pattern = strings.Replace(pattern, fmt.Sprintf(`<get %s>`, name), strings.ToLower(value), -1)
		}
	}

	// Filter in <input> and <reply> tags.
	giveup = 0
	pattern = strings.Replace(pattern, "<input>", "<input1>", -1)
	pattern = strings.Replace(pattern, "<reply>", "<reply1>", -1)
	for strings.Index(pattern, "<input") > -1 || strings.Index(pattern, "<reply") > -1 {
		giveup++
		if giveup > 50 {
			break
		}

		for i := 1; i <= sessions.HistorySize; i++ {
			inputPattern := fmt.Sprintf("<input%d>", i)
			replyPattern := fmt.Sprintf("<reply%d>", i)
			history, err := rs.sessions.GetHistory(username)
			if err == nil {
				pattern = strings.Replace(pattern, inputPattern, history.Input[i-1], -1)
				pattern = strings.Replace(pattern, replyPattern, history.Reply[i-1], -1)
			} else {
				pattern = strings.Replace(pattern, inputPattern, UNDEFINED, -1)
				pattern = strings.Replace(pattern, replyPattern, UNDEFINED, -1)
			}
		}
	}

	// Recover escaped Unicode symbols.
	if rs.UTF8 && strings.Index(pattern, `\u`) > -1 {
		// TODO: make this more general
		pattern = strings.Replace(pattern, `\u0040`, "@", -1)
	}

	return pattern
}

/*
processTags processes tags in a reply element.

Params:

	username: The name of the user.
	message: The user's message.
	reply: The reply element to process tags on.
	st: Array of matched stars in the trigger.
	bst: Array of matched bot stars in a %Previous.
	step: Recursion depth counter.
*/
func (rs *RiveScript) processTags(username string, message string, reply string, st []string, bst []string, step uint) string {
	// Prepare the stars and botstars.
	stars := []string{""}
	stars = append(stars, st...)
	botstars := []string{""}
	botstars = append(botstars, bst...)
	if len(stars) == 1 {
		stars = append(stars, UNDEFINED)
	}
	if len(botstars) == 1 {
		botstars = append(botstars, UNDEFINED)
	}

	// Turn arrays into randomized sets.
	match := reReplyArray.FindStringSubmatch(reply)
	var giveup uint
	for len(match) > 0 {
		giveup++
		if giveup > rs.Depth {
			rs.warn("Infinite loop interpolating arrays into reply!")
			break
		}

		name := match[1]
		var result string
		if value, ok := rs.array[name]; ok {
			result = "{random}" + strings.Join(value, "|") + "{/random}"
		} else {
			// Dummy it out so we can reinsert it, as-is, later.
			result = "\x00@" + name + "\x00"
		}

		reply = strings.Replace(reply, "(@"+name+")", result, -1)
		match = reReplyArray.FindStringSubmatch(reply)
	}

	// Re-insert dummied out (non-existant) arrays from the above block.
	reply = regReplace(reply, "\x00@([A-Za-z0-9_]+)\x00", "(@$1)")

	// Tag shortcuts.
	reply = strings.Replace(reply, "<person>", "{person}<star>{/person}", -1)
	reply = strings.Replace(reply, "<@>", "{@<star>}", -1)
	reply = strings.Replace(reply, "<formal>", "{formal}<star>{/formal}", -1)
	reply = strings.Replace(reply, "<sentence>", "{sentence}<star>{/sentence}", -1)
	reply = strings.Replace(reply, "<uppercase>", "{uppercase}<star>{/uppercase}", -1)
	reply = strings.Replace(reply, "<lowercase>", "{lowercase}<star>{/lowercase}", -1)

	// Weight and star tags.
	reply = reWeight.ReplaceAllString(reply, "") // Remove {weight} tags.
	reply = strings.Replace(reply, "<star>", stars[1], -1)
	reply = strings.Replace(reply, "<botstar>", botstars[1], -1)
	for i := 1; i < len(stars); i++ {
		reply = strings.Replace(reply, fmt.Sprintf("<star%d>", i), stars[i], -1)
	}
	for i := 1; i < len(botstars); i++ {
		reply = strings.Replace(reply, fmt.Sprintf("<botstar%d>", i), botstars[i], -1)
	}

	// <input> and <reply>
	reply = strings.Replace(reply, "<input>", "<input1>", -1)
	reply = strings.Replace(reply, "<reply>", "<reply1>", -1)
	history, err := rs.sessions.GetHistory(username)
	if err == nil {
		for i := 1; i <= sessions.HistorySize; i++ {
			reply = strings.Replace(reply, fmt.Sprintf("<input%d>", i), history.Input[i-1], -1)
			reply = strings.Replace(reply, fmt.Sprintf("<reply%d>", i), history.Reply[i-1], -1)
		}
	}

	// <id> and escape codes.
	reply = strings.Replace(reply, "<id>", username, -1)
	reply = strings.Replace(reply, `\s`, " ", -1)
	reply = strings.Replace(reply, `\n`, "\n", -1)
	reply = strings.Replace(reply, `\#`, "#", -1)

	// {random}
	match = reRandom.FindStringSubmatch(reply)
	giveup = 0
	for len(match) > 0 {
		giveup++
		if giveup > rs.Depth {
			rs.warn("Infinite loop looking for random tag!")
			break
		}

		var random []string
		text := match[1]
		if strings.Index(text, "|") > -1 {
			random = strings.Split(text, "|")
		} else {
			random = strings.Split(text, " ")
		}

		output := ""
		if len(random) > 0 {
			output = random[rs.randomInt(len(random))]
		}

		reply = strings.Replace(reply, fmt.Sprintf("{random}%s{/random}", text), output, -1)
		match = reRandom.FindStringSubmatch(reply)
	}

	// Person substitution and string formatting.
	formats := []string{"person", "formal", "sentence", "uppercase", "lowercase"}
	for _, format := range formats {
		formatRegexp := regexp.MustCompile(fmt.Sprintf(`\{%s\}(.+?)\{/%s\}`, format, format))
		match = formatRegexp.FindStringSubmatch(reply)
		giveup = 0
		for len(match) > 0 {
			giveup++
			if giveup > rs.Depth {
				rs.warn("Infinite loop looking for %s tag!", format)
				break
			}

			content := match[1]
			var replace string
			if format == "person" {
				replace = rs.substitute(content, rs.person, rs.sorted.person)
			} else {
				replace = stringFormat(format, content)
			}

			reply = strings.Replace(reply, fmt.Sprintf("{%s}%s{/%s}", format, content, format), replace, -1)
			match = formatRegexp.FindStringSubmatch(reply)
		}
	}

	// Handle all variable-related tags with an iterative regexp approach to
	// allow for nesting of tags in arbitrary ways (think <set a=<get b>>)
	// Dummy out the <call> tags first, because we don't handle them here.
	reply = strings.Replace(reply, "<call>", "{__call__}", -1)
	reply = strings.Replace(reply, "</call>", "{/__call__}", -1)
	for {
		// Look for tags that don't contain any other tags inside them.
		matcher := reAnytag.FindStringSubmatch(reply)
		if len(matcher) == 0 {
			break // No tags left!
		}

		match := matcher[1]
		parts := strings.Split(match, " ")
		tag := strings.ToLower(parts[0])
		data := ""
		if len(parts) > 1 {
			data = strings.Join(parts[1:], " ")
		}
		insert := ""

		// Handle the various types of tags.
		if tag == "bot" || tag == "env" {
			// <bot> and <env> work similarly
			var target map[string]string
			if tag == "bot" {
				target = rs.vars
			} else {
				target = rs.global
			}

			if strings.Index(data, "=") > -1 {
				// Assigning the value.
				parts := strings.Split(data, "=")
				rs.say("Assign %s variable %s = %s", tag, parts[0], parts[1])
				target[parts[0]] = parts[1]
			} else {
				// Getting a bot/env variable.
				if _, ok := target[data]; ok {
					insert = target[data]
				} else {
					insert = UNDEFINED
				}
			}
		} else if tag == "set" {
			// <set> user vars
			parts := strings.Split(data, "=")
			if len(parts) > 1 {
				rs.say("Set uservar %s = %s", parts[0], parts[1])
				rs.sessions.Set(username, map[string]string{parts[0]: parts[1]})
			} else {
				rs.warn("Malformed <set> tag: %s", match)
			}
		} else if tag == "add" || tag == "sub" || tag == "mult" || tag == "div" {
			// Math operator tags.
			parts := strings.Split(data, "=")
			name := parts[0]
			strValue := parts[1]

			// Initialize the variable?
			var origStr string
			origStr, err = rs.sessions.Get(username, name)
			if err != nil {
				rs.sessions.Set(username, map[string]string{name: "0"})
				origStr = "0"
			}

			// Sanity check.
			var value int
			value, err = strconv.Atoi(strValue)
			abort := false
			if err != nil {
				insert = fmt.Sprintf("[ERR: Math can't %s non-numeric value %s]", tag, strValue)
				abort = true
			}
			var orig int
			orig, err = strconv.Atoi(origStr)
			if err != nil {
				insert = fmt.Sprintf("[ERR: Math can't %s non-numeric user variable %s]", tag, name)
				abort = true
			}

			if !abort {
				result := orig
				if tag == "add" {
					result += value
				} else if tag == "sub" {
					result -= value
				} else if tag == "mult" {
					result *= value
				} else if tag == "div" {
					if value == 0 {
						insert = "[ERR: Can't Divide By Zero]"
					} else {
						result /= value
					}
				}

				if len(insert) == 0 {
					// Save it to their account.
					rs.sessions.Set(username, map[string]string{name: strconv.Itoa(result)})
				}
			}
		} else if tag == "get" {
			// <get> user vars
			insert, err = rs.sessions.Get(username, data)
			if err != nil {
				insert = UNDEFINED
			}
		} else {
			// Unrecognized tag; preserve it.
			insert = fmt.Sprintf("\x00%s\x01", match)
		}

		reply = strings.Replace(reply, fmt.Sprintf("<%s>", match), insert, -1)
	}

	// Recover mangled HTML-like tags.
	reply = strings.Replace(reply, "\x00", "<", -1)
	reply = strings.Replace(reply, "\x01", ">", -1)

	// Topic setter.
	match = reTopic.FindStringSubmatch(reply)
	giveup = 0
	for len(match) > 0 {
		giveup++
		if giveup > rs.Depth {
			rs.warn("Infinite loop looking for topic tag!")
			break
		}

		name := match[1]
		rs.sessions.Set(username, map[string]string{"topic": name})
		reply = strings.Replace(reply, fmt.Sprintf("{topic=%s}", name), "", -1)
		match = reTopic.FindStringSubmatch(reply)
	}

	// Inline redirector.
	match = reRedirect.FindStringSubmatch(reply)
	giveup = 0
	for len(match) > 0 {
		giveup++
		if giveup > rs.Depth {
			rs.warn("Infinite loop looking for redirect tag!")
			break
		}

		target := match[1]
		rs.say("Inline redirection to: %s", target)
		subreply, err := rs.getReply(username, strings.TrimSpace(target), false, step+1)
		if err != nil {
			subreply = err.Error()
		}
		reply = strings.Replace(reply, fmt.Sprintf("{@%s}", target), subreply, -1)
		match = reRedirect.FindStringSubmatch(reply)
	}

	// Object caller.
	reply = strings.Replace(reply, "{__call__}", "<call>", -1)
	reply = strings.Replace(reply, "{/__call__}", "</call>", -1)
	match = reCall.FindStringSubmatch(reply)
	giveup = 0
	for len(match) > 0 {
		giveup++
		if giveup > rs.Depth {
			rs.warn("Infinite loop looking for call tags!")
			break
		}

		text := strings.TrimSpace(match[1])
		parts := strings.Split(text, " ")
		obj := parts[0]
		args := []string{}
		if len(parts) > 1 {
			args = parts[1:]
		}

		// Do we know this object?
		var output string
		if _, ok := rs.subroutines[obj]; ok {
			// It exists as a native Go macro.
			output = rs.subroutines[obj](rs, args)
		} else if _, ok := rs.objlangs[obj]; ok {
			lang := rs.objlangs[obj]
			output = rs.handlers[lang].Call(obj, args)
		} else {
			output = "[ERR: Object Not Found]"
		}

		reply = strings.Replace(reply, fmt.Sprintf("<call>%s</call>", match[1]), output, -1)
		match = reCall.FindStringSubmatch(reply)
	}

	return reply
}

// substitute applies a substitution to an input message.
func (rs *RiveScript) substitute(message string, subs map[string]string, sorted []string) string {
	// Safety checking.
	if len(subs) == 0 {
		return message
	}

	// Make placeholders each time we substitute something.
	ph := []string{}
	pi := 0

	// fmt.Printf("Running substitutitons on input message: %s\n", message)

	for _, pattern := range sorted {
		result := subs[pattern]
		qm := quotemeta(pattern)
		// fmt.Printf("Pattern: %s; Result: %s\n", pattern, result)

		// Make a placeholder.
		ph = append(ph, result)
		placeholder := fmt.Sprintf("\x00%d\x00", pi)
		pi++

		// Run substitutions.
		message = regReplace(message, fmt.Sprintf(`^%s$`, qm), placeholder)
		message = regReplace(message, fmt.Sprintf(`^%s(\W+)`, qm), fmt.Sprintf("%s$1", placeholder))
		message = regReplace(message, fmt.Sprintf(`(\W+)%s(\W+)`, qm), fmt.Sprintf("$1%s$2", placeholder))
		message = regReplace(message, fmt.Sprintf(`(\W+)%s$`, qm), fmt.Sprintf("$1%s", placeholder))
	}

	// Convert the placeholders back in.
	var tries uint
	for strings.Index(message, "\x00") > -1 {
		tries++
		if tries > rs.Depth {
			rs.warn("Too many loops in substitution placeholders!")
			break
		}

		match := rePlaceholder.FindStringSubmatch(message)
		if len(match) > 0 {
			i, _ := strconv.Atoi(match[1])
			result := ph[i]
			message = strings.Replace(message, fmt.Sprintf("\x00%d\x00", i), result, -1)
		}
	}

	// fmt.Printf("Final message: %s", message)

	return message
}
