package ligno

import (
	"fmt"
	"strconv"
	"sync"
)

// Level represents rank of message importance.
// Log records can contain level and filters can decide not to process
// records based on this level.
type Level uint

// Levels of log records. Additional can be created, these are just defaults
// offered by library.
const (
	NOTSET   Level = iota
	DEBUG          = iota * 10
	INFO           = iota * 10
	WARNING        = iota * 10
	ERROR          = iota * 10
	CRITICAL       = iota * 10
)

var (
	// level2Name is map from level to name of known level names.
	level2Name = map[Level]string{
		NOTSET:   "NOTSET",
		DEBUG:    "DEBUG",
		INFO:     "INFO",
		WARNING:  "WARNING",
		ERROR:    "ERROR",
		CRITICAL: "CRITICAL",
	}
	// name2level is map from name to level of known level names.
	name2Level = map[string]Level{
		"NOTSET":   NOTSET,
		"DEBUG":    DEBUG,
		"INFO":     INFO,
		"WARNING":  WARNING,
		"ERROR":    ERROR,
		"CRITICAL": CRITICAL,
	}
	mu sync.RWMutex
)

// getLevelName returns name of provided level.
// If provided level does not exist, empty string is returned.
func getLevelName(level Level) (name string) {
	mu.RLock()
	defer mu.RUnlock()
	return level2Name[level]
}

// getLevelFromName returns level from provided name.
// If level with provided name is not found, NOTSET is returned.
func getLevelFromName(name string) (level Level) {
	mu.RLock()
	defer mu.RUnlock()
	return name2Level[name]
}

// AddLevel add new level to system with provided name and rank.
// It is forbidden to register levels that already exist.
func AddLevel(name string, rank uint) (Level, error) {
	mu.Lock()
	defer mu.Unlock()
	l := Level(rank)
	if _, ok := name2Level[name]; ok {
		return NOTSET, fmt.Errorf("level with name '%s' already exists", name)
	}
	if _, ok := level2Name[l]; ok {
		return NOTSET, fmt.Errorf("level with rank '%d' already exists", rank)
	}
	level2Name[l] = name
	name2Level[name] = l
	return l, nil
}

// String returns level's string representation.
func (l Level) String() string {
	if ll, ok := level2Name[l]; ok {
		return ll
	}
	return fmt.Sprintf("Level(%d)", l)
}

// MarshalJSON returns levels JSON representation (implementation of json.Marshaler)
func (l Level) MarshalJSON() ([]byte, error) {
	fmt.Println("In MarshalJSON")
	return []byte(fmt.Sprintf("%q", l.String())), nil
}

// UnmarshalJSON recreates level from JSON representation (implementation of json.Unmarshaler)
func (l *Level) UnmarshalJSON(b []byte) error {
	levelStr, _ := strconv.Unquote(string(b))

	level, ok := name2Level[levelStr]
	if ok {
		fmt.Printf("Found %s in map, value: %d\n", levelStr, level)
	}
	*l = level
	return nil
}
