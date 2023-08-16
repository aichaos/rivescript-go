package rivescript

import "errors"

// The types of errors returned by RiveScript.
var (
	ErrDeepRecursion    = errors.New("deep recursion detected")
	ErrRepliesNotSorted = errors.New("replies not sorted")
	ErrNoDefaultTopic   = errors.New("no default topic 'random' was found")
	ErrNoTriggerMatched = errors.New("no trigger matched")
	ErrNoReplyFound     = errors.New("the trigger matched but yielded no reply")
)
