package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
)

// FilterBar handles filter input.
type FilterBar struct {
	textInput textinput.Model
	visible   bool
	focused   bool
	width     int
	height    int
	theme     theme.Theme
}

// NewFilterBar creates a new FilterBar.
func NewFilterBar(theme theme.Theme) FilterBar {
	ti := textinput.New()
	ti.Placeholder = "Filter expression (e.g., .level >= 30)"
	ti.Prompt = "F:"
	ti.PromptStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	ti.TextStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(theme.Colors().Highlight)
	ti.Cursor.TextStyle = lipgloss.NewStyle().Foreground(theme.Colors().Background)

	return FilterBar{
		textInput: ti,
		visible:   false,
		focused:   false,
		theme:     theme,
	}
}

// Show shows the filter bar.
func (m *FilterBar) Show() {
	m.visible = true
	m.focused = true
	m.textInput.Focus()
}

// Hide hides the filter bar.
func (m *FilterBar) Hide() {
	m.visible = false
	m.focused = false
	m.textInput.Blur()
}

// IsVisible returns true if the filter bar is visible.
func (m *FilterBar) IsVisible() bool {
	return m.visible
}

// IsFocused returns true if the filter bar is focused.
func (m *FilterBar) IsFocused() bool {
	return m.focused
}

// SetValue sets the filter value.
func (m *FilterBar) SetValue(value string) {
	m.textInput.SetValue(value)
}

// GetValue returns the current filter value.
func (m *FilterBar) GetValue() string {
	return m.textInput.Value()
}

// Clear clears the filter input.
func (m *FilterBar) Clear() {
	m.textInput.Reset()
}

// SetSize sets the dimensions of the filter bar.
func (m *FilterBar) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.textInput.Width = width - 20
}

// SetTheme sets the theme.
func (m *FilterBar) SetTheme(theme theme.Theme) {
	m.theme = theme
	m.textInput.PromptStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	m.textInput.TextStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	m.textInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	m.textInput.Cursor.Style = lipgloss.NewStyle().Foreground(theme.Colors().Highlight)
	m.textInput.Cursor.TextStyle = lipgloss.NewStyle().Foreground(theme.Colors().Background)
}

// Update updates the filter bar.
func (m FilterBar) Update(msg tea.Msg) (FilterBar, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the filter bar.
func (m FilterBar) View() string {
	if !m.visible {
		return ""
	}
	return m.textInput.View()
}

// Focus sets focus to the filter bar.
func (m *FilterBar) Focus() {
	m.focused = true
	m.textInput.Focus()
}

// Blur removes focus from the filter bar.
func (m *FilterBar) Blur() {
	m.focused = false
	m.textInput.Blur()
}

// Reset resets the filter bar.
func (m *FilterBar) Reset() {
	m.textInput.Reset()
	m.visible = false
	m.focused = false
}

// GetExpression returns the filter expression.
func (m *FilterBar) GetExpression() string {
	return strings.TrimSpace(m.textInput.Value())
}

// HasExpression returns true if there's a filter expression.
func (m *FilterBar) HasExpression() bool {
	return m.GetExpression() != ""
}
