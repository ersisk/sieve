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
	theme         theme.Theme
	files         []string
	filteredFiles []string
	selected      int
	visible       bool
	width         int
	height        int
	offset        int
	searchMode    bool
	searchQuery   string
}

// FileSelectedMsg is sent when a file is selected.
type FileSelectedMsg struct {
	Path string
}

// NewFilePicker creates a new file picker.
func NewFilePicker(theme theme.Theme) FilePicker {
	return FilePicker{
		theme:         theme,
		files:         []string{},
		filteredFiles: []string{},
		selected:      0,
		visible:       false,
		width:         80,
		height:        24,
		offset:        0,
		searchMode:    false,
		searchQuery:   "",
	}
}

// SetFiles sets the list of files to display.
func (f *FilePicker) SetFiles(files []string) {
	f.files = files
	f.filteredFiles = files
	f.selected = 0
	f.offset = 0
	f.searchQuery = ""
	f.searchMode = false
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
		if f.searchMode {
			return f.handleSearchInput(msg)
		}
		return f.handleNormalInput(msg)
	}

	return f, nil
}

func (f *FilePicker) handleSearchInput(msg tea.KeyMsg) (*FilePicker, tea.Cmd) {
	switch msg.Type {
	case tea.KeyEscape:
		f.searchMode = false
		f.searchQuery = ""
		f.filterFiles()
	case tea.KeyEnter:
		f.searchMode = false
		if len(f.filteredFiles) > 0 {
			return f, func() tea.Msg {
				return FileSelectedMsg{Path: f.filteredFiles[f.selected]}
			}
		}
	case tea.KeyBackspace:
		if len(f.searchQuery) > 0 {
			f.searchQuery = f.searchQuery[:len(f.searchQuery)-1]
			f.filterFiles()
		}
	case tea.KeyRunes:
		f.searchQuery += string(msg.Runes)
		f.filterFiles()
	case tea.KeySpace:
		f.searchQuery += " "
		f.filterFiles()
	}
	return f, nil
}

func (f *FilePicker) handleNormalInput(msg tea.KeyMsg) (*FilePicker, tea.Cmd) {
	switch msg.String() {
	case "/":
		f.searchMode = true
		f.searchQuery = ""
	case "up", "k":
		f.selectPrev()
	case "down", "j":
		f.selectNext()
	case "g":
		f.selectFirst()
	case "G":
		f.selectLast()
	case "enter":
		if len(f.filteredFiles) > 0 {
			return f, func() tea.Msg {
				return FileSelectedMsg{Path: f.filteredFiles[f.selected]}
			}
		}
	case "q", "esc", "ctrl+c":
		return f, func() tea.Msg {
			return QuitMsg{}
		}
	}
	return f, nil
}

func (f *FilePicker) filterFiles() {
	if f.searchQuery == "" {
		f.filteredFiles = f.files
		f.selected = 0
		f.offset = 0
		return
	}

	query := strings.ToLower(f.searchQuery)
	filtered := make([]string, 0)
	for _, file := range f.files {
		if strings.Contains(strings.ToLower(file), query) {
			filtered = append(filtered, file)
		}
	}
	f.filteredFiles = filtered
	f.selected = 0
	f.offset = 0
}

