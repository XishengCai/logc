package pkg

import (
	"github.com/emicklei/go-restful"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sync"
	"time"
)

type Logger struct {
	// The logs are `io.Copy`'d to this in a mutex. It's common to set this to a
	// file, or leave it default which is `os.Stderr`. You can also set this to
	// something more adventurous, such as logging to Kafka.
	Out io.Writer

	// Hooks for the logger instance. These allow firing events based on logging
	// levels and log entries. For example, to send errors to an error tracking
	// service, log to StatsD or dump the core on fatal errors.
	Hooks LevelHooks

	// All log entries pass through the formatter before logged to Out. The
	// include formatters are `TextFormatter` and `JSONFormatter` for which
	// TextFormatter is the default. In development(when a TTY is attached) it
	// logs with colors, but to a file it wouldn't. You can easily implement your
	// own that implements the `Formatter` interface, see the `README` or included
	// formatters for examples.
	Formatter Formatter

	// Flag for whether to log caller info(off by default)
	ReportCaller bool

	// The logging level the logger should log at. This is typically (add defaults
	// to) `loc.Info`, which allows Info(), Warn(), Error() and Fatal() to be
	// logged.
	Level Level

	// Used to sync writing to the log. Locking is enabled by Default.
	mu  MutexWrap

	// Reusable empty entry
	entryPool sync.Pool

	// Function to exit the application, default to `os.Exit()`
	ExitFunc exitFunc

}

type exitFunc func(int)

type MutexWrap struct {
	lock sync.Mutex
	disabled bool
}

func(mw *MutexWrap) Lock() {
	if !mw.disabled {
		mw.lock.Lock()
	}
}

func (mw *MutexWrap) Unlock() {
	if !mw.disabled {
		mw.lock.Unlock()
	}
}

func (mw *MutexWrap) Disable() {
	mw.disabled = true
}

// Create a new logger. Configuration should be set by changing `Formatter`,
// `Out` and `Hooks` directly on the default logger instance. You can also just
// instantiate your own:
//
// 	var log = &Logger {
//		Out: os.Stderr,
// 		Formatter: new(JSONFormatter),
// 		Hooks: make(LevelHooks),
// 		Level: logc.DebugLevel,
//  }
//
// It's recommended to make this a global instance called `log`.
func New() *Logger {
	return &Logger {
			Out:	os.Stderr,
			Formatter:	new(TextFormatter),
			Hooks:	make(LevelHooks),
			Level:	InfoLevel,
			ExitFunc:os.Exit,
			ReportCaller: false,
	}
}

func (logger *Logger) newEntry() *Entry {
	entry, ok := logger.entryPool.Get().(*Entry)
	if ok {
		return entry
	}
	return NewEntry(logger)
}

func (logger *Logger) releaseEntry(entry *Entry) {
	entry.Data = map[string]interface{}{}
	logger.entryPool.Put(entry)
}

// Adds a filed to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry
// If you want multiple fields, use `WithFields`
func (logger *Logger) WithField(key string, value interface{}) *Entry{
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WitheField(key, value)
}

// Adds a struct of fields to the log entry. All it does it call `WithField` for
// each `Field`.
func (logger *Logger)WithFields(fields Fields) *Entry{
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WitheFields(fields)
}

// Add an error as single field to the log entry. All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WitheError(err error) *Entry{
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WithError(err)
}

// Overrides the time of the log entry.
func (logger *Logger)WithTime(t time.Time) *Entry {
	entry := logger.newEntry()
	defer logger.releaseEntry(entry)
	return entry.WitheTime(t)
}

func (logger *Logger)Tracef(format string, args ...interface{}) {
	if logger.Is
}
