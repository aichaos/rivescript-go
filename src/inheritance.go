package rivescript

import "fmt"

/* Topic inheritance functions.

These are helper functions to assist with topic inheritance and includes.
*/

/*
getTopicTriggers recursively scans topics and collects triggers therein.

This function scans through a topic and collects its triggers, along with the
triggers belonging to any topic that's inherited by or included by the parent
topic. Some triggers will come out with an {inherits} tag to signify
inheritance depth.

Params:

	topic: The name of the topic to scan through
	thats: Whether to get only triggers that have %Previous.
		`false` returns all triggers.

Each "trigger" returned from this function is actually an array, where index
0 is the trigger text and index 1 is the pointer to the trigger's data within
the original topic structure.
*/
func (rs *RiveScript) getTopicTriggers(topic string, thats bool) []sortedTriggerEntry {
	return rs._getTopicTriggers(topic, thats, 0, 0, false)
}

/*
_getTopicTriggers implements the bulk of the logic for getTopicTriggers.

Additional private parameters used:
- depth: Recursion depth counter.
- inheritance: Inheritance counter.
- inherited: Inherited status.

Important info about the depth vs. inheritance params to this function:
depth increments by 1 each time this function recursively calls itself.
inheritance only increments by 1 when this topic inherits another topic.

This way, `> topic alpha includes beta inherits gamma` will have this effect:
- alpha and beta's triggers are combined together into one matching pool,
- and then those triggers have higher priority than gamma's.

The inherited option is true if this is a recursive call, from a topic that
inherits other topics. This forces the {inherits} tag to be added to the
triggers. This only applies when the topic 'includes' another topic.
*/
func (rs *RiveScript) _getTopicTriggers(topic string, thats bool, depth uint, inheritance int, inherited bool) []sortedTriggerEntry {
	// Break if we're in too deep.
	if depth > rs.Depth {
		rs.warn("Deep recursion while scanning topic inheritance!")
		return []sortedTriggerEntry{}
	}

	/*
		Keep in mind here that there is a difference between 'includes' and
		'inherits' -- topics that inherit other topics are able to OVERRIDE
		triggers that appear in the inherited topic. This means that if the top
		topic has a trigger of simply '*', then NO triggers are capable of
		matching in ANY inherited topic, because even though * has the lowest
		priority, it has an automatic priority over all inherited topics.

		The getTopicTriggers method takes this into account. All topics that
		inherit other topics will have their triggers prefixed with a fictional
		{inherits} tag, which would start at {inherits=0} and increment of this
		topic has other inheriting topics. So we can use this tag to make sure
		topics that inherit things will have their triggers always be on top of
		the stack, from inherits=0 to inherits=n.
	*/
	rs.say("Collecting trigger list for topic %s (depth=%d; inheritance=%d; inherited=%v)",
		topic, depth, inheritance, inherited)

	// Collect an array of triggers to return.
	triggers := []sortedTriggerEntry{}

	// Get those that exist in this topic directly.
	inThisTopic := []sortedTriggerEntry{}

	if _, ok := rs.topics[topic]; ok {
		for _, trigger := range rs.topics[topic].triggers {
			if !thats {
				// All triggers.
				entry := sortedTriggerEntry{trigger.trigger, trigger}
				inThisTopic = append(inThisTopic, entry)
			} else {
				// Only triggers that have %Previous.
				if trigger.previous != "" {
					inThisTopic = append(inThisTopic, sortedTriggerEntry{trigger.previous, trigger})
				}
			}
		}
	}

	// Does this topic include others?
	if _, ok := rs.includes[topic]; ok {
		for includes := range rs.includes[topic] {
			rs.say("Topic %s includes %s", topic, includes)
			triggers = append(triggers, rs._getTopicTriggers(includes, thats, depth+1, inheritance+1, false)...)
		}
	}

	// Does this topic inherit others?
	if _, ok := rs.inherits[topic]; ok {
		for inherits := range rs.inherits[topic] {
			rs.say("Topic %s inherits %s", topic, inherits)
			triggers = append(triggers, rs._getTopicTriggers(inherits, thats, depth+1, inheritance+1, true)...)
		}
	}

	// Collect the triggers for *this* topic. If this topic inherits any other
	// topics, it means that this topic's triggers have higher priority than
	// those in any inherited topics. Enforce this with an {inherits} tag.
	if len(rs.inherits[topic]) > 0 || inherited {
		for _, trigger := range inThisTopic {
			rs.say("Prefixing trigger with {inherits=%d} %s", inheritance, trigger.trigger)
			label := fmt.Sprintf("{inherits=%d}%s", inheritance, trigger.trigger)
			triggers = append(triggers, sortedTriggerEntry{label, trigger.pointer})
		}
	} else {
		for _, trigger := range inThisTopic {
			triggers = append(triggers, sortedTriggerEntry{trigger.trigger, trigger.pointer})
		}
	}

	return triggers
}

/*
getTopicTree returns an array of every topic related to a topic (all the
topics it inherits or includes, plus all the topics included or inherited
by those topics, and so on). The array includes the original topic, too.
*/
func (rs *RiveScript) getTopicTree(topic string, depth uint) []string {
	// Break if we're in too deep.
	if depth > rs.Depth {
		rs.warn("Deep recursion while scanning topic tree!")
		return []string{}
	}

	// Collect an array of all topics.
	topics := []string{topic}

	for includes := range rs.includes[topic] {
		topics = append(topics, rs.getTopicTree(includes, depth+1)...)
	}
	for inherits := range rs.inherits[topic] {
		topics = append(topics, rs.getTopicTree(inherits, depth+1)...)
	}

	return topics
}
