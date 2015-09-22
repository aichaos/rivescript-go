package rivescript

// Interface for object macro handlers.

type MacroInterface interface {
	Load(name string, code []string)
	Call(name string, fields []string) string
}
