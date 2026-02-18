package logentry

import (
	"strconv"
	"strings"
)

// Level represents the severity level of a log entry.
type Level int

const (
	Unknown Level = iota
	Debug
	Info
	Warn
	Error
	Fatal
)

// String returns the uppercase string representation of a Level.
func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	case Fatal:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// ParseLevel parses a string or numeric value into a Level.
// Supports: full names ("INFO", "info"), single chars ("I", "i"),
// and numeric Bunyan-style levels (10=trace, 20=debug, 30=info, 40=warn, 50=error, 60=fatal).
func ParseLevel(s string) Level {
	// Try numeric first
	if n, err := strconv.Atoi(s); err == nil {
		return parseBunyanLevel(n)
	}

	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "DEBUG", "D", "TRACE", "T":
		return Debug
	case "INFO", "I", "INFORMATION":
		return Info
	case "WARN", "W", "WARNING":
		return Warn
	case "ERROR", "E", "ERR":
		return Error
	case "FATAL", "F", "CRITICAL", "CRIT", "PANIC":
		return Fatal
	default:
		return Unknown
	}
}

// parseBunyanLevel converts a numeric Bunyan-style level to a Level.
func parseBunyanLevel(n int) Level {
	switch {
	case n <= 10:
		return Debug // trace
	case n <= 20:
		return Debug
	case n <= 30:
		return Info
	case n <= 40:
		return Warn
	case n <= 50:
		return Error
	case n <= 60:
		return Fatal
	default:
		return Fatal
	}
}
