package rivescript

// Miscellaneous utility functions.

import (
	"regexp"
	"strings"
)

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
			wc += 1
		}
	}

	return wc
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
