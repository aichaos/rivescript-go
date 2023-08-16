package rivescript

/*
For my own sanity while programming the code, these structs mirror the data
in the 'ast' subpackage but uses non-exported fields for the bot's own use.

The logic is as follows:

- The parser subpackage becomes a stand-alone Go module that third party
  developers can use to make their own applications around the RiveScript
  scripting language itself. To that end, it exports a public AST tree.
- In RiveScript's parse() function, it uses the public parser package and
  gets back an AST tree full of exported fields. It doesn't need these fields
  to be exported, and it copies them into internal fields of similar names.
- I don't want to use the exported AST names directly because it makes the
  code become a Russian Roulette of capital or non-capital names.

An example of how unwieldy the code would be if I use the direct AST types:

	rs.thats[topic].Triggers[trigger.Trigger].Previous[trigger.Previous]
	                ^                ^        ^                ^

If the ast package structs are updated, update the mappings in this package too.
*/

type astTopic struct {
	triggers []*astTrigger
}

type astTrigger struct {
	trigger   string
	reply     []string
	condition []string
	redirect  string
	previous  string
}
