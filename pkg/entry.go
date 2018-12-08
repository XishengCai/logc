package pkg

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
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
	callerInitOnce sync.Once
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
var ErrorKey = "error"

// An entry is the final or intermediate logc entry. It contains all
// the fields passed with WithField{,s}. It's finally logged when Trace, Debug,
// Info, Warn, Error, Fatal or Panic is called on it. These objects can be
// reused and passed around as much as you wish to avoid field duplication.

type Entry struct {
	Logger *Logger

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

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry{
	return entry.WitheField(ErrorKey, err)
}

// Add a single field to the Entry
func (entry *Entry) WitheField(key string, value interface{}) *Entry{
	return entry.WitheFields(Fields{key: value})
}

// Add a map of fields to the Entry
func(entry *Entry) WitheFields(fields Fields) *Entry {
	data := make(Fields, len(entry.Data) + len(fields))
	for k, v := range entry.Data {
		data[k] = v
	}

	var field_err string
	for k, v := range fields {
		if t := reflect.TypeOf(v); t != nil && t.Kind() == reflect.Func {
			field_err = fmt.Sprintf("can not add fields %q", k)
			if entry.err != "" {
				field_err = entry.err + "," + field_err
			}
		}else {
			data[k] = v
		}
	}
	return &Entry{Logger: entry.Logger, Data: data, Time: entry.Time, err: field_err}
}


// Overrides the time of the Entry.
func(entry *Entry) WitheTime(t time.Time) *Entry {
	return &Entry{Logger: entry.Logger, Data: entry.Data, Time: t}
}


// getPackageName reduce a fully qualified function name to the package name
// There really ought to be a better way...
func getPackageName(f string) string {
	for {
		lastPeriod := strings.LastIndex(f, ".")
		lastSlash := strings.LastIndex(f, "/")
		if lastPeriod > lastSlash {
			f = f[:lastPeriod]
		}else {
			break
		}
	}
	return f
}

// getCaller retrieves the name of the first non-logc calling calling function
func getCaller() *runtime.Frame {
	// Retrict the lookback frames to avoid runway lookups
	pcs := make([]uintptr, maximumCallerDepth)
	depth := runtime.Callers(minimumCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	// cache this package's full-qualified name
	callerInitOnce.Do(func(){
		logcPackage = getPackageName(runtime.FuncForPC(pcs[0]).Name())

		// now that we have the cache, we can skip a minimum count of known-logc functions
		// xxx this is dubious, the number of frames may very store an entry in a logger interface
		minimumCallerDepth = knownLogcFrames
	})

	for f, again := frames.Next(); again; f, again = frames.Next() {
		pkg := getPackageName(f.Function)

		// If the caller isn't part of this package, we're done
		if pkg != logcPackage {
			return &f
		}
	}

	// if we got here, we failed to fine the caller's context
	return nil
}

func (entry Entry) HasCaller() (has bool) {
	return entry.Logger != nil &&
		entry.Logger.ReportCaller &&
		entry.Caller != nil
}

// This function is not declared with a pointer value because otherwise
// race conditions will occur when using multiple goroutines.
func(entry Entry) log(level Level, msg string) {
	var buffer *bytes.Buffer

	// Default to now, but allow users to override if they want.
	//
	// We don't have to worry about polluting future calls to Entry#log()
	// with this assignment because this function is declared with a
	// non-point receiver.
	if entry.Time.IsZero() {
		entry.Time = time.Now()
	}

	entry.Level = level
	entry.Message = msg
	if entry.Logger.ReportCaller {
		entry.Caller = getCaller()
	}

	entry.fireHooks()

	buffer = bufferPool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer bufferPool.Put(buffer)
	entry.Buffer = buffer

	entry.write()

	entry.Buffer = nil

	// To avoid Entry#log() returning a value that only would mak sense for
	// panic() to use in Entry#Panic(), we avoid the allocation by checking
	// directly here.
	if level <= PanicLevel() {
		panic(&entry)
	}
}

func (entry *Entry) write() {
	entry.Logger.mu.Lock()

}
