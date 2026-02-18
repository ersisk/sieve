package filter

import (
	"testing"

	"github.com/ersanisk/sieve/pkg/logentry"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "simple equality",
			input:   ".level == 30",
			wantErr: false,
		},
		{
			name:    "inequality",
			input:   ".level != 40",
			wantErr: false,
		},
		{
			name:    "greater than",
			input:   ".level > 30",
			wantErr: false,
		},
		{
			name:    "less than",
			input:   ".level < 50",
			wantErr: false,
		},
		{
			name:    "greater equal",
			input:   ".level >= 30",
			wantErr: false,
		},
		{
			name:    "less equal",
			input:   ".level <= 50",
			wantErr: false,
		},
		{
			name:    "string equality",
			input:   `.service == "api"`,
			wantErr: false,
		},
		{
			name:    "single quoted string",
			input:   `.service == 'api'`,
			wantErr: false,
		},
		{
			name:    "contains operator",
			input:   `.message contains "error"`,
			wantErr: false,
		},
		{
			name:    "matches operator",
			input:   `.message matches ".*error.*"`,
			wantErr: false,
		},
		{
			name:    "and operator",
			input:   `.level == 30 and .service == "api"`,
			wantErr: false,
		},
		{
			name:    "or operator",
			input:   `.level == 30 or .level == 40`,
			wantErr: false,
		},
		{
			name:    "not operator",
			input:   "not .level == 50",
			wantErr: false,
		},
		{
			name:    "complex expression",
			input:   `.level >= 30 and .service == "api"`,
			wantErr: false,
		},
		{
			name:    "numeric literal",
			input:   `.port == 8080`,
			wantErr: false,
		},
		{
			name:    "float literal",
			input:   `.duration_ms == 45.5`,
			wantErr: false,
		},
		{
			name:    "boolean literal",
			input:   `.success == true`,
			wantErr: false,
		},

		{
			name:    "empty expression",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid operator",
			input:   ".level === 30",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && expr == nil {
				t.Errorf("Parse() returned nil expression for valid input")
			}
		})
	}
}

