package rivescript

/*
	NOTE: This module is a wrapper around the bulk of the actual source code
	under the 'src/' subpackage. This gives multiple benefits such as keeping
	the root of the git repo as tidy as possible (low number of source files)
	and keeping the public facing, official API in one small place in the code.

	Everything exported from the 'src' subpackage should not be used directly
	by third party developers. A lot of the symbols from the src package must
	be exported to get this wrapper program to work (and keep a nice looking
	module import path), but only this public facing API module should be used.
*/

import (
	"github.com/aichaos/rivescript-go/macro"
	"github.com/aichaos/rivescript-go/sessions"
	"github.com/aichaos/rivescript-go/sessions/memory"
	"github.com/aichaos/rivescript-go/src"
)

// VERSION describes the module version.
const VERSION string = "0.2.0"

// RiveScript represents an individual chatbot instance.
type RiveScript struct {
	rs *rivescript.RiveScript
}

/*
New creates a new RiveScript instance.

A RiveScript instance represents one chat bot personality; it has its own
replies and its own memory of user data. You could make multiple bots in the
same program, each with its own replies loaded from different sources.
*/
func New(cfg *Config) *RiveScript {
	bot := &RiveScript{
		rs: rivescript.New(),
	}

	// If no config was given, default to the BasicConfig.
	if cfg == nil {
		cfg = &Config{
			Strict: true,
			Depth:  50,
		}
	}

	// If no session manager configured, default to the in-memory one.
	if cfg.SessionManager == nil {
		cfg.SessionManager = memory.New()
	}

	// Default depth if not given is 50.
	if cfg.Depth <= 0 {
		cfg.Depth = 50
	}

	bot.rs.Configure(cfg.Debug, cfg.Strict, cfg.UTF8, cfg.Depth, memory.New())

	return bot
}

// Version returns the RiveScript library version.
func (rs *RiveScript) Version() string {
	return VERSION
}

// SetUnicodePunctuation allows you to override the text of the unicode
// punctuation regexp. Provide a string literal that will validate in
// `regexp.MustCompile()`
func (rs *RiveScript) SetUnicodePunctuation(value string) {
	rs.rs.SetUnicodePunctuation(value)
}

////////////////////////////////////////////////////////////////////////////////
////// Loading and Parsing Functions ///////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

/*
LoadFile loads a single RiveScript source file from disk.

Parameters

	path: Path to a RiveScript source file.
*/
func (rs *RiveScript) LoadFile(path string) error {
	return rs.rs.LoadFile(path)
}

/*
LoadDirectory loads multiple RiveScript documents from a folder on disk.

Parameters

	path: Path to the directory on disk
	extensions...: List of file extensions to filter on, default is
	               '.rive' and '.rs'
*/
func (rs *RiveScript) LoadDirectory(path string, extensions ...string) error {
	return rs.rs.LoadDirectory(path, extensions...)
}

/*
Stream loads RiveScript code from a text buffer.

Parameters

	code: Raw source code of a RiveScript document, with line breaks after
	      each line.
*/
func (rs *RiveScript) Stream(code string) error {
	return rs.rs.Stream(code)
}

/*
SortReplies sorts the reply structures in memory for optimal matching.

After you have finished loading your RiveScript code, call this method to
populate the various sort buffers. This is absolutely necessary for reply
matching to work efficiently!

If the bot has loaded no topics, or if it ends up with no sorted triggers at
the end, it will return an error saying such. This usually means the bot didn't
load any RiveScript code, for example because it looked in the wrong directory.
*/
func (rs *RiveScript) SortReplies() error {
	return rs.rs.SortReplies()
}

////////////////////////////////////////////////////////////////////////////////
////// Public Configuration Functions //////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

/*
SetHandler sets a custom language handler for RiveScript object macros.

Parameters

	lang: What your programming language is called, e.g. "javascript"
	handler: An implementation of macro.MacroInterface.
*/
func (rs *RiveScript) SetHandler(name string, handler macro.MacroInterface) {
	rs.rs.SetHandler(name, handler)
}

/*
RemoveHandler removes an object macro language handler.

If the handler has already loaded object macros, they will be deleted from
the bot along with the handler.

Parameters

	lang: The programming language for the handler to remove.
*/
func (rs *RiveScript) RemoveHandler(lang string) {
	rs.rs.RemoveHandler(lang)
}

/*
SetSubroutine defines a Go object macro from your program.

Parameters

	name: The name of your subroutine for the `<call>` tag in RiveScript.
	fn: A function with a prototype `func(*RiveScript, []string) string`
*/
func (rs *RiveScript) SetSubroutine(name string, fn rivescript.Subroutine) {
	rs.rs.SetSubroutine(name, fn)
}

/*
DeleteSubroutine removes a Go object macro.

Parameters

	name: The name of the object macro to be deleted.
*/
func (rs *RiveScript) DeleteSubroutine(name string) {
	rs.rs.DeleteSubroutine(name)
}

