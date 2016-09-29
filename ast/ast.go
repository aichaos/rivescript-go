/*
Package ast defines the Abstract Syntax Tree for RiveScript.

The tree looks like this (in JSON-style syntax):

	{
		"Begin": {
			"Global": {}, // Global vars
			"Var": {},    // Bot variables
			"Sub": {},    // Substitution map
			"Person": {}, // Person substitution map
			"Array": {},  // Arrays
		},
		"Topics": {},
		"Objects": [],
	}
*/
package ast

// Type Root represents the root of the AST tree.
type Root struct {
	Begin   Begin
	Topics  map[string]*Topic
	Objects []*Object
}

// Type Begin represents the "begin block" style data (configuration).
type Begin struct {
	Global map[string]string
	Var    map[string]string
	Sub    map[string]string
	Person map[string]string
	Array  map[string][]string // Map of string (names) to arrays-of-strings
}

// Type Topic represents a topic of conversation.
type Topic struct {
	Triggers []*Trigger
	Includes map[string]bool
	Inherits map[string]bool
}

// Type Trigger has a trigger pattern and all the subsequent handlers for it.
type Trigger struct {
	Trigger   string
	Reply     []string
	Condition []string
	Redirect  string
	Previous  string
}

// Type Object contains source code of dynamically parsed object macros.
type Object struct {
	Name     string
	Language string
	Code     []string
}

// New creates a new, empty, abstract syntax tree.
func New() *Root {
	ast := new(Root)

	// Initialize all the structures.
	ast.Begin.Global = map[string]string{}
	ast.Begin.Var = map[string]string{}
	ast.Begin.Sub = map[string]string{}
	ast.Begin.Person = map[string]string{}
	ast.Begin.Array = map[string][]string{}

	// Initialize the 'random' topic.
	ast.Topics = map[string]*Topic{}
	ast.AddTopic("random")

	// Objects
	ast.Objects = []*Object{}

	return ast
}

// AddTopic sets up the AST tree for a new topic and gets it ready for
// triggers to be added.
func (ast *Root) AddTopic(name string) {
	ast.Topics[name] = new(Topic)
	ast.Topics[name].Triggers = []*Trigger{}
	ast.Topics[name].Includes = map[string]bool{}
	ast.Topics[name].Inherits = map[string]bool{}
}
