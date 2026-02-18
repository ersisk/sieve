package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
	"github.com/ersanisk/sieve/pkg/logentry"
)

// highlightQuery highlights all occurrences of query in text using a bold+color style.
func highlightQuery(text, query string, hlStyle lipgloss.Style) string {
	if query == "" {
		return text
	}
	lower := strings.ToLower(text)
	lowerQ := strings.ToLower(query)
	if !strings.Contains(lower, lowerQ) {
		return text
	}

	var result strings.Builder
	for len(text) > 0 {
		idx := strings.Index(strings.ToLower(text), lowerQ)
		if idx < 0 {
			result.WriteString(text)
			break
		}
		result.WriteString(text[:idx])
		result.WriteString(hlStyle.Render(text[idx : idx+len(query)]))
		text = text[idx+len(query):]
	}
	return result.String()
}

// LogView displays log entries with virtual scrolling.
type LogView struct {
	entries     []logentry.Entry
	offset      int
	height      int
	selected    int
	width       int
	theme       theme.Theme
	lineNumbers bool
	expanded    map[int]bool
	searchQuery string
}

// NewLogView creates a new LogView.
func NewLogView(theme theme.Theme) LogView {
	return LogView{
		offset:      0,
		height:      0,
		selected:    0,
		theme:       theme,
		lineNumbers: true,
		expanded:    make(map[int]bool),
	}
}

// SetEntries sets the log entries.
func (m *LogView) SetEntries(entries []logentry.Entry) {
	m.entries = entries
	if m.selected >= len(entries) {
		m.selected = len(entries) - 1
	}
	if len(entries) == 0 {
		m.selected = 0
	}
}

// GetEntries returns the current entries.
func (m *LogView) GetEntries() []logentry.Entry {
	return m.entries
}

// GetSelected returns the selected entry and its index.
func (m *LogView) GetSelected() (logentry.Entry, int) {
	if m.selected >= 0 && m.selected < len(m.entries) {
		return m.entries[m.selected], m.selected
	}
	return logentry.Entry{}, -1
}

// SetSelected sets the selected entry by index.
func (m *LogView) SetSelected(index int) {
	if index < 0 {
		m.selected = 0
	} else if index >= len(m.entries) {
		m.selected = len(m.entries) - 1
	} else {
		m.selected = index
	}
	m.ensureVisible()
}

// ScrollUp scrolls up by the specified amount.
func (m *LogView) ScrollUp(amount int) {
	m.selected -= amount
	if m.selected < 0 {
		m.selected = 0
	}
	m.ensureVisible()
}

// ScrollDown scrolls down by the specified amount.
func (m *LogView) ScrollDown(amount int) {
	m.selected += amount
	if m.selected >= len(m.entries) {
		m.selected = len(m.entries) - 1
	}
	m.ensureVisible()
}

// ScrollToTop scrolls to the top.
func (m *LogView) ScrollToTop() {
	m.selected = 0
	m.offset = 0
}

// ScrollToBottom scrolls to the bottom.
func (m *LogView) ScrollToBottom() {
	m.selected = len(m.entries) - 1
	m.ensureVisible()
}

// ScrollToLine scrolls to a specific line.
func (m *LogView) ScrollToLine(line int) {
	m.SetSelected(line)
}

// ScrollUpOne scrolls up by one line.
func (m *LogView) ScrollUpOne() {
	m.ScrollUp(1)
}

// ScrollDownOne scrolls down by one line.
func (m *LogView) ScrollDownOne() {
	m.ScrollDown(1)
}

// ScrollPageUp scrolls up by one page.
func (m *LogView) ScrollPageUp() {
	m.ScrollUp(m.height)
}

// ScrollPageDown scrolls down by one page.
func (m *LogView) ScrollPageDown() {
	m.ScrollDown(m.height)
}

// ensureVisible ensures the selected entry is visible.
func (m *LogView) ensureVisible() {
	if m.selected < m.offset {
		m.offset = m.selected
	} else if m.selected >= m.offset+m.height {
		m.offset = m.selected - m.height + 1
	}
	if m.offset < 0 {
		m.offset = 0
	}
}

// SetSize sets the dimensions of the view.
func (m *LogView) SetSize(width, height int) {
	m.width = width
	m.height = height
	if m.height <= 0 {
		m.height = 1
	}
	m.ensureVisible()
}

// GetSize returns the dimensions of the view.
func (m *LogView) GetSize() (int, int) {
	return m.width, m.height
}

// SetTheme sets the theme.
func (m *LogView) SetTheme(theme theme.Theme) {
	m.theme = theme
}

// SetSearchQuery sets the current search query for highlighting.
func (m *LogView) SetSearchQuery(query string) {
	m.searchQuery = query
}

// ToggleLineNumbers toggles line numbers.
func (m *LogView) ToggleLineNumbers() {
	m.lineNumbers = !m.lineNumbers
}

