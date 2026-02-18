package app

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/ersanisk/sieve/internal/filter"
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

// searchCmd performs search on entries.
func searchCmd(query string, entries []logentry.Entry) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(10 * time.Millisecond)
		return ui.RefreshMsg{}
	}
}

// filterCmd applies filter to entries.
func filterCmd(expr string, entries []logentry.Entry) tea.Cmd {
	return func() tea.Msg {
		_, err := filter.Parse(expr)
		if err != nil {
			return ui.ErrorMsg{Error: fmt.Errorf("failed to parse filter: %w", err)}
		}

		time.Sleep(10 * time.Millisecond)
		return ui.RefreshMsg{}
	}
}

// tickCmd returns a tick command for animations.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return ui.TickMsg{Time: t}
	})
}
