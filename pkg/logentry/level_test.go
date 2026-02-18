package logentry

import "testing"

func TestParseLevel(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Level
	}{
		// Full uppercase names
		{name: "INFO uppercase", input: "INFO", want: Info},
		{name: "DEBUG uppercase", input: "DEBUG", want: Debug},
		{name: "WARN uppercase", input: "WARN", want: Warn},
		{name: "ERROR uppercase", input: "ERROR", want: Error},
		{name: "FATAL uppercase", input: "FATAL", want: Fatal},

		// Full lowercase names
		{name: "info lowercase", input: "info", want: Info},
		{name: "debug lowercase", input: "debug", want: Debug},
		{name: "warn lowercase", input: "warn", want: Warn},
		{name: "error lowercase", input: "error", want: Error},
		{name: "fatal lowercase", input: "fatal", want: Fatal},

		// Single character
		{name: "I single char", input: "I", want: Info},
		{name: "D single char", input: "D", want: Debug},
		{name: "W single char", input: "W", want: Warn},
		{name: "E single char", input: "E", want: Error},
		{name: "F single char", input: "F", want: Fatal},

		// Single char lowercase
		{name: "i lowercase char", input: "i", want: Info},
		{name: "d lowercase char", input: "d", want: Debug},

		// Aliases
		{name: "WARNING alias", input: "WARNING", want: Warn},
		{name: "ERR alias", input: "ERR", want: Error},
		{name: "CRITICAL alias", input: "CRITICAL", want: Fatal},
		{name: "CRIT alias", input: "CRIT", want: Fatal},
		{name: "PANIC alias", input: "PANIC", want: Fatal},
		{name: "TRACE alias", input: "TRACE", want: Debug},
		{name: "INFORMATION alias", input: "INFORMATION", want: Info},

		// Numeric Bunyan-style
		{name: "Bunyan 10 trace", input: "10", want: Debug},
		{name: "Bunyan 20 debug", input: "20", want: Debug},
		{name: "Bunyan 30 info", input: "30", want: Info},
		{name: "Bunyan 40 warn", input: "40", want: Warn},
		{name: "Bunyan 50 error", input: "50", want: Error},
		{name: "Bunyan 60 fatal", input: "60", want: Fatal},

		// Whitespace
		{name: "with leading space", input: "  INFO", want: Info},
		{name: "with trailing space", input: "INFO  ", want: Info},

		// Unknown
		{name: "empty string", input: "", want: Unknown},
		{name: "garbage", input: "xyz", want: Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ParseLevel(tt.input)
			if got != tt.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{Unknown, "UNKNOWN"},
		{Debug, "DEBUG"},
		{Info, "INFO"},
		{Warn, "WARN"},
		{Error, "ERROR"},
		{Fatal, "FATAL"},
		{Level(99), "UNKNOWN"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.level.String()
			if got != tt.want {
				t.Errorf("Level(%d).String() = %q, want %q", tt.level, got, tt.want)
			}
		})
	}
}
