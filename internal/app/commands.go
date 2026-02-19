package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ersanisk/sieve/internal/parser"
	"github.com/ersanisk/sieve/internal/ui"
	"github.com/ersanisk/sieve/pkg/logentry"
)

// loadFileCmd loads a file and parses its contents.
func loadFileCmd(path string) tea.Cmd {
	return func() tea.Msg {
		parser := parser.NewParser()

		file, err := os.Open(path)
		if err != nil {
			return ui.ErrorMsg{Error: fmt.Errorf("failed to open file: %w", err)}
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				// Log error but don't override the main error
				fmt.Fprintf(os.Stderr, "Warning: failed to close file: %v\n", closeErr)
			}
		}()

		entries, err := parser.ParseLines(file)
		if err != nil {
			return ui.ErrorMsg{Error: fmt.Errorf("failed to parse file: %w", err)}
		}

		return ui.FileLoadedMsg{Path: path, Entries: entries}
	}
}

// tickCmd returns a tick command for animations.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return ui.TickMsg{Time: t}
	})
}

// clearInfoCmd returns a command that clears info/error messages after a delay.
func clearInfoCmd(delay time.Duration) tea.Cmd {
	return tea.Tick(delay, func(t time.Time) tea.Msg {
		return ui.ClearInfoMsg{}
	})
}

// NewLinesMsg is sent when new lines are appended to the followed file.
type NewLinesMsg struct {
	Entries []logentry.Entry
}

// followCmd reads new lines appended to the file since lastSize.
func followCmd(path string, lastSize int64, p *parser.Parser) tea.Cmd {
	return func() tea.Msg {
		info, err := os.Stat(path)
		if err != nil {
			return nil
		}
		if info.Size() <= lastSize {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer func() {
			_ = f.Close()
		}()

		if _, err := f.Seek(lastSize, io.SeekStart); err != nil {
			return nil
		}

		entries, err := p.ParseLines(f)
		if err != nil || len(entries) == 0 {
			return nil
		}

		return NewLinesMsg{Entries: entries}
	}
}

// findLogFilesCmd searches for .log files in the given directory.
func findLogFilesCmd(dir string) tea.Cmd {
	return func() tea.Msg {
		if dir == "" {
			dir = "."
		}

		// Convert to absolute path for consistency
		absDir, err := filepath.Abs(dir)
		if err != nil {
			absDir = dir
		}

		var logFiles []string

		err = filepath.Walk(absDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors and continue
			}

			// Skip directories
			if info.IsDir() {
				// Don't recurse into hidden directories or common non-log directories
				if strings.HasPrefix(info.Name(), ".") ||
					info.Name() == "node_modules" ||
					info.Name() == "vendor" {
					return filepath.SkipDir
				}
				return nil
			}

			// Check if file has .log extension
			if strings.HasSuffix(strings.ToLower(info.Name()), ".log") {
				// Always use absolute path for consistency
				absPath, err := filepath.Abs(path)
				if err != nil {
					absPath = path
				}
				logFiles = append(logFiles, absPath)
			}

			return nil
		})

		if err != nil {
			return ui.ErrorMsg{Error: fmt.Errorf("failed to search for log files: %w", err)}
		}

		// Sort files alphabetically
		sort.Strings(logFiles)

		return ui.LogFilesFoundMsg{Files: logFiles}
	}
}
