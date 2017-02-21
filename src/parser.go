package rivescript

// parse loads the RiveScript code into the bot's memory.
func (rs *RiveScript) parse(path string, lines []string) error {
	rs.say("Parsing code...")

	// Get the abstract syntax tree of this file.
	AST, err := rs.parser.Parse(path, lines)
	if err != nil {
		return err
	}

	// Get all of the "begin" type variables
	for k, v := range AST.Begin.Global {
		if v == UNDEFTAG {
			delete(rs.global, k)
		} else {
			rs.global[k] = v
		}
	}
	for k, v := range AST.Begin.Var {
		if v == UNDEFTAG {
			delete(rs.vars, k)
		} else {
			rs.vars[k] = v
		}
	}
	for k, v := range AST.Begin.Sub {
		if v == UNDEFTAG {
			delete(rs.sub, k)
		} else {
			rs.sub[k] = v
		}
	}
	for k, v := range AST.Begin.Person {
		if v == UNDEFTAG {
			delete(rs.person, k)
		} else {
			rs.person[k] = v
		}
	}
	for k, v := range AST.Begin.Array {
		rs.array[k] = v
	}

	// Consume all the parsed triggers.
	for topic, data := range AST.Topics {
		// Keep a map of the topics that are included/inherited under this topic.
		if _, ok := rs.includes[topic]; !ok {
			rs.includes[topic] = map[string]bool{}
		}
		if _, ok := rs.inherits[topic]; !ok {
			rs.inherits[topic] = map[string]bool{}
		}

		// Merge in the topic inclusions/inherits.
		for included := range data.Includes {
			rs.includes[topic][included] = true
		}
		for inherited := range data.Inherits {
			rs.inherits[topic][inherited] = true
		}

		// Initialize the topic structure.
		if _, ok := rs.topics[topic]; !ok {
			rs.topics[topic] = new(astTopic)
			rs.topics[topic].triggers = []*astTrigger{}
		}

		// Consume the AST triggers into the brain.
		for _, trig := range data.Triggers {
			// Convert this AST trigger into an internal astmap trigger.
			trigger := new(astTrigger)
			trigger.trigger = trig.Trigger
			trigger.reply = trig.Reply
			trigger.condition = trig.Condition
			trigger.redirect = trig.Redirect
			trigger.previous = trig.Previous

			rs.topics[topic].triggers = append(rs.topics[topic].triggers, trigger)
		}
	}

	// Load all the parsed objects.
	for _, object := range AST.Objects {
		// Have a language handler for this?
		if _, ok := rs.handlers[object.Language]; ok {
			rs.say("Loading object macro %s (%s)", object.Name, object.Language)
			rs.handlers[object.Language].Load(object.Name, object.Code)
			rs.objlangs[object.Name] = object.Language
		}
	}

	return nil
}
