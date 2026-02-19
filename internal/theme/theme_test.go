package theme

import (
	"testing"

	"github.com/ersanisk/sieve/pkg/logentry"
)

func TestGetKnownTheme(t *testing.T) {
	tests := []struct {
		name     string
		wantName string
	}{
		{"monokai", "monokai"},
		{"dracula", "dracula"},
		{"gruvbox", "gruvbox"},
		{"nord", "nord"},
		{"kanagawa", "kanagawa"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			theme := Get(tt.name)
			if theme.Name() != tt.wantName {
				t.Errorf("Get(%q).Name() = %q, want %q", tt.name, theme.Name(), tt.wantName)
			}
		})
	}
}

func TestGetUnknownThemeFallback(t *testing.T) {
	theme := Get("nonexistent")
	if theme.Name() != "kanagawa" {
		t.Errorf("Get(nonexistent).Name() = %q, want kanagawa", theme.Name())
	}
}

func TestThemeLevelStyles(t *testing.T) {
	levels := []logentry.Level{
		logentry.Debug,
		logentry.Info,
		logentry.Warn,
		logentry.Error,
		logentry.Fatal,
		logentry.Unknown,
	}

	for _, name := range Names() {
		theme := Get(name)
		t.Run(name, func(t *testing.T) {
			for _, level := range levels {
				style := theme.LevelStyle(level)
				// Verify the style can render without panicking
				_ = style.Render(level.String())
			}
		})
	}
}

func TestThemeStyles(t *testing.T) {
	for _, name := range Names() {
		theme := Get(name)
		t.Run(name, func(t *testing.T) {
			// Verify all style methods work without panicking
			_ = theme.TimestampStyle().Render("2024-01-15T10:00:00Z")
			_ = theme.KeyStyle().Render("key")
			_ = theme.ValueStyle().Render("value")
			_ = theme.StatusBarStyle().Render("status")
			_ = theme.BorderStyle().Render("border")
			_ = theme.HighlightStyle().Render("highlight")
		})
	}
}

func TestNames(t *testing.T) {
	names := Names()
	if len(names) < 5 {
		t.Errorf("Names() returned %d themes, want at least 5", len(names))
	}

	expected := map[string]bool{"monokai": false, "dracula": false, "gruvbox": false, "nord": false, "kanagawa": false}
	for _, name := range names {
		if _, ok := expected[name]; ok {
			expected[name] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("Names() missing theme %q", name)
		}
	}
}