// View renders the file picker.
func (f *FilePicker) View() string {
	if !f.visible {
		return ""
	}

	colors := f.theme.Colors()

	// Calculate container dimensions
	containerWidth := f.width - 8
	if containerWidth < 40 {
		containerWidth = 40
	}
	if containerWidth > 100 {
		containerWidth = 100
	}

	contentWidth := containerWidth - 6

	var content strings.Builder

	// Header with icon and title
	headerStyle := lipgloss.NewStyle().
		Foreground(colors.Info).
		Bold(true)

	titleText := " Select Log File"
	if len(f.filteredFiles) > 0 {
		titleText = fmt.Sprintf(" Select Log File (%d/%d)", f.selected+1, len(f.filteredFiles))
	}
	content.WriteString(headerStyle.Render(titleText))
	content.WriteString("\n")

	// Separator line
	separatorStyle := lipgloss.NewStyle().
		Foreground(colors.Info).
		Faint(true)
	content.WriteString(separatorStyle.Render(strings.Repeat("─", contentWidth)))
	content.WriteString("\n")

	// Search bar (always show, highlight when active)
	searchBarStyle := lipgloss.NewStyle().
		Padding(0, 1).
		Width(contentWidth)

	if f.searchMode {
		searchBarStyle = searchBarStyle.
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colors.Warn)
		searchText := f.searchQuery
		if searchText == "" {
			searchText = "type to search..."
		}
		content.WriteString(searchBarStyle.Render("/ " + searchText + "█"))
	} else {
		searchBarStyle = searchBarStyle.
			Foreground(colors.Foreground).
			Faint(true)
		if f.searchQuery != "" {
			content.WriteString(searchBarStyle.Render("/ " + f.searchQuery + " (press / to edit)"))
		} else {
			content.WriteString(searchBarStyle.Render("Press / to search"))
		}
	}
	content.WriteString("\n\n")

	// Calculate visible range for file list
	visibleHeight := f.height - 14
	if visibleHeight < 3 {
		visibleHeight = 3
	}

	// Adjust offset to keep selected item visible
	if f.selected < f.offset {
		f.offset = f.selected
	}
	if f.selected >= f.offset+visibleHeight {
		f.offset = f.selected - visibleHeight + 1
	}

	// File list
	if len(f.filteredFiles) == 0 {
		emptyStyle := lipgloss.NewStyle().
			Foreground(colors.Warn).
			Italic(true).
			Padding(1, 0)
		if f.searchQuery != "" {
			content.WriteString(emptyStyle.Render("No files match your search"))
		} else {
			content.WriteString(emptyStyle.Render("No .log files found"))
		}
		content.WriteString("\n")
	} else {
		start := f.offset
		end := f.offset + visibleHeight
		if end > len(f.filteredFiles) {
			end = len(f.filteredFiles)
		}

		for i := start; i < end; i++ {
			isSelected := i == f.selected

			// File icon and path
			icon := "  "
			if isSelected {
				icon = " "
			}

			line := f.formatPath(f.filteredFiles[i])
			maxLen := contentWidth - 6
			if len(line) > maxLen {
				// Truncate from the middle for better readability
				if maxLen > 20 {
					left := (maxLen - 3) / 2
					right := maxLen - 3 - left
					line = line[:left] + "..." + line[len(line)-right:]
				} else {
					line = line[:maxLen-3] + "..."
				}
			}

			var itemStyle lipgloss.Style
			if isSelected {
				itemStyle = lipgloss.NewStyle().
					Foreground(colors.Background).
					Background(colors.Info).
					Bold(true).
					Width(contentWidth).
					Padding(0, 1)
			} else {
				itemStyle = lipgloss.NewStyle().
					Foreground(colors.Foreground).
					Width(contentWidth).
					Padding(0, 1)
			}

			content.WriteString(itemStyle.Render(icon + line))
			content.WriteString("\n")
		}

		// Scroll indicator
		if len(f.filteredFiles) > visibleHeight {
			scrollInfo := fmt.Sprintf("  %d more above", f.offset)
			scrollInfoBottom := fmt.Sprintf("  %d more below", len(f.filteredFiles)-end)

			scrollStyle := lipgloss.NewStyle().
				Foreground(colors.Foreground).
				Faint(true).
				Italic(true)

			if f.offset > 0 {
				content.WriteString(scrollStyle.Render(scrollInfo))
				content.WriteString("\n")
			}
			if end < len(f.filteredFiles) {
				content.WriteString(scrollStyle.Render(scrollInfoBottom))
				content.WriteString("\n")
			}
		}
	}

	// Footer with keybindings
	content.WriteString("\n")
	footerStyle := lipgloss.NewStyle().
		Foreground(colors.Foreground).
		Faint(true)

	var footerText string
	if f.searchMode {
		footerText = "type to filter  enter confirm  esc cancel"
	} else {
		footerText = "j/k navigate  / search  enter select  q quit"
	}
	content.WriteString(footerStyle.Render(footerText))

	// Container with border
	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colors.Info).
		Padding(1, 2).
		Width(containerWidth)

	// Center the container on screen
	return lipgloss.Place(
		f.width,
		f.height,
		lipgloss.Center,
		lipgloss.Center,
		containerStyle.Render(content.String()),
	)
}

func (f *FilePicker) selectNext() {
	if len(f.filteredFiles) == 0 {
		return
	}
	if f.selected < len(f.filteredFiles)-1 {
		f.selected++
	}
}

func (f *FilePicker) selectPrev() {
	if len(f.filteredFiles) == 0 {
		return
	}
	if f.selected > 0 {
		f.selected--
	}
}

func (f *FilePicker) selectFirst() {
	if len(f.filteredFiles) == 0 {
		return
	}
	f.selected = 0
}

func (f *FilePicker) selectLast() {
	if len(f.filteredFiles) == 0 {
		return
	}
	f.selected = len(f.filteredFiles) - 1
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
