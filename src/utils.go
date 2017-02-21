package rivescript

// Miscellaneous utility functions.

import (
	"fmt"
	"regexp"
	"strings"
)

// randomInt gets a random number using RiveScript's internal RNG.
func (rs *RiveScript) randomInt(max int) int {
	rs.randomLock.Lock()
	defer rs.randomLock.Unlock()
	return rs.rng.Intn(max)
}

// wordCount counts the number of real words in a string.
func wordCount(pattern string, all bool) int {
	var words []string
	if all {
		words = strings.Fields(pattern) // Splits at whitespaces
	} else {
		words = regSplit(pattern, `[\s\*\#\_\|]+`)
	}

	wc := 0
	for _, word := range words {
		if len(word) > 0 {
			wc++
		}
	}

	return wc
}

// stripNasties strips special characters out of a string.
func stripNasties(pattern string) string {
	return reNasties.ReplaceAllString(pattern, "")
}

// isAtomic tells you whether a string is atomic or not.
func isAtomic(pattern string) bool {
	// Atomic triggers don't contain any wildcards or parenthesis or anything of
	// the sort. We don't need to test the full character set, just left brackets
	// will do.
	specials := []string{"*", "#", "_", "(", "[", "<", "@"}
	for _, special := range specials {
		if strings.Index(pattern, special) > -1 {
			return false
		}
	}
	return true
}

// stringFormat formats a string.
func stringFormat(format string, text string) string {
	if format == "uppercase" {
		return strings.ToUpper(text)
	} else if format == "lowercase" {
		return strings.ToLower(text)
	} else if format == "sentence" {
		if len(text) > 1 {
			return strings.ToUpper(text[0:1]) + strings.ToLower(text[1:])
		}
		return strings.ToUpper(text)
	} else if format == "formal" {
		words := strings.Split(text, " ")
		result := []string{}
		for _, word := range words {
			if len(word) > 1 {
				result = append(result, strings.ToUpper(word[0:1])+strings.ToLower(word[1:]))
			} else {
				result = append(result, strings.ToUpper(word))
			}
		}
		return strings.Join(result, " ")
	}
	return text
}

// quotemeta escapes a string for use in a regular expression.
func quotemeta(pattern string) string {
	unsafe := `\.+*?[^]$(){}=!<>|:`
	for _, char := range strings.Split(unsafe, "") {
		pattern = strings.Replace(pattern, char, fmt.Sprintf("\\%s", char), -1)
	}
	return pattern
}

// Sort a list of strings by length. Callable like:
// sort.Sort(byLength(strings)) where strings is a []string type.
// https://gobyexample.com/sorting-by-functions
type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

// regSplit splits a string using a regular expression.
// http://stackoverflow.com/questions/4466091/split-string-using-regular-expression-in-go
func regSplit(text string, delimiter string) []string {
	reg := regexp.MustCompile(delimiter)
	indexes := reg.FindAllStringIndex(text, -1)
	lastStart := 0
	result := make([]string, len(indexes)+1)
	for i, element := range indexes {
		result[i] = text[lastStart:element[0]]
		lastStart = element[1]
	}
	result[len(indexes)] = text[lastStart:len(text)]
	return result
}

/*
regReplace quickly replaces a string using a regular expression.

This is a convenience function to do a RegExp-based find/replace in a
JavaScript-like fashion. Example usage:

message = regReplace(message, `hello (.+?)`, "goodbye $1")

Params:

	input: The input string to run the substitution against.
	pattern: Literal string for a regular expression pattern.
	result: String to substitute the result out for. You can use capture group
	        placeholders like $1 in this string.
*/
func regReplace(input string, pattern string, result string) string {
	reg := regexp.MustCompile(pattern)
	match := reg.FindStringSubmatch(input)
	input = reg.ReplaceAllString(input, result)
	if len(match) > 1 {
		for i := range match[1:] {
			input = strings.Replace(input, fmt.Sprintf("$%d", i), match[i], -1)
		}
	}
	return input
}
