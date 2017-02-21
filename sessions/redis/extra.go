package redis

// NOTE: This file contains added functions above and beyond the SessionManager
// implementation.

import (
	"encoding/json"
	"fmt"

	"github.com/aichaos/rivescript-go/sessions"
)

// key generates a key name to use in Redis.
func (s *Session) key(username string) string {
	return s.prefix + username
}

// frozenKey generates the 'frozen' key name to use in Redis.
func (s *Session) frozenKey(username string) string {
	return s.frozenPrefix + username
}

// getRedis gets a UserData out of the Redis cache.
func (s *Session) getRedis(username string) (*sessions.UserData, error) {
	data, err := s.getRedisFrozen(username, false)
	return data, err
}

// getRedisFrozen is the implementation behind getRedis and allows for the
// key to be overridden with the 'frozen' version.
func (s *Session) getRedisFrozen(username string, frozen bool) (*sessions.UserData, error) {
	var key string
	if frozen {
		key = s.frozenKey(username)
	} else {
		key = s.key(username)
	}

	// Check Redis for the key.
	value, err := s.client.Get(key).Result()
	if err != nil {
		return nil, fmt.Errorf(
			`no data for username "%s": %s`,
			username, err,
		)
	}

	// Decode the JSON.
	var user *sessions.UserData
	err = json.Unmarshal([]byte(value), &user)
	if err != nil {
		return nil, fmt.Errorf(
			`JSON unmarshal error for username "%s": %s`,
			username, err,
		)
	}

	return user, nil
}

// putRedis puts a UserData into the Redis cache.
func (s *Session) putRedis(username string, data *sessions.UserData) {
	s.putRedisFrozen(username, data, false)
}

// putRedisFrozen is the implementation behind putRedis and allows for the
// key to be overridden with the 'frozen' version.
func (s *Session) putRedisFrozen(username string, data *sessions.UserData, frozen bool) error {
	// Which key to use?
	var key string
	if frozen {
		key = s.frozenKey(username)
	} else {
		key = s.key(username)
	}

	encoded, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	err = s.client.Set(key, string(encoded), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

// defaultSession initializes the default session variables for a user.
func defaultSession() *sessions.UserData {
	return &sessions.UserData{
		Variables: map[string]string{
			"topic": "random",
		},
		LastMatch: "",
		History:   sessions.NewHistory(),
	}
}
