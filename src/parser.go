package src

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
		if v == "<undef>" {
			delete(rs.global, k)
		} else {
			rs.global[k] = v
		}
	}
	for k, v := range AST.Begin.Var {
		if v == "<undef>" {
			delete(rs.var_, k)
		} else {
			rs.var_[k] = v
		}
	}
	for k, v := range AST.Begin.Sub {
		if v == "<undef>" {
			delete(rs.sub, k)
		} else {
			rs.sub[k] = v
		}
	}
	for k, v := range AST.Begin.Person {
		if v == "<undef>" {
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
		for included, _ := range data.Includes {
			rs.includes[topic][included] = true
		}
		for inherited, _ := range data.Inherits {
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
			trigger.reply   = trig.Reply
			trigger.condition = trig.Condition
			trigger.redirect  = trig.Redirect
			trigger.previous  = trig.Previous

			rs.topics[topic].triggers = append(rs.topics[topic].triggers, trigger)

			// Does this one have a %Previous? If so, make a pointer to this
			// exact trigger in rs.thats
			if trigger.previous != "" {
				// Initialize the structure first.
				if _, ok := rs.thats[topic]; !ok {
					rs.thats[topic] = new(thatTopic)
					rs.say("%q", rs.thats[topic])
					rs.thats[topic].triggers = map[string]*thatTrigger{}
				}
				if _, ok := rs.thats[topic].triggers[trigger.trigger]; !ok {
					rs.say("%q", rs.thats[topic].triggers[trigger.trigger])
					rs.thats[topic].triggers[trigger.trigger] = new(thatTrigger)
					rs.thats[topic].triggers[trigger.trigger].previous = map[string]*astTrigger{}
				}
				rs.thats[topic].triggers[trigger.trigger].previous[trigger.previous] = trigger
			}
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
