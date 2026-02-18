package app

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/filter"
	"github.com/ersanisk/sieve/internal/theme"
	"github.com/ersanisk/sieve/internal/ui"
	"github.com/ersanisk/sieve/pkg/logentry"
)

type Model struct {
	logView       ui.LogView
	statusBar     ui.StatusBar
	searchBar     ui.SearchBar
	filterBar     ui.FilterBar
	sidebar       ui.Sidebar
	help          ui.Help
	treeView      ui.TreeView
	dashboard     ui.Dashboard
	logDetail     *ui.LogDetail
	keyMap        KeyMap
	theme         theme.Theme
	entries       []logentry.Entry
	filtered      []logentry.Entry
	selectedEntry logentry.Entry
	mode          string
	loading       bool
	loadingMsg    string
	filePath      string
	followMode    bool
	levelFilter   logentry.Level
	filter        *filter.CompiledFilter
	filterExpr    string
}

func NewModel(filePath string, themeName string) Model {
	theme := getTheme(themeName)

	return Model{
		logView:     ui.NewLogView(theme),
		statusBar:   ui.NewStatusBar(theme),
		searchBar:   ui.NewSearchBar(theme),
		filterBar:   ui.NewFilterBar(theme),
		sidebar:     ui.NewSidebar(theme),
		help:        ui.NewHelp(theme),
		treeView:    ui.NewTreeView(theme),
		dashboard:   ui.NewDashboard(theme),
		logDetail:   ui.NewLogDetail(theme),
		keyMap:      DefaultKeyMap(),
		theme:       theme,
		entries:     []logentry.Entry{},
		filtered:    []logentry.Entry{},
		mode:        "view",
		loading:     false,
		filePath:    filePath,
		followMode:  false,
		levelFilter: logentry.Unknown,
	}
}

func getTheme(name string) theme.Theme {
	if t := theme.Get(name); t != nil {
		return t
	}
	return theme.Get("default")
}

func levelToInt(level logentry.Level) int {
	switch level {
	case logentry.Debug:
		return 10
	case logentry.Info:
		return 30
	case logentry.Warn:
		return 40
	case logentry.Error:
		return 50
	case logentry.Fatal:
		return 60
	default:
		return 0
	}
}

func (m Model) Init() tea.Cmd {
	if m.filePath != "" {
		return tea.Batch(tickCmd(), loadFileCmd(m.filePath))
	}
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.filterBar.IsFocused() && msg.Type == tea.KeyEnter {
			m.filterBar.Hide()
			return m, tea.Batch(tickCmd(), tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return ui.FilterSubmitMsg{}
			}))
		}
		if m.filterBar.IsFocused() && msg.Type == tea.KeyEsc {
			m.filterBar.Hide()
			return m, tickCmd()
		}
		if m.searchBar.IsFocused() && msg.Type == tea.KeyEnter {
			m.searchBar.Hide()
			return m, tea.Batch(tickCmd(), tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return ui.SearchSubmitMsg{}
			}))
		}
		if m.searchBar.IsFocused() && msg.Type == tea.KeyEsc {
			m.searchBar.Hide()
			return m, tickCmd()
		}
		if m.searchBar.IsFocused() {
			var cmd tea.Cmd
			m.searchBar, cmd = m.searchBar.Update(msg)
			return m, tea.Batch(cmd, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return ui.SearchInputMsg{Query: m.searchBar.GetValue()}
			}))
		}
		if m.filterBar.IsFocused() {
			var cmd tea.Cmd
			m.filterBar, cmd = m.filterBar.Update(msg)
			return m, tea.Batch(cmd, tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
				return ui.FilterInputMsg{Expression: m.filterBar.GetValue()}
			}))
		}
		return m.handleKey(msg)
	case tea.MouseMsg:
		return m, nil
	case tea.WindowSizeMsg:
		m.handleResize(msg)
		return m, tickCmd()
	case ui.RefreshMsg:
		return m, tickCmd()
	case ui.LoadingStartedMsg:
		m.loading = true
		m.loadingMsg = msg.Message
		return m, tickCmd()
	case ui.LoadingFinishedMsg:
		m.loading = false
		m.loadingMsg = ""
		m.statusBar.SetTotalLines(msg.Count)
		return m, tickCmd()
	case ui.QuitMsg:
		return m, tea.Quit
	case ui.ErrorMsg:
		fmt.Fprintf(os.Stderr, "Error: %v\n", msg.Error)
		return m, tea.Quit
	case ui.FileLoadedMsg:
		m.entries = msg.Entries
		m.filtered = msg.Entries
		m.logView.SetEntries(msg.Entries)
		m.statusBar.SetFilePath(msg.Path)
		m.statusBar.SetTotalLines(len(msg.Entries))
		m.loading = false
		return m, tickCmd()
	case ui.SearchInputMsg:
		m.searchBar.SetValue(msg.Query)
		return m, tickCmd()
	case ui.SearchSubmitMsg:
		return m, tickCmd()
	case ui.FilterInputMsg:
		m.filterBar.SetValue(msg.Expression)
		return m, tickCmd()
	case ui.FilterSubmitMsg:
		return m.applyFilter(m.filterBar.GetValue())
	case ui.SetLevelFilterMsg:
		m.levelFilter = msg.Level
		return m, tickCmd()
	case ui.ToggleHelpMsg:
		m.help.Show()
		m.mode = "help"
		return m, tickCmd()
	case ui.ToggleSidebarMsg:
		if m.sidebar.IsVisible() {
			m.sidebar.Hide()
		} else {
			m.sidebar.Show()
		}
		return m, tickCmd()
	case ui.ToggleDashboardMsg:
		if m.dashboard.IsVisible() {
			m.dashboard.Hide()
		} else {
			m.dashboard.Show()
		}
		return m, tickCmd()
	case ui.ToggleFollowMsg:
		m.followMode = !m.followMode
		m.statusBar.SetFollowing(m.followMode)
		return m, tickCmd()
	case ui.ScrollUpMsg:
		if m.help.IsVisible() {
			m.help, _ = m.help.Update(msg)
		} else {
			m.logView.ScrollUp(msg.Amount)
			m.updateSelectedEntry()
		}
		return m, tickCmd()
	case ui.ScrollDownMsg:
		if m.help.IsVisible() {
			m.help, _ = m.help.Update(msg)
		} else {
			m.logView.ScrollDown(msg.Amount)
			m.updateSelectedEntry()
		}
		return m, tickCmd()
	case ui.ScrollToTopMsg:
		m.logView.ScrollToTop()
		m.updateSelectedEntry()
		return m, tickCmd()
	case ui.ScrollToBottomMsg:
		m.logView.ScrollToBottom()
		m.updateSelectedEntry()
		return m, tickCmd()
	}

	return m, cmd
}

