package logentry

import (
	"testing"
	"time"
)

func TestEntryGetField(t *testing.T) {
	tests := []struct {
		name   string
		entry  Entry
		key    string
		want   any
		wantOK bool
	}{
		{
			name: "existing field",
			entry: Entry{
				Fields: map[string]any{"service": "auth", "status": 200},
			},
			key:    "service",
			want:   "auth",
			wantOK: true,
		},
		{
			name: "non-existing field",
			entry: Entry{
				Fields: map[string]any{"service": "auth"},
			},
			key:    "missing",
			want:   nil,
			wantOK: false,
		},
		{
			name:   "nil fields map",
			entry:  Entry{},
			key:    "anything",
			want:   nil,
			wantOK: false,
		},
		{
			name: "numeric field",
			entry: Entry{
				Fields: map[string]any{"status": float64(404)},
			},
			key:    "status",
			want:   float64(404),
			wantOK: true,
		},
		{
			name: "nil value field",
			entry: Entry{
				Fields: map[string]any{"data": nil},
			},
			key:    "data",
			want:   nil,
			wantOK: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := tt.entry.GetField(tt.key)
			if ok != tt.wantOK {
				t.Errorf("GetField(%q) ok = %v, want %v", tt.key, ok, tt.wantOK)
			}
			if got != tt.want {
				t.Errorf("GetField(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestEntryZeroValue(t *testing.T) {
	var e Entry
	if e.Level != Unknown {
		t.Errorf("zero Entry.Level = %v, want Unknown", e.Level)
	}
	if e.Message != "" {
		t.Errorf("zero Entry.Message = %q, want empty", e.Message)
	}
	if !e.Timestamp.IsZero() {
		t.Errorf("zero Entry.Timestamp = %v, want zero", e.Timestamp)
	}
	if e.IsJSON {
		t.Error("zero Entry.IsJSON = true, want false")
	}
	if e.Fields != nil {
		t.Errorf("zero Entry.Fields = %v, want nil", e.Fields)
	}
}

func TestEntryFields(t *testing.T) {
	ts := time.Date(2024, 1, 15, 10, 30, 0, 0, time.UTC)
	e := Entry{
		Level:     Error,
		Message:   "connection refused",
		Timestamp: ts,
		Caller:    "server.go:42",
		Fields:    map[string]any{"host": "localhost", "port": float64(8080)},
		Raw:       `{"level":"error","msg":"connection refused"}`,
		Line:      10,
		IsJSON:    true,
	}

	if e.Level != Error {
		t.Errorf("Level = %v, want Error", e.Level)
	}
	if e.Message != "connection refused" {
		t.Errorf("Message = %q, want %q", e.Message, "connection refused")
	}
	if !e.Timestamp.Equal(ts) {
		t.Errorf("Timestamp = %v, want %v", e.Timestamp, ts)
	}
	if e.Caller != "server.go:42" {
		t.Errorf("Caller = %q, want %q", e.Caller, "server.go:42")
	}
	if e.Line != 10 {
		t.Errorf("Line = %d, want 10", e.Line)
	}
	if !e.IsJSON {
		t.Error("IsJSON = false, want true")
	}

	v, ok := e.GetField("host")
	if !ok || v != "localhost" {
		t.Errorf("GetField(host) = %v, %v; want localhost, true", v, ok)
	}
}
