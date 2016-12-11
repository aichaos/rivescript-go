// Package null provides a session manager that has no memory.
package null

import "github.com/aichaos/rivescript-go/sessions"

// Type NullStore implements a memory store that has no memory.
//
// It's mostly useful for the unit tests. With this memory store in place,
// RiveScript is unable to maintain any user variables at all.
type NullStore struct{}

// New creates a new NullStore.
func New() *NullStore {
	return new(NullStore)
}

func (s *NullStore) Init(username string) *sessions.UserData {
	return nullSession()
}

func (s *NullStore) Set(username string, vars map[string]string) {}

func (s *NullStore) AddHistory(username, input, reply string) {}

func (s *NullStore) SetLastMatch(username, trigger string) {}

func (s *NullStore) Get(username string, name string) (string, error) {
	return "undefined", nil
}

func (s *NullStore) GetAny(username string) (*sessions.UserData, error) {
	return nullSession(), nil
}

func (s *NullStore) GetAll() map[string]*sessions.UserData {
	return map[string]*sessions.UserData{}
}

func (s *NullStore) GetLastMatch(username string) (string, error) {
	return "", nil
}

func (s *NullStore) GetHistory(username string) (*sessions.History, error) {
	return sessions.NewHistory(), nil
}

func (s *NullStore) Clear(username string) {}

func (s *NullStore) ClearAll() {}

func (s *NullStore) Freeze(username string) error {
	return nil
}

func (s *NullStore) Thaw(username string, action sessions.ThawAction) error {
	return nil
}

func nullSession() *sessions.UserData {
	return &sessions.UserData{
		Variables: map[string]string{},
		History:   sessions.NewHistory(),
		LastMatch: "",
	}
}
