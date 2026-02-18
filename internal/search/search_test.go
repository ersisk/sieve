package search

import (
	"testing"

	"github.com/ersanisk/sieve/pkg/logentry"
)

func TestFuzzyMatch(t *testing.T) {
	entries := []logentry.Entry{
		{
			Level:   logentry.Info,
			Message: "Server starting successfully",
			Caller:  "main.go:42",
			Fields: map[string]any{
				"service": "api",
				"port":    8080,
			},
		},
		{
			Level:   logentry.Error,
			Message: "Connection refused to database",
			Caller:  "db.go:123",
			Fields: map[string]any{
				"service": "db",
				"host":    "localhost:5432",
			},
		},
		{
			Level:   logentry.Warn,
			Message: "Slow query detected",
			Caller:  "query.go:45",
			Fields: map[string]any{
				"service":     "api",
				"duration_ms": 1523.5,
			},
		},
	}

	tests := []struct {
		name  string
		query string
		want  int
	}{
		{"match message exact", "Server starting", 1},
		{"match message partial", "connection", 1},
		{"match caller", "db.go", 1},
		{"match field", "api", 2},
		{"match numeric", "8080", 1},
		{"match float", "1523", 1},
		{"no match", "nonexistent", 0},
		{"empty query", "", 0},
		{"case insensitive", "DATABASE", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := FuzzyMatch(entries, tt.query)
			if len(results) != tt.want {
				t.Errorf("FuzzyMatch() got %d results, want %d", len(results), tt.want)
			}
		})
	}
}

func TestSmartMatch(t *testing.T) {
	entries := []logentry.Entry{
		{
			Level:   logentry.Info,
			Message: "Server starting",
			Caller:  "main.go:42",
			Fields:  map[string]any{"service": "api"},
		},
		{
			Level:   logentry.Error,
			Message: "Connection error",
			Caller:  "db.go:123",
			Fields:  map[string]any{"service": "db"},
		},
	}

	tests := []struct {
		name  string
		query string
		want  int
	}{
		{"match in message", "Server", 1},
		{"match in caller", "db.go", 1},
		{"match in field", "api", 1},
		{"case insensitive", "connection", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results := SmartMatch(entries, tt.query)
			if len(results) != tt.want {
				t.Errorf("SmartMatch() got %d results, want %d", len(results), tt.want)
			}
		})
	}
}

func TestTokenizeQuery(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  []string
	}{
		{
			name:  "single word",
			query: "error",
			want:  []string{"error"},
		},
		{
			name:  "multiple words",
			query: "connection error database",
			want:  []string{"connection", "error", "database"},
		},
		{
			name:  "quoted phrase",
			query: `"connection error" database`,
			want:  []string{"connection error", "database"},
		},
		{
			name:  "multiple quotes",
			query: `"server error" "db timeout"`,
			want:  []string{"server error", "db timeout"},
		},
		{
			name:  "mixed",
			query: `error "connection failed" warning`,
			want:  []string{"error", "connection failed", "warning"},
		},
		{
			name:  "extra spaces",
			query: "  error   warning  ",
			want:  []string{"error", "warning"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TokenizeQuery(tt.query)
			if len(got) != len(tt.want) {
				t.Errorf("TokenizeQuery() got %d tokens, want %d", len(got), len(tt.want))
				return
			}
			for i, token := range got {
				if token != tt.want[i] {
					t.Errorf("TokenizeQuery()[%d] = %v, want %v", i, token, tt.want[i])
				}
			}
		})
	}
}

func TestRegexMatch(t *testing.T) {
	entries := []logentry.Entry{
		{
			Level:   logentry.Info,
			Message: "Server starting on port 8080",
			Caller:  "main.go:42",
			Fields:  map[string]any{"service": "api"},
		},
		{
			Level:   logentry.Error,
			Message: "Connection refused at localhost:5432",
			Caller:  "db.go:123",
			Fields:  map[string]any{"service": "db"},
		},
	}

	tests := []struct {
		name    string
		pattern string
		want    int
		wantErr bool
	}{
		{"match number", `port \d+`, 1, false},
		{"match text", "Connection", 1, false},
		{"match caller", `db\.go:\d+`, 1, false},
		{"no match", "nonexistent", 0, false},
		{"invalid regex", "[", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := RegexMatch(entries, tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegexMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(results) != tt.want {
				t.Errorf("RegexMatch() got %d results, want %d", len(results), tt.want)
			}
		})
	}
}