func (m Model) View() string {
	if m.loading {
		return m.renderLoading()
	}

	if m.help.IsVisible() {
		return m.help.View()
	}

	if m.dashboard.IsVisible() {
		return m.dashboard.View()
	}

	if m.logDetail.IsVisible() {
		return m.renderMain() + "\n" + m.logDetail.View()
	}

	main := m.renderMain()
	return main
}

func (m Model) handleKey(msg tea.KeyMsg) (Model, tea.Cmd) {
	if m.logDetail.IsVisible() {
		if msg.Type == tea.KeyEsc {
			m.logDetail.Hide()
		}
		return m, tickCmd()
	}

	if m.help.IsVisible() {
		if msg.Type == tea.KeyEsc {
			m.help.Hide()
			m.mode = "view"
			return m, tickCmd()
		}
		switch msg.String() {
		case m.keyMap.ScrollUp.key.String():
			return m, tea.Batch(tickCmd(), func() tea.Msg { return ui.ScrollUpMsg{Amount: 1} })
		case m.keyMap.ScrollDown.key.String():
			return m, tea.Batch(tickCmd(), func() tea.Msg { return ui.ScrollDownMsg{Amount: 1} })
		}
		return m, tickCmd()
	}

	switch msg.String() {
	case m.keyMap.Quit.key.String():
		return m, tea.Quit
	case m.keyMap.ForceQuit.key.String():
		return m, tea.Quit
	case m.keyMap.ScrollUp.key.String():
		m.logView.ScrollUpOne()
		m.updateSelectedEntry()
		return m, tickCmd()
	case m.keyMap.ScrollDown.key.String():
		m.logView.ScrollDownOne()
		m.updateSelectedEntry()
		return m, tickCmd()
	case m.keyMap.ScrollToTop.key.String():
		m.logView.ScrollToTop()
		m.updateSelectedEntry()
		return m, tickCmd()
	case m.keyMap.ScrollToBottom.key.String():
		m.logView.ScrollToBottom()
		m.updateSelectedEntry()
		return m, tickCmd()
	case m.keyMap.ScrollPageUp.key.String():
		m.logView.ScrollPageUp()
		m.updateSelectedEntry()
		return m, tickCmd()
	case m.keyMap.ScrollPageDown.key.String():
		m.logView.ScrollPageDown()
		m.updateSelectedEntry()
		return m, tickCmd()
	case m.keyMap.Search.key.String():
		m.searchBar.Show()
		m.mode = "search"
		return m, tickCmd()
	case m.keyMap.Filter.key.String():
		m.filterBar.Show()
		m.mode = "filter"
		return m, tickCmd()
	case m.keyMap.ToggleHelp.key.String():
		m.help.Show()
		m.mode = "help"
		return m, tickCmd()
	case m.keyMap.ToggleSidebar.key.String():
		if m.sidebar.IsVisible() {
			m.sidebar.Hide()
		} else {
			m.sidebar.Show()
		}
		return m, tickCmd()
	case m.keyMap.ToggleDashboard.key.String():
		if m.dashboard.IsVisible() {
			m.dashboard.Hide()
		} else {
			m.dashboard.Show()
		}
		return m, tickCmd()
	case m.keyMap.ToggleFollow.key.String():
		m.followMode = !m.followMode
		m.statusBar.SetFollowing(m.followMode)
		return m, tickCmd()
	case m.keyMap.LevelDebug.key.String():
		m.levelFilter = logentry.Debug
		return m, tickCmd()
	case m.keyMap.LevelInfo.key.String():
		m.levelFilter = logentry.Info
		return m, tickCmd()
	case m.keyMap.LevelWarn.key.String():
		m.levelFilter = logentry.Warn
		return m, tickCmd()
	case m.keyMap.LevelError.key.String():
		m.levelFilter = logentry.Error
		return m, tickCmd()
	case m.keyMap.LevelFatal.key.String():
		m.levelFilter = logentry.Fatal
		return m, tickCmd()
	case m.keyMap.LevelNone.key.String():
		m.levelFilter = logentry.Unknown
		return m, tickCmd()
	case m.keyMap.Refresh.key.String():
		return m, tea.Batch(tickCmd(), loadFileCmd(m.filePath), tickCmd())
	case m.keyMap.Expand.key.String():
		m.logDetail.Show(m.selectedEntry)
		return m, tickCmd()
	}

	return m, tickCmd()
}

