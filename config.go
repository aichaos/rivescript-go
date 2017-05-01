package rivescript

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/aichaos/rivescript-go/macro"
	"github.com/aichaos/rivescript-go/sessions"
)

/*
Config provides options to configure the RiveScript bot.

Create a pointer to this type and send it to the New() constructor to change
the default settings. You only need to provide settings you want to override;
the zero-values of all the options are handled appropriately by the RiveScript
library.

The default values are documented below.
*/
type Config struct {
	// Debug enables verbose logging to standard output. Default false.
	Debug bool

	// Strict enables strict syntax checking, where a syntax error in RiveScript
	// code is considered fatal at parse time. Default true.
	Strict bool

	// UTF8 enables UTF-8 mode within the bot. Default false.
	//
	// When UTF-8 mode is enabled, triggers in the RiveScript source files are
	// allowed to contain foreign characters. Additionally, the user's incoming
	// messages are left *mostly* intact, so that they send messages with
	// foreign characters to the bot.
	UTF8 bool

	// Depth controls the global limit for recursive functions within
	// RiveScript. Default 50.
	Depth uint

	// Random number seed, if you'd like to customize it. The default is for
	// RiveScript to choose its own seed, `time.Now().UnixNano()`
	Seed int64

	// SessionManager is an implementation of the same name for managing user
	// variables for the bot. The default is the in-memory session handler.
	SessionManager sessions.SessionManager
}

// WithUTF8 provides a Config object that enables UTF-8 mode.
func WithUTF8() *Config {
	return &Config{
		UTF8: true,
	}
}

/*
SetHandler sets a custom language handler for RiveScript object macros.

Parameters

	lang: What your programming language is called, e.g. "javascript"
	handler: An implementation of macro.MacroInterface.
*/
func (rs *RiveScript) SetHandler(lang string, handler macro.MacroInterface) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	rs.handlers[lang] = handler
}

/*
RemoveHandler removes an object macro language handler.

If the handler has already loaded object macros, they will be deleted from
the bot along with the handler.

Parameters

	lang: The programming language for the handler to remove.
*/
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

/*
SetSubroutine defines a Go object macro from your program.

Parameters

	name: The name of your subroutine for the `<call>` tag in RiveScript.
	fn: A function with a prototype `func(*RiveScript, []string) string`
*/
func (rs *RiveScript) SetSubroutine(name string, fn Subroutine) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	rs.subroutines[name] = fn
}

/*
DeleteSubroutine removes a Go object macro.

Parameters

	name: The name of the object macro to be deleted.
*/
func (rs *RiveScript) DeleteSubroutine(name string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	delete(rs.subroutines, name)
}

/*
SetGlobal sets a global variable.

This is equivalent to `! global` in RiveScript. Set the value to `undefined`
to delete a global.
*/
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

/*
SetVariable sets a bot variable.

This is equivalent to `! var` in RiveScript. Set the value to `undefined`
to delete a bot variable.
*/
func (rs *RiveScript) SetVariable(name, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == UNDEFINED {
		delete(rs.vars, name)
	} else {
		rs.vars[name] = value
	}
}

/*
SetSubstitution sets a substitution pattern.

This is equivalent to `! sub` in RiveScript. Set the value to `undefined`
to delete a substitution.
*/
func (rs *RiveScript) SetSubstitution(name, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == UNDEFINED {
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
func (rs *RiveScript) SetPerson(name, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == UNDEFINED {
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
func (rs *RiveScript) SetUservar(username, name, value string) {
	rs.sessions.Set(username, map[string]string{
		name: value,
	})
}

/*
SetUservars sets a map of variables for a user.

Set multiple user variables by providing a map[string]string of key/value pairs.
Equivalent to calling `SetUservar()` for each pair in the map.
*/
func (rs *RiveScript) SetUservars(username string, data map[string]string) {
	rs.sessions.Set(username, data)
}

/*
GetGlobal gets a global variable.

This is equivalent to `<env name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
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

/*
GetVariable gets a bot variable.

This is equivalent to `<bot name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (rs *RiveScript) GetVariable(name string) (string, error) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if _, ok := rs.vars[name]; ok {
		return rs.vars[name], nil
	}
	return UNDEFINED, fmt.Errorf("bot variable %s not found", name)
}

/*
GetUservar gets a user variable.

This is equivalent to `<get name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (rs *RiveScript) GetUservar(username, name string) (string, error) {
	return rs.sessions.Get(username, name)
}

/*
GetUservars gets all the variables for a user.

This returns a `map[string]string` containing all the user's variables.
*/
func (rs *RiveScript) GetUservars(username string) (*sessions.UserData, error) {
	return rs.sessions.GetAny(username)
}

/*
GetAllUservars gets all the variables for all the users.

This returns a map of username (strings) to `map[string]string` of their
variables.
*/
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

/*
FreezeUservars freezes the variable state of a user.

This will clone and preserve the user's entire variable state, so that it
can be restored later with `ThawUservars()`.
*/
func (rs *RiveScript) FreezeUservars(username string) error {
	return rs.sessions.Freeze(username)
}

/*
ThawUservars unfreezes a user's variables.

The `action` can be one of the following:
* thaw: Restore the variables and delete the frozen copy.
* discard: Don't restore the variables, just delete the frozen copy.
* keep: Keep the frozen copy after restoring.
*/
func (rs *RiveScript) ThawUservars(username string, action sessions.ThawAction) error {
	return rs.sessions.Thaw(username, action)
}

// LastMatch returns the user's last matched trigger.
func (rs *RiveScript) LastMatch(username string) (string, error) {
	return rs.sessions.GetLastMatch(username)
}

/*
CurrentUser returns the current user's ID.

This is only useful from within an object macro, to get the ID of the user who
invoked the macro. This value is set at the beginning of `Reply()` and unset
at the end, so this function will return empty outside of a reply context.
*/
func (rs *RiveScript) CurrentUser() (string, error) {
	if rs.inReplyContext {
		return rs.currentUser, nil
	}
	return "", errors.New("CurrentUser() can only be called inside a reply context")
}
