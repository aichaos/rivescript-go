package rivescript

// Interface for object macro handlers.

type MacroInterface interface {
	Load(name string, code []string)
	Call(name string, fields []string) string
}

// A built-in object macro for Golang functions. Note that Go isn't a dynamic
// language and it can't compile code dynamically, so this interface is purely
// for the back-end SetSubroutine() method.
type golangHandler struct {
	functions map[string]*Subroutine
}

func (self golangHandler) loadFunction(name string, fn Subroutine) {
	// TODO: implementation
}
