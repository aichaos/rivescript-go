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

// Root represents the root of the AST tree.
type Root struct {
	Begin   Begin             `json:"begin"`
	Topics  map[string]*Topic `json:"topics"`
	Objects []*Object         `json:"objects"`
}

// Begin represents the "begin block" style data (configuration).
type Begin struct {
	Global map[string]string   `json:"global"`
	Var    map[string]string   `json:"var"`
	Sub    map[string]string   `json:"sub"`
	Person map[string]string   `json:"person"`
	Array  map[string][]string `json:"array"` // Map of string (names) to arrays-of-strings
}

// Topic represents a topic of conversation.
type Topic struct {
	Triggers []*Trigger      `json:"triggers"`
	Includes map[string]bool `json:"includes"`
	Inherits map[string]bool `json:"inherits"`
}

// Trigger has a trigger pattern and all the subsequent handlers for it.
type Trigger struct {
	Trigger   string   `json:"trigger"`
	Reply     []string `json:"reply"`
	Condition []string `json:"condition"`
	Redirect  string   `json:"redirect"`
	Previous  string   `json:"previous"`
}

// Object contains source code of dynamically parsed object macros.
type Object struct {
	Name     string   `json:"name"`
	Language string   `json:"language"`
	Code     []string `json:"code"`
}

// New creates a new, empty, abstract syntax tree.
func New() *Root {
	ast := &Root{
		// Initialize all the structures.
		Begin: Begin{
			Global: map[string]string{},
			Var:    map[string]string{},
			Sub:    map[string]string{},
			Person: map[string]string{},
			Array:  map[string][]string{},
		},
		Topics:  map[string]*Topic{},
		Objects: []*Object{},
	}

	// Initialize the 'random' topic.
	ast.AddTopic("random")

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
