// Package memory provides the default in-memory session store with skv backing the memory stuff up
package skv

import (
	"fmt"
	"strings"
	"sync"

	"github.com/aichaos/rivescript-go/sessions"
	"github.com/rapidloop/skv"
)

// Type MemoryStore implements the default in-memory session store for
// RiveScript.
type MemoryStore struct {
	lock   sync.Mutex
	users  map[string]*sessions.UserData
	frozen map[string]*sessions.UserData
	store  *skv.KVStore
	dbfile string
}

// New creates a new MemoryStore.
func New(dbfile string) (ms *MemoryStore, err error) {
	db, err := skv.Open(dbfile)
	if err != nil {
		return nil, err
	}
	ms = &MemoryStore{
		store:  db,
		dbfile: dbfile,
		users:  map[string]*sessions.UserData{},
		frozen: map[string]*sessions.UserData{},
	}
	return ms, nil
}

// init makes sure a username exists in the memory store.
func (s *MemoryStore) Init(username string) *sessions.UserData {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.users[username]; !ok {
		// check if it's in the db!
		var val *sessions.UserData
		err := s.store.Get(username, &val)
		if err != nil {
			s.users[username] = defaultSession()
		} else {
			s.users[username] = val
		}
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
	err := s.store.Put(username, s.users[username])
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	return
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
	s.store.Put(username, data)
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
		// check if it's in the db!
		var val *sessions.UserData
		err := s.store.Get(username, &val)
		if err != nil {
			return "", fmt.Errorf(`no data for username "%s"`, username)
		} else {
			s.users[username] = val
		}
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
		// check if it's in the db!
		var val *sessions.UserData
		err := s.store.Get(username, &val)
		if err != nil {
			return val, err
		} else {
			s.users[username] = val
		}
	}
	return cloneUser(s.users[username]), nil
}

// GetAll gets all data for all users. in memory
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

// GetLastMatch returns the last matched trigger for the user. In memory
func (s *MemoryStore) GetLastMatch(username string) (string, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	data, ok := s.users[username]
	if !ok {
		return "", fmt.Errorf(`no data for username "%s"`, username)
	}
	return data.LastMatch, nil
}

// GetHistory gets the user's history. In memory not from db.
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
	// delete from db
	err := s.store.Delete(username)
	if err != nil {
		return
	}
	return
}

// ClearAll resets all user data for all users.
func (s *MemoryStore) ClearAll() {
	s.lock.Lock()
	defer s.lock.Unlock()
	for key, user := range s.users {
		err := s.store.Delete(key)
		if err != nil {
			// Do this for all users. User gets one error back and a string with all the errors.
			fmt.Printf("%v %v %v\n", user, err)
		}
	}
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
