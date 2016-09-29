package src

// Public API Configuration Methods

import (
	"errors"
	"github.com/aichaos/rivescript-go/macro"
)

func (rs *RiveScript) SetHandler(lang string, handler macro.MacroInterface) {
	rs.handlers[lang] = handler
}

func (rs *RiveScript) RemoveHandler(lang string) {
	delete(rs.handlers, lang)
}

func (rs *RiveScript) SetSubroutine(name string, fn Subroutine) {
	rs.subroutines[name] = fn
}

func (rs *RiveScript) DeleteSubroutine(name string) {
	delete(rs.subroutines, name)
}

func (rs *RiveScript) SetGlobal(name string, value string) {
	if value == "undefined" {
		delete(rs.global, name)
	} else {
		rs.global[name] = value
	}
}

func (rs *RiveScript) SetVariable(name string, value string) {
	if value == "undefined" {
		delete(rs.var_, name)
	} else {
		rs.var_[name] = value
	}
}

func (rs *RiveScript) SetSubstitution(name string, value string) {
	if value == "undefined" {
		delete(rs.sub, name)
	} else {
		rs.sub[name] = value
	}
}

func (rs *RiveScript) SetPerson(name string, value string) {
	if value == "undefined" {
		delete(rs.person, name)
	} else {
		rs.person[name] = value
	}
}

func (rs *RiveScript) SetUservar(username string, name string, value string) {
	// Initialize the user?
	if _, ok := rs.users[username]; !ok {
		rs.users[username] = newUser()
	}

	if value == "undefined" {
		delete(rs.users[username].data, name)
	} else {
		rs.users[username].data[name] = value
	}
}

func (rs *RiveScript) SetUservars(username string, data map[string]string) {
	// Initialize the user?
	if _, ok := rs.users[username]; !ok {
		rs.users[username] = newUser()
	}

	for key, value := range data {
		if value == "undefined" {
			delete(rs.users[username].data, key)
		} else {
			rs.users[username].data[key] = value
		}
	}
}

func (rs *RiveScript) GetGlobal(name string) (string, error) {
	if _, ok := rs.global[name]; ok {
		return rs.global[name], nil
	}
	return "undefined", errors.New("Global variable not found.")
}

func (rs *RiveScript) GetVariable(name string) (string, error) {
	if _, ok := rs.var_[name]; ok {
		return rs.var_[name], nil
	}
	return "undefined", errors.New("Variable not found.")
}

func (rs *RiveScript) GetUservar(username string, name string) (string, error) {
	if _, ok := rs.users[username]; ok {
		if _, ok := rs.users[username].data[name]; ok {
			return rs.users[username].data[name], nil
		}
	}
	return "undefined", errors.New("User variable not found.")
}

func (rs *RiveScript) GetUservars(username string) (map[string]string, error) {
	if _, ok := rs.users[username]; ok {
		return rs.users[username].data, nil
	}
	return map[string]string{}, errors.New("Username not found.")
}

func (rs *RiveScript) GetAllUservars() map[string]map[string]string {
	result := map[string]map[string]string{}
	for username, data := range rs.users {
		result[username] = data.data
	}
	return result
}

func (rs *RiveScript) ClearUservars(username string) {
	delete(rs.users, username)
}

func (rs *RiveScript) ClearAllUservars() {
	for username, _ := range rs.users {
		delete(rs.users, username)
	}
}

func (rs *RiveScript) FreezeUservars(username string) error {
	if _, ok := rs.users[username]; ok {
		delete(rs.freeze, username) // Always start fresh
		rs.freeze[username] = newUser()

		for key, value := range rs.users[username].data {
			rs.freeze[username].data[key] = value
		}

		for i, entry := range rs.users[username].inputHistory {
			rs.freeze[username].inputHistory[i] = entry
		}
		for i, entry := range rs.users[username].replyHistory {
			rs.freeze[username].replyHistory[i] = entry
		}
		return nil
	}
	return errors.New("Username not found.")
}

func (rs *RiveScript) ThawUservars(username string, action string) error {
	if _, ok := rs.freeze[username]; ok {
		// What are we doing?
		if action == "thaw" {
			rs.ClearUservars(username)
			rs.users[username] = rs.freeze[username]
			delete(rs.freeze, username)
		} else if action == "discard" {
			delete(rs.freeze, username)
		} else if action == "keep" {
			delete(rs.users, username) // Always start fresh
			rs.users[username] = newUser()

			for key, value := range rs.freeze[username].data {
				rs.users[username].data[key] = value
			}

			for i, entry := range rs.freeze[username].inputHistory {
				rs.users[username].inputHistory[i] = entry
			}
			for i, entry := range rs.freeze[username].replyHistory {
				rs.users[username].replyHistory[i] = entry
			}
		} else {
			return errors.New("Unsupported thaw action. Valid options are: thaw, discard, keep.")
		}
		return nil
	}
	return errors.New("Username not found.")
}

func (rs *RiveScript) LastMatch(username string) (string, error) {
	if _, ok := rs.users[username]; ok {
		return rs.users[username].lastMatch, nil
	}
	return "", errors.New("Username not found.")
}

func (rs *RiveScript) CurrentUser() string {
	return rs.currentUser
}
