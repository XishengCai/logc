package pkg

import (
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	nocolor = 0
	red		= 31
	green 	= 32
	yellow	= 33
	blue	= 36
	gray	= 37
)

var (
	baseTimestamp time.Time
	emptyFieldMap FieldMap
)

func init() {
	baseTimestamp = time.Now()
}

// TextFormatter formats logs into text
type TextFormatter struct {
	// Set to true to bypass checking for a TTY before outputting colors.
	ForceColors bool

	// Force disabling colors
	DisableColors bool

	// Override coloring based on CLICOLOR and CLIOLOR_FORCE. - https://bixense.com/clicolors/
	EnvironmentOverrideColors bool

	// Disable timestamp logging. useful when output is redirected to logging
	// system that already adds timestamps.
	DisableTimestamp  bool

	// Enable logging the full timestamp when a TTY is attached of just
	// the time passed since beginning of execution.
	FullTimestamp bool

	// TimestampFormat to use for display when a full timestamp is printed
	TimestampFormat string

	// The fields are sorted by default for a consistent output. For applications
	// that log extremely frequently and don't use the JSON formatter this may not
	// be desired.
	DisableSorting bool

	// The keys sorting function, when uninitialized it uses sort.Strings.
	SortingFunc func([]string)

	// Disables the truncation of the level text to 4 characters.
	DisableLevelTruncation bool

	// QuoteEmptyFields will wrap empty fields quotes if true
	QuoteEmptyFields bool

	// Whether the logger's out is to a terminal
	isTerminal bool

	// FieldMap allows users to customize the names of keys for default fields.
	// As an example:
	// formatter := &TextFormatter {
	//		FieldMap: FieldMap{
	//      	FieldKeyTime: "@timestamp",
	//  		FieldKeyLevel: "@level",
	// 			FieldKeyMsg: 	"@message"}}
	FieldMap 	FieldMap

	terminalInitOnce sync.Once
}

func (f *TextFormatter) init(entry *Entry) {
	if entry.Logger != nil {

	}
}

func (f *TextFormatter) printColored(b *bytes.Buffer, entry *Entry, keys []string, data Fields, timeoutFormat string) {
	var levelColor int
	switch entry.Level {
	case DebugLevel, TraceLevel:
		levelColor = gray
	case WarnLevel:
		levelColor = yellow
	case ErrorLevel, FatalLevel, PanicLevel:
		levelColor = red
	default:
		levelColor = blue
	}

	levelText := strings.ToUpper(entry.Level.String())
	if !f.DisableLevelTruncation {
		levelText = levelText[0:4]
	}

	// Remove a single newline if ti already exists in the message to keep
	// the behavior of logrus text_formatter the same as the stdlib log package
	entry.Message = strings.TrimSuffix(entry.Message, "\n")

	caller := ""

	if entry.HasCaller() {
		caller = fmt.Sprintf("%s:%d %s()",
			entry.Caller.File, entry.Caller.Line, entry.Caller.Function)
	}

	if f.DisableTimestamp {
		fmt.Fprintf(b,  "\x1b[%dm%s\x1b[0m%s %-44s ", levelColor, levelText, caller, entry.Message)
	} else if !f.FullTimestamp {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%04d]%s %-44s ", levelColor, levelText, int(entry.Time.Sub(baseTimestamp)/time.Second), caller, entry.Message)
	} else {
		fmt.Fprintf(b, "\x1b[%dm%s\x1b[0m[%s]%s %-44s ", levelColor, levelText, entry.Time.Format(timestampFormat), caller, entry.Message)
	}

	for _, k := range keys {
		v := data[k]
		fmt.Fprintf(b, " \x1b[%dm%s\x1b[0m=", levelColor, k)
		f.appendValue(b, v)
	}

}

func (f *TextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}){
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *TextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}

func (f *TextFormatter) needQuoting(text string) bool {
	if f.QuoteEmptyFields && len(text) == 0 {
		return true
	}

	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'X') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' ||
			ch == '/' || ch == '@' || ch == '^' ||
			ch == '+') {
				return true
		}
	}
	return false
}