package rivescript

// deprecated.go is where functions that are deprecated move to.

import (
	"fmt"
	"os"
)

// common function to put the deprecated note.
func deprecated(name, since string) {
	fmt.Fprintf(
		os.Stderr,
		"Use of 'rivescript.%s()' is deprecated since v%s (this is v%s)\n",
		name,
		since,
		VERSION,
	)
}

// SetDebug enables or disable debug mode.
func (self *RiveScript) SetDebug(value bool) {
	deprecated("SetDebug", "0.1.0")
	self.rs.Debug = value
}

// GetDebug tells you the current status of the debug mode.
func (self *RiveScript) GetDebug() bool {
	deprecated("GetDebug", "0.1.0")
	return self.rs.Debug
}

// SetUTF8 enables or disabled UTF-8 mode.
func (self *RiveScript) SetUTF8(value bool) {
	deprecated("SetUTF8", "0.1.0")
	self.rs.UTF8 = value
}

// GetUTF8 returns the current status of UTF-8 mode.
func (self *RiveScript) GetUTF8() bool {
	deprecated("GetUTF8", "0.1.0")
	return self.rs.UTF8
}

// SetDepth lets you override the recursion depth limit (default 50).
func (self *RiveScript) SetDepth(value uint) {
	deprecated("SetDepth", "0.1.0")
	self.rs.Depth = value
}

// GetDepth returns the current recursion depth limit.
func (self *RiveScript) GetDepth() uint {
	deprecated("GetDepth", "0.1.0")
	return self.rs.Depth
}

// SetStrict enables strict syntax checking when parsing RiveScript code.
func (self *RiveScript) SetStrict(value bool) {
	deprecated("SetStrict", "0.1.0")
	self.rs.Strict = value
}

// GetStrict returns the strict syntax check setting.
func (self *RiveScript) GetStrict() bool {
	deprecated("GetStrict", "0.1.0")
	return self.rs.Strict
}
