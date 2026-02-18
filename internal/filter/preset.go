package filter

import (
	"fmt"

	"github.com/ersanisk/sieve/pkg/logentry"
)

// Preset represents a predefined filter configuration.
type Preset struct {
	Name        string
	Description string
	Expression  string
}

// Presets is a collection of common filter presets.
var Presets = []Preset{
	{
		Name:        "errors",
		Description: "Show only error and fatal logs",
		Expression:  ".level >= 50",
	},
	{
		Name:        "errors-and-warnings",
		Description: "Show errors and warnings",
		Expression:  ".level >= 40",
	},
	{
		Name:        "debug",
		Description: "Show all logs including debug",
		Expression:  ".level >= 10",
	},
	{
		Name:        "production",
		Description: "Show info, warn, error, fatal",
		Expression:  ".level >= 30",
	},
}

// GetPreset retrieves a preset by name.
func GetPreset(name string) (*Preset, bool) {
	for _, p := range Presets {
		if p.Name == name {
			return &p, true
		}
	}
	return nil, false
}

// PresetFilter represents a compiled preset filter.
type PresetFilter struct {
	compiled *CompiledFilter
}

// NewPresetFilter creates a new filter from a preset name.
func NewPresetFilter(name string) (*PresetFilter, error) {
	preset, ok := GetPreset(name)
	if !ok {
		return nil, fmt.Errorf("preset not found: %s", name)
	}

	expr, err := Parse(preset.Expression)
	if err != nil {
		return nil, fmt.Errorf("failed to parse preset expression: %w", err)
	}

	compiled, err := Compile(expr)
	if err != nil {
		return nil, fmt.Errorf("failed to compile preset: %w", err)
	}

	return &PresetFilter{compiled: compiled}, nil
}

// Evaluate evaluates the preset filter against an entry.
func (pf *PresetFilter) Evaluate(entry logentry.Entry) (bool, error) {
	return pf.compiled.Evaluate(entry)
}

// LevelToPreset returns the preset name for a given minimum level.
func LevelToPreset(level logentry.Level) string {
	switch {
	case level >= logentry.Fatal:
		return "errors"
	case level >= logentry.Error:
		return "errors"
	case level >= logentry.Warn:
		return "errors-and-warnings"
	case level >= logentry.Info:
		return "production"
	case level >= logentry.Debug:
		return "debug"
	default:
		return "production"
	}
}
