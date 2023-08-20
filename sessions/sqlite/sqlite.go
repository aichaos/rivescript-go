package sqlite

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/aichaos/rivescript-go/sessions"
	_ "modernc.org/sqlite"
)

var schema string = `PRAGMA journal_mode = WAL;
PRAGMA synchronous = normal;
PRAGMA foreign_keys = on;
PRAGMA encoding = "UTF-8";
BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS "users" (
	"id" INTEGER,
	"username" TEXT UNIQUE,
	"last_match" TEXT,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "user_variables" (
	"id" INTEGER,
	"user_id" INTEGER NOT NULL,
	"key" TEXT NOT NULL,
	"value" TEXT,
	PRIMARY KEY("id" AUTOINCREMENT),
	UNIQUE("user_id", "key")
);
CREATE TABLE IF NOT EXISTS "history" (
	"id" INTEGER,
	"user_id" INTEGER NOT NULL,
	"input" TEXT NOT NULL,
	"reply" TEXT NOT NULL,
	"timestamp" INTEGER NOT NULL DEFAULT (CAST(strftime('%s', 'now') AS INTEGER)),
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "frozen_user" (
	"id" INTEGER,
	"user_id" INTEGER NOT NULL,
	"data" TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE IF NOT EXISTS "local_storage" (
	"id" INTEGER,
	"key" TEXT NOT NULL,
	"value" TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE VIEW IF NOT EXISTS v_user_variables AS
SELECT users.username AS username,
	user_variables.key,
	user_variables.value
FROM users,
	user_variables
WHERE users.id = user_variables.user_id;
COMMIT;`

type Client struct {
	lock sync.Mutex
	db   *sql.DB
}

// New creates a new Client.
func New(filename string) (*Client, error) {
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(0)

	_, err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	return &Client{
		db: db,
	}, nil
}

func (s *Client) Close() error {
	return s.db.Close()
}

// init makes sure a username exists in the memory store.
func (s *Client) Init(username string) *sessions.UserData {
	user, err := s.GetAny(username)
	if err != nil {
		func() {

			s.lock.Lock()
			defer s.lock.Unlock()

			tx, _ := s.db.Begin()
			stmt, _ := tx.Prepare(`INSERT OR IGNORE INTO users (username, last_match) VALUES (?,"");`)
			defer stmt.Close()
			stmt.Exec(username)
			tx.Commit()
		}()

		s.Set(username, map[string]string{
			"topic": "random",
		})

		return &sessions.UserData{
			Variables: map[string]string{
				"topic": "random",
			},
			LastMatch: "",
			History:   sessions.NewHistory(),
		}
	}
	return user
} // Init()

// Set a user variable.
func (s *Client) Set(username string, vars map[string]string) {
	s.Init(username)

	s.lock.Lock()
	defer s.lock.Unlock()

	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare(`INSERT OR REPLACE INTO user_variables (user_id, key, value) VALUES ((SELECT id FROM users WHERE username = ?), ?, ?);`)
	defer stmt.Close()
	for k, v := range vars {
		stmt.Exec(username, k, v)
	}
	tx.Commit()
}

// AddHistory adds history items.
func (s *Client) AddHistory(username, input, reply string) {
	s.Init(username)

	s.lock.Lock()
	defer s.lock.Unlock()

	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare(`INSERT INTO history (user_id, input,reply)VALUES((SELECT id FROM users WHERE username = ?),?,?);`)
	defer stmt.Close()
	stmt.Exec(username, input, reply)
	tx.Commit()
}

// SetLastMatch sets the user's last matched trigger.
func (s *Client) SetLastMatch(username, trigger string) {
	s.Init(username)

	s.lock.Lock()
	defer s.lock.Unlock()

	tx, _ := s.db.Begin()
	stmt, _ := tx.Prepare(`UPDATE users SET last_match = ? WHERE username = ?;`)
	defer stmt.Close()
	stmt.Exec(trigger, username)
	tx.Commit()
}

// Get a user variable.
func (s *Client) Get(username, name string) (string, error) {
	var value string
	row := s.db.QueryRow(`SELECT value FROM user_variables WHERE user_id = (SELECT id FROM users WHERE username = ?) AND key = ?;`, username, name)
	switch err := row.Scan(&value); err {
	case sql.ErrNoRows:
		return "", fmt.Errorf("no %s variable found for user %s", name, username)
	case nil:
		return value, nil
	default:
		return "", fmt.Errorf("unknown sql error")
	}
}

