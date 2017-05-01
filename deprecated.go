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
		Version,
	)
}

// SetDebug enables or disable debug mode.
func (rs *RiveScript) SetDebug(value bool) {
	deprecated("SetDebug", "0.1.0")
	rs.Debug = value
}

// GetDebug tells you the current status of the debug mode.
func (rs *RiveScript) GetDebug() bool {
	deprecated("GetDebug", "0.1.0")
	return rs.Debug
}

// SetUTF8 enables or disabled UTF-8 mode.
func (rs *RiveScript) SetUTF8(value bool) {
	deprecated("SetUTF8", "0.1.0")
	rs.UTF8 = value
}

// GetUTF8 returns the current status of UTF-8 mode.
func (rs *RiveScript) GetUTF8() bool {
	deprecated("GetUTF8", "0.1.0")
	return rs.UTF8
}

// SetDepth lets you override the recursion depth limit (default 50).
func (rs *RiveScript) SetDepth(value uint) {
	deprecated("SetDepth", "0.1.0")
	rs.Depth = value
}

// GetDepth returns the current recursion depth limit.
func (rs *RiveScript) GetDepth() uint {
	deprecated("GetDepth", "0.1.0")
	return rs.Depth
}

// SetStrict enables strict syntax checking when parsing RiveScript code.
func (rs *RiveScript) SetStrict(value bool) {
	deprecated("SetStrict", "0.1.0")
	rs.Strict = value
}

// GetStrict returns the strict syntax check setting.
func (rs *RiveScript) GetStrict() bool {
	deprecated("GetStrict", "0.1.0")
	return rs.Strict
}
