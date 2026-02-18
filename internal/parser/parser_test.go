package parser

import (
	"strings"
	"testing"

	"github.com/ersanisk/sieve/pkg/logentry"
)

func TestParseLine_ValidJSON(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name      string
		input     string
		lineNum   int
		wantLevel logentry.Level
		wantMsg   string
		wantJSON  bool
	}{
		{
			name:      "basic info log",
			input:     `{"level":"info","msg":"hello world"}`,
			lineNum:   1,
			wantLevel: logentry.Info,
			wantMsg:   "hello world",
			wantJSON:  true,
		},
		{
			name:      "error log with timestamp",
			input:     `{"level":"error","msg":"failed","ts":"2024-01-15T10:00:00Z"}`,
			lineNum:   5,
			wantLevel: logentry.Error,
			wantMsg:   "failed",
			wantJSON:  true,
		},
		{
			name:      "warn log with nested fields",
			input:     `{"level":"warn","msg":"slow query","duration_ms":1523}`,
			lineNum:   10,
			wantLevel: logentry.Warn,
			wantMsg:   "slow query",
			wantJSON:  true,
		},
		{
			name:      "debug log",
			input:     `{"level":"debug","msg":"debugging","user_id":123}`,
			lineNum:   15,
			wantLevel: logentry.Debug,
			wantMsg:   "debugging",
			wantJSON:  true,
		},
		{
			name:      "fatal log",
			input:     `{"level":"fatal","msg":"crash","code":500}`,
			lineNum:   20,
			wantLevel: logentry.Fatal,
			wantMsg:   "crash",
			wantJSON:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.ParseLine(tt.input, tt.lineNum)
			if got.Level != tt.wantLevel {
				t.Errorf("ParseLine() Level = %v, want %v", got.Level, tt.wantLevel)
			}
			if got.Message != tt.wantMsg {
				t.Errorf("ParseLine() Message = %v, want %v", got.Message, tt.wantMsg)
			}
			if got.IsJSON != tt.wantJSON {
				t.Errorf("ParseLine() IsJSON = %v, want %v", got.IsJSON, tt.wantJSON)
			}
			if got.Line != tt.lineNum {
				t.Errorf("ParseLine() Line = %v, want %v", got.Line, tt.lineNum)
			}
		})
	}
}

func TestParseLine_LevelVariations(t *testing.T) {
	p := NewParser()

	tests := []struct {
		input     string
		wantLevel logentry.Level
	}{
		{`{"level":"INFO"}`, logentry.Info},
		{`{"level":"info"}`, logentry.Info},
		{`{"lvl":"INFO"}`, logentry.Info},
		{`{"severity":"WARN"}`, logentry.Warn},
		{`{"priority":"ERROR"}`, logentry.Error},
		{`{"level":30}`, logentry.Info},
		{`{"level":50}`, logentry.Error},
		{`{"level":"D"}`, logentry.Debug},
		{`{"level":"I"}`, logentry.Info},
		{`{"level":"E"}`, logentry.Error},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := p.ParseLine(tt.input, 1)
			if got.Level != tt.wantLevel {
				t.Errorf("ParseLine() Level = %v, want %v", got.Level, tt.wantLevel)
			}
		})
	}
}

func TestParseLine_MessageVariations(t *testing.T) {
	p := NewParser()

	tests := []struct {
		input   string
		wantMsg string
	}{
		{`{"msg":"hello"}`, "hello"},
		{`{"message":"world"}`, "world"},
		{`{"text":"test"}`, "test"},
		{`{"msg":"multi\nline"}`, "multi\nline"},
		{`{"msg":123}`, "123.00"},
		{`{"msg":45.67}`, "45.67"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := p.ParseLine(tt.input, 1)
			if got.Message != tt.wantMsg {
				t.Errorf("ParseLine() Message = %v, want %v", got.Message, tt.wantMsg)
			}
		})
	}
}

func TestParseLine_TimestampVariations(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name      string
		input     string
		wantValid bool
	}{
		{
			name:      "RFC3339",
			input:     `{"time":"2024-01-15T10:00:00Z"}`,
			wantValid: true,
		},
		{
			name:      "RFC3339Nano",
			input:     `{"timestamp":"2024-01-15T10:00:00.123456789Z"}`,
			wantValid: true,
		},
		{
			name:      "standard format",
			input:     `{"ts":"2024-01-15 10:00:00"}`,
			wantValid: true,
		},
		{
			name:      "Unix timestamp float",
			input:     `{"time":1705298400.123}`,
			wantValid: true,
		},
		{
			name:      "Unix timestamp int",
			input:     `{"ts":1705298400}`,
			wantValid: true,
		},
		{
			name:      "invalid format",
			input:     `{"time":"not a date"}`,
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.ParseLine(tt.input, 1)
			isValid := !got.Timestamp.IsZero()
			if isValid != tt.wantValid {
				t.Errorf("ParseLine() Timestamp valid = %v, want %v", isValid, tt.wantValid)
			}
		})
	}
}

