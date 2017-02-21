// Package sessions provides the interface and default session store for
// RiveScript.
package sessions

/*
Interface SessionManager describes a session manager for user variables
in RiveScript.

The session manager keeps track of getting and setting user variables,
for example when the `<set>` or `<get>` tags are used in RiveScript
or when API functions like `SetUservar()` are called.

By default RiveScript stores user sessions in memory and provides methods
to export and import them (e.g. to persist them when the bot shuts down
so they can be reloaded). If you'd prefer a more 'active' session storage,
for example one that puts user variables into a database or cache, you can
create your own session manager that implements this interface.
*/
type SessionManager interface {
	// Init makes sure a username has a session (creates one if not). It returns
	// the pointer to the user data in either case.
	Init(username string) *UserData

	// Set user variables from a map.
	Set(username string, vars map[string]string)

	// AddHistory adds input and reply to the user's history.
	AddHistory(username, input, reply string)

	// SetLastMatch sets the last matched trigger.
	SetLastMatch(username, trigger string)

	// Get a user variable.
	Get(username string, key string) (string, error)

	// Get all variables for a user.
	GetAny(username string) (*UserData, error)

	// Get all variables about all users.
	GetAll() map[string]*UserData

	// GetLastMatch returns the last trigger the user matched.
	GetLastMatch(username string) (string, error)

	// GetHistory returns the user's history.
	GetHistory(username string) (*History, error)

	// Clear all variables for a given user.
	Clear(username string)

	// Clear all variables for all users.
	ClearAll()

	// Freeze makes a snapshot of a user's variables.
	Freeze(string) error

	// Thaw unfreezes a snapshot of a user's variables and returns an error
	// if the user had no frozen variables.
	Thaw(username string, ThawAction ThawAction) error
}

// HistorySize is the number of entries stored in the history.
const HistorySize int = 9

// UserData is a container for user variables.
type UserData struct {
	Variables map[string]string `json:"vars"`
	LastMatch string            `json:"lastMatch"`
	*History  `json:"history"`
}

// History keeps track of recent input and reply history.
type History struct {
	Input []string `json:"input"`
	Reply []string `json:"reply"`
}

// NewHistory creates a new History object with the history arrays filled out.
func NewHistory() *History {
	h := &History{
		Input: []string{},
		Reply: []string{},
	}

	for i := 0; i < HistorySize; i++ {
		h.Input = append(h.Input, "undefined")
		h.Reply = append(h.Reply, "undefined")
	}

	return h
}

// Type ThawAction describes the action for the `Thaw()` method.
type ThawAction int

// Valid options for ThawAction.
const (
	// Thaw means to restore the user variables and erase the frozen copy.
	Thaw = iota

	// Discard means to cancel the frozen copy and not restore them.
	Discard

	// Keep means to restore the user variables and still keep the frozen copy.
	Keep
)
