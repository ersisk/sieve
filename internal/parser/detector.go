package parser

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

// Format represents the detected log format type.
type Format int

const (
	FormatUnknown Format = iota
	FormatJSON
	FormatJSONLines
	FormatMixed
	FormatPlain
)

func (f Format) String() string {
	switch f {
	case FormatJSON:
		return "JSON"
	case FormatJSONLines:
		return "JSONLines"
	case FormatMixed:
		return "Mixed"
	case FormatPlain:
		return "Plain"
	default:
		return "Unknown"
	}
}

const sampleSize = 100

// DetectFormat analyzes the input and returns the detected format type.
func DetectFormat(r io.Reader) Format {
	var lines []string
	scanner := bufio.NewScanner(r)

	for i := 0; i < sampleSize && scanner.Scan(); i++ {
		line := scanner.Text()
		if line != "" {
			lines = append(lines, line)
		}
	}

	if len(lines) == 0 {
		return FormatUnknown
	}

	jsonCount := 0
	for _, line := range lines {
		if isValidJSON(line) {
			jsonCount++
		}
	}

	jsonRatio := float64(jsonCount) / float64(len(lines))

	switch {
	case jsonRatio == 1.0:
		if len(lines) == 1 {
			return FormatJSON
		}
		return FormatJSONLines
	case jsonRatio > 0.5:
		return FormatMixed
	case jsonRatio > 0:
		return FormatMixed
	default:
		return FormatPlain
	}
}

// isValidJSON checks if a string is valid JSON.
func isValidJSON(s string) bool {
	trimmed := strings.TrimSpace(s)
	var js json.RawMessage
	return json.Unmarshal([]byte(trimmed), &js) == nil
}
