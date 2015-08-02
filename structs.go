package rivescript

// Miscellaneous structures
type UserData struct {
	data map[string]string
}

// This is like astTopic but is just for %Previous mapping
type thatTopic struct {
	triggers map[string]*thatTrigger
}

type thatTrigger struct {
	previous map[string]*astTrigger
}
