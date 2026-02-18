package logentry

import "time"

// Entry represents a single parsed log entry.
type Entry struct {
	Level     Level
	Message   string
	Timestamp time.Time
	Caller    string
	Fields    map[string]any
	Raw       string
	Line      int
	IsJSON    bool
}

// GetField returns the value of a field by key and a boolean indicating existence.
func (e Entry) GetField(key string) (any, bool) {
	if e.Fields == nil {
		return nil, false
	}
	v, ok := e.Fields[key]
	return v, ok
}
