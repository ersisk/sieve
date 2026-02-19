package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
	"github.com/ersanisk/sieve/pkg/logentry"
)

// LogDetail displays detailed information about a log entry.
type LogDetail struct {
	visible  bool
	entry    logentry.Entry
	theme    theme.Theme
	width    int
	height   int
	viewport viewport.Model
}

// NewLogDetail creates a new LogDetail modal.
func NewLogDetail(theme theme.Theme) *LogDetail {
	return &LogDetail{
		visible: false,
		theme:   theme,
	}
}

// Show shows the log detail modal for the given entry.
func (m *LogDetail) Show(entry logentry.Entry) {
	m.entry = entry
	m.visible = true
	m.updateContent()
	m.viewport.GotoTop()
}

// Hide hides the log detail modal.
func (m *LogDetail) Hide() {
	m.visible = false
}

// IsVisible returns whether the modal is visible.
func (m *LogDetail) IsVisible() bool {
	return m.visible
}

// SetSize sets the dimensions of the modal.
func (m *LogDetail) SetSize(width, height int) {
	m.width = width
	m.height = height

	modalWidth := min(m.width-8, 100)
	modalHeight := min(m.height-8, 50)

	m.viewport.Width = modalWidth
	m.viewport.Height = modalHeight
	m.updateContent()
}

// Update handles events for the log detail modal.
func (m *LogDetail) Update(msg any) {
	// Ideally LogDetail should be a proper tea.Model, but for now we wrap it.
	// We only care about scrolling keys.
	newViewport, _ := m.viewport.Update(msg)
	m.viewport = newViewport
}

// View renders the log detail modal.
func (m LogDetail) View() string {
	if !m.visible {
		return ""
	}

	if m.width == 0 || m.height == 0 {
		return ""
	}

	// Ensure viewport size is correct if not set
	if m.viewport.Width == 0 {
		modalWidth := min(m.width-4, 80)
		modalHeight := min(m.height-4, 30)
		m.viewport.Width = modalWidth
		m.viewport.Height = modalHeight
	}

	content := m.viewport.View()

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Colors().Highlight).
		Padding(1, 2).
		Width(m.viewport.Width).
		Height(m.viewport.Height).
		Background(m.theme.Colors().Background).
		Foreground(m.theme.Colors().Foreground)

	centered := lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		modalStyle.Render(content),
	)

	return centered
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (m *LogDetail) updateContent() {
	if m.viewport.Width == 0 {
		return
	}
	content := m.renderContent(m.viewport.Width, m.viewport.Height)
	m.viewport.SetContent(content)
}

func (m LogDetail) renderContent(width, height int) string {
	var builder strings.Builder

	// Header
	builder.WriteString(m.renderHeader(width))
	builder.WriteString("\n")
	builder.WriteString(m.renderSeparator(width, "═"))
	builder.WriteString("\n\n")

	// Main info section
	builder.WriteString(m.renderLevel(width))
	builder.WriteString("\n")
	builder.WriteString(m.renderTimestamp(width))
	builder.WriteString("\n")
	builder.WriteString(m.renderMessage(width))
	builder.WriteString("\n\n")

	// Fields section
	if len(m.entry.Fields) > 0 {
		builder.WriteString(m.renderSeparator(width, "─"))
		builder.WriteString("\n")
		builder.WriteString(m.renderFields(width, height))
		builder.WriteString("\n")
	}

	// Raw JSON section
	if m.entry.Raw != "" {
		builder.WriteString(m.renderSeparator(width, "─"))
		builder.WriteString("\n")
		builder.WriteString(m.renderRawJSON(width))
	}

	return builder.String()
}

func (m LogDetail) renderHeader(width int) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Colors().Highlight).
		Width(width)

	return headerStyle.Render("📋 Log Entry Details")
}

func (m LogDetail) renderSeparator(width int, char string) string {
	separatorStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Border)
	return separatorStyle.Render(strings.Repeat(char, width))
}

func (m LogDetail) renderLevel(width int) string {
	var color lipgloss.Color
	switch m.entry.Level {
	case logentry.Debug:
		color = m.theme.Colors().Debug
	case logentry.Info:
		color = m.theme.Colors().Info
	case logentry.Warn:
		color = m.theme.Colors().Warn
	case logentry.Error:
		color = m.theme.Colors().Error
	case logentry.Fatal:
		color = m.theme.Colors().Fatal
	default:
		color = m.theme.Colors().Foreground
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Key).
		Bold(true).
		Width(12)

	levelStyle := lipgloss.NewStyle().
		Background(color).
		Foreground(m.theme.Colors().Background).
		Bold(true).
		Padding(0, 2)

	return labelStyle.Render("LEVEL") + levelStyle.Render(m.entry.Level.String())
}

func (m LogDetail) renderTimestamp(width int) string {
	if m.entry.Timestamp.IsZero() {
		return ""
	}

	labelStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Key).
		Bold(true).
		Width(12)

	timestampStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Timestamp)

	return labelStyle.Render("TIME") + timestampStyle.Render(m.entry.Timestamp.Format("2006-01-02 15:04:05.000"))
}

func (m LogDetail) renderMessage(width int) string {
	labelStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Key).
		Bold(true).
		Width(12)

	messageStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground).
		Width(width - 12)

	return labelStyle.Render("MESSAGE") + messageStyle.Render(m.entry.Message)
}

func (m LogDetail) renderFields(width, height int) string {
	if len(m.entry.Fields) == 0 {
		return ""
	}

	var builder strings.Builder
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Colors().Highlight)
	builder.WriteString(headerStyle.Render("FIELDS"))
	builder.WriteString("\n\n")

	keyStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Key).
		Width(18)
	valueStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Value)

	// Sort keys for stable ordering
	keys := make([]string, 0, len(m.entry.Fields))
	for k := range m.entry.Fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := m.entry.Fields[key]
		valueStr := m.formatValue(value)
		builder.WriteString("  ")
		builder.WriteString(keyStyle.Render(key + ":"))
		builder.WriteString(" ")
		builder.WriteString(valueStyle.Render(valueStr))
		builder.WriteString("\n")
	}

	return builder.String()
}

func (m LogDetail) renderRawJSON(width int) string {
	if m.entry.Raw == "" {
		return ""
	}

	var builder strings.Builder
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Colors().Highlight)

	builder.WriteString(headerStyle.Render("RAW JSON"))
	builder.WriteString("\n\n")

	rawStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Value).
		Width(width - 4) // Padding için

	builder.WriteString("  ")
	builder.WriteString(rawStyle.Render(m.entry.Raw))

	return builder.String()
}

func (m LogDetail) formatValue(value any) string {
	switch v := value.(type) {
	case string:
		if len(v) > 50 {
			return "\"" + v[:47] + "...\""
		}
		return "\"" + v + "\""
	case float64:
		return fmt.Sprintf("%.6g", v)
	case int:
		return fmt.Sprintf("%d", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case map[string]any:
		return fmt.Sprintf("{%d fields}", len(v))
	case []any:
		return fmt.Sprintf("[%d items]", len(v))
	default:
		return fmt.Sprintf("%v", v)
	}
}
