package pkg

import (
	"time"
)

// Default key names for the default fields
const (
	defaultTimestampFormat = time.RFC3339
	FieldKeyMsg = "msg"
	FieldKeyLevel = "level"
	FieldKeyTime = "time"
	FieldKeyLogcError = "logc_error"
	FieldKeyFunc = "func"
	FieldKeyField = "file"
)

// The Formatter interface is used to implement a custom Formatter. It takes an
// `Entry`. It exposes all the fields, including the default ones:

// *`entry.Data["msg"]`. The message passed form Info, Warn, Error..
// *`entry.Data["time"]`. The timestamp
// *`entry.Data["level"]. The level the entry was logged at.

// Any additional fields added with `WitheField` or `WithFields` are also in
// `entry.Data`. Format is expected to return an array of bytes which are then
// logged to `logger.Out`.
type Formatter interface {
	Format(*Entry) ([]byte, error)
}

// This is to not silently overwrite `time`, `msg`, `func` and `level` fields when
// dumping it. If this code wasn't there doing:
//
// logc.WithField("level", 1).Info("hello")
// Would just silently drop the user provided level. Instead with this code
// it'll logged as:

// {"level": "info", "fields.level":1, "msg": "helo", "time": ".."}

// It's not exported because it's still using Data in an opinionated way. It's to
// avoid code duplication between the tow default formatters.
