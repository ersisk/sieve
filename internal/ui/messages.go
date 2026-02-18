package ui

import (
	"github.com/ersanisk/sieve/pkg/logentry"
)

// LogSelectedMsg is sent when a log entry is selected.
type LogSelectedMsg struct {
	Entry logentry.Entry
	Index  int
}

// ScrollUpMsg is sent to scroll up.
type ScrollUpMsg struct {
	Amount int
}

// ScrollDownMsg is sent to scroll down.
type ScrollDownMsg struct {
	Amount int
}

// ScrollToTopMsg is sent to scroll to top.
type ScrollToTopMsg struct{}

// ScrollToLineMsg is sent to scroll to a specific line.
type ScrollToLineMsg struct {
	Line int
}

// SearchInputMsg is sent when search input changes.
type SearchInputMsg struct {
	Query string
}

// SearchSubmitMsg is sent when search is submitted.
type SearchSubmitMsg struct{}

// SearchClearMsg is sent to clear search.
type SearchClearMsg struct{}

// SearchNextMsg is sent to go to next search result.
type SearchNextMsg struct{}

// SearchPrevMsg is sent to go to previous search result.
type SearchPrevMsg struct{}

// FilterInputMsg is sent when filter input changes.
type FilterInputMsg struct {
	Expression string
}

// FilterSubmitMsg is sent when filter is submitted.
type FilterSubmitMsg struct{}

// FilterClearMsg is sent to clear filter.
type FilterClearMsg struct{}

// ToggleHelpMsg is sent to toggle help overlay.
type ToggleHelpMsg struct{}

// ToggleSidebarMsg is sent to toggle sidebar.
type ToggleSidebarMsg struct{}

// ToggleDashboardMsg is sent to toggle dashboard.
type ToggleDashboardMsg struct{}

// ToggleFollowMsg is sent to toggle follow mode.
type ToggleFollowMsg struct{}

// SetLevelFilterMsg is sent to set level filter.
type SetLevelFilterMsg struct {
	Level logentry.Level
}

// FilterSetPresetMsg is sent to set a filter preset.
type FilterSetPresetMsg struct {
	Preset string
}

// ResizeMsg is sent when window is resized.
type ResizeMsg struct {
	Width  int
	Height int
}

// EnterModeMsg is sent to enter a mode.
type EnterModeMsg struct {
	Mode string
}

// ExitModeMsg is sent to exit current mode.
type ExitModeMsg struct{}

// QuitMsg is sent to quit.
type QuitMsg struct{}

// RefreshMsg is sent to refresh view.
type RefreshMsg struct{}

// FileLoadedMsg is sent when a file is loaded.
type FileLoadedMsg struct {
	Path    string
	Entries []logentry.Entry
}

// EntryFocusedMsg is sent when an entry is focused.
type EntryFocusedMsg struct {
	Entry logentry.Entry
	Index int
}

// ExpandNodeMsg is sent to expand a tree node.
type ExpandNodeMsg struct {
	Path string
}

// CollapseNodeMsg is sent to collapse a tree node.
type CollapseNodeMsg struct{}

// ExpandAllMsg is sent to expand all tree nodes.
type ExpandAllMsg struct{}

// CollapseAllMsg is sent to collapse all tree nodes.
type CollapseAllMsg struct{}

// CopyEntryMsg is sent to copy entry to clipboard.
type CopyEntryMsg struct {
	Entry logentry.Entry
}

// BookmarkMsg is sent to bookmark an entry.
type BookmarkMsg struct {
	Index int
}

// GoToBookmarkMsg is sent to go to a bookmark.
type GoToBookmarkMsg struct {
	Index int
}

// TickMsg is sent periodically for animations.
type TickMsg struct {
	Time time.Time
}

// LoadingStartedMsg is sent when loading starts.
type LoadingStartedMsg struct {
	Message string
}

// LoadingFinishedMsg is sent when loading finishes.
type LoadingFinishedMsg struct {
	Count int
}

// ErrorMsg is sent when an error occurs.
type ErrorMsg struct {
	Error error
}

// LoadFileMsg is sent to load a file.
type LoadFileMsg struct {
	Path string
}
}
