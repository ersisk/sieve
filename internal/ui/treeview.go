package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/ersanisk/sieve/internal/theme"
)

// TreeView displays a hierarchical tree view of JSON data.
type TreeView struct {
	visible  bool
	width    int
	height   int
	data     any
	expanded map[string]bool
	theme    theme.Theme
	selected int
}

// TreeNode represents a node in the tree.
type TreeNode struct {
	Key    string
	Value  any
	IsLeaf bool
	Level  int
	Path   string
}

// NewTreeView creates a new TreeView.
func NewTreeView(theme theme.Theme) TreeView {
	return TreeView{
		visible:  false,
		theme:    theme,
		expanded: make(map[string]bool),
	}
}

// Show shows the tree view.
func (m *TreeView) Show() {
	m.visible = true
}

// Hide hides the tree view.
func (m *TreeView) Hide() {
	m.visible = false
}

// IsVisible returns true if the tree view is visible.
func (m *TreeView) IsVisible() bool {
	return m.visible
}

// SetData sets the data to display.
func (m *TreeView) SetData(data any) {
	m.data = data
	m.expanded = make(map[string]bool)
}

// GetData returns the current data.
func (m *TreeView) GetData() any {
	return m.data
}

// SetSize sets the dimensions of the tree view.
func (m *TreeView) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// GetSize returns the dimensions of the tree view.
func (m *TreeView) GetSize() (int, int) {
	return m.width, m.height
}

// SetTheme sets the theme.
func (m *TreeView) SetTheme(theme theme.Theme) {
	m.theme = theme
}

// ToggleNode toggles expansion of a node.
func (m *TreeView) ToggleNode(path string) {
	if _, ok := m.expanded[path]; ok {
		delete(m.expanded, path)
	} else {
		m.expanded[path] = true
	}
}

// IsExpanded returns true if a node is expanded.
func (m *TreeView) IsExpanded(path string) bool {
	return m.expanded[path]
}

// ExpandAll expands all nodes.
func (m *TreeView) ExpandAll() {
	m.expanded = make(map[string]bool)
	nodes := m.buildTree(m.data, "", 0)
	for _, node := range nodes {
		if !node.IsLeaf {
			m.expanded[node.Path] = true
		}
	}
}

// CollapseAll collapses all nodes.
func (m *TreeView) CollapseAll() {
	m.expanded = make(map[string]bool)
}

// View renders the tree view.
func (m TreeView) View() string {
	if !m.visible || m.data == nil {
		return ""
	}

	containerStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(m.theme.Colors().Border).
		Width(m.width).
		Height(m.height)

	content := m.renderContent()
	return containerStyle.Render(content)
}

// renderContent renders the tree content.
func (m TreeView) renderContent() string {
	nodes := m.buildTree(m.data, "", 0)

	if len(nodes) == 0 {
		style := lipgloss.NewStyle().
			Foreground(m.theme.Colors().Foreground).
			Italic(true).
			Align(lipgloss.Center).
			Height(m.height - 2)
		return style.Render("No data to display")
	}

	var builder strings.Builder
	for _, node := range nodes {
		builder.WriteString(m.renderNode(node))
	}

	return builder.String()
}

// buildTree builds a tree from data.
func (m TreeView) buildTree(data any, path string, level int) []TreeNode {
	var nodes []TreeNode

	switch v := data.(type) {
	case map[string]any:
		keys := m.getSortedKeys(v)
		for _, key := range keys {
			nodePath := path + "." + key
			node := TreeNode{
				Key:    key,
				Value:  v[key],
				IsLeaf: m.isLeaf(v[key]),
				Level:  level,
				Path:   nodePath,
			}
			nodes = append(nodes, node)
		}
	case []any:
		for i, val := range v {
			nodePath := fmt.Sprintf("%s[%d]", path, i)
			node := TreeNode{
				Key:    fmt.Sprintf("[%d]", i),
				Value:  val,
				IsLeaf: m.isLeaf(val),
				Level:  level,
				Path:   nodePath,
			}
			nodes = append(nodes, node)
		}
	}

	return nodes
}

// renderNode renders a single tree node.
func (m TreeView) renderNode(node TreeNode) string {
	keyStyle := m.theme.KeyStyle()
	valueStyle := m.theme.ValueStyle()

	indent := strings.Repeat("  ", node.Level)

	if node.IsLeaf {
		keyText := keyStyle.Render(fmt.Sprintf("%s├─ %s", indent, node.Key))
		valueText := valueStyle.Render(fmt.Sprintf(": %v", node.Value))
		return fmt.Sprintf("%s %s\n", keyText, valueText)
	}

	isExpanded := m.expanded[node.Path]
	prefix := "+"
	if isExpanded {
		prefix = "-"
	}

	keyText := keyStyle.Render(fmt.Sprintf("%s%s %s", indent, prefix, node.Key))
	return fmt.Sprintf("%s\n", keyText)
}

// isLeaf checks if a value is a leaf node.
func (m TreeView) isLeaf(value any) bool {
	switch value.(type) {
	case map[string]any, []any:
		return false
	}
	return true
}

// getSortedKeys returns sorted keys from a map.
func (m TreeView) getSortedKeys(mmap map[string]any) []string {
	keys := make([]string, 0, len(mmap))
	for key := range mmap {
		keys = append(keys, key)
	}

	return keys
}

// GetSelected returns the selected node.
func (m *TreeView) GetSelected() *TreeNode {
	nodes := m.buildTree(m.data, "", 0)
	if m.selected >= 0 && m.selected < len(nodes) {
		return &nodes[m.selected]
	}
	return nil
}

// SetSelected sets the selected node.
func (m *TreeView) SetSelected(index int) {
	if index < 0 {
		m.selected = 0
	} else {
		m.selected = index
	}
}

// GetSelectedPath returns the path of the selected node.
func (m *TreeView) GetSelectedPath() string {
	if node := m.GetSelected(); node != nil {
		return node.Path
	}
	return ""
}

// ToggleSelected toggles expansion of the selected node.
func (m *TreeView) ToggleSelected() {
	if node := m.GetSelected(); node != nil && !node.IsLeaf {
		m.ToggleNode(node.Path)
	}
}
