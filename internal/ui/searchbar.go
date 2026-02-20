package ui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
)

// SearchBar handles search input.
type SearchBar struct {
	textInput textinput.Model
	visible   bool
	focused   bool
	width     int
	theme     theme.Theme
}

// NewSearchBar creates a new SearchBar.
func NewSearchBar(theme theme.Theme) SearchBar {
	ti := textinput.New()
	ti.Placeholder = "Search..."
	ti.Prompt = "/ "
	ti.CharLimit = 200
	ti.PromptStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	ti.TextStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	ti.Cursor.Style = lipgloss.NewStyle().Foreground(theme.Colors().Highlight)
	ti.Cursor.TextStyle = lipgloss.NewStyle().Foreground(theme.Colors().Background)

	return SearchBar{
		textInput: ti,
		visible:   false,
		focused:   false,
		theme:     theme,
	}
}

// Show shows the search bar.
func (m *SearchBar) Show() {
	m.visible = true
	m.focused = true
	m.textInput.Focus()
	m.textInput.Reset()
}

// Hide hides the search bar.
func (m *SearchBar) Hide() {
	m.visible = false
	m.focused = false
	m.textInput.Blur()
}

// IsVisible returns true if the search bar is visible.
func (m *SearchBar) IsVisible() bool {
	return m.visible
}

// IsFocused returns true if the search bar is focused.
func (m *SearchBar) IsFocused() bool {
	return m.focused
}

// SetValue sets the search value.
func (m *SearchBar) SetValue(value string) {
	m.textInput.SetValue(value)
}

// GetValue returns the current search value.
func (m *SearchBar) GetValue() string {
	return m.textInput.Value()
}

// Clear clears the search input.
func (m *SearchBar) Clear() {
	m.textInput.Reset()
}

// SetSize sets the dimensions of the search bar.
func (m *SearchBar) SetSize(width, height int) {
	m.width = width
	// Reserve space for prompt (2 chars: "/ ") and some padding
	if width > 4 {
		m.textInput.Width = width - 4
	}
}

// SetTheme sets the theme.
func (m *SearchBar) SetTheme(theme theme.Theme) {
	m.theme = theme
	m.textInput.PromptStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	m.textInput.TextStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	m.textInput.PlaceholderStyle = lipgloss.NewStyle().Foreground(theme.Colors().Foreground)
	m.textInput.Cursor.Style = lipgloss.NewStyle().Foreground(theme.Colors().Highlight)
	m.textInput.Cursor.TextStyle = lipgloss.NewStyle().Foreground(theme.Colors().Background)
}

// Update updates the search bar.
func (m SearchBar) Update(msg tea.Msg) (SearchBar, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	// ESC tuşunu textinput'a göndermeden önce kontrol et
	if keyMsg, ok := msg.(tea.KeyMsg); ok {
		if keyMsg.Type == tea.KeyEsc {
			// ESC tuşunu consume etme, parent'a geçmesine izin ver
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// View renders the search bar.
func (m SearchBar) View() string {
	if !m.visible {
		return ""
	}
	return m.textInput.View()
}

// Focus sets focus to the search bar.
func (m *SearchBar) Focus() {
	m.focused = true
	m.textInput.Focus()
}

// Blur removes focus from the search bar.
func (m *SearchBar) Blur() {
	m.focused = false
	m.textInput.Blur()
}

// Reset resets the search bar.
func (m *SearchBar) Reset() {
	m.textInput.Reset()
	m.visible = false
	m.focused = false
}

// Width returns the current width.
func (m *SearchBar) Width() int {
	return m.width
}

// GetQuery returns the query with leading slash removed.
func (m *SearchBar) GetQuery() string {
	query := strings.TrimPrefix(m.textInput.Value(), "/")
	return strings.TrimSpace(query)
}

// HasQuery returns true if there's a query.
func (m *SearchBar) HasQuery() bool {
	return m.GetQuery() != ""
}
