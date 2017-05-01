package rivescript

import (
	"fmt"
	re "regexp"
	"strconv"
	"strings"
)

/*
Reply fetches a reply from the bot for a user's message.

Parameters

	username: The name of the user requesting a reply.
	message: The user's message.
*/
func (rs *RiveScript) Reply(username, message string) (string, error) {
	rs.say("Asked to reply to [%s] %s", username, message)
	var err error

	// Initialize a user profile for this user?
	rs.sessions.Init(username)

	// Store the current user's ID.
	rs.inReplyContext = true
	rs.currentUser = username

	// Format their message.
	message = rs.formatMessage(message, false)
	var reply string

	// If the BEGIN block exists, consult it first.
	if _, ok := rs.topics["__begin__"]; ok {
		var begin string
		begin, err = rs.getReply(username, "request", true, 0)
		if err != nil {
			return "", err
		}

		// OK to continue?
		if strings.Index(begin, "{ok}") > -1 {
			reply, err = rs.getReply(username, message, false, 0)
			if err != nil {
				return "", err
			}
			begin = strings.NewReplacer("{ok}", reply).Replace(begin)
		}

		reply = begin
		reply = rs.processTags(username, message, reply, []string{}, []string{}, 0)
	} else {
		reply, err = rs.getReply(username, message, false, 0)
		if err != nil {
			return "", err
		}
	}

	// Save their message history.
	rs.sessions.AddHistory(username, message, reply)

	// Unset the current user's ID.
	rs.currentUser = ""
	rs.inReplyContext = false

	return reply, nil
}

