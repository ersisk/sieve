package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
)

// Help displays help information and key bindings.
type Help struct {
	visible      bool
	width        int
	height       int
	scrollOffset int
	theme        theme.Theme
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
	m.scrollOffset = 0
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

// Update handles keyboard input for help.
func (m Help) Update(msg tea.Msg) (Help, tea.Cmd) {
	switch msg := msg.(type) {
	case ScrollUpMsg:
		m.scrollUp(msg.Amount)
	case ScrollDownMsg:
		m.scrollDown(msg.Amount)
	}
	return m, nil
}

// scrollUp scrolls up by the specified amount.
func (m *Help) scrollUp(amount int) {
	m.scrollOffset -= amount
	if m.scrollOffset < 0 {
		m.scrollOffset = 0
	}
}

// scrollDown scrolls down by the specified amount.
func (m *Help) scrollDown(amount int) {
	contentLines := m.getContentLines()
	maxOffset := len(contentLines) - (m.height - 4)
	if maxOffset < 0 {
		maxOffset = 0
	}
	m.scrollOffset += amount
	if m.scrollOffset > maxOffset {
		m.scrollOffset = maxOffset
	}
}

// getContentLines returns the content as lines.
func (m Help) getContentLines() []string {
	content := m.renderContent()
	return strings.Split(content, "\n")
}

// View renders the help overlay.
func (m Help) View() string {
	if !m.visible {
		return ""
	}

	if m.width == 0 || m.height == 0 {
		return ""
	}

	content := m.renderContent()
	lines := strings.Split(content, "\n")

	// Border + padding alır 4 satır (top border + padding + bottom padding + bottom border)
	innerHeight := m.height - 4
	if innerHeight < 1 {
		innerHeight = 1
	}

	totalLines := len(lines)
	maxOffset := totalLines - innerHeight
	if maxOffset < 0 {
		maxOffset = 0
	}

	scrollOffset := m.scrollOffset
	if scrollOffset > maxOffset {
		scrollOffset = maxOffset
	}

	end := scrollOffset + innerHeight
	if end > totalLines {
		end = totalLines
	}

	visibleLines := lines[scrollOffset:end]

	// İçeriği tam genişliğe pad et
	innerWidth := m.width - 6 // border + padding
	if innerWidth < 0 {
		innerWidth = 0
	}
	for i, line := range visibleLines {
		if len(line) < innerWidth {
			visibleLines[i] = line + strings.Repeat(" ", innerWidth-len(line))
		}
	}

	scrollIndicator := ""
	if totalLines > innerHeight {
		pct := 0
		if maxOffset > 0 {
			pct = scrollOffset * 100 / maxOffset
		}
		scrollIndicator = fmt.Sprintf(" [%d%%] j/k: scroll", pct)
	}

	visibleContent := strings.Join(visibleLines, "\n")

	titleLine := fmt.Sprintf("Sieve Help%s", scrollIndicator)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Colors().Highlight).
		Padding(0, 2).
		Width(m.width - 2).
		Background(m.theme.Colors().Background).
		Foreground(m.theme.Colors().Foreground)

	titleStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Highlight).
		Bold(true).
		Width(m.width - 6).
		Align(lipgloss.Center)

	body := titleStyle.Render(titleLine) + "\n\n" + visibleContent

	return borderStyle.Render(body)
}

// renderContent renders the help content.
func (m Help) renderContent() string {
	var builder strings.Builder

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
		{"f", "Filter expression"},
		{"Esc", "Clear filter/Cancel"},
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
		{"Enter", "View log details"},
		{"Tab", "Toggle sidebar"},
		{"d", "Toggle dashboard"},
		{"?", "Toggle help"},
		{"Esc", "Close overlay / exit mode"},
	}))

	builder.WriteString("\n")
	builder.WriteString(m.renderSection("File & Program", []keyBinding{
		{"r", "Toggle sort order"},
		{"R", "Refresh file"},
		{"F", "Toggle follow mode"},
		{"Ctrl+C", "Force quit"},
		{"q", "Quit"},
	}))

	builder.WriteString("\n")
	builder.WriteString(m.renderFilterExamples())

	return builder.String()
}

func (m Help) renderFilterExamples() string {
	exampleStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Value).
		Italic(true)

	return `Filter Examples:
  .level >= 40              Show warn/error/fatal logs
  .service == "api"         Filter by service name
  .msg contains "error"      Search in message
  .duration_ms > 1000        Numeric comparison
  .status >= 500            HTTP status filter

Operators:
  ==, !=, >, <, >=, <=    Comparison
  contains                  Text contains
  matches                   Regex match
  and, or, not             Logical` + "\n" + exampleStyle.Render("Press Enter to apply filter, Esc to cancel")
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