func TestRegexCaseInsensitiveMatch(t *testing.T) {
	entries := []logentry.Entry{
		{Message: "Server STARTING"},
		{Message: "connection ERROR"},
	}

	results, err := RegexCaseInsensitiveMatch(entries, "starting")
	if err != nil {
		t.Fatalf("RegexCaseInsensitiveMatch() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("RegexCaseInsensitiveMatch() got %d results, want 1", len(results))
	}
}

func TestRegexFieldMatch(t *testing.T) {
	entries := []logentry.Entry{
		{
			Message: "Server starting",
			Caller:  "main.go:42",
			Fields:  map[string]any{"service": "api", "port": 8080},
		},
	}

	tests := []struct {
		name    string
		field   string
		pattern string
		want    int
		wantErr bool
	}{
		{"match message field", "message", "Server", 1, false},
		{"match caller field", "caller", `main\.go`, 1, false},
		{"match custom field", "service", "api", 1, false},
		{"no match", "service", "db", 0, false},
		{"invalid regex", "message", "[", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := RegexFieldMatch(entries, tt.field, tt.pattern)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegexFieldMatch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && len(results) != tt.want {
				t.Errorf("RegexFieldMatch() got %d results, want %d", len(results), tt.want)
			}
		})
	}
}

func TestRegexMultiMatch(t *testing.T) {
	entries := []logentry.Entry{
		{Message: "Server starting error", Fields: map[string]any{"service": "api"}},
		{Message: "Connection error", Fields: map[string]any{"service": "db"}},
		{Message: "Server starting", Fields: map[string]any{"service": "cache"}},
	}

	results, err := RegexMultiMatch(entries, []string{"Server", "error"})
	if err != nil {
		t.Fatalf("RegexMultiMatch() error = %v", err)
	}
	if len(results) != 3 {
		t.Errorf("RegexMultiMatch() got %d results, want 3", len(results))
	}
}

func TestRegexExcludeMatch(t *testing.T) {
	entries := []logentry.Entry{
		{Message: "Server starting"},
		{Message: "Connection error"},
		{Message: "Server stopped"},
	}

	results, err := RegexExcludeMatch(entries, "error")
	if err != nil {
		t.Fatalf("RegexExcludeMatch() error = %v", err)
	}
	if len(results) != 2 {
		t.Errorf("RegexExcludeMatch() got %d results, want 2", len(results))
	}
}

func TestRegexAndMatch(t *testing.T) {
	entries := []logentry.Entry{
		{Message: "Server starting error", Caller: "main.go:42"},
		{Message: "Connection error", Caller: "db.go:123"},
		{Message: "Server starting", Caller: "main.go:42"},
	}

	results, err := RegexAndMatch(entries, []string{"Server", "error"})
	if err != nil {
		t.Fatalf("RegexAndMatch() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("RegexAndMatch() got %d results, want 1", len(results))
	}
}

func TestSearchResult_Sorting(t *testing.T) {
	entries := []logentry.Entry{
		{Message: "exact match", Level: logentry.Info},
		{Message: "partial", Level: logentry.Info},
		{Message: "nomatch", Level: logentry.Info},
	}

	results := FuzzyMatch(entries, "exact")
	if len(results) != 1 {
		t.Fatalf("FuzzyMatch() got %d results, want 1", len(results))
	}

	if results[0].Score <= 0 {
		t.Errorf("FuzzyMatch() Score = %v, want > 0", results[0].Score)
	}
}

func TestValueToString(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  string
		ok    bool
	}{
		{"string", "test", "test", true},
		{"int", 42, "42", true},
		{"float", 3.14, "3.14", true},
		{"float .0", 3.0, "3", true},
		{"bool true", true, "true", true},
		{"bool false", false, "false", true},
		{"nil", nil, "", false},
		{"map", map[string]any{}, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := valueToString(tt.value)
			if ok != tt.ok {
				t.Errorf("valueToString() ok = %v, want %v", ok, tt.ok)
				return
			}
			if ok && got != tt.want {
				t.Errorf("valueToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
