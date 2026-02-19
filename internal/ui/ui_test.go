package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
	"github.com/ersanisk/sieve/pkg/logentry"
)

func TestLogView_NewLogView(t *testing.T) {
	theme := &MockTheme{}
	view := NewLogView(theme)

	if view.width != 0 {
		t.Errorf("NewLogView() width = %d, want 0", view.width)
	}

	if view.height != 0 {
		t.Errorf("NewLogView() height = %d, want 0", view.height)
	}

	if view.selected != 0 {
		t.Errorf("NewLogView() selected = %d, want 0", view.selected)
	}

	if !view.lineNumbers {
		t.Error("NewLogView() lineNumbers = false, want true")
	}
}

func TestLogView_SetEntries(t *testing.T) {
	theme := &MockTheme{}
	view := NewLogView(theme)

	entries := []logentry.Entry{
		{Level: logentry.Info, Message: "test 1"},
		{Level: logentry.Error, Message: "test 2"},
	}

	view.SetEntries(entries)

	if len(view.entries) != 2 {
		t.Errorf("SetEntries() got %d entries, want 2", len(view.entries))
	}
}

func TestLogView_ScrollDown(t *testing.T) {
	theme := &MockTheme{}
	view := NewLogView(theme)
	view.height = 10

	entries := make([]logentry.Entry, 20)
	for i := range entries {
		entries[i] = logentry.Entry{Level: logentry.Info, Message: "test"}
	}

	view.SetEntries(entries)
	view.ScrollDown(5)

	if view.selected != 5 {
		t.Errorf("ScrollDown(5) selected = %d, want 5", view.selected)
	}
}

func TestLogView_ScrollUp(t *testing.T) {
	theme := &MockTheme{}
	view := NewLogView(theme)

	entries := []logentry.Entry{
		{Level: logentry.Info, Message: "test 1"},
		{Level: logentry.Info, Message: "test 2"},
	}

	view.SetEntries(entries)
	view.ScrollDown(1)
	view.ScrollUp(1)

	if view.selected != 0 {
		t.Errorf("ScrollUp(1) selected = %d, want 0", view.selected)
	}
}

func TestLogView_ScrollToTop(t *testing.T) {
	theme := &MockTheme{}
	view := NewLogView(theme)

	entries := []logentry.Entry{
		{Level: logentry.Info, Message: "test 1"},
		{Level: logentry.Info, Message: "test 2"},
	}

	view.SetEntries(entries)
	view.ScrollDown(1)
	view.ScrollToTop()

	if view.selected != 0 {
		t.Errorf("ScrollToTop() selected = %d, want 0", view.selected)
	}

	if view.offset != 0 {
		t.Errorf("ScrollToTop() offset = %d, want 0", view.offset)
	}
}

func TestLogView_ScrollToBottom(t *testing.T) {
	theme := &MockTheme{}
	view := NewLogView(theme)
	view.height = 5

	entries := make([]logentry.Entry, 10)
	for i := range entries {
		entries[i] = logentry.Entry{Level: logentry.Info, Message: "test"}
	}

	view.SetEntries(entries)
	view.ScrollToBottom()

	if view.selected != 9 {
		t.Errorf("ScrollToBottom() selected = %d, want 9", view.selected)
	}
}

func TestLogView_GetSelected(t *testing.T) {
	theme := &MockTheme{}
	view := NewLogView(theme)

	entry, index := view.GetSelected()
	if index != -1 {
		t.Errorf("GetSelected() index = %d, want -1", index)
	}

	entries := []logentry.Entry{
		{Level: logentry.Info, Message: "test 1"},
		{Level: logentry.Error, Message: "test 2"},
	}
	view.SetEntries(entries)

	entry, index = view.GetSelected()
	if index != 0 {
		t.Errorf("GetSelected() index = %d, want 0", index)
	}
	if entry.Message != "test 1" {
		t.Errorf("GetSelected() entry = %v, want test 1", entry)
	}
}

func TestStatusBar_NewStatusBar(t *testing.T) {
	theme := &MockTheme{}
	bar := NewStatusBar(theme)

	if bar.mode != "view" {
		t.Errorf("NewStatusBar() mode = %s, want view", bar.mode)
	}

	if bar.following {
		t.Error("NewStatusBar() following = true, want false")
	}
}

func TestStatusBar_SetFilePath(t *testing.T) {
	theme := &MockTheme{}
	bar := NewStatusBar(theme)

	bar.SetFilePath("/path/to/file.log")
	if bar.filePath != "/path/to/file.log" {
		t.Errorf("SetFilePath() path = %s, want /path/to/file.log", bar.filePath)
	}
}

func TestStatusBar_SetTotalLines(t *testing.T) {
	theme := &MockTheme{}
	bar := NewStatusBar(theme)

	bar.SetTotalLines(100)
	if bar.totalLines != 100 {
		t.Errorf("SetTotalLines() total = %d, want 100", bar.totalLines)
	}
}

func TestSearchBar_NewSearchBar(t *testing.T) {
	theme := &MockTheme{}
	bar := NewSearchBar(theme)

	if bar.IsVisible() {
		t.Error("NewSearchBar() visible = true, want false")
	}

	if bar.IsFocused() {
		t.Error("NewSearchBar() focused = true, want false")
	}
}

func TestSearchBar_Show(t *testing.T) {
	theme := &MockTheme{}
	bar := NewSearchBar(theme)

	bar.Show()
	if !bar.IsVisible() {
		t.Error("Show() visible = false, want true")
	}

	if !bar.IsFocused() {
		t.Error("Show() focused = false, want true")
	}
}

