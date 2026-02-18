package parser

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/ersanisk/sieve/pkg/logentry"
)

// Parser parses log lines into logentry.Entry objects.
type Parser struct{}

// NewParser creates a new Parser instance.
func NewParser() *Parser {
	return &Parser{}
}

// ParseLine parses a single log line into an Entry.
func (p *Parser) ParseLine(raw string, lineNum int) logentry.Entry {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return logentry.Entry{
			Raw:  raw,
			Line: lineNum,
		}
	}

	var fields map[string]any
	err := json.Unmarshal([]byte(trimmed), &fields)
	if err != nil {
		return logentry.Entry{
			Level:   logentry.Unknown,
			Message: raw,
			Raw:     raw,
			Line:    lineNum,
			IsJSON:  false,
		}
	}

	entry := logentry.Entry{
		Level:     p.parseLevel(fields),
		Message:   p.parseMessage(fields),
		Timestamp: p.parseTimestamp(fields),
		Caller:    p.parseCaller(fields),
		Fields:    fields,
		Raw:       raw,
		Line:      lineNum,
		IsJSON:    true,
	}

	return entry
}

// ParseLines reads from a reader and parses all lines into Entries.
func (p *Parser) ParseLines(r io.Reader) ([]logentry.Entry, error) {
	var entries []logentry.Entry

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read: %w", err)
	}

	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		if i == len(lines)-1 && line == "" {
			continue
		}
		entries = append(entries, p.ParseLine(line, i+1))
	}

	return entries, nil
}

// parseLevel extracts the log level from fields.
// Checks common field names: level, lvl, severity, priority.
func (p *Parser) parseLevel(fields map[string]any) logentry.Level {
	levelKeys := []string{"level", "lvl", "severity", "priority"}

	for _, key := range levelKeys {
		if v, ok := fields[key]; ok {
			switch val := v.(type) {
			case string:
				return logentry.ParseLevel(val)
			case float64:
				return logentry.ParseLevel(fmt.Sprintf("%.0f", val))
			case int:
				return logentry.ParseLevel(fmt.Sprintf("%d", val))
			}
		}
	}

	return logentry.Unknown
}

// parseMessage extracts the message from fields.
// Checks common field names: msg, message, text.
func (p *Parser) parseMessage(fields map[string]any) string {
	msgKeys := []string{"msg", "message", "text"}

	for _, key := range msgKeys {
		if v, ok := fields[key]; ok {
			switch val := v.(type) {
			case string:
				return val
			case float64:
				return fmt.Sprintf("%.2f", val)
			case int:
				return fmt.Sprintf("%d", val)
			}
		}
	}

	return ""
}

// parseTimestamp extracts and parses the timestamp from fields.
// Checks common field names: time, timestamp, ts, @timestamp.
// Supports RFC3339 and Unix timestamps.
func (p *Parser) parseTimestamp(fields map[string]any) time.Time {
	timeKeys := []string{"time", "timestamp", "ts", "@timestamp"}

	for _, key := range timeKeys {
		if v, ok := fields[key]; ok {
			switch val := v.(type) {
			case string:
				if t, err := time.Parse(time.RFC3339, val); err == nil {
					return t
				}
				if t, err := time.Parse(time.RFC3339Nano, val); err == nil {
					return t
				}
				if t, err := time.Parse("2006-01-02 15:04:05", val); err == nil {
					return t
				}
				if t, err := time.Parse("2006-01-02T15:04:05.999999999Z", val); err == nil {
					return t
				}
			case float64:
				if val > 1e9 && val < 2e9 {
					return time.Unix(int64(val), 0)
				}
				return time.Unix(0, int64(val*1e9))
			case int:
				if val > 1e9 && val < 2e9 {
					return time.Unix(int64(val), 0)
				}
				return time.Unix(0, int64(val)*1e9)
			}
		}
	}

	return time.Time{}
}

// parseCaller extracts the caller/location from fields.
// Checks common field names: caller, source, file, location.
func (p *Parser) parseCaller(fields map[string]any) string {
	callerKeys := []string{"caller", "source", "file", "location"}

	for _, key := range callerKeys {
		if v, ok := fields[key]; ok {
			if val, ok := v.(string); ok {
				return val
			}
		}
	}

	return ""
}
