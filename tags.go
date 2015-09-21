package rivescript

// Tag processing functions.

import (
	"fmt"
	"strconv"
	"strings"
)

// formatMessage formats a user's message for safe processing.
func (rs RiveScript) formatMessage(msg string, botReply bool) string {
	// Lowercase it.
	msg = strings.ToLower(msg)

	// Run substitutions and sanitize what's left.
	msg = rs.substitute(msg, rs.sub, rs.sorted.sub)

	// In UTF-8 mode, only strip metacharacters and HTML brackets (to protect
	// against obvious XSS attacks).
	if rs.UTF8 {
		msg = re_meta.ReplaceAllString(msg, "")
		msg = rs.UnicodePunctuation.ReplaceAllString(msg, "")

		// For the bot's reply, also strip common punctuation.
		if botReply {
			msg = re_symbols.ReplaceAllString(msg, "")
		}
	} else {
		// For everything else, strip all non-alphanumerics.
		msg = stripNasties(msg)
	}

	return msg
}

// triggerRegexp prepares a trigger pattern for the regular expression engine.
func (rs RiveScript) triggerRegexp(username string, pattern string) string {
	// If the trigger is simply '*' then the * needs to become (.*?)
	// to match the blank string too.
	pattern = re_zerowidthstar.ReplaceAllString(pattern, "<zerowidthstar>")

	// Simple replacements.
	pattern = strings.Replace(pattern, "*", `(.+?)`, -1)
	pattern = strings.Replace(pattern, "#", `(\d+?)`, -1)
	pattern = strings.Replace(pattern, "_", `(\w+?)`, -1)
	pattern = re_weight.ReplaceAllString(pattern, "") // Remove {weight} tags
	pattern = strings.Replace(pattern, "<zerowidthstar>", `(.*?)`, -1)

	// UTF-8 mode special characters.
	if rs.UTF8 {
		// Literal @ symbols (like in an e-mail address) conflict with arrays.
		pattern = strings.Replace(pattern, "@", `\u0040`, -1)
	}

	// Optionals.
	match := re_optional.FindStringSubmatch(pattern)
	giveup := 0
	for len(match) > 0 {
		giveup++
		if giveup > rs.Depth {
			rs.warn("Infinite loop when trying to process optionals in a trigger!")
			return ""
		}

		parts := strings.Split(match[1], "|")
		opts := []string{}
		for _, p := range parts {
			opts = append(opts, fmt.Sprintf(`\s*%s\s*`, p))
		}
		opts = append(opts, `\s*`)

		// If this optional had a star or anything in it, make it non-matching.
		pipes := strings.Join(opts, "|")
		pipes = strings.Replace(pipes, `(.+?)`, `(?:.+?)`, -1)
		pipes = strings.Replace(pipes, `(\d+?)`, `(?:\d+?)`, -1)
		pipes = strings.Replace(pipes, `(\d+?)`, `(?:\w+?)`, -1)

		pattern = regReplace(pattern, fmt.Sprintf(`\s*\[%s\]\s*`, quotemeta(match[1])), fmt.Sprintf("(?:%s)", pipes))
		match = re_optional.FindStringSubmatch(pattern)
	}

	// _ wildcards can't match numbers! Quick note on why I did it this way:
	// the initial replacement above (_ => (\w+?)) needs to be \w because the
	// square brackets in [A-Za-z] will confuse the optionals logic just above.
	// So then we switch it back down here.
	pattern = strings.Replace(pattern, `\w`, "[A-Za-z]", -1)

	// Filter in arrays.
	giveup = 0
	for strings.Index(pattern, "@") > -1 {
		giveup++
		if giveup > rs.Depth {
			break
		}

		match := re_array.FindStringSubmatch(pattern)
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

		match := re_botvar.FindStringSubmatch(pattern)
		if len(match) > 0 {
			name := match[1]
			rep := ""
			if _, ok := rs.var_[name]; ok {
				rep = stripNasties(rs.var_[name])
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

		match := re_uservar.FindStringSubmatch(pattern)
		if len(match) > 0 {
			name := match[1]
			rep := "undefined"
			if _, ok := rs.users[username].data[name]; ok {
				rep = rs.users[username].data[name]
			}
			pattern = strings.Replace(pattern, fmt.Sprintf(`<get %s>`, name), strings.ToLower(rep), -1)
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

		for i := 1; i <= 9; i++ {
			inputPattern := fmt.Sprintf("<input%d>", i)
			replyPattern := fmt.Sprintf("<reply%d>", i)
			pattern = strings.Replace(pattern, inputPattern, rs.users[username].inputHistory[i], -1)
			pattern = strings.Replace(pattern, replyPattern, rs.users[username].replyHistory[i], -1)
		}
	}

	// Recover escaped Unicode symbols.
	if rs.UTF8 && strings.Index(pattern, `\u`) > -1 {
		// TODO: make this more general
		pattern = strings.Replace(pattern, `\u0040`, "@", -1)
	}

	return pattern
}

// substitute applies a substitution to an input message.
func (rs RiveScript) substitute(message string, subs map[string]string, sorted []string) string {
	// Safety checking.
	if len(subs) == 0 {
		rs.warn("You forgot to call sortReplies()!")
		return ""
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
		// fmt.Printf("BEFORE: %s\n", message)
		message = regReplace(message, fmt.Sprintf(`^%s$`, qm), placeholder)
		message = regReplace(message, fmt.Sprintf(`^%s(\W+)`, qm), fmt.Sprintf("%s$1", placeholder))
		message = regReplace(message, fmt.Sprintf(`(\W+)%s(\W+)`, qm), fmt.Sprintf("$1%s$2", placeholder))
		message = regReplace(message, fmt.Sprintf(`(\W+)%s$`, qm), fmt.Sprintf("$1%s", placeholder))
		// fmt.Printf("AFTER: %s\n", message)
	}

	// Convert the placeholders back in.
	tries := 0
	for strings.Index(message, "\x00") > -1 {
		tries++
		if tries > rs.Depth {
			rs.warn("Too many loops in substitution placeholders!")
			break
		}

		match := re_placeholder.FindStringSubmatch(message)
		if len(match) > 0 {
			i, _ := strconv.Atoi(match[1])
			result := ph[i]
			message = strings.Replace(message, fmt.Sprintf("\x00%d\x00", i), result, -1)
		}
	}

	// fmt.Printf("Final message: %s", message)

	return message
}