func TestSearchBar_Hide(t *testing.T) {
	theme := &MockTheme{}
	bar := NewSearchBar(theme)
	bar.Show()
	bar.Hide()

	if bar.IsVisible() {
		t.Error("Hide() visible = true, want false")
	}

	if bar.IsFocused() {
		t.Error("Hide() focused = true, want false")
	}
}

func TestFilterBar_NewFilterBar(t *testing.T) {
	theme := &MockTheme{}
	bar := NewFilterBar(theme)

	if bar.IsVisible() {
		t.Error("NewFilterBar() visible = true, want false")
	}
}

func TestHelp_NewHelp(t *testing.T) {
	theme := &MockTheme{}
	help := NewHelp(theme)

	if help.IsVisible() {
		t.Error("NewHelp() visible = true, want false")
	}
}

func TestHelp_Show(t *testing.T) {
	theme := &MockTheme{}
	help := NewHelp(theme)
	help.Show()

	if !help.IsVisible() {
		t.Error("Show() visible = false, want true")
	}
}

func TestHelp_Hide(t *testing.T) {
	theme := &MockTheme{}
	help := NewHelp(theme)
	help.Show()
	help.Hide()

	if help.IsVisible() {
		t.Error("Hide() visible = true, want false")
	}
}

func TestSidebar_NewSidebar(t *testing.T) {
	theme := &MockTheme{}
	sidebar := NewSidebar(theme)

	if sidebar.IsVisible() {
		t.Error("NewSidebar() visible = true, want false")
	}
}

func TestTreeView_NewTreeView(t *testing.T) {
	theme := &MockTheme{}
	tree := NewTreeView(theme)

	if tree.IsVisible() {
		t.Error("NewTreeView() visible = true, want false")
	}
}

func TestDashboard_NewDashboard(t *testing.T) {
	theme := &MockTheme{}
	dash := NewDashboard(theme)

	if dash.IsVisible() {
		t.Error("NewDashboard() visible = true, want false")
	}
}

func TestDashboard_SetEntries(t *testing.T) {
	theme := &MockTheme{}
	dash := NewDashboard(theme)

	entries := []logentry.Entry{
		{Level: logentry.Info, Fields: map[string]any{"key": "value"}},
		{Level: logentry.Error},
	}
	dash.SetEntries(entries)

	if len(dash.entries) != 2 {
		t.Errorf("SetEntries() got %d entries, want 2", len(dash.entries))
	}

	if dash.levelCounts[logentry.Info] != 1 {
		t.Errorf("SetEntries() info count = %d, want 1", dash.levelCounts[logentry.Info])
	}

	if dash.levelCounts[logentry.Error] != 1 {
		t.Errorf("SetEntries() error count = %d, want 1", dash.levelCounts[logentry.Error])
	}
}

// MockTheme is a mock implementation of theme.Theme for testing.
type MockTheme struct{}

func (m *MockTheme) Name() string {
	return "mock"
}

func (m *MockTheme) Colors() theme.ThemeColors {
	return theme.ThemeColors{
		Debug:      "#0000ff",
		Info:       "#00ff00",
		Warn:       "#ffff00",
		Error:      "#ff0000",
		Fatal:      "#ff00ff",
		Timestamp:  "#888888",
		Key:        "#00ffff",
		Value:      "#ffffff",
		Background: "#000000",
		Foreground: "#ffffff",
		StatusBar:  "#333333",
		StatusText: "#ffffff",
		Border:     "#666666",
		Highlight:  "#444444",
	}
}

func (m *MockTheme) LevelStyle(level logentry.Level) lipgloss.Style {
	return lipgloss.NewStyle()
}

func (m *MockTheme) TimestampStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (m *MockTheme) KeyStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (m *MockTheme) ValueStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (m *MockTheme) StatusBarStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (m *MockTheme) BorderStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (m *MockTheme) HighlightStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (m *MockTheme) ErrorStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func (m *MockTheme) InfoStyle() lipgloss.Style {
	return lipgloss.NewStyle()
}

func TestLogDetail_Show(t *testing.T) {
	theme := &MockTheme{}
	detail := NewLogDetail(theme)
	detail.SetSize(100, 30)

	entry := logentry.Entry{
		Level:   logentry.Info,
		Message: "Test Message",
		Fields: map[string]any{
			"b_field": "value2",
			"a_field": "value1",
		},
		Raw: `{"level":"info","msg":"Test Message","a_field":"value1","b_field":"value2"}`,
	}

	detail.Show(entry)

	if !detail.IsVisible() {
		t.Error("Show() visible = false, want true")
	}

	// Verify content through viewport
	content := detail.viewport.View()

	// Check for Raw JSON header
	if !contains(content, "Raw JSON:") {
		t.Error("View() missing 'Raw JSON:' header")
	}

	// Check for Raw JSON content
	if !contains(content, entry.Raw) {
		t.Error("View() missing raw JSON content")
	}

	// Check sorting: a_field should appear before b_field
	idxA := indexOf(content, "a_field")
	idxB := indexOf(content, "b_field")

	if idxA == -1 || idxB == -1 {
		t.Error("View() missing fields")
	} else if idxA > idxB {
		t.Errorf("Fields not sorted: a_field index %d > b_field index %d", idxA, idxB)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func indexOf(s, substr string) int {
	return strings.Index(s, substr)
}