/*
getReply is the internal logic behind Reply().

Parameters

	username: The name of the user requesting a reply.
	message: The user's message.
	isBegin: Whether this reply is for the "BEGIN Block" context or not.
	step: Recursion depth counter.
*/
func (rs *RiveScript) getReply(username string, message string, isBegin bool, step uint) (string, error) {
	// Needed to sort replies?
	if len(rs.sorted.topics) == 0 {
		rs.warn("You forgot to call SortReplies()!")
		return "", ErrRepliesNotSorted
	}

	// Collect data on this user.
	topic, err := rs.sessions.Get(username, "topic")
	if err != nil {
		topic = "random"
	}
	stars := []string{}
	thatStars := []string{} // For %Previous
	var reply string

	// Avoid letting them fall into a missing topic.
	if _, ok := rs.topics[topic]; !ok {
		rs.warn("User %s was in an empty topic named '%s'", username, topic)
		rs.sessions.Set(username, map[string]string{"topic": "random"})
		topic = "random"
	}

	// Avoid deep recursion.
	if step > rs.Depth {
		return "", ErrDeepRecursion
	}

	// Are we in the BEGIN block?
	if isBegin {
		topic = "__begin__"
	}

	// More topic sanity checking.
	if _, ok := rs.topics[topic]; !ok {
		// This was handled before, which would mean topic=random and it doesn't
		// exist. Serious issue!
		return "", ErrNoDefaultTopic
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
		allTopics := []string{topic}
		if len(rs.includes[topic]) > 0 || len(rs.inherits[topic]) > 0 {
			// Get ALL the topics!
			allTopics = rs.getTopicTree(topic, 0)
		}

		// Scan them all.
		for _, top := range allTopics {
			rs.say("Checking topic %s for any %%Previous's.", top)

			if len(rs.sorted.thats[top]) > 0 {
				rs.say("There's a %%Previous in this topic!")

				// Get the bot's last reply to the user.
				history, _ := rs.sessions.GetHistory(username)
				lastReply := history.Reply[0]

				// Format the bot's reply the same way as the human's.
				lastReply = rs.formatMessage(lastReply, true)
				rs.say("Bot's last reply: %s", lastReply)

				// See if it's a match.
				for _, trig := range rs.sorted.thats[top] {
					pattern := trig.pointer.previous
					botside := rs.triggerRegexp(username, pattern)
					rs.say("Try to match lastReply (%s) to %s (%s)", lastReply, pattern, botside)

					// Match?
					matcher := re.MustCompile(fmt.Sprintf("^%s$", botside))
					match := matcher.FindStringSubmatch(lastReply)
					if len(match) > 0 {
						// Huzzah! See if OUR message is right too...
						rs.say("Bot side matched!")

						// Collect the bot stars.
						thatStars = []string{}
						if len(match) > 1 {
							for i := range match[1:] {
								thatStars = append(thatStars, match[i+1])
							}
						}

						// Compare the triggers to the user's message.
						userSide := trig.pointer
						regexp := rs.triggerRegexp(username, userSide.trigger)
						rs.say("Try to match %s against %s (%s)", message, userSide.trigger, regexp)

						// If the trigger is atomic, we don't need to deal with the regexp engine.
						isMatch := false
						if isAtomic(userSide.trigger) {
							if message == regexp {
								isMatch = true
							}
						} else {
							matcher := re.MustCompile(fmt.Sprintf("^%s$", regexp))
							match := matcher.FindStringSubmatch(message)
							if len(match) > 0 {
								isMatch = true

								// Get the user's message stars.
								if len(match) > 1 {
									for i := range match[1:] {
										stars = append(stars, match[i+1])
									}
								}
							}
						}

						// Was it a match?
						if isMatch {
							// Keep the trigger pointer.
							matched = userSide
							foundMatch = true
							matchedTrigger = userSide.trigger
							break
						}
					}
				}
			}
		}
	}

	// Search their topic for a match to their trigger.
	if !foundMatch {
		rs.say("Searching their topic for a match...")
		for _, trig := range rs.sorted.topics[topic] {
			pattern := trig.trigger
			regexp := rs.triggerRegexp(username, pattern)
			rs.say("Try to match \"%s\" against %s (%s)", message, pattern, regexp)

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
						for i := range match[1:] {
							stars = append(stars, match[i+1])
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
	rs.sessions.SetLastMatch(username, matchedTrigger)

	// Did we match?
	if foundMatch {
		for range []int{0} { // A single loop so we can break out early
			// See if there are any hard redirects.
			if len(matched.redirect) > 0 {
				rs.say("Redirecting us to %s", matched.redirect)
				redirect := matched.redirect
				redirect = rs.processTags(username, message, redirect, stars, thatStars, 0)
				redirect = strings.ToLower(redirect)
				rs.say("Pretend user said: %s", redirect)
				reply, err = rs.getReply(username, redirect, isBegin, step+1)
				if err != nil {
					return "", err
				}
				break
			}

			// Check the conditionals.
			for _, row := range matched.condition {
				halves := strings.Split(row, "=>")
				if len(halves) == 2 {
					condition := reCondition.FindStringSubmatch(strings.TrimSpace(halves[0]))
					if len(condition) > 0 {
						left := strings.TrimSpace(condition[1])
						eq := condition[2]
						right := strings.TrimSpace(condition[3])
						potreply := strings.TrimSpace(halves[1]) // Potential reply

						// Process tags all around
						left = rs.processTags(username, message, left, stars, thatStars, step)
						right = rs.processTags(username, message, right, stars, thatStars, step)

						// Defaults?
						if len(left) == 0 {
							left = UNDEFINED
						}
						if len(right) == 0 {
							right = UNDEFINED
						}

						rs.say("Check if %s %s %s", left, eq, right)

						// Validate it.
						passed := false
						if eq == "eq" || eq == "==" {
							if left == right {
								passed = true
							}
						} else if eq == "ne" || eq == "!=" || eq == "<>" {
							if left != right {
								passed = true
							}
						} else {
							// Dealing with numbers here.
							iLeft, errLeft := strconv.Atoi(left)
							iRight, errRight := strconv.Atoi(right)
							if errLeft == nil && errRight == nil {
								if eq == "<" && iLeft < iRight {
									passed = true
								} else if eq == "<=" && iLeft <= iRight {
									passed = true
								} else if eq == ">" && iLeft > iRight {
									passed = true
								} else if eq == ">=" && iLeft >= iRight {
									passed = true
								}
							} else {
								rs.warn("Failed to evaluate numeric condition!")
							}
						}

						if passed {
							reply = potreply
							break
						}
					}
				}
			}

			// Have our reply yet?
			if len(reply) > 0 {
				break
			}

			// Process weights in the replies.
			bucket := []string{}
			for _, rep := range matched.reply {
				match := reWeight.FindStringSubmatch(rep)
				if len(match) > 0 {
					weight, _ := strconv.Atoi(match[1])
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
				reply = bucket[rs.randomInt(len(bucket))]
			}
			break
		}
	}

	// Still no reply?? Give up with the fallback error replies.
	if !foundMatch {
		return "", ErrNoTriggerMatched
	} else if len(reply) == 0 {
		return "", ErrNoReplyFound
	}

	rs.say("Reply: %s", reply)

	// Process tags for the BEGIN block.
	if isBegin {
		// The BEGIN block can set {topic} and user vars.

		// Topic setter
		match := reTopic.FindStringSubmatch(reply)
		var giveup uint
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

		// Set user vars
		match = reSet.FindStringSubmatch(reply)
		giveup = 0
		for len(match) > 0 {
			giveup++
			if giveup > rs.Depth {
				rs.warn("Infinite loop looking for set tag!")
				break
			}
			name := match[1]
			value := match[2]
			rs.sessions.Set(username, map[string]string{name: value})
			reply = strings.Replace(reply, fmt.Sprintf("<set %s=%s>", name, value), "", -1)
			match = reSet.FindStringSubmatch(reply)
		}
	} else {
		reply = rs.processTags(username, message, reply, stars, thatStars, 0)
	}

	return reply, nil
}
