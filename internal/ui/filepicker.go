package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
)

// FilePicker is a file picker component.
type FilePicker struct {
	theme    theme.Theme
	files    []string
	selected int
	visible  bool
	width    int
	height   int
	offset   int
}

// FileSelectedMsg is sent when a file is selected.
type FileSelectedMsg struct {
	Path string
}

// NewFilePicker creates a new file picker.
func NewFilePicker(theme theme.Theme) FilePicker {
	return FilePicker{
		theme:    theme,
		files:    []string{},
		selected: 0,
		visible:  false,
		width:    80,
		height:   24,
		offset:   0,
	}
}

// SetFiles sets the list of files to display.
func (f *FilePicker) SetFiles(files []string) {
	f.files = files
	f.selected = 0
	f.offset = 0
}

// Show shows the file picker.
func (f *FilePicker) Show() {
	f.visible = true
}

// Hide hides the file picker.
func (f *FilePicker) Hide() {
	f.visible = false
}

// IsVisible returns whether the file picker is visible.
func (f FilePicker) IsVisible() bool {
	return f.visible
}

// SetSize sets the dimensions of the file picker.
func (f *FilePicker) SetSize(width, height int) {
	f.width = width
	f.height = height
}

// Update handles messages.
func (f *FilePicker) Update(msg tea.Msg) (*FilePicker, tea.Cmd) {
	if !f.visible {
		return f, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			f.selectPrev()
		case "down", "j":
			f.selectNext()
		case "g":
			f.selectFirst()
		case "G":
			f.selectLast()
		case "enter":
			if len(f.files) > 0 {
				return f, func() tea.Msg {
					return FileSelectedMsg{Path: f.files[f.selected]}
				}
			}
		case "q", "esc", "ctrl+c":
			return f, func() tea.Msg {
				return QuitMsg{}
			}
		}
	}

	return f, nil
}

// View renders the file picker.
func (f *FilePicker) View() string {
	if !f.visible {
		return ""
	}

	var sb strings.Builder

	// Title
	titleStyle := lipgloss.NewStyle().
		Foreground(f.theme.Colors().Foreground).
		Background(f.theme.Colors().Info).
		Bold(true).
		Padding(0, 1).
		Width(f.width)

	title := "Select a Log File"
	if len(f.files) > 0 {
		title = fmt.Sprintf("Select a Log File (%d/%d)", f.selected+1, len(f.files))
	}
	sb.WriteString(titleStyle.Render(title))
	sb.WriteString("\n\n")

	// Calculate visible range
	visibleHeight := f.height - 6 // Reserve space for title, footer, etc.
	if visibleHeight < 1 {
		visibleHeight = 1
	}

	// Adjust offset to keep selected item visible
	if f.selected < f.offset {
		f.offset = f.selected
	}
	if f.selected >= f.offset+visibleHeight {
		f.offset = f.selected - visibleHeight + 1
	}

	// File list
	if len(f.files) == 0 {
		noFilesStyle := lipgloss.NewStyle().
			Foreground(f.theme.Colors().Foreground).
			Italic(true).
			Padding(1, 2)
		sb.WriteString(noFilesStyle.Render("No .log files found in this directory"))
	} else {
		start := f.offset
		end := f.offset + visibleHeight
		if end > len(f.files) {
			end = len(f.files)
		}

		for i := start; i < end; i++ {
			var style lipgloss.Style
			if i == f.selected {
				style = lipgloss.NewStyle().
					Foreground(f.theme.Colors().Background).
					Background(f.theme.Colors().Info).
					Bold(true).
					Padding(0, 1).
					Width(f.width - 4)
			} else {
				style = lipgloss.NewStyle().
					Foreground(f.theme.Colors().Foreground).
					Padding(0, 1).
					Width(f.width - 4)
			}

			line := f.formatPath(f.files[i])
			if len(line) > f.width-6 {
				// Truncate from the middle for better readability
				if f.width > 20 {
					left := (f.width - 12) / 2
					right := (f.width - 12) - left
					line = line[:left] + "..." + line[len(line)-right:]
				} else {
					line = line[:f.width-9] + "..."
				}
			}
			sb.WriteString("  ")
			sb.WriteString(style.Render(line))
			sb.WriteString("\n")
		}
	}

	// Footer with help
	sb.WriteString("\n")
	footerStyle := lipgloss.NewStyle().
		Foreground(f.theme.Colors().Foreground).
		Faint(true).
		Width(f.width)

	footer := fmt.Sprintf("↑/k up • ↓/j down • g top • G bottom • enter select • q/esc quit | %d files", len(f.files))
	sb.WriteString(footerStyle.Render(footer))

	// Center the content
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(f.theme.Colors().Info).
		Padding(1, 2).
		Width(f.width)

	return lipgloss.Place(
		f.width,
		f.height,
		lipgloss.Center,
		lipgloss.Center,
		containerStyle.Render(sb.String()),
	)
}

func (f *FilePicker) selectNext() {
	if len(f.files) == 0 {
		return
	}
	if f.selected < len(f.files)-1 {
		f.selected++
	}
}

func (f *FilePicker) selectPrev() {
	if len(f.files) == 0 {
		return
	}
	if f.selected > 0 {
		f.selected--
	}
}

func (f *FilePicker) selectFirst() {
	if len(f.files) == 0 {
		return
	}
	f.selected = 0
}

func (f *FilePicker) selectLast() {
	if len(f.files) == 0 {
		return
	}
	f.selected = len(f.files) - 1
}

// SetTheme sets the theme.
func (f *FilePicker) SetTheme(theme theme.Theme) {
	f.theme = theme
}

// formatPath formats a file path for display (replaces home dir with ~)
func (f *FilePicker) formatPath(path string) string {
	// Get current working directory
	cwd, err := os.Getwd()
	if err == nil {
		// Try to make it relative to cwd
		rel, err := filepath.Rel(cwd, path)
		if err == nil && !strings.HasPrefix(rel, "..") {
			// If it's within cwd, use relative path
			if rel == "." {
				return filepath.Base(path)
			}
			return rel
		}
	}

	// Otherwise, try to replace home directory with ~
	home, err := os.UserHomeDir()
	if err == nil && strings.HasPrefix(path, home) {
		return "~" + path[len(home):]
	}

	return path
}
