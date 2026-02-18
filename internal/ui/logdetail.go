package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
	"github.com/ersanisk/sieve/pkg/logentry"
)

// LogDetail displays detailed information about a log entry.
type LogDetail struct {
	visible bool
	entry   logentry.Entry
	theme   theme.Theme
	width   int
	height  int
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
}

// View renders the log detail modal.
func (m LogDetail) View() string {
	if !m.visible {
		return ""
	}

	if m.width == 0 || m.height == 0 {
		return ""
	}

	modalWidth := min(m.width-4, 80)
	modalHeight := min(m.height-4, 30)

	if modalWidth < 20 {
		return ""
	}

	content := m.renderContent(modalWidth, modalHeight)

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Colors().Highlight).
		Padding(1, 2).
		Width(modalWidth).
		Height(modalHeight).
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

func (m LogDetail) renderContent(width, height int) string {
	var builder strings.Builder

	builder.WriteString(m.renderHeader(width))
	builder.WriteString("\n\n")
	builder.WriteString(m.renderLevel(width))
	builder.WriteString(m.renderTimestamp(width))
	builder.WriteString(m.renderMessage(width))
	builder.WriteString("\n")
	builder.WriteString(m.renderFields(width, height))

	return builder.String()
}

func (m LogDetail) renderHeader(width int) string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(m.theme.Colors().Highlight).
		Width(width)

	return headerStyle.Render("Log Entry Details")
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

	levelStyle := lipgloss.NewStyle().
		Background(color).
		Foreground(m.theme.Colors().Background).
		Bold(true).
		Padding(0, 1)

	return fmt.Sprintf("Level: %s", levelStyle.Render(m.entry.Level.String()))
}

func (m LogDetail) renderTimestamp(width int) string {
	if m.entry.Timestamp.IsZero() {
		return ""
	}
	timestampStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Timestamp)
	return fmt.Sprintf("Time:  %s", timestampStyle.Render(m.entry.Timestamp.Format("2006-01-02 15:04:05.000")))
}

func (m LogDetail) renderMessage(width int) string {
	messageStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground)
	return fmt.Sprintf("Msg:   %s", messageStyle.Render(m.entry.Message))
}

func (m LogDetail) renderFields(width, height int) string {
	if len(m.entry.Fields) == 0 {
		return ""
	}

	var builder strings.Builder
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Underline(true).
		Foreground(m.theme.Colors().Key)
	builder.WriteString(headerStyle.Render("Fields:"))
	builder.WriteString("\n")

	keyStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Key).
		Width(15)
	valueStyle := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Value)

	for key, value := range m.entry.Fields {
		valueStr := m.formatValue(value)
		builder.WriteString(keyStyle.Render(key + ": "))
		builder.WriteString(valueStyle.Render(valueStr))
		builder.WriteString("\n")
	}

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
