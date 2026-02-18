package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
	"github.com/ersanisk/sieve/pkg/logentry"
)

// Dashboard displays statistics and metrics about log entries.
type Dashboard struct {
	visible     bool
	width       int
	height      int
	entries     []logentry.Entry
	levelCounts map[logentry.Level]int
	linesPerSec float64
	theme       theme.Theme
}

// NewDashboard creates a new Dashboard.
func NewDashboard(theme theme.Theme) Dashboard {
	return Dashboard{
		visible:     false,
		theme:       theme,
		levelCounts: make(map[logentry.Level]int),
	}
}

// Show shows the dashboard.
func (m *Dashboard) Show() {
	m.visible = true
}

// Hide hides the dashboard.
func (m *Dashboard) Hide() {
	m.visible = false
}

// IsVisible returns true if the dashboard is visible.
func (m *Dashboard) IsVisible() bool {
	return m.visible
}

// SetEntries sets the log entries and updates statistics.
func (m *Dashboard) SetEntries(entries []logentry.Entry) {
	m.entries = entries
	m.levelCounts = make(map[logentry.Level]int)

	for _, entry := range entries {
		m.levelCounts[entry.Level]++
	}
}

// GetEntries returns the current entries.
func (m *Dashboard) GetEntries() []logentry.Entry {
	return m.entries
}

// SetLinesPerSec sets the lines per second metric.
func (m *Dashboard) SetLinesPerSec(lps float64) {
	m.linesPerSec = lps
}

// SetSize sets the dimensions of the dashboard.
func (m *Dashboard) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetSize returns the dimensions of the dashboard.
func (m *Dashboard) GetSize() (int, int) {
	return m.width, m.height
}

// SetTheme sets the theme.
func (m *Dashboard) SetTheme(theme theme.Theme) {
	m.theme = theme
}

// View renders the dashboard.
func (m Dashboard) View() string {
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

// renderContent renders the dashboard content.
func (m Dashboard) renderContent() string {
	titleStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground).
		Bold(true).
		Underline(true)

	var builder strings.Builder
	builder.WriteString(titleStyle.Render("Dashboard"))
	builder.WriteString("\n\n")

	builder.WriteString(m.renderLevelDistribution())
	builder.WriteString("\n\n")
	builder.WriteString(m.renderMetrics())

	return builder.String()
}

// renderLevelDistribution renders the level distribution.
func (m Dashboard) renderLevelDistribution() string {
	if len(m.entries) == 0 {
		style := lipgloss.NewStyle().
			Foreground(m.theme.Colors().Foreground).
			Italic(true)
		return style.Render("No data available")
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground).
		Bold(true)

	var builder strings.Builder
	builder.WriteString(headerStyle.Render("Level Distribution"))
	builder.WriteString("\n")

	total := len(m.entries)
	levels := []logentry.Level{
		logentry.Debug,
		logentry.Info,
		logentry.Warn,
		logentry.Error,
		logentry.Fatal,
	}

	for _, level := range levels {
		count := m.levelCounts[level]
		if count > 0 {
			percentage := float64(count) / float64(total) * 100
			bar := m.renderBar(count, total, m.theme.LevelStyle(level))
			builder.WriteString(fmt.Sprintf("%-8s %s %5d (%5.1f%%)\n",
				level.String()+":", bar, count, percentage))
		}
	}

	return builder.String()
}

// renderMetrics renders general metrics.
func (m Dashboard) renderMetrics() string {
	headerStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground).
		Bold(true)

	valueStyle := m.theme.ValueStyle()

	var builder strings.Builder
	builder.WriteString(headerStyle.Render("Metrics"))
	builder.WriteString("\n")

	totalEntries := len(m.entries)
	builder.WriteString(fmt.Sprintf("%-20s %s\n", "Total Entries:",
		valueStyle.Render(fmt.Sprintf("%d", totalEntries))))

	if m.linesPerSec > 0 {
		builder.WriteString(fmt.Sprintf("%-20s %s\n", "Lines/Second:",
			valueStyle.Render(fmt.Sprintf("%.2f", m.linesPerSec))))
	}

	uniqueFields := m.countUniqueFields()
	builder.WriteString(fmt.Sprintf("%-20s %s\n", "Unique Fields:",
		valueStyle.Render(fmt.Sprintf("%d", uniqueFields))))

	return builder.String()
}

// renderBar renders a progress bar.
func (m Dashboard) renderBar(count, total int, style lipgloss.Style) string {
	if total == 0 {
		return ""
	}

	maxWidth := 30
	percentage := float64(count) / float64(total)
	barWidth := int(float64(maxWidth) * percentage)
	if barWidth > maxWidth {
		barWidth = maxWidth
	}

	bar := strings.Repeat("█", barWidth)
	empty := strings.Repeat("░", maxWidth-barWidth)

	return style.Render(bar + lipgloss.NewStyle().Foreground(m.theme.Colors().Background).Render(empty))
}

// countUniqueFields counts the number of unique fields across all entries.
func (m Dashboard) countUniqueFields() int {
	fields := make(map[string]bool)

	for _, entry := range m.entries {
		for key := range entry.Fields {
			fields[key] = true
		}
	}

	return len(fields)
}

// GetLevelCount returns the count for a specific level.
func (m *Dashboard) GetLevelCount(level logentry.Level) int {
	return m.levelCounts[level]
}

// GetTotalCount returns the total count of all entries.
func (m *Dashboard) GetTotalCount() int {
	return len(m.entries)
}

// GetHighestLevel returns the highest log level.
func (m *Dashboard) GetHighestLevel() logentry.Level {
	highest := logentry.Unknown

	for level := range m.levelCounts {
		if level > highest {
			highest = level
		}
	}

	return highest
}
