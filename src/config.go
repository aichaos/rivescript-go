package rivescript

// Public API Configuration Methods

import (
	"errors"

	"github.com/aichaos/rivescript-go/macro"
	"github.com/aichaos/rivescript-go/sessions"
)

func (rs *RiveScript) SetHandler(lang string, handler macro.MacroInterface) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	rs.handlers[lang] = handler
}

func (rs *RiveScript) RemoveHandler(lang string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	delete(rs.handlers, lang)
}

func (rs *RiveScript) SetSubroutine(name string, fn Subroutine) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	rs.subroutines[name] = fn
}

func (rs *RiveScript) DeleteSubroutine(name string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	delete(rs.subroutines, name)
}

func (rs *RiveScript) SetGlobal(name string, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == "undefined" {
		delete(rs.global, name)
	} else {
		rs.global[name] = value
	}
}

func (rs *RiveScript) SetVariable(name string, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == "undefined" {
		delete(rs.var_, name)
	} else {
		rs.var_[name] = value
	}
}

func (rs *RiveScript) SetSubstitution(name string, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == "undefined" {
		delete(rs.sub, name)
	} else {
		rs.sub[name] = value
	}
}

func (rs *RiveScript) SetPerson(name string, value string) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if value == "undefined" {
		delete(rs.person, name)
	} else {
		rs.person[name] = value
	}
}

func (rs *RiveScript) SetUservar(username string, name string, value string) {
	rs.sessions.Set(username, map[string]string{
		name: value,
	})
}

func (rs *RiveScript) SetUservars(username string, data map[string]string) {
	rs.sessions.Set(username, data)
}

func (rs *RiveScript) GetGlobal(name string) (string, error) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if _, ok := rs.global[name]; ok {
		return rs.global[name], nil
	}
	return "undefined", errors.New("Global variable not found.")
}

func (rs *RiveScript) GetVariable(name string) (string, error) {
	rs.cLock.Lock()
	defer rs.cLock.Unlock()

	if _, ok := rs.var_[name]; ok {
		return rs.var_[name], nil
	}
	return "undefined", errors.New("Variable not found.")
}

func (rs *RiveScript) GetUservar(username string, name string) (string, error) {
	return rs.sessions.Get(username, name)
}

func (rs *RiveScript) GetUservars(username string) (*sessions.UserData, error) {
	return rs.sessions.GetAny(username)
}

func (rs *RiveScript) GetAllUservars() map[string]*sessions.UserData {
	return rs.sessions.GetAll()
}

func (rs *RiveScript) ClearUservars(username string) {
	rs.sessions.Clear(username)
}

func (rs *RiveScript) ClearAllUservars() {
	rs.sessions.ClearAll()
}

func (rs *RiveScript) FreezeUservars(username string) error {
	return rs.sessions.Freeze(username)
}

func (rs *RiveScript) ThawUservars(username string, action sessions.ThawAction) error {
	return rs.sessions.Thaw(username, action)
}

func (rs *RiveScript) LastMatch(username string) (string, error) {
	return rs.sessions.GetLastMatch(username)
}

func (rs *RiveScript) CurrentUser() string {
	return rs.currentUser
}
