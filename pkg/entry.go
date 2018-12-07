package pkg

import (
	"bytes"
	"log"
	"runtime"
	"sync"
	"time"
)

var (
	bufferPool *sync.Pool

	// qualified package name, cached at first use
	logcPackage string

	// Positions in the call stack when tracing to report the calling method
	minimumCallerDepth int

	// Use for caller information initialisation
	callerInitOnce sync.Mutex
)

const (
	maximumCallerDepth int = 25
	knownLogcFrames    int = 4
)

func init() {
	bufferPool = &sync.Pool{
		New: func() interface{} {
			return new(bytes.Buffer)
		},
	}

	// start at bottom of the stack before the package-name cache is primed
	minimumCallerDepth  = 1
}

// Defines the key when adding errors using WithError.
var Errorkey = "error"

// An entry is the final or intermediate logc entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.

type Entry struct {
	Logger *log.Logger

	// Contains all the fields set by the user.
	Data Fields

	// Time at which the log entry was created
	Time time.Time

	// Level the log entry was logged at: Trace, Debug, Info, Warn, Error, Fatal or Panic
	// This field will be set n entry firing and the value will be equal to the one in Logger struct field.
	Level Level

	// Calling method, with package name
	Caller *runtime.Frame

	// Message passed to Trace, Debug, Info, Warn, Error, Fatal or Panic
	Message string

	// When formatter is called in entry.log(), a Buffer may be set to entry
	Buffer *bytes.Buffer

	// err may contain a field formatting error
	err string
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		Logger: logger,
		// Default is three fields, plus one optional. Give a little extra room
		Data: make(Fields, 6),
	}
}

func (entry *Entry) String() (string, error){
	serialized, err := entry.Logger.Formatter.Format(entry)
	if err != nil {
		return "", err
	}
	str := string(serialized)
	return str, nil
}