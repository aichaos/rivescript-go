package rivescript

/* "Abstract Syntax Tree" of parsed RiveScript code.

The tree looks like this (in JSON-style syntax):

{
	"begin": {
		"global": {} // Global vars
		"var_": {}   // Bot vars
		"sub": {}    // Substitutions
		"person": {} // Person substitutions
		"array": {}  // Arrays
	}
	"topics": {}
	"objects" []
}
*/

type astRoot struct {
	begin   astBegin
	topics  map[string]*astTopic
	objects []*astObject
}

type astBegin struct {
	global map[string]string
	var_   map[string]string
	sub    map[string]string
	person map[string]string
	array  map[string][]string // Map of string (names) to arrays-of-strings
}

type astTopic struct {
	triggers []*astTrigger
	includes map[string]bool
	inherits map[string]bool
}

type astTrigger struct {
	trigger   string
	reply     []string
	condition []string
	redirect  string
	previous  string
}

type astObject struct {
	name     string
	language string
	code     []string
}

func newAST() *astRoot {
	ast := new(astRoot)

	// Initialize all the structures.
	ast.begin.global = map[string]string{}
	ast.begin.var_ = map[string]string{}
	ast.begin.sub = map[string]string{}
	ast.begin.person = map[string]string{}
	ast.begin.array = map[string][]string{}

	// Initialize the 'random' topic.
	ast.topics = map[string]*astTopic{}
	ast = initTopic(ast, "random")

	// Objects
	ast.objects = []*astObject{}

	return ast
}

// initTopic sets up the AST tree for a new topic and gets it ready for
// triggers to be added.
func initTopic(ast *astRoot, name string) *astRoot {
	ast.topics[name] = new(astTopic)
	ast.topics[name].triggers = []*astTrigger{}
	ast.topics[name].includes = map[string]bool{}
	ast.topics[name].inherits = map[string]bool{}
	return ast
}