func (m *Model) handleResize(msg tea.WindowSizeMsg) {
	width, height := msg.Width, msg.Height

	m.logView.SetSize(width, height-2)
	m.statusBar.SetSize(width, 2)
	m.searchBar.SetSize(width, 3)
	m.filterBar.SetSize(width, 3)
	m.help.SetSize(width, height)
	m.dashboard.SetSize(width/2, height/2)
	m.sidebar.SetSize(width/3, height-2)
	m.treeView.SetSize(width/3, height-2)
	m.logDetail.SetSize(width, height)

	m.statusBar.SetFilePath(m.filePath)
	m.statusBar.SetTotalLines(len(m.filtered))
}

func (m *Model) updateSelectedEntry() {
	entry, index := m.logView.GetSelected()
	if index >= 0 {
		m.selectedEntry = entry
		m.sidebar.SetEntry(entry)
		m.statusBar.SetSelected(index)
	}
}

func (m Model) renderMain() string {
	width, height := m.logView.GetSize()

	if width == 0 || height == 0 {
		return ""
	}

	m.logView.SetSize(width, height-2)
	m.statusBar.SetSize(width, 2)

	var result string

	if m.searchBar.IsVisible() {
		result += m.searchBar.View() + "\n"
	}

	if m.filterBar.IsVisible() {
		result += m.filterBar.View() + "\n"
	}

	result += m.logView.View()
	result += "\n"
	result += m.statusBar.View()

	return result
}

func (m Model) renderLoading() string {
	style := lipgloss.NewStyle().
		Foreground(m.theme.Colors().Foreground).
		Align(lipgloss.Center).
		Bold(true)

	if m.loadingMsg != "" {
		return style.Render(m.loadingMsg)
	}
	return style.Render("Loading...")
}

func (m Model) applyFilter(expr string) (Model, tea.Cmd) {
	if expr == "" {
		m.filter = nil
		m.filterExpr = ""
		m.filtered = m.entries
		m.logView.SetEntries(m.filtered)
		m.statusBar.SetTotalLines(len(m.filtered))
		return m, tickCmd()
	}

	parsed, err := filter.Parse(expr)
	if err != nil {
		m.statusBar.SetError(fmt.Sprintf("Filter error: %v", err))
		return m, tickCmd()
	}

	compiled, err := filter.Compile(parsed)
	if err != nil {
		m.statusBar.SetError(fmt.Sprintf("Filter error: %v", err))
		return m, tickCmd()
	}

	m.filter = compiled
	m.filterExpr = expr

	var filtered []logentry.Entry
	for _, entry := range m.entries {
		matches, err := compiled.Evaluate(entry)
		if err == nil && matches {
			filtered = append(filtered, entry)
		}
	}

	m.filtered = filtered
	m.logView.SetEntries(m.filtered)
	m.statusBar.SetTotalLines(len(m.filtered))
	m.statusBar.SetInfo(fmt.Sprintf("Filtered: %d/%d entries", len(filtered), len(m.entries)))

	return m, tickCmd()
}
