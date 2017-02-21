// Package redis implements a Redis backed session manager for RiveScript.
package redis

// NOTE: This source file contains the implementation of a SessionManager.

import (
	"fmt"
	"strings"

	"github.com/aichaos/rivescript-go/sessions"
	redis "gopkg.in/redis.v5"
)

// Config allows for configuring the Redis instance and key prefix.
type Config struct {
	// The key prefix to use in Redis. For example, with a username of 'alice',
	// the Redis key might be 'rivescript/alice'.
	//
	// The default prefix is 'rivescript/'
	Prefix string

	// The key used to prefix frozen user variables (those created by
	// `Freeze()`). The default is `frozen:<prefix>`
	FrozenPrefix string

	// Settings for the Redis client.
	Redis *redis.Options
}

// Session wraps a Redis client connection.
type Session struct {
	prefix       string
	frozenPrefix string
	client       *redis.Client
}

// New creates a new Redis session instance.
func New(options *Config) *Session {
	// No options given?
	if options == nil {
		options = &Config{}
	}

	// Default prefix is 'rivescript/'
	if options.Prefix == "" {
		options.Prefix = "rivescript/"
	}
	if options.FrozenPrefix == "" {
		options.FrozenPrefix = "frozen:" + options.Prefix
	}

	// Default options for Redis if none provided.
	if options.Redis == nil {
		options.Redis = &redis.Options{
			Addr: "localhost:6379",
			DB:   0,
		}
	}

	return &Session{
		prefix:       options.Prefix,
		frozenPrefix: options.FrozenPrefix,
		client:       redis.NewClient(options.Redis),
	}
}

// Init makes sure that a username has a session (creates one if not), and
// returns the pointer to it in any event.
func (s *Session) Init(username string) *sessions.UserData {
	// See if they have any data in Redis, and return it if so.
	user, err := s.getRedis(username)
	if err == nil {
		return user
	}

	// Create the default session.
	user = defaultSession()

	// Put them in Redis.
	s.putRedis(username, user)
	return user
}

// Set puts a user variable into Redis.
func (s *Session) Set(username string, vars map[string]string) {
	data := s.Init(username)

	for key, value := range vars {
		data.Variables[key] = value
	}

	s.putRedis(username, data)
}

// AddHistory adds to a user's history data.
func (s *Session) AddHistory(username, input, reply string) {
	data := s.Init(username)

	// Pop, unshift, pop, unshift.
	data.History.Input = data.History.Input[:len(data.History.Input)-1]
	data.History.Input = append([]string{strings.TrimSpace(input)}, data.History.Input...)
	data.History.Reply = data.History.Reply[:len(data.History.Reply)-1]
	data.History.Reply = append([]string{strings.TrimSpace(reply)}, data.History.Reply...)

	s.putRedis(username, data)
}

// SetLastMatch sets the user's last matched trigger.
func (s *Session) SetLastMatch(username, trigger string) {
	data := s.Init(username)
	data.LastMatch = trigger
	s.putRedis(username, data)
}

// Get a user variable out of Redis.
func (s *Session) Get(username, name string) (string, error) {
	data, err := s.getRedis(username)
	if err != nil {
		return "", err
	}

	value, ok := data.Variables[name]
	if !ok {
		return "", fmt.Errorf(`variable "%s" for user "%s" not set`, name, username)
	}
	return value, nil
}

// GetAny returns all variables about a user.
func (s *Session) GetAny(username string) (*sessions.UserData, error) {
	// Check redis.
	data, err := s.getRedis(username)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetAll gets all data for all users.
func (s *Session) GetAll() map[string]*sessions.UserData {
	result := map[string]*sessions.UserData{}

	keys, err := s.client.Keys(s.prefix + "*").Result()
	if err != nil {
		return result
	}

	for _, key := range keys {
		username := strings.Replace(key, s.prefix, "", 1)
		data, _ := s.GetAny(username)
		result[username] = data
	}

	return result
}

// GetLastMatch retrieves the user's last matched trigger.
func (s *Session) GetLastMatch(username string) (string, error) {
	data, err := s.getRedis(username)
	if err != nil {
		return "", err
	}

	return data.LastMatch, nil
}

// GetHistory gets the user's history.
func (s *Session) GetHistory(username string) (*sessions.History, error) {
	data, err := s.getRedis(username)
	if err != nil {
		return nil, err
	}

	return data.History, nil
}

// Clear deletes all variables about a user.
func (s *Session) Clear(username string) {
	s.client.Del(s.prefix + username)
}

// ClearAll resets all user data for all users.
func (s *Session) ClearAll() {
	// List all the users.
	keys, err := s.client.Keys(s.prefix + "*").Result()
	if err != nil {
		return
	}

	// Delete them all.
	s.client.Del(keys...)
}

// Freeze makes a snapshot of user variables.
func (s *Session) Freeze(username string) error {
	data, err := s.getRedis(username)
	if err != nil {
		return err
	}

	// Duplicate it into the frozen Redis key.
	return s.putRedisFrozen(username, data, true)
}

// Thaw restores user variables from a snapshot.
func (s *Session) Thaw(username string, action sessions.ThawAction) error {
	frozen, err := s.getRedisFrozen(username, true)
	if err != nil {
		return fmt.Errorf(`no frozen data for username "%s": %s`, username, err)
	}

	// Which type of thaw action are they using?
	switch action {
	case sessions.Thaw:
		// Thaw means to restore the frozen copy and then delete the copy.
		s.Clear(username)
		s.putRedis(username, frozen)
		s.client.Del(s.frozenKey(username))
	case sessions.Discard:
		// Discard means to just delete the frozen copy, do not restore it.
		s.client.Del(s.frozenKey(username))
	case sessions.Keep:
		// Keep restores from the frozen copy, but keeps the frozen copy.
		s.Clear(username)
		s.putRedis(username, frozen)
	default:
		return fmt.Errorf(`can't thaw data for username "%s": invalid thaw action`, username)
	}

	return nil
}
