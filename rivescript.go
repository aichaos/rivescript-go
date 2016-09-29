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
	"github.com/aichaos/rivescript-go/src"
)

const VERSION string = "0.0.3"

type RiveScript struct {
	rs *src.RiveScript
}

func New() *RiveScript {
	bot := new(RiveScript)
	bot.rs = src.New()
	return bot
}

// Version returns the RiveScript library version.
func (self *RiveScript) Version() string {
	return VERSION
}

// SetDebug enables or disable debug mode.
func (self *RiveScript) SetDebug(value bool) {
	self.rs.Debug = value
}

// GetDebug tells you the current status of the debug mode.
func (self *RiveScript) GetDebug() bool {
	return self.rs.Debug
}

// SetUTF8 enables or disabled UTF-8 mode.
func (self *RiveScript) SetUTF8(value bool) {
	self.rs.UTF8 = value
}

// GetUTF8 returns the current status of UTF-8 mode.
func (self *RiveScript) GetUTF8() bool {
	return self.rs.UTF8
}

// SetUnicodePunctuation allows you to override the text of the unicode
// punctuation regexp. Provide a string literal that will validate in
// `regexp.MustCompile()`
func (self *RiveScript) SetUnicodePunctuation(value string) {
	self.rs.SetUnicodePunctuation(value)
}

// SetDepth lets you override the recursion depth limit (default 50).
func (self *RiveScript) SetDepth(value int) {
	self.rs.Depth = value
}

// GetDepth returns the current recursion depth limit.
func (self *RiveScript) GetDepth() int {
	return self.rs.Depth
}

// SetStrict enables strict syntax checking when parsing RiveScript code.
func (self *RiveScript) SetStrict(value bool) {
	self.rs.Strict = value
}

// GetStrict returns the strict syntax check setting.
func (self *RiveScript) GetStrict() bool {
	return self.rs.Strict
}

////////////////////////////////////////////////////////////////////////////////
////// Loading and Parsing Functions ///////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

/*
LoadFile loads a single RiveScript source file from disk.

Parameters

	path: Path to a RiveScript source file.
*/
func (self *RiveScript) LoadFile(path string) error {
	return self.rs.LoadFile(path)
}

/*
LoadDirectory loads multiple RiveScript documents from a folder on disk.

Parameters

	path: Path to the directory on disk
	extensions...: List of file extensions to filter on, default is
	               '.rive' and '.rs'
*/
func (self *RiveScript) LoadDirectory(path string, extensions ...string) error {
	return self.rs.LoadDirectory(path, extensions...)
}

/*
Stream loads RiveScript code from a text buffer.

Parameters

	code: Raw source code of a RiveScript document, with line breaks after
	      each line.
*/
func (self *RiveScript) Stream(code string) error {
	return self.rs.Stream(code)
}

