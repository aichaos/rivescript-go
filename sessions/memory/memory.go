// Package memory provides the default in-memory session store.
package memory

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aichaos/rivescript-go/sessions"
)

// Type MemoryStore implements the default in-memory session store for
// RiveScript.
type MemoryStore struct {
	lock   sync.Mutex
	users  map[string]*sessions.UserData
	frozen map[string]*sessions.UserData
}

// New creates a new MemoryStore.
func New() *MemoryStore {
	return &MemoryStore{
		users:  map[string]*sessions.UserData{},
		frozen: map[string]*sessions.UserData{},
	}
}

// init makes sure a username exists in the memory store.
func (s *MemoryStore) Init(username string) *sessions.UserData {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.users[username]; !ok {
		s.users[username] = defaultSession()
	}
	return s.users[username]
}

// Set a user variable.
func (s *MemoryStore) Set(username string, vars map[string]string) {
	s.Init(username)
	s.lock.Lock()
	defer s.lock.Unlock()

	for k, v := range vars {
		s.users[username].Variables[k] = v
	}
}

// AddHistory adds history items.
func (s *MemoryStore) AddHistory(username, input, reply string) {
	data := s.Init(username)
	s.lock.Lock()
	defer s.lock.Unlock()

	data.History.Input = data.History.Input[:len(data.History.Input)-1]                    // Pop
	data.History.Input = append([]string{strings.TrimSpace(input)}, data.History.Input...) // Unshift
	data.History.Reply = data.History.Reply[:len(data.History.Reply)-1]                    // Pop
	data.History.Reply = append([]string{strings.TrimSpace(reply)}, data.History.Reply...) // Unshift
}

// SetLastMatch sets the user's last matched trigger.
func (s *MemoryStore) SetLastMatch(username, trigger string) {
	data := s.Init(username)
	s.lock.Lock()
	defer s.lock.Unlock()
	data.LastMatch = trigger
}

// Get a user variable.
func (s *MemoryStore) Get(username, name string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.users[username]; !ok {
		return "", fmt.Errorf(`no data for username "%s"`, username)
	}

	value, ok := s.users[username].Variables[name]
	if !ok {
		return "", fmt.Errorf(`variable "%s" for user "%s" not set`, name, username)
	}

	return value, nil
}

// GetAny gets all variables for a user.
func (s *MemoryStore) GetAny(username string) (*sessions.UserData, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.users[username]; !ok {
		return &sessions.UserData{}, fmt.Errorf(`no data for username "%s"`, username)
	}
	return cloneUser(s.users[username]), nil
}

// GetAll gets all data for all users.
func (s *MemoryStore) GetAll() map[string]*sessions.UserData {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Make safe copies of all our structures.
	var result map[string]*sessions.UserData
	for k, v := range s.users {
		result[k] = cloneUser(v)
	}
	return result
}

// GetLastMatch returns the last matched trigger for the user,
func (s *MemoryStore) GetLastMatch(username string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	data, ok := s.users[username]
	if !ok {
		return "", fmt.Errorf(`no data for username "%s"`, username)
	}
	return data.LastMatch, nil
}

// GetHistory gets the user's history.
func (s *MemoryStore) GetHistory(username string) (*sessions.History, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	data, ok := s.users[username]
	if !ok {
		return nil, fmt.Errorf(`no data for username "%s"`, username)
	}
	return data.History, nil
}

// Clear data for a user.
func (s *MemoryStore) Clear(username string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.users, username)
}

// ClearAll resets all user data for all users.
func (s *MemoryStore) ClearAll() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.users = make(map[string]*sessions.UserData)
	s.frozen = make(map[string]*sessions.UserData)
}

// Freeze makes a snapshot of user variables.
func (s *MemoryStore) Freeze(username string) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	data, ok := s.users[username]
	if !ok {
		return fmt.Errorf(`no data for username %s`, username)
	}

	s.frozen[username] = cloneUser(data)
	return nil
}

// Thaw restores from a snapshot.
func (s *MemoryStore) Thaw(username string, action sessions.ThawAction) error {
	s.lock.Lock()
	defer s.lock.Unlock()

	frozen, ok := s.frozen[username]
	if !ok {
		return fmt.Errorf(`no frozen data for username "%s"`, username)
	}

	if action == sessions.Thaw {
		s.users[username] = cloneUser(frozen)
		delete(s.frozen, username)
	} else if action == sessions.Discard {
		delete(s.frozen, username)
	} else if action == sessions.Keep {
		s.users[username] = cloneUser(frozen)
	}

	return nil
}

// cloneUser makes a safe clone of a UserData.
func cloneUser(data *sessions.UserData) *sessions.UserData {
	new := defaultSession()

	// Copy user variables.
	for k, v := range data.Variables {
		new.Variables[k] = v
	}

	// Copy history.
	for i := 0; i < sessions.HistorySize; i++ {
		new.History.Input[i] = data.History.Input[i]
		new.History.Reply[i] = data.History.Reply[i]
	}

	return new
}

// defaultSession initializes the default session variables for a user.
// This mostly just means the topic is set to "random" and structs
// are initialized.
func defaultSession() *sessions.UserData {
	return &sessions.UserData{
		Variables: map[string]string{
			"topic": "random",
		},
		LastMatch: "",
		History:   sessions.NewHistory(),
	}
}
