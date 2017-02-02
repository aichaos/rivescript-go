package rivescript

import "errors"

// The types of errors returned by RiveScript.
var (
	ErrDeepRecursion    = errors.New("Deep Recursion Detected")
	ErrRepliesNotSorted = errors.New("Replies Not Sorted")
	ErrNoDefaultTopic   = errors.New("No default topic 'random' was found")
	ErrNoTriggerMatched = errors.New("No Trigger Matched")
	ErrNoReplyFound     = errors.New("The trigger matched but yielded no reply")
)