// ToggleExpanded toggles expansion of the selected entry.
func (m *LogView) ToggleExpanded() {
	if m.selected >= 0 && m.selected < len(m.entries) {
		if _, ok := m.expanded[m.selected]; ok {
			delete(m.expanded, m.selected)
		} else {
			m.expanded[m.selected] = true
		}
	}
}

// IsExpanded returns true if the entry is expanded.
func (m *LogView) IsExpanded(index int) bool {
	return m.expanded[index]
}

// GetOffset returns the current scroll offset.
func (m *LogView) GetOffset() int {
	return m.offset
}

// GetTotalLines returns the total number of lines.
func (m *LogView) GetTotalLines() int {
	return len(m.entries)
}

// View renders the log view.
// Only renders visible entries for virtual scrolling.
func (m LogView) View() string {
	if len(m.entries) == 0 {
		return m.renderEmpty()
	}

	visibleEnd := m.offset + m.height
	if visibleEnd > len(m.entries) {
		visibleEnd = len(m.entries)
	}

	var builder strings.Builder

	for i := m.offset; i < visibleEnd; i++ {
		if i >= len(m.entries) {
			break
		}
		builder.WriteString(m.renderEntry(i))
	}

	return builder.String()
}

// renderEntry renders a single log entry.
func (m LogView) renderEntry(index int) string {
	entry := m.entries[index]
	isSelected := index == m.selected
	isExpanded := m.expanded[index]

	var line strings.Builder

	if m.lineNumbers {
		lineNumStyle := m.theme.TimestampStyle()
		if isSelected {
			lineNumStyle = lineNumStyle.Background(m.theme.Colors().Highlight).Foreground(m.theme.Colors().Background).Bold(true)
		}
		line.WriteString(lineNumStyle.Render(fmt.Sprintf("%5d ", index+1)))
	}

	entryStyle := m.theme.LevelStyle(entry.Level)
	if isSelected {
		entryStyle = entryStyle.Background(m.theme.Colors().Highlight).Foreground(m.theme.Colors().Background)
	}

	timestamp := ""
	if !entry.Timestamp.IsZero() {
		timestampStyle := m.theme.TimestampStyle()
		if isSelected {
			timestampStyle = timestampStyle.Background(m.theme.Colors().Highlight).Foreground(m.theme.Colors().Background)
		}
		timestamp = timestampStyle.Render(fmt.Sprintf("[%s] ", entry.Timestamp.Format("15:04:05")))
	}

	level := entryStyle.Render(fmt.Sprintf("%-5s ", entry.Level.String()))

	messageStyle := lipgloss.NewStyle().Foreground(m.theme.Colors().Foreground)
	if isSelected {
		messageStyle = messageStyle.Background(m.theme.Colors().Highlight).Foreground(m.theme.Colors().Background).Bold(true)
	}

	rawMsg := truncateText(entry.Message, m.width-40)
	var message string
	if m.searchQuery != "" && !isSelected {
		hlStyle := lipgloss.NewStyle().
			Foreground(m.theme.Colors().Background).
			Background(m.theme.Colors().Warn).
			Bold(true)
		message = messageStyle.Render(highlightQuery(rawMsg, m.searchQuery, hlStyle))
	} else {
		message = messageStyle.Render(rawMsg)
	}

	line.WriteString(timestamp)
	line.WriteString(level)
	line.WriteString(message)

	if isExpanded {
		line.WriteString("\n")
		line.WriteString(m.renderExpandedFields(entry, isSelected))
	}

	line.WriteString("\n")
	return line.String()
}

// renderExpandedFields renders the fields of an expanded entry.
func (m LogView) renderExpandedFields(entry logentry.Entry, isSelected bool) string {
	if len(entry.Fields) == 0 {
		return ""
	}

	keyStyle := m.theme.KeyStyle()
	valueStyle := m.theme.ValueStyle()
	if isSelected {
		keyStyle = keyStyle.Bold(true)
		valueStyle = valueStyle.Bold(true)
	}

	var builder strings.Builder
	for key, value := range entry.Fields {
		valueStr := formatValue(value)
		line := fmt.Sprintf("  %s: %s", keyStyle.Render(key), valueStyle.Render(valueStr))
		builder.WriteString(line)
		builder.WriteString("\n")
	}

	return builder.String()
}

// renderEmpty renders the empty state.
func (m LogView) renderEmpty() string {
	style := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground).
		Bold(true).
		Align(lipgloss.Center).
		Width(m.width).
		Height(m.height)

	return style.Render("No log entries")
}

// formatValue formats a value for display.
func formatValue(value any) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("\"%s\"", v)
	case float64:
		return fmt.Sprintf("%.2f", v)
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case map[string]any:
		return fmt.Sprintf("{%d fields}", len(v))
	default:
		return fmt.Sprintf("%v", v)
	}
}

// truncateText truncates text to fit within width.
func truncateText(text string, width int) string {
	if width <= 0 {
		return ""
	}
	runes := []rune(text)
	if len(runes) <= width {
		return string(runes)
	}
	if width <= 3 {
		return string(runes[:width])
	}
	return string(runes[:width-3]) + "..."
}
