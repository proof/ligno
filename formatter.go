package ligno

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Formatter is interface for converting log record to string representation.
type Formatter interface {
	Format(record Record) string
}

// DefaultFormatter converts log record to simple string for printing.
type DefaultFormatter struct{}

// defaultTimeFormat is formatting string for time for DefaultFormatter
const defaultTimeFormat = "2006-01-02 15:05:06.0000"

// Format converts provided log record to format suitable for printing in one line.
// String produced resembles traditional log message.
func (df *DefaultFormatter) Format(record Record) string {
	time := record.Time().Format(defaultTimeFormat)
	delete(record, TimeKey)
	level := record.Level()
	delete(record, LevelKey)
	event := record.Event()
	delete(record, EventKey)
	var buff bytes.Buffer

	keys := make([]string, 0, len(record))
	for k := range record {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for i := 0; i < len(keys); i++ {
		k := keys[i]
		v := strconv.Quote(fmt.Sprintf("%+v", record[k]))
		if strings.IndexFunc(k, needsQuote) >= 0 || k == "" {
			k = strconv.Quote(k)
		}
		buff.WriteString(fmt.Sprintf("%s=%+v", k, v))
		if i < len(keys)-1 {
			buff.WriteString(" ")
		}
	}
	return fmt.Sprintf("%-25s %-10s %-15s [%s]", time, level, event, buff.String())
}

// defaultFormatter is instance of DefaultFormatter.
var defaultFormatter = &DefaultFormatter{}

// Needs quote determines if provided rune is such that word that contains this
// rune needs to be quoted.
func needsQuote(r rune) bool {
	return r == ' ' || r == '"' || r == '\\' || r == '=' ||
		!unicode.IsPrint(r)
}

// JSONFormatter is simple formatter that only marshals log record to json.
type JSONFormatter struct{}

// Format returns JSON representation of provided record.
func (jf *JSONFormatter) Format(record Record) string {
	d, err := json.MarshalIndent(record, "", "    ")
	if err != nil {
		panic(err)
	}
	return string(d)
}
