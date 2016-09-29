// Package macros exports types relevant to object macros.
package macro

// MacroInterface is the interface for a Go object macro handler.
//
// Here, "object macro handler" means Go code is handling object macros for a
// foreign programming language, for example JavaScript.
type MacroInterface interface {
	Load(name string, code []string)
	Call(name string, fields []string) string
}