func TestCompiledFilter_Evaluate(t *testing.T) {
	tests := []struct {
		name    string
		expr    string
		entry   logentry.Entry
		want    bool
		wantErr bool
	}{
		{
			name:  "level equality - info",
			expr:  ".level == 30",
			entry: logentry.Entry{Level: logentry.Info, Fields: map[string]any{"level": 30}},
			want:  true,
		},
		{
			name:  "level equality - not info",
			expr:  ".level == 30",
			entry: logentry.Entry{Level: logentry.Error, Fields: map[string]any{"level": 50}},
			want:  false,
		},
		{
			name:  "level greater equal",
			expr:  ".level >= 40",
			entry: logentry.Entry{Level: logentry.Warn, Fields: map[string]any{"level": 40}},
			want:  true,
		},
		{
			name:  "level greater equal - below",
			expr:  ".level >= 40",
			entry: logentry.Entry{Level: logentry.Info, Fields: map[string]any{"level": 30}},
			want:  false,
		},
		{
			name:  "string equality - match",
			expr:  `.service == "api"`,
			entry: logentry.Entry{Fields: map[string]any{"service": "api"}},
			want:  true,
		},
		{
			name:  "string equality - no match",
			expr:  `.service == "auth"`,
			entry: logentry.Entry{Fields: map[string]any{"service": "api"}},
			want:  false,
		},
		{
			name:  "contains - match",
			expr:  `.message contains "error"`,
			entry: logentry.Entry{Message: "connection error occurred"},
			want:  true,
		},
		{
			name:  "contains - no match",
			expr:  `.message contains "warning"`,
			entry: logentry.Entry{Message: "connection error occurred"},
			want:  false,
		},
		{
			name:  "and - both true",
			expr:  `.level == 30 and .service == "api"`,
			entry: logentry.Entry{Level: logentry.Info, Fields: map[string]any{"level": 30, "service": "api"}},
			want:  true,
		},
		{
			name:  "and - one false",
			expr:  `.level == 30 and .service == "auth"`,
			entry: logentry.Entry{Level: logentry.Info, Fields: map[string]any{"level": 30, "service": "api"}},
			want:  false,
		},
		{
			name:  "or - one true",
			expr:  `.level == 30 or .level == 40`,
			entry: logentry.Entry{Level: logentry.Info, Fields: map[string]any{"level": 30}},
			want:  true,
		},
		{
			name:  "or - both false",
			expr:  `.level == 30 or .level == 40`,
			entry: logentry.Entry{Level: logentry.Debug, Fields: map[string]any{"level": 20}},
			want:  false,
		},
		{
			name:  "not - true",
			expr:  "not .level == 50",
			entry: logentry.Entry{Level: logentry.Info, Fields: map[string]any{"level": 30}},
			want:  true,
		},
		{
			name:  "not - false",
			expr:  "not .level == 50",
			entry: logentry.Entry{Level: logentry.Error, Fields: map[string]any{"level": 50}},
			want:  false,
		},
		{
			name:  "numeric field",
			expr:  `.port == 8080`,
			entry: logentry.Entry{Fields: map[string]any{"port": 8080}},
			want:  true,
		},
		{
			name:  "float field",
			expr:  `.duration_ms == 45.5`,
			entry: logentry.Entry{Fields: map[string]any{"duration_ms": 45.5}},
			want:  true,
		},
		{
			name:  "boolean field",
			expr:  `.success == true`,
			entry: logentry.Entry{Fields: map[string]any{"success": true}},
			want:  true,
		},
		{
			name:  "missing field",
			expr:  `.nonexistent == "value"`,
			entry: logentry.Entry{},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expr, err := Parse(tt.expr)
			if err != nil {
				t.Fatalf("Parse() failed: %v", err)
			}

			filter, err := Compile(expr)
			if err != nil {
				t.Fatalf("Compile() failed: %v", err)
			}

			got, err := filter.Evaluate(tt.entry)
			if (err != nil) != tt.wantErr {
				t.Errorf("Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestByLevel(t *testing.T) {
	tests := []struct {
		name     string
		minLevel logentry.Level
		entries  []logentry.Entry
		want     []bool
	}{
		{
			name:     "error level filter",
			minLevel: logentry.Error,
			entries: []logentry.Entry{
				{Level: logentry.Info, Fields: map[string]any{"level": 30}},
				{Level: logentry.Warn, Fields: map[string]any{"level": 40}},
				{Level: logentry.Error, Fields: map[string]any{"level": 50}},
				{Level: logentry.Fatal, Fields: map[string]any{"level": 60}},
			},
			want: []bool{false, false, true, true},
		},
		{
			name:     "warn level filter",
			minLevel: logentry.Warn,
			entries: []logentry.Entry{
				{Level: logentry.Info, Fields: map[string]any{"level": 30}},
				{Level: logentry.Warn, Fields: map[string]any{"level": 40}},
				{Level: logentry.Error, Fields: map[string]any{"level": 50}},
			},
			want: []bool{false, true, true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := ByLevel(tt.minLevel)

			for i, entry := range tt.entries {
				got, err := filter.Evaluate(entry)
				if err != nil {
					t.Errorf("Evaluate() error = %v", err)
					continue
				}
				if got != tt.want[i] {
					t.Errorf("Evaluate()[%d] = %v, want %v", i, got, tt.want[i])
				}
			}
		})
	}
}

func TestByValue(t *testing.T) {
	tests := []struct {
		name  string
		field string
		value any
		op    Operator
		entry logentry.Entry
		want  bool
	}{
		{
			name:  "equality match",
			field: "service",
			value: "api",
			op:    OpEqual,
			entry: logentry.Entry{Fields: map[string]any{"service": "api"}},
			want:  true,
		},
		{
			name:  "equality no match",
			field: "service",
			value: "auth",
			op:    OpEqual,
			entry: logentry.Entry{Fields: map[string]any{"service": "api"}},
			want:  false,
		},
		{
			name:  "greater than match",
			field: "port",
			value: 80,
			op:    OpGreater,
			entry: logentry.Entry{Fields: map[string]any{"port": 8080}},
			want:  true,
		},
		{
			name:  "greater than no match",
			field: "port",
			value: 9000,
			op:    OpGreater,
			entry: logentry.Entry{Fields: map[string]any{"port": 8080}},
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := ByValue(tt.field, tt.value, tt.op)
			got, err := filter.Evaluate(tt.entry)
			if err != nil {
				t.Errorf("Evaluate() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreset(t *testing.T) {
	tests := []struct {
		name    string
		preset  string
		wantErr bool
	}{
		{"errors preset", "errors", false},
		{"errors-and-warnings preset", "errors-and-warnings", false},
		{"debug preset", "debug", false},
		{"non-existent preset", "nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			preset, ok := GetPreset(tt.preset)
			if tt.wantErr {
				if ok {
					t.Errorf("GetPreset() returned true for non-existent preset")
				}
				return
			}
			if !ok {
				t.Errorf("GetPreset() returned false for existing preset")
				return
			}
			if preset.Name != tt.preset {
				t.Errorf("GetPreset() Name = %v, want %v", preset.Name, tt.preset)
			}
		})
	}
}

func TestNewPresetFilter(t *testing.T) {
	tests := []struct {
		name    string
		preset  string
		wantErr bool
	}{
		{"errors preset", "errors", false},
		{"errors-and-warnings preset", "errors-and-warnings", false},
		{"non-existent preset", "nonexistent", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter, err := NewPresetFilter(tt.preset)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPresetFilter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && filter != nil {
				t.Errorf("NewPresetFilter() returned filter for error case")
			}
		})
	}
}

func TestOperator_String(t *testing.T) {
	tests := []struct {
		op  Operator
		str string
	}{
		{OpEqual, "=="},
		{OpNotEqual, "!="},
		{OpGreater, ">"},
		{OpLess, "<"},
		{OpGreaterEqual, ">="},
		{OpLessEqual, "<="},
		{OpContains, "contains"},
		{OpMatches, "matches"},
		{OpAnd, "and"},
		{OpOr, "or"},
		{OpNot, "not"},
	}

	for _, tt := range tests {
		t.Run(tt.str, func(t *testing.T) {
			if tt.op.String() != tt.str {
				t.Errorf("Operator.String() = %v, want %v", tt.op.String(), tt.str)
			}
		})
	}
}
