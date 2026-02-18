package ui

import (
	"fmt"
	"strings"

	"github.com/ersanisk/sieve/internal/theme"
	"github.com/ersanisk/sieve/pkg/logentry"
)

// StatusBar displays status information at the bottom of the screen.
type StatusBar struct {
	filePath    string
	totalLines  int
	selected    int
	filter      string
	mode        string
	following   bool
	levelFilter logentry.Level
	width       int
	height      int
	theme       theme.Theme
}

// NewStatusBar creates a new StatusBar.
func NewStatusBar(theme theme.Theme) StatusBar {
	return StatusBar{
		mode:        "view",
		theme:       theme,
		following:   false,
		levelFilter: logentry.Unknown,
	}
}

// SetFilePath sets the file path.
func (m *StatusBar) SetFilePath(path string) {
	m.filePath = path
}

// SetTotalLines sets the total number of lines.
func (m *StatusBar) SetTotalLines(count int) {
	m.totalLines = count
}

// SetSelected sets the selected line number.
func (m *StatusBar) SetSelected(index int) {
	m.selected = index + 1
}

// SetFilter sets the active filter expression.
func (m *StatusBar) SetFilter(filter string) {
	m.filter = filter
}

// SetMode sets the current mode.
func (m *StatusBar) SetMode(mode string) {
	m.mode = mode
}

// SetFollowing sets the follow mode state.
func (m *StatusBar) SetFollowing(following bool) {
	m.following = following
}

// SetLevelFilter sets the level filter.
func (m *StatusBar) SetLevelFilter(level logentry.Level) {
	m.levelFilter = level
}

// SetSize sets the dimensions of the status bar.
func (m *StatusBar) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetSize returns the dimensions of the status bar.
func (m *StatusBar) GetSize() (int, int) {
	return m.width, m.height
}

// SetTheme sets the theme.
func (m *StatusBar) SetTheme(theme theme.Theme) {
	m.theme = theme
}

// View renders the status bar.
func (m StatusBar) View() string {
	style := m.theme.StatusBarStyle()
	content := m.renderContent()
	return style.Width(m.width).Height(m.height).Render(content)
}

// renderContent renders the status bar content.
func (m StatusBar) renderContent() string {
	var parts []string

	filePart := m.renderFilePart()
	if filePart != "" {
		parts = append(parts, filePart)
	}

	linesPart := m.renderLinesPart()
	if linesPart != "" {
		parts = append(parts, linesPart)
	}

	filterPart := m.renderFilterPart()
	if filterPart != "" {
		parts = append(parts, filterPart)
	}

	modePart := m.renderModePart()
	if modePart != "" {
		parts = append(parts, modePart)
	}

	return strings.Join(parts, " │ ")
}

// renderFilePart renders the file information.
func (m StatusBar) renderFilePart() string {
	if m.filePath == "" {
		return ""
	}

	shortPath := m.filePath
	if len(shortPath) > 20 {
		shortPath = "..." + shortPath[len(shortPath)-20:]
	}

	return fmt.Sprintf("📄 %s", shortPath)
}

// renderLinesPart renders the line information.
func (m StatusBar) renderLinesPart() string {
	if m.totalLines == 0 {
		return ""
	}
	return fmt.Sprintf("📝 %d/%d", m.selected, m.totalLines)
}

// renderFilterPart renders the filter information.
func (m StatusBar) renderFilterPart() string {
	var filterInfo []string

	if m.filter != "" {
		shortFilter := m.filter
		if len(shortFilter) > 15 {
			shortFilter = shortFilter[:15] + "..."
		}
		filterInfo = append(filterInfo, fmt.Sprintf("🔍 %s", shortFilter))
	}

	if m.levelFilter != logentry.Unknown {
		filterInfo = append(filterInfo, fmt.Sprintf("🏷️ %s", m.levelFilter.String()))
	}

	if m.following {
		filterInfo = append(filterInfo, "👁️ FOLLOW")
	}

	return strings.Join(filterInfo, " ")
}

// renderModePart renders the mode information.
func (m StatusBar) renderModePart() string {
	if m.mode == "" {
		return ""
	}
	return fmt.Sprintf("⚙️ %s", strings.ToUpper(m.mode))
}
