package theme

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/ersanisk/sieve/pkg/logentry"
)

// ThemeColors holds all color definitions for a theme.
type ThemeColors struct {
	Debug      lipgloss.Color
	Info       lipgloss.Color
	Warn       lipgloss.Color
	Error      lipgloss.Color
	Fatal      lipgloss.Color
	Timestamp  lipgloss.Color
	Key        lipgloss.Color
	Value      lipgloss.Color
	Background lipgloss.Color
	Foreground lipgloss.Color
	StatusBar  lipgloss.Color
	StatusText lipgloss.Color
	Border     lipgloss.Color
	Highlight  lipgloss.Color
}

// Theme defines the styling interface for the application.
type Theme interface {
	Name() string
	Colors() ThemeColors
	LevelStyle(level logentry.Level) lipgloss.Style
	TimestampStyle() lipgloss.Style
	KeyStyle() lipgloss.Style
	ValueStyle() lipgloss.Style
	StatusBarStyle() lipgloss.Style
	BorderStyle() lipgloss.Style
	HighlightStyle() lipgloss.Style
	ErrorStyle() lipgloss.Style
	InfoStyle() lipgloss.Style
}

// BaseTheme provides a default Theme implementation using ThemeColors.
type BaseTheme struct {
	ThemeName   string
	ThemeColors ThemeColors
}

func (t BaseTheme) Name() string {
	return t.ThemeName
}

func (t BaseTheme) Colors() ThemeColors {
	return t.ThemeColors
}

func (t BaseTheme) LevelStyle(level logentry.Level) lipgloss.Style {
	var color lipgloss.Color
	switch level {
	case logentry.Debug:
		color = t.ThemeColors.Debug
	case logentry.Info:
		color = t.ThemeColors.Info
	case logentry.Warn:
		color = t.ThemeColors.Warn
	case logentry.Error:
		color = t.ThemeColors.Error
	case logentry.Fatal:
		color = t.ThemeColors.Fatal
	default:
		color = t.ThemeColors.Foreground
	}
	return lipgloss.NewStyle().Foreground(color).Bold(true)
}

func (t BaseTheme) TimestampStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.ThemeColors.Timestamp)
}

func (t BaseTheme) KeyStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.ThemeColors.Key)
}

func (t BaseTheme) ValueStyle() lipgloss.Style {
	return lipgloss.NewStyle().Foreground(t.ThemeColors.Value)
}

func (t BaseTheme) StatusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.ThemeColors.StatusBar).
		Foreground(t.ThemeColors.StatusText).
		Padding(0, 1)
}

func (t BaseTheme) BorderStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(t.ThemeColors.Border)
}

func (t BaseTheme) HighlightStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Background(t.ThemeColors.Highlight).
		Foreground(t.ThemeColors.Background)
}

func (t BaseTheme) ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.ThemeColors.Error).
		Bold(true)
}

func (t BaseTheme) InfoStyle() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(t.ThemeColors.Info).
		Bold(true)
}
