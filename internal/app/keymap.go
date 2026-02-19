package app

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// KeyMap defines keyboard shortcuts for the application.
type KeyMap struct {
	Quit            keyBinding
	ForceQuit       keyBinding
	ScrollUp        keyBinding
	ScrollDown      keyBinding
	ScrollLeft      keyBinding
	ScrollRight     keyBinding
	ScrollPageUp    keyBinding
	ScrollPageDown  keyBinding
	ScrollToTop     keyBinding
	ScrollToBottom  keyBinding
	Search          keyBinding
	SearchNext      keyBinding
	SearchPrev      keyBinding
	Filter          keyBinding
	ClearFilter     keyBinding
	ToggleHelp      keyBinding
	ToggleSidebar   keyBinding
	ToggleDashboard keyBinding
	ToggleFollow    keyBinding
	LevelDebug      keyBinding
	LevelInfo       keyBinding
	LevelWarn       keyBinding
	LevelError      keyBinding
	LevelFatal      keyBinding
	LevelNone       keyBinding
	Expand          keyBinding
	Collapse        keyBinding
	Copy            keyBinding
	RefreshFile     keyBinding
	ToggleSort      keyBinding
}

// keyBinding represents a single keyboard binding.
type keyBinding struct {
	key   tea.Key
	help  string
	style lipgloss.Style
}

// DefaultKeyMap returns the default keyboard mapping.
func DefaultKeyMap() KeyMap {
	keyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Bold(true)

	return KeyMap{
		Quit:            keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'q'}}, help: "Quit", style: keyStyle},
		ForceQuit:       keyBinding{key: tea.Key{Type: tea.KeyCtrlC}, help: "Force Quit (Ctrl+C)", style: keyStyle},
		ScrollUp:        keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'k'}}, help: "Scroll Up", style: keyStyle},
		ScrollDown:      keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'j'}}, help: "Scroll Down", style: keyStyle},
		ScrollLeft:      keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'h'}}, help: "Scroll Left", style: keyStyle},
		ScrollRight:     keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'l'}}, help: "Scroll Right", style: keyStyle},
		ScrollPageUp:    keyBinding{key: tea.Key{Type: tea.KeyPgUp}, help: "Page Up (PgUp)", style: keyStyle},
		ScrollPageDown:  keyBinding{key: tea.Key{Type: tea.KeyPgDown}, help: "Page Down (PgDn)", style: keyStyle},
		ScrollToTop:     keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'g'}}, help: "Go to Top", style: keyStyle},
		ScrollToBottom:  keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'G'}}, help: "Go to Bottom", style: keyStyle},
		Search:          keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'/'}}, help: "Search", style: keyStyle},
		SearchNext:      keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'n'}}, help: "Next Search Result", style: keyStyle},
		SearchPrev:      keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'N'}}, help: "Previous Search Result", style: keyStyle},
		Filter:          keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'f'}}, help: "Filter", style: keyStyle},
		ClearFilter:     keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'F'}}, help: "Clear Filter", style: keyStyle},
		ToggleHelp:      keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'?'}}, help: "Toggle Help", style: keyStyle},
		ToggleSidebar:   keyBinding{key: tea.Key{Type: tea.KeyTab}, help: "Toggle Sidebar (Tab)", style: keyStyle},
		ToggleDashboard: keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'d'}}, help: "Toggle Dashboard", style: keyStyle},
		ToggleFollow:    keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'F'}}, help: "Toggle Follow Mode", style: keyStyle},
		LevelDebug:      keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'1'}}, help: "Filter Debug", style: keyStyle},
		LevelInfo:       keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'2'}}, help: "Filter Info", style: keyStyle},
		LevelWarn:       keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'3'}}, help: "Filter Warn", style: keyStyle},
		LevelError:      keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'4'}}, help: "Filter Error", style: keyStyle},
		LevelFatal:      keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'5'}}, help: "Filter Fatal", style: keyStyle},
		LevelNone:       keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'0'}}, help: "Clear Level Filter", style: keyStyle},
		Expand:          keyBinding{key: tea.Key{Type: tea.KeyEnter}, help: "Expand Entry (Enter)", style: keyStyle},
		Collapse:        keyBinding{key: tea.Key{Type: tea.KeyEsc}, help: "Collapse/Close (Esc)", style: keyStyle},
		Copy:            keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'c'}}, help: "Copy Entry", style: keyStyle},
		RefreshFile:     keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'R'}}, help: "Refresh File (Shift+R)", style: keyStyle},
		ToggleSort:      keyBinding{key: tea.Key{Type: tea.KeyRunes, Runes: []rune{'r'}}, help: "Toggle Sort", style: keyStyle},
	}
}

// ShortHelp returns a short help string for keybindings.
func (k KeyMap) ShortHelp() string {
	return "q: quit | ?: help | /: search | f: filter | d: dashboard | r: sort | R: refresh"
}

// FullHelp returns a full help string for keybindings.
func (k KeyMap) FullHelp() string {
	return `Navigation:
  j/↓         Scroll down
  k/↑         Scroll up
  g            Go to top
  G            Go to bottom
  PgDn/Space  Page down
  PgUp         Page up

Search & Filter:
  /            Search
  n/N          Next/Prev search result
  f            Filter
  F            Clear filter

Level Filter:
  1-5          Filter by level (debug/error/etc)
  0            Clear level filter

View & Actions:
  Enter        Expand/Collapse entry
  Tab          Toggle sidebar
  d            Toggle dashboard
  Esc          Close overlay / exit mode

File & Program:
  r            Toggle Sort (asc/desc)
  R            Refresh file
  q            Quit
  ?            Toggle this help`
}
