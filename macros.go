package rivescript

// Interface for object macro handlers.

type MacroHandler interface {
	Load(name string, code string)
	Call(name string, fields []string) string
}

// A built-in object macro for Golang functions. Note that Go isn't a dynamic
// language and it can't compile code dynamically, so this interface is purely
// for the back-end SetSubroutine() method.
type GolangHandler struct {
	functions map[string]*Subroutine
}

func (self GolangHandler) LoadFunction(name string, fn Subroutine) {
	// TODO: implementation
}