/*
SetGlobal sets a global variable.

This is equivalent to `! global` in RiveScript. Set the value to `undefined`
to delete a global.
*/
func (rs *RiveScript) SetGlobal(name, value string) {
	rs.rs.SetGlobal(name, value)
}

/*
GetGlobal gets a global variable.

This is equivalent to `<env name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (rs *RiveScript) GetGlobal(name string) (string, error) {
	return rs.rs.GetGlobal(name)
}

/*
SetVariable sets a bot variable.

This is equivalent to `! var` in RiveScript. Set the value to `undefined`
to delete a bot variable.
*/
func (rs *RiveScript) SetVariable(name, value string) {
	rs.rs.SetVariable(name, value)
}

/*
GetVariable gets a bot variable.

This is equivalent to `<bot name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (rs *RiveScript) GetVariable(name string) (string, error) {
	return rs.rs.GetVariable(name)
}

/*
SetSubstitution sets a substitution pattern.

This is equivalent to `! sub` in RiveScript. Set the value to `undefined`
to delete a substitution.
*/
func (rs *RiveScript) SetSubstitution(name, value string) {
	rs.rs.SetSubstitution(name, value)
}

/*
SetPerson sets a person substitution pattern.

This is equivalent to `! person` in RiveScript. Set the value to `undefined`
to delete a person substitution.
*/
func (rs *RiveScript) SetPerson(name, value string) {
	rs.rs.SetPerson(name, value)
}

/*
SetUservar sets a variable for a user.

This is equivalent to `<set>` in RiveScript. Set the value to `undefined`
to delete a substitution.
*/
func (rs *RiveScript) SetUservar(username, name, value string) {
	rs.rs.SetUservar(username, name, value)
}

/*
SetUservars sets a map of variables for a user.

Set multiple user variables by providing a map[string]string of key/value pairs.
Equivalent to calling `SetUservar()` for each pair in the map.
*/
func (rs *RiveScript) SetUservars(username string, data map[string]string) {
	rs.rs.SetUservars(username, data)
}

/*
GetUservar gets a user variable.

This is equivalent to `<get name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (rs *RiveScript) GetUservar(username, name string) (string, error) {
	return rs.rs.GetUservar(username, name)
}

/*
GetUservars gets all the variables for a user.

This returns a `map[string]string` containing all the user's variables.
*/
func (rs *RiveScript) GetUservars(username string) (*sessions.UserData, error) {
	return rs.rs.GetUservars(username)
}

/*
GetAllUservars gets all the variables for all the users.

This returns a map of username (strings) to `map[string]string` of their
variables.
*/
func (rs *RiveScript) GetAllUservars() map[string]*sessions.UserData {
	return rs.rs.GetAllUservars()
}

// ClearAllUservars clears all variables for all users.
func (rs *RiveScript) ClearAllUservars() {
	rs.rs.ClearAllUservars()
}

// ClearUservars clears all a user's variables.
func (rs *RiveScript) ClearUservars(username string) {
	rs.rs.ClearUservars(username)
}

/*
FreezeUservars freezes the variable state of a user.

This will clone and preserve the user's entire variable state, so that it
can be restored later with `ThawUservars()`.
*/
func (rs *RiveScript) FreezeUservars(username string) error {
	return rs.rs.FreezeUservars(username)
}

/*
ThawUservars unfreezes a user's variables.

The `action` can be one of the following:
* thaw: Restore the variables and delete the frozen copy.
* discard: Don't restore the variables, just delete the frozen copy.
* keep: Keep the frozen copy after restoring.
*/
func (rs *RiveScript) ThawUservars(username string, action sessions.ThawAction) error {
	return rs.rs.ThawUservars(username, action)
}

// LastMatch returns the user's last matched trigger.
func (rs *RiveScript) LastMatch(username string) (string, error) {
	return rs.rs.LastMatch(username)
}

/*
CurrentUser returns the current user's ID.

This is only useful from within an object macro, to get the ID of the user who
invoked the macro. This value is set at the beginning of `Reply()` and unset
at the end, so this function will return empty outside of a reply context.
*/
func (rs *RiveScript) CurrentUser() (string, error) {
	return rs.rs.CurrentUser()
}

////////////////////////////////////////////////////////////////////////////////
////// Reply Fetching Functions ////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

/*
Reply fetches a reply from the bot for a user's message.

Parameters

	username: The name of the user requesting a reply.
	message: The user's message.
*/
func (rs *RiveScript) Reply(username, message string) (string, error) {
	reply, err := rs.rs.Reply(username, message)
	return reply, err
}

////////////////////////////////////////////////////////////////////////////////
////// Debugging Functions /////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// DumpSorted is a debug method which dumps the sort tree from the bot's memory.
func (rs *RiveScript) DumpSorted() {
	rs.rs.DumpSorted()
}

// DumpTopics is a debug method which dumps the topic structure from the bot's memory.
func (rs *RiveScript) DumpTopics() {
	rs.rs.DumpTopics()
}
