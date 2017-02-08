package rivescript

// Public API Configuration Methods

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aichaos/rivescript-go/macro"
	"github.com/aichaos/rivescript-go/sessions"
)

// SetHandler registers a handler for foreign language object macros.
func (rs *RiveScript) SetHandler(lang string, handler macro.MacroInterface) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	rs.handlers[lang] = handler
}

// RemoveHandler deletes support for a foreign language object macro.
func (rs *RiveScript) RemoveHandler(lang string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	// Purge all loaded objects for this handler.
	for name, language := range rs.objlangs {
		if language == lang {
			delete(rs.objlangs, name)
		}
	}

	// And delete the handler itself.
	delete(rs.handlers, lang)
}

// SetSubroutine defines a Go function to handle an object macro for RiveScript.
func (rs *RiveScript) SetSubroutine(name string, fn Subroutine) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	rs.subroutines[name] = fn
}

// DeleteSubroutine deletes a Go object macro handler.
func (rs *RiveScript) DeleteSubroutine(name string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	delete(rs.subroutines, name)
}

// SetGlobal configures a global variable in RiveScript.
func (rs *RiveScript) SetGlobal(name, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	// Special globals that reconfigure the interpreter.
	if name == "debug" {
		switch strings.ToLower(value) {
		case "true", "t", "on", "yes":
			rs.Debug = true
		default:
			rs.Debug = false
		}
	} else if name == "depth" {
		depth, err := strconv.Atoi(value)
		if err != nil {
			rs.warn("Can't set global `depth` to `%s`: %s\n", value, err)
		} else {
			rs.Depth = uint(depth)
		}
	}

	if value == UNDEFINED {
		delete(rs.global, name)
	} else {
		rs.global[name] = value
	}
}

// SetVariable configures a bot variable in RiveScript.
func (rs *RiveScript) SetVariable(name, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == UNDEFINED {
		delete(rs.vars, name)
	} else {
		rs.vars[name] = value
	}
}

// SetSubstitution sets a substitution pattern.
func (rs *RiveScript) SetSubstitution(name, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == UNDEFINED {
		delete(rs.sub, name)
	} else {
		rs.sub[name] = value
	}
}

// SetPerson sets a person substitution.
func (rs *RiveScript) SetPerson(name, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == UNDEFINED {
		delete(rs.person, name)
	} else {
		rs.person[name] = value
	}
}

// SetUservar sets a user variable to a value.
func (rs *RiveScript) SetUservar(username, name, value string) {
	rs.sessions.Set(username, map[string]string{
		name: value,
	})
}

// SetUservars sets many user variables at a time.
func (rs *RiveScript) SetUservars(username string, data map[string]string) {
	rs.sessions.Set(username, data)
}

// GetGlobal retrieves the value of a global variable.
func (rs *RiveScript) GetGlobal(name string) (string, error) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	// Special globals.
	if name == "debug" {
		return fmt.Sprintf("%v", rs.Debug), nil
	} else if name == "depth" {
		return strconv.Itoa(int(rs.Depth)), nil
	}

	if _, ok := rs.global[name]; ok {
		return rs.global[name], nil
	}
	return UNDEFINED, fmt.Errorf("global variable %s not found", name)
}

// GetVariable retrieves the value of a bot variable.
func (rs *RiveScript) GetVariable(name string) (string, error) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if _, ok := rs.vars[name]; ok {
		return rs.vars[name], nil
	}
	return UNDEFINED, fmt.Errorf("bot variable %s not found", name)
}

// GetUservar retrieves the value of a user variable.
func (rs *RiveScript) GetUservar(username, name string) (string, error) {
	return rs.sessions.Get(username, name)
}

// GetUservars retrieves all variables about a user.
func (rs *RiveScript) GetUservars(username string) (*sessions.UserData, error) {
	return rs.sessions.GetAny(username)
}

// GetAllUservars gets all variables about all users.
func (rs *RiveScript) GetAllUservars() map[string]*sessions.UserData {
	return rs.sessions.GetAll()
}

// ClearUservars deletes all the variables that belong to a user.
func (rs *RiveScript) ClearUservars(username string) {
	rs.sessions.Clear(username)
}

// ClearAllUservars deletes all variables for all users.
func (rs *RiveScript) ClearAllUservars() {
	rs.sessions.ClearAll()
}

// FreezeUservars takes a snapshot of a user's variables.
func (rs *RiveScript) FreezeUservars(username string) error {
	return rs.sessions.Freeze(username)
}

// ThawUservars restores a snapshot of user variables.
func (rs *RiveScript) ThawUservars(username string, action sessions.ThawAction) error {
	return rs.sessions.Thaw(username, action)
}

// LastMatch returns the last match of the user.
func (rs *RiveScript) LastMatch(username string) (string, error) {
	return rs.sessions.GetLastMatch(username)
}

// CurrentUser returns the current user and can only be called from within an
// object macro context.
func (rs *RiveScript) CurrentUser() (string, error) {
	if rs.inReplyContext {
		return rs.currentUser, nil
	}
	return "", errors.New("CurrentUser() can only be called inside a reply context")
}
