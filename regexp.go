package rivescript

// Common regular expressions.

import "regexp"

var re_weight = regexp.MustCompile(`\{weight=(\d+)\}`)
var re_inherits = regexp.MustCompile(`\{inherits=(\d+)\}`)
var re_meta = regexp.MustCompile(`[\<>]+`)
var re_symbols = regexp.MustCompile(`[.?,!;:@#$%^&*()]+`)
var re_nasties = regexp.MustCompile(`[^A-Za-z0-9 ]`)
var re_zerowidthstar = regexp.MustCompile(`^\*$`)
var re_optional = regexp.MustCompile(`\[(.+?)\]`)
var re_array = regexp.MustCompile(`@(.+?)\b`)
var re_botvar = regexp.MustCompile(`<bot (.+?)>`)
var re_uservar = regexp.MustCompile(`<get (.+?)>`)
var re_input = regexp.MustCompile(`<input([1-9])>`)
var re_reply = regexp.MustCompile(`<reply([1-9])>`)
var re_placeholder = regexp.MustCompile(`\x00(\d+)\x00`)
