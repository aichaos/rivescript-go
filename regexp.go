package rivescript

// Common regular expressions.

import "regexp"

var re_weight = regexp.MustCompile(`\{weight=(\d+)\}`)
var re_inherits = regexp.MustCompile(`\{inherits=(\d+)\}`)
