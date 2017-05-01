package rivescript

import "regexp"

// Commonly used regular expressions.
var (
	reWeight        = regexp.MustCompile(`\s*\{weight=(\d+)\}\s*`)
	reInherits      = regexp.MustCompile(`\{inherits=(\d+)\}`)
	reMeta          = regexp.MustCompile(`[\<>]+`)
	reSymbols       = regexp.MustCompile(`[.?,!;:@#$%^&*()]+`)
	reNasties       = regexp.MustCompile(`[^A-Za-z0-9 ]`)
	reZerowidthstar = regexp.MustCompile(`^\*$`)
	reOptional      = regexp.MustCompile(`\[(.+?)\]`)
	reArray         = regexp.MustCompile(`@(.+?)\b`)
	reReplyArray    = regexp.MustCompile(`\(@([A-Za-z0-9_]+)\)`)
	reBotvars       = regexp.MustCompile(`<bot (.+?)>`)
	reUservars      = regexp.MustCompile(`<get (.+?)>`)
	reInput         = regexp.MustCompile(`<input([1-9])>`)
	reReply         = regexp.MustCompile(`<reply([1-9])>`)
	reRandom        = regexp.MustCompile(`\{random\}(.+?)\{/random\}`)

	// Self-contained tags like <set> that contain no nested tag.
	reAnytag = regexp.MustCompile(`<([^<]+?)>`)

	reTopic     = regexp.MustCompile(`\{topic=(.+?)\}`)
	reRedirect  = regexp.MustCompile(`\{@(.+?)\}`)
	reCall      = regexp.MustCompile(`<call>(.+?)</call>`)
	reCondition = regexp.MustCompile(`^(.+?)\s+(==|eq|!=|ne|<>|<|<=|>|>=)\s+(.*?)$`)
	reSet       = regexp.MustCompile(`<set (.+?)=(.+?)>`)

	// Placeholders used during substitutions.
	rePlaceholder = regexp.MustCompile(`\x00(\d+)\x00`)
)
