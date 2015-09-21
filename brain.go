package rivescript

import (
	"fmt"
	"math/rand"
	re "regexp"
	"strconv"
	"strings"
)

// Brain logic for RiveScript

/*
Reply fetches a reply from the bot for a user's message.

Params:
- username: The name of the user requesting a reply.
- message: The user's message.
*/
func (rs RiveScript) Reply(username string, message string) string {
	rs.say("Asked to reply to [%s] %s", username, message)

	// Initialize a user profile for this user?
	if _, ok := rs.users[username]; !ok {
		rs.users[username] = NewUser()
	}

	// Store the current user's ID.
	rs.currentUser = username

	// Format their message.
	message = rs.formatMessage(message, false)
	reply := ""

	// If the BEGIN block exists, consult it first.
	if _, ok := rs.topics["__begin__"]; ok {
		begin := rs.getReply(username, "request", true, 0)

		// OK to continue?
		if strings.Index(begin, "{ok}") > -1 {
			reply = rs.getReply(username, message, false, 0)
			begin = strings.NewReplacer("{ok}", reply).Replace(begin)
		}

		reply = begin
		// reply = rs.processTags(reply)
	} else {
		reply = rs.getReply(username, message, false, 0)
	}

	// Save their message history.
	user := rs.users[username]
	user.inputHistory = user.inputHistory[:len(user.inputHistory)-1]    // Pop
	user.inputHistory = append([]string{message}, user.inputHistory...) // Unshift
	user.replyHistory = user.replyHistory[:len(user.replyHistory)-1]    // Pop
	user.replyHistory = append([]string{reply}, user.replyHistory...)   // Unshift

	// Unset the current user's ID.
	rs.currentUser = ""

	return reply
}

/*
getReply is the internal logic behind Reply().

Params:
- username: The name of the user requesting a reply.
- message: The user's message.
- isBegin: Whether this reply is for the "BEGIN Block" context or not.
- step: Recursion depth counter.
*/
func (rs RiveScript) getReply(username string, message string, isBegin bool, step int) string {
	// Needed to sort replies?
	if len(rs.sorted.topics) == 0 {
		rs.warn("You forgot to call SortReplies()!")
		return "ERR: Replies Not Sorted"
	}

	// Collect data on this user.
	topic := rs.users[username].data["topic"]
	stars := []string{}
	thatStars := []string{} // For %Previous
	reply := ""

	// Avoid letting them fall into a missing topic.
	if _, ok := rs.topics[topic]; !ok {
		rs.warn("User %s was in an empty topic named '%s'", username, topic)
		rs.users[username].data["topic"] = "random"
		topic = "random"
	}

	// Avoid deep recursion.
	if step > rs.Depth {
		return "ERR: Deep Recursion Detected"
	}

	// Are we in the BEGIN block?
	if isBegin {
		topic = "__begin__"
	}

	// More topic sanity checking.
	if _, ok := rs.topics[topic]; !ok {
		// This was handled before, which would mean topic=random and it doesn't
		// exist. Serious issue!
		return "ERR: No default topic 'random' was found!"
	}

	// Create a pointer for the matched data when we find it.
	var matched *astTrigger
	matchedTrigger := ""
	foundMatch := false

	// See if there were any %Previous's in this topic, or any topic related to
	// it. This should only be done the first time -- not during a recursive
	// redirection. This is because in a redirection, "lastReply" is still gonna
	// be the same as it was the first time, resulting in an infinite loop!
	if step == 0 {
		// TODO
	}

	// Search their topic for a match to their trigger.
	if !foundMatch {
		rs.say("Searching their topic for a match...")
		for _, trig := range rs.sorted.topics[topic] {
			pattern := trig.trigger
			regexp := rs.triggerRegexp(username, pattern)
			rs.say("Try to match \"%s\" against %s", message, pattern)

			// If the trigger is atomic, we don't need to bother with the regexp engine.
			isMatch := false
			if isAtomic(pattern) && message == regexp {
				isMatch = true
			} else {
				// Non-atomic triggers always need the regexp.
				matcher := re.MustCompile(fmt.Sprintf("^%s$", regexp))
				match := matcher.FindStringSubmatch(message)
				if len(match) > 0 {
					// The regexp matched!
					isMatch = true

					// Collect the stars.
					if len(match) > 1 {
						for i, _ := range match[1:] {
							stars = append(stars, match[i])
						}
					}
				}
			}

			// A match somehow?
			if isMatch {
				rs.say("Found a match!")

				// Keep the pointer to this trigger's data.
				matched = trig.pointer
				foundMatch = true
				matchedTrigger = pattern
				break
			}
		}
	}

	// Store what trigger they matched on.
	rs.users[username].lastMatch = matchedTrigger

	// Did we match?
	if foundMatch {
		for range []int{0} { // A single loop so we can break out early
			// See if there are any hard redirects.
			if len(matched.redirect) > 0 {
				rs.say("Redirecting us to %s", matched.redirect)
				redirect := matched.redirect
				//redirect = rs.processTags() TODO
				rs.say("Pretend user said: %s", redirect)
				reply = rs.getReply(username, redirect, isBegin, step+1)
				break
			}

			// Check the conditionals. TODO

			// Have our reply yet?
			if len(reply) > 0 {
				break
			}

			// Process weights in the replies.
			bucket := []string{}
			for _, rep := range matched.reply {
				weight := 1
				match := re_weight.FindStringSubmatch(rep)
				if len(match) > 0 {
					weight , _ = strconv.Atoi(match[1])
					if weight <= 0 {
						rs.warn("Can't have a weight <= 0!")
						weight = 1
					}

					for i := weight; i > 0; i-- {
						bucket = append(bucket, rep)
					}
				} else {
					bucket = append(bucket, rep)
				}
			}

			// Get a random reply.
			if len(bucket) > 0 {
				rs.say("BUCKET LENGTH: %d", len(bucket))
				reply = bucket[rand.Intn(len(bucket))]
			}
			break
		}
	}

	// Still no reply?? Give up with the fallback error replies.
	if !foundMatch {
		reply = "ERR: No Reply Matched"
	} else if len(reply) == 0 {
		reply = "ERR: No Reply Found"
	}

	rs.say("Reply: %s", reply)

	// Process tags for the BEGIN block.
	if isBegin {
		// TODO
	} else {
		// reply = processTags TODO
	}

	_ = thatStars

	return reply
}