// GetAny gets all variables for a user.
func (s *Client) GetAny(username string) (*sessions.UserData, error) {
	history, err := s.GetHistory(username)
	if err != nil {
		return nil, err
	}
	last_match, err := s.GetLastMatch(username)
	if err != nil {
		return nil, err
	}

	var variables map[string]string = make(map[string]string)
	rows, err := s.db.Query(`SELECT key,value FROM user_variables WHERE user_id = (SELECT id FROM users WHERE username = ?);`, username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var key, value string
	for rows.Next() {
		err = rows.Scan(&key, &value)
		if err != nil {
			continue
		}
		variables[key] = value
	}

	return &sessions.UserData{
		History:   history,
		LastMatch: last_match,
		Variables: variables,
	}, nil
}

// GetAll gets all data for all users.
func (s *Client) GetAll() map[string]*sessions.UserData {
	var users []string = make([]string, 0)
	rows, _ := s.db.Query(`SELECT username FROM users;`)
	defer rows.Close()
	var user string
	for rows.Next() {
		rows.Scan(&user)
		users = append(users, user)
	}

	var usersmap map[string]*sessions.UserData = make(map[string]*sessions.UserData)
	for _, user := range users {
		u, _ := s.GetAny(user)
		usersmap[user] = u
	}

	return usersmap
}

// GetLastMatch returns the last matched trigger for the user,
func (s *Client) GetLastMatch(username string) (string, error) {
	var last_match string
	row := s.db.QueryRow(`SELECT last_match FROM users WHERE username = ?;`, username)
	switch err := row.Scan(&last_match); err {
	case sql.ErrNoRows:
		return "", fmt.Errorf("no last match found for user %s", username)
	case nil:
		return last_match, nil
	default:
		return "", fmt.Errorf("unknown sql error: %s", err)
	}
}

// GetHistory gets the user's history.
func (s *Client) GetHistory(username string) (*sessions.History, error) {
	data := &sessions.History{
		Input: []string{},
		Reply: []string{},
	}

	for i := 0; i < sessions.HistorySize; i++ {
		data.Input = append(data.Input, "undefined")
		data.Reply = append(data.Reply, "undefined")
	}

	rows, err := s.db.Query("SELECT input,reply FROM history WHERE user_id = (SELECT id FROM users WHERE username = ?) ORDER BY timestamp ASC LIMIT 10;", username)
	if err != nil {
		return data, err
	}
	defer rows.Close()
	for rows.Next() {
		var input, reply string
		err := rows.Scan(&input, &reply)
		if err != nil {
			log.Println("[ERROR]", err)
			continue
		}
		data.Input = data.Input[:len(data.Input)-1]                            // Pop
		data.Input = append([]string{strings.TrimSpace(input)}, data.Input...) // Unshift
		data.Reply = data.Reply[:len(data.Reply)-1]                            // Pop
		data.Reply = append([]string{strings.TrimSpace(reply)}, data.Reply...) // Unshift

	}

	return data, nil
}

// Clear data for a user.
func (s *Client) Clear(username string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx, _ := s.db.Begin()
	tx.Exec(`DELETE FROM user_variables WHERE user_id = (SELECT id FROM users WHERE username = ?);`, username)
	tx.Exec(`DELETE FROM history WHERE user_id = (SELECT id FROM users WHERE username = ?);`, username)

	tx.Exec(`DELETE FROM users WHERE username = ?;`, username)
	tx.Commit()
}

// ClearAll resets all user data for all users.
func (s *Client) ClearAll() {
	s.lock.Lock()
	defer s.lock.Unlock()

	tx, _ := s.db.Begin()
	tx.Exec(`DELETE FROM user_variables;`)
	tx.Exec(`DELETE FROM history;`)

	s.db.Exec(`DELETE FROM users;`)
	tx.Commit()
}

// Freeze makes a snapshot of user variables.
func (s *Client) Freeze(username string) error {
	user := s.Init(username)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}

	s.lock.Lock()
	defer s.lock.Unlock()

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare(`INSERT OR REPLACE INTO frozen_user (user_id, data)VALUES((SELECT id FROM users WHERE username = ?), ?);`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, string(data))
	if err != nil {
		return err
	}

	return tx.Commit()
}

// Thaw restores from a snapshot.
func (s *Client) Thaw(username string, action sessions.ThawAction) error {
	user, err := func(u string) (sessions.UserData, error) {
		var data string
		var reply sessions.UserData
		row := s.db.QueryRow(`SELECT data FROM frozen_user WHERE user_id = (SELECT id FROM users WHERE username = ?);`, username)
		switch err := row.Scan(&data); err {
		case sql.ErrNoRows:
			return sessions.UserData{}, fmt.Errorf("no rows found")
		case nil:
			err = json.Unmarshal([]byte(data), &reply)
			if err != nil {
				return sessions.UserData{}, err
			}
			return reply, nil
		default:
			return sessions.UserData{}, fmt.Errorf("unknown sql error")
		}
	}(username)
	if err != nil {
		return fmt.Errorf("no data for snapshot for user %s", username)
	}

	switch action {
	case sessions.Thaw:
		if err := func() error {
			s.lock.Lock()
			defer s.lock.Unlock()

			_, err = s.db.Exec(`DELETE FROM frozen_user WHERE user_id = (SELECT id FROM users WHERE username = ?);`, username)
			if err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}

		s.Clear(username)
		s.Set(username, user.Variables)
		s.SetLastMatch(username, user.LastMatch)
		for i := len(user.History.Input) - 1; i >= 0; i-- {
			s.AddHistory(username, user.History.Input[i], user.History.Reply[i])
		}

		return nil
	case sessions.Discard:
		s.lock.Lock()
		defer s.lock.Unlock()

		_, err = s.db.Exec(`DELETE FROM frozen_user WHERE user_id = (SELECT id FROM users WHERE username = ?);`, username)
		if err != nil {
			return err
		}
	case sessions.Keep:
		s.Clear(username)
		s.Set(username, user.Variables)
		s.SetLastMatch(username, user.LastMatch)
		for i := range user.History.Input {
			s.AddHistory(username, user.History.Input[i], user.History.Reply[i])
		}
		return nil
	default:
		return fmt.Errorf("something went wrong")
	}

	return nil
}
