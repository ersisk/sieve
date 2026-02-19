package app

import (
	"fmt"
	"io"
	"os"
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
		defer file.Close()

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
		defer f.Close()

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
