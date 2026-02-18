package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
)

// Help displays help information and key bindings.
type Help struct {
	visible bool
	width   int
	height  int
	theme   theme.Theme
}

// NewHelp creates a new Help overlay.
func NewHelp(theme theme.Theme) Help {
	return Help{
		visible: false,
		theme:   theme,
	}
}

// Show shows the help overlay.
func (m *Help) Show() {
	m.visible = true
}

// Hide hides the help overlay.
func (m *Help) Hide() {
	m.visible = false
}

// IsVisible returns true if the help overlay is visible.
func (m *Help) IsVisible() bool {
	return m.visible
}

// SetSize sets the dimensions of the help overlay.
func (m *Help) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetSize returns the dimensions of the help overlay.
func (m *Help) GetSize() (int, int) {
	return m.width, m.height
}

// SetTheme sets the theme.
func (m *Help) SetTheme(theme theme.Theme) {
	m.theme = theme
}

// View renders the help overlay.
func (m Help) View() string {
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

// renderContent renders the help content.
func (m Help) renderContent() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground).
		Bold(true).
		Underline(true)

	var builder strings.Builder

	builder.WriteString(titleStyle.Render("Sieve - Help"))
	builder.WriteString("\n\n")

	builder.WriteString(m.renderSection("Navigation", []keyBinding{
		{"j / ↓", "Scroll down"},
		{"k / ↑", "Scroll up"},
		{"g", "Go to top"},
		{"G", "Go to bottom"},
		{"PgDn / Space", "Page down"},
		{"PgUp", "Page up"},
	}))

	builder.WriteString("\n")
	builder.WriteString(m.renderSection("Search & Filter", []keyBinding{
		{"/", "Search"},
		{"n", "Next search result"},
		{"N", "Previous search result"},
		{"f", "Filter"},
		{"F", "Clear filter"},
	}))

	builder.WriteString("\n")
	builder.WriteString(m.renderSection("Level Filter", []keyBinding{
		{"1", "Debug"},
		{"2", "Info"},
		{"3", "Warn"},
		{"4", "Error"},
		{"5", "Fatal"},
		{"0", "No filter"},
	}))

	builder.WriteString("\n")
	builder.WriteString(m.renderSection("View & Actions", []keyBinding{
		{"Enter", "Expand/Collapse entry"},
		{"Tab", "Toggle sidebar"},
		{"d", "Toggle dashboard"},
		{"Esc", "Close overlay / exit mode"},
	}))

	builder.WriteString("\n")
	builder.WriteString(m.renderSection("File & Program", []keyBinding{
		{"r", "Refresh"},
		{"f (in tail mode)", "Toggle follow"},
		{"q", "Quit"},
		{"?", "Toggle this help"},
	}))

	return builder.String()
}

// renderSection renders a help section.
func (m Help) renderSection(title string, bindings []keyBinding) string {
	headerStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground).
		Bold(true)

	keyStyle := m.theme.KeyStyle()
	descStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground)

	var builder strings.Builder
	builder.WriteString(headerStyle.Render(title))
	builder.WriteString(":\n")

	for _, binding := range bindings {
		keyText := keyStyle.Render(fmt.Sprintf("%-10s", binding.key))
		descText := descStyle.Render(binding.description)
		builder.WriteString(fmt.Sprintf("  %s %s\n", keyText, descText))
	}

	return builder.String()
}

// keyBinding represents a keyboard binding.
type keyBinding struct {
	key         string
	description string
}