/*
SortReplies sorts the reply structures in memory for optimal matching.

After you have finished loading your RiveScript code, call this method to
populate the various sort buffers. This is absolutely necessary for reply
matching to work efficiently!
*/
func (self *RiveScript) SortReplies() {
	self.rs.SortReplies()
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
func (self *RiveScript) SetHandler(name string, handler macro.MacroInterface) {
	self.rs.SetHandler(name, handler)
}

/*
RemoveHandler removes an object macro language handler.

Parameters

	lang: The programming language for the handler to remove.
*/
func (self *RiveScript) RemoveHandler(lang string) {
	self.rs.RemoveHandler(lang)
}

/*
SetSubroutine defines a Go object macro from your program.

Parameters

	name: The name of your subroutine for the `<call>` tag in RiveScript.
	fn: A function with a prototype `func(*RiveScript, []string) string`
*/
func (self *RiveScript) SetSubroutine(name string, fn src.Subroutine) {
	self.rs.SetSubroutine(name, fn)
}

/*
DeleteSubroutine removes a Go object macro.

Parameters

	name: The name of the object macro to be deleted.
*/
func (self *RiveScript) DeleteSubroutine(name string) {
	self.rs.DeleteSubroutine(name)
}

/*
SetGlobal sets a global variable.

This is equivalent to `! global` in RiveScript. Set the value to `undefined`
to delete a global.
*/
func (self *RiveScript) SetGlobal(name, value string) {
	self.rs.SetGlobal(name, value)
}

/*
GetGlobal gets a global variable.

This is equivalent to `<env name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (self *RiveScript) GetGlobal(name string) (string, error) {
	return self.rs.GetGlobal(name)
}

/*
SetVariable sets a bot variable.

This is equivalent to `! var` in RiveScript. Set the value to `undefined`
to delete a bot variable.
*/
func (self *RiveScript) SetVariable(name, value string) {
	self.rs.SetVariable(name, value)
}

/*
GetVariable gets a bot variable.

This is equivalent to `<bot name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (self *RiveScript) GetVariable(name string) (string, error) {
	return self.rs.GetVariable(name)
}

/*
SetSubstitution sets a substitution pattern.

This is equivalent to `! sub` in RiveScript. Set the value to `undefined`
to delete a substitution.
*/
func (self *RiveScript) SetSubstitution(name, value string) {
	self.rs.SetSubstitution(name, value)
}

/*
SetPerson sets a person substitution pattern.

This is equivalent to `! person` in RiveScript. Set the value to `undefined`
to delete a person substitution.
*/
func (self *RiveScript) SetPerson(name, value string) {
	self.rs.SetPerson(name, value)
}

/*
SetUservar sets a variable for a user.

This is equivalent to `<set>` in RiveScript. Set the value to `undefined`
to delete a substitution.
*/
func (self *RiveScript) SetUservar(username, name, value string) {
	self.rs.SetUservar(username, name, value)
}

/*
SetUservars sets a map of variables for a user.

Set multiple user variables by providing a map[string]string of key/value pairs.
Equivalent to calling `SetUservar()` for each pair in the map.
*/
func (self *RiveScript) SetUservars(username string, data map[string]string) {
	self.rs.SetUservars(username, data)
}

/*
GetUservar gets a user variable.

This is equivalent to `<get name>` in RiveScript. Returns `undefined` if the
variable isn't defined.
*/
func (self *RiveScript) GetUservar(username, name string) (string, error) {
	return self.rs.GetUservar(username, name)
}

/*
GetUservars gets all the variables for a user.

This returns a `map[string]string` containing all the user's variables.
*/
func (self *RiveScript) GetUservars(username string) (map[string]string, error) {
	return self.rs.GetUservars(username)
}

/*
GetAllUservars gets all the variables for all the users.

This returns a map of username (strings) to `map[string]string` of their
variables.
*/
func (self *RiveScript) GetAllUservars() map[string]map[string]string {
	return self.rs.GetAllUservars()
}

// ClearAllUservars clears all variables for all users.
func (self *RiveScript) ClearAllUservars() {
	self.rs.ClearAllUservars()
}

// ClearUservars clears all a user's variables.
func (self *RiveScript) ClearUservars(username string) {
	self.rs.ClearUservars(username)
}

/*
FreezeUservars freezes the variable state of a user.

This will clone and preserve the user's entire variable state, so that it
can be restored later with `ThawUservars()`.
*/
func (self *RiveScript) FreezeUservars(username string) error {
	return self.rs.FreezeUservars(username)
}

/*
ThawUservars unfreezes a user's variables.

The `action` can be one of the following:
* thaw: Restore the variables and delete the frozen copy.
* discard: Don't restore the variables, just delete the frozen copy.
* keep: Keep the frozen copy after restoring.
*/
func (self *RiveScript) ThawUservars(username, action string) error {
	return self.rs.ThawUservars(username, action)
}

// LastMatch returns the user's last matched trigger.
func (self *RiveScript) LastMatch(username string) (string, error) {
	return self.rs.LastMatch(username)
}

/*
CurrentUser returns the current user's ID.

This is only useful from within an object macro, to get the ID of the user who
invoked the macro. This value is set at the beginning of `Reply()` and unset
at the end, so this function will return empty outside of a reply context.
*/
func (self *RiveScript) CurrentUser() string {
	return self.rs.CurrentUser()
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
func (self *RiveScript) Reply(username, message string) string {
	return self.rs.Reply(username, message)
}

////////////////////////////////////////////////////////////////////////////////
////// Debugging Functions /////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////

// DumpSorted is a debug method which dumps the sort tree from the bot's memory.
func (self *RiveScript) DumpSorted() {
	self.rs.DumpSorted()
}

// DumpTopics is a debug method which dumps the topic structure from the bot's memory.
func (self *RiveScript) DumpTopics() {
	self.rs.DumpTopics()
}
