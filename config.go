package rivescript

// Public API Configuration Methods

import (
	"errors"
)

/*
SetHandler sets a custom language handler for RiveScript object macros.
*/
func (rs *RiveScript) SetHandler(lang string, handler MacroInterface) {
	rs.handlers[lang] = handler
}

/*
DeleteHandler removes an object macro language handler.
*/
func (rs *RiveScript) RemoveHandler(lang string) {
	delete(rs.handlers, lang)
}

// Subroutine is a Golang function type for defining an object macro in Go.
type Subroutine func(*RiveScript, []string) string

/*
SetSubroutine defines a Go object macro from your program.

Params:

	name: The name of your subroutine for the `<call>` tag in RiveScript.
	fn: A function with a prototype `func(*RiveScript, []string) string`
*/
func (rs *RiveScript) SetSubroutine(name string, fn Subroutine) {
	rs.subroutines[name] = fn
}

/*
DeleteSubroutine removes a Go object macro.
*/
func (rs *RiveScript) DeleteSubroutine(name string) {
	delete(rs.subroutines, name)
}

/*
SetGlobal sets a global variable.

This is equivalent to `! global` in RiveScript. Set the value to `undefined`
to delete a global.
*/
func (rs *RiveScript) SetGlobal(name string, value string) {
	if value == "undefined" {
		delete(rs.global, name)
	} else {
		rs.global[name] = value
	}
}

/*
SetVariable sets a bot variable.

This is equivalent to `! var` in RiveScript. Set the value to `undefined`
to delete a bot variable.
*/
func (rs *RiveScript) SetVariable(name string, value string) {
	if value == "undefined" {
		delete(rs.var_, name)
	} else {
		rs.var_[name] = value
	}
}

/*
SetSubstitution sets a substitution pattern.

This is equivalent to `! sub` in RiveScript. Set the value to `undefined`
to delete a substitution.
*/
func (rs *RiveScript) SetSubstitution(name string, value string) {
	if value == "undefined" {
		delete(rs.sub, name)
	} else {
		rs.sub[name] = value
	}
}

/*
SetPerson sets a person substitution pattern.

This is equivalent to `! person` in RiveScript. Set the value to `undefined`
to delete a person substitution.
*/
func (rs *RiveScript) SetPerson(name string, value string) {
	if value == "undefined" {
		delete(rs.person, name)
	} else {
		rs.person[name] = value
	}
}

/*
SetUservar sets a variable for a user.

This is equivalent to `<set>` in RiveScript. Set the value to `undefined`
to delete a substitution.
*/
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

/*
SetUservars sets a map of variables for a user.

Set multiple user variables by providing a map[string]string of key/value pairs.
Equivalent to calling `SetUservar()` for each pair in the map.
*/
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

/*
GetVariable gets a bot variable.

This is equivalent to `<bot name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (rs *RiveScript) GetVariable(name string) (string, error) {
	if _, ok := rs.var_[name]; ok {
		return rs.var_[name], nil
	}
	return "undefined", errors.New("Variable not found.")
}

/*
GetUservar gets a user variable.

This is equivalent to `<get name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (rs *RiveScript) GetUservar(username string, name string) (string, error) {
	if _, ok := rs.users[username]; ok {
		if _, ok := rs.users[username].data[name]; ok {
			return rs.users[username].data[name], nil
		}
	}
	return "undefined", errors.New("User variable not found.")
}

/*
GetUservars gets all the variables for a user.

This returns a `map[string]string` containing all the user's variables.
*/
func (rs *RiveScript) GetUservars(username string) (map[string]string, error) {
	if _, ok := rs.users[username]; ok {
		return rs.users[username].data, nil
	}
	return map[string]string{}, errors.New("Username not found.")
}

/*
GetAllUservars gets all the variables for all the users.

This returns a map of username (strings) to `map[string]string` of their
variables.
*/
func (rs *RiveScript) GetAllUservars() map[string]map[string]string {
	result := map[string]map[string]string{}
	for username, data := range rs.users {
		result[username] = data.data
	}
	return result
}

/*
ClearUservars clears all a user's variables.
*/
func (rs *RiveScript) ClearUservars(username string) {
	delete(rs.users, username)
}

/*
ClearAllUservars clears all variables for all users.
*/
func (rs *RiveScript) ClearAllUservars() {
	for username, _ := range rs.users {
		delete(rs.users, username)
	}
}

/*
FreezeUservars freezes the variable state of a user.

This will clone and preserve the user's entire variable state, so that it
can be restored later with `ThawUservars()`.
*/
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

/*
ThawUservars unfreezes a user's variables.

The `action` can be one of the following:
* thaw: Restore the variables and delete the frozen copy.
* discard: Don't restore the variables, just delete the frozen copy.
* keep: Keep the frozen copy after restoring.
*/
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

/*
LastMatch returns the user's last matched trigger.
*/
func (rs *RiveScript) LastMatch(username string) (string, error) {
	if _, ok := rs.users[username]; ok {
		return rs.users[username].lastMatch, nil
	}
	return "", errors.New("Username not found.")
}

/*
CurrentUser returns the current user's ID.
*/
func (rs *RiveScript) CurrentUser() string {
	return rs.currentUser
}