func TestParseLine_CallerVariations(t *testing.T) {
	p := NewParser()

	tests := []struct {
		input      string
		wantCaller string
	}{
		{`{"caller":"main.go:42"}`, "main.go:42"},
		{`{"source":"app/handler.go:123"}`, "app/handler.go:123"},
		{`{"file":"/app/main.go:42"}`, "/app/main.go:42"},
		{`{"location":"service:80"}`, "service:80"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := p.ParseLine(tt.input, 1)
			if got.Caller != tt.wantCaller {
				t.Errorf("ParseLine() Caller = %v, want %v", got.Caller, tt.wantCaller)
			}
		})
	}
}

func TestParseLine_MalformedJSON(t *testing.T) {
	p := NewParser()

	tests := []struct {
		name  string
		input string
	}{
		{"unclosed brace", `{invalid`},
		{"missing quotes", `{level:info}`},
		{"trailing comma", `{"level":"info",}`},
		{"random text", `just some text`},
		{"empty string", ``},
		{"only spaces", `   `},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := p.ParseLine(tt.input, 1)
			if got.IsJSON {
				t.Errorf("ParseLine() IsJSON = true, want false for malformed input")
			}
			if got.Level != logentry.Unknown {
				t.Errorf("ParseLine() Level = %v, want Unknown for malformed input", got.Level)
			}
		})
	}
}

func TestParseLine_NestedJSON(t *testing.T) {
	p := NewParser()

	input := `{"level":"info","msg":"nested","data":{"users":[{"id":1,"name":"Alice"},{"id":2,"name":"Bob"}],"total":2}}`
	got := p.ParseLine(input, 1)

	if !got.IsJSON {
		t.Errorf("ParseLine() IsJSON = false, want true")
	}

	data, ok := got.Fields["data"]
	if !ok {
		t.Fatal("ParseLine() missing 'data' field")
	}

	dataMap, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("ParseLine() data field is not map[string]any, got %T", data)
	}

	total, ok := dataMap["total"].(float64)
	if !ok || total != 2 {
		t.Errorf("ParseLine() data.total = %v, want 2", total)
	}
}

func TestParseLine_FieldAccess(t *testing.T) {
	p := NewParser()

	input := `{"level":"info","msg":"test","service":"api","request_id":"req-123","duration_ms":45}`
	got := p.ParseLine(input, 1)

	service, ok := got.GetField("service")
	if !ok {
		t.Error("ParseLine() missing 'service' field")
	}
	if service != "api" {
		t.Errorf("ParseLine() service = %v, want 'api'", service)
	}

	requestID, ok := got.GetField("request_id")
	if !ok {
		t.Error("ParseLine() missing 'request_id' field")
	}
	if requestID != "req-123" {
		t.Errorf("ParseLine() request_id = %v, want 'req-123'", requestID)
	}

	duration, ok := got.GetField("duration_ms")
	if !ok {
		t.Error("ParseLine() missing 'duration_ms' field")
	}
	if duration != float64(45) {
		t.Errorf("ParseLine() duration_ms = %v, want 45", duration)
	}

	_, ok = got.GetField("nonexistent")
	if ok {
		t.Error("ParseLine() GetField returned true for nonexistent key")
	}
}

func TestParseLines(t *testing.T) {
	p := NewParser()

	input := `{"level":"info","msg":"line1"}
{"level":"debug","msg":"line2"}
{"level":"error","msg":"line3"}`

	r := strings.NewReader(input)
	entries, err := p.ParseLines(r)

	if err != nil {
		t.Fatalf("ParseLines() error = %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("ParseLines() got %d entries, want 3", len(entries))
	}

	if entries[0].Level != logentry.Info || entries[0].Message != "line1" {
		t.Errorf("ParseLines()[0] = %+v, want Level=Info, Message=line1", entries[0])
	}
	if entries[1].Level != logentry.Debug || entries[1].Message != "line2" {
		t.Errorf("ParseLines()[1] = %+v, want Level=Debug, Message=line2", entries[1])
	}
	if entries[2].Level != logentry.Error || entries[2].Message != "line3" {
		t.Errorf("ParseLines()[2] = %+v, want Level=Error, Message=line3", entries[2])
	}
}

func TestParseLines_EmptyInput(t *testing.T) {
	p := NewParser()

	r := strings.NewReader("")
	entries, err := p.ParseLines(r)

	if err != nil {
		t.Fatalf("ParseLines() error = %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("ParseLines() got %d entries, want 0", len(entries))
	}
}

func TestParseLines_MixedFormat(t *testing.T) {
	p := NewParser()

	input := `{"level":"info","msg":"json line"}
plain text line
{"level":"warn","msg":"another json"}`

	r := strings.NewReader(input)
	entries, err := p.ParseLines(r)

	if err != nil {
		t.Fatalf("ParseLines() error = %v", err)
	}

	if len(entries) != 3 {
		t.Fatalf("ParseLines() got %d entries, want 3", len(entries))
	}

	if !entries[0].IsJSON {
		t.Errorf("ParseLines()[0] IsJSON = false, want true")
	}
	if entries[1].IsJSON {
		t.Errorf("ParseLines()[1] IsJSON = true, want false")
	}
	if !entries[2].IsJSON {
		t.Errorf("ParseLines()[2] IsJSON = false, want true")
	}
}
