package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
	"github.com/ersanisk/sieve/pkg/logentry"
)

// Sidebar displays detailed information about a selected entry.
type Sidebar struct {
	visible  bool
	width    int
	height   int
	entry    logentry.Entry
	theme    theme.Theme
	expanded map[string]bool
}

// NewSidebar creates a new Sidebar.
func NewSidebar(theme theme.Theme) Sidebar {
	return Sidebar{
		visible:  false,
		theme:    theme,
		expanded: make(map[string]bool),
	}
}

// Show shows the sidebar.
func (m *Sidebar) Show() {
	m.visible = true
}

// Hide hides the sidebar.
func (m *Sidebar) Hide() {
	m.visible = false
}

// IsVisible returns true if the sidebar is visible.
func (m *Sidebar) IsVisible() bool {
	return m.visible
}

// SetEntry sets the entry to display.
func (m *Sidebar) SetEntry(entry logentry.Entry) {
	m.entry = entry
	m.expanded = make(map[string]bool)
}

// GetEntry returns the current entry.
func (m *Sidebar) GetEntry() logentry.Entry {
	return m.entry
}

// SetSize sets the dimensions of the sidebar.
func (m *Sidebar) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetSize returns the dimensions of the sidebar.
func (m *Sidebar) GetSize() (int, int) {
	return m.width, m.height
}

// SetTheme sets the theme.
func (m *Sidebar) SetTheme(theme theme.Theme) {
	m.theme = theme
}

// ToggleField toggles expansion of a field.
func (m *Sidebar) ToggleField(key string) {
	if _, ok := m.expanded[key]; ok {
		delete(m.expanded, key)
	} else {
		m.expanded[key] = true
	}
}

// View renders the sidebar.
func (m Sidebar) View() string {
	if !m.visible {
		return ""
	}

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Colors().Border).
		Width(m.width).
		Height(m.height)

	content := m.renderContent()
	return containerStyle.Render(content)
}

// renderContent renders the sidebar content.
func (m Sidebar) renderContent() string {
	if m.entry.Raw == "" {
		style := lipgloss.NewStyle().
			Foreground(m.theme.Colors().Foreground).
			Align(lipgloss.Center).
			Height(m.height - 2)
		return style.Render("No entry selected")
	}

	var builder strings.Builder
	builder.WriteString(m.renderHeader())
	builder.WriteString("\n\n")
	builder.WriteString(m.renderFields())

	return builder.String()
}

// renderHeader renders the entry header.
func (m Sidebar) renderHeader() string {
	levelStyle := m.theme.LevelStyle(m.entry.Level)
	level := levelStyle.Render(fmt.Sprintf("[%s]", m.entry.Level.String()))

	timestampStyle := m.theme.TimestampStyle()
	timestamp := ""
	if !m.entry.Timestamp.IsZero() {
		timestamp = timestampStyle.Render(fmt.Sprintf(" %s", m.entry.Timestamp.Format("2006-01-02 15:04:05")))
	}

	return level + timestamp
}

// renderFields renders the entry fields.
func (m Sidebar) renderFields() string {
	if len(m.entry.Fields) == 0 {
		style := lipgloss.NewStyle().
			Foreground(m.theme.Colors().Foreground).
			Italic(true)
		return style.Render("No fields")
	}

	var builder strings.Builder
	for key, value := range m.entry.Fields {
		builder.WriteString(m.renderField(key, value, 0))
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderField renders a single field.
func (m Sidebar) renderField(key string, value any, indent int) string {
	keyStyle := m.theme.KeyStyle()
	keyText := keyStyle.Render(fmt.Sprintf("%*s%s", indent, "", key))

	valueText := m.renderValue(value, indent+2)
	return fmt.Sprintf("%s: %s", keyText, valueText)
}

// renderValue renders a field value.
func (m Sidebar) renderValue(value any, indent int) string {
	valueStyle := m.theme.ValueStyle()

	switch v := value.(type) {
	case string:
		return valueStyle.Render(fmt.Sprintf("\"%s\"", v))
	case float64:
		return valueStyle.Render(fmt.Sprintf("%.2f", v))
	case int:
		return valueStyle.Render(fmt.Sprintf("%d", v))
	case bool:
		return valueStyle.Render(fmt.Sprintf("%t", v))
	case map[string]any:
		return m.renderMap(v, indent)
	case []any:
		return m.renderArray(v, indent)
	default:
		return valueStyle.Render(fmt.Sprintf("%v", v))
	}
}

// renderMap renders a map value.
func (m Sidebar) renderMap(mmap map[string]any, indent int) string {
	if len(mmap) == 0 {
		return "{}"
	}

	var builder strings.Builder
	builder.WriteString("{")

	i := 0
	for key, val := range mmap {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString("\n")
		builder.WriteString(m.renderField(key, val, indent+2))
		i++
	}

	builder.WriteString("\n")
	for j := 0; j < indent; j++ {
		builder.WriteString(" ")
	}
	builder.WriteString("}")

	return builder.String()
}

// renderArray renders an array value.
func (m Sidebar) renderArray(arr []any, indent int) string {
	if len(arr) == 0 {
		return "[]"
	}

	var builder strings.Builder
	builder.WriteString("[")

	for i, val := range arr {
		if i > 0 {
			builder.WriteString(", ")
		}
		builder.WriteString(m.renderValue(val, indent+2))
	}

	builder.WriteString("]")
	return builder.String()
}
