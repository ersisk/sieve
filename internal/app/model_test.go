package app

import (
	"sort"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ersanisk/sieve/pkg/logentry"
)

func TestToggleSort(t *testing.T) {
	model := Model{
		sortOrder: SortAsc,
	}

	if model.sortOrder != SortAsc {
		t.Errorf("Expected initial sort order to be ascending, got %v", model.sortOrder)
	}

	model.toggleSort()
	if model.sortOrder != SortDesc {
		t.Errorf("Expected sort order to be descending after toggle, got %v", model.sortOrder)
	}

	model.toggleSort()
	if model.sortOrder != SortAsc {
		t.Errorf("Expected sort order to be ascending after second toggle, got %v", model.sortOrder)
	}
}

func TestSortLogic(t *testing.T) {
	now := time.Now()
	later := now.Add(time.Hour)

	entries := []logentry.Entry{
		{Timestamp: later, Message: "second"},
		{Timestamp: now, Message: "first"},
		{Timestamp: now.Add(time.Minute * 30), Message: "middle"},
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	if entries[0].Message != "first" {
		t.Errorf("Expected first message 'first' when sorted ascending, got '%s'", entries[0].Message)
	}
	if entries[1].Message != "middle" {
		t.Errorf("Expected second message 'middle' when sorted ascending, got '%s'", entries[1].Message)
	}
	if entries[2].Message != "second" {
		t.Errorf("Expected third message 'second' when sorted ascending, got '%s'", entries[2].Message)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})

	if entries[0].Message != "second" {
		t.Errorf("Expected first message 'second' when sorted descending, got '%s'", entries[0].Message)
	}
	if entries[1].Message != "middle" {
		t.Errorf("Expected second message 'middle' when sorted descending, got '%s'", entries[1].Message)
	}
	if entries[2].Message != "first" {
		t.Errorf("Expected third message 'first' when sorted descending, got '%s'", entries[2].Message)
	}
}

func TestFilterBarESCKey(t *testing.T) {
	model := NewModel("", "kanagawa", false)

	// Filter bar'ı aç
	model.filterBar.Show()

	if !model.filterBar.IsFocused() {
		t.Error("FilterBar should be focused after Show()")
	}

	if !model.filterBar.IsVisible() {
		t.Error("FilterBar should be visible after Show()")
	}

	// ESC tuşu gönder
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := model.Update(escMsg)
	model = newModel.(Model)

	if model.filterBar.IsFocused() {
		t.Error("FilterBar should not be focused after ESC")
	}

	if model.filterBar.IsVisible() {
		t.Error("FilterBar should not be visible after ESC")
	}

	if model.mode != "view" {
		t.Errorf("Mode should be 'view' after ESC, got '%s'", model.mode)
	}
}

func TestSearchBarESCKey(t *testing.T) {
	model := NewModel("", "kanagawa", false)

	// Search bar'ı aç
	model.searchBar.Show()

	if !model.searchBar.IsFocused() {
		t.Error("SearchBar should be focused after Show()")
	}

	// ESC tuşu gönder
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := model.Update(escMsg)
	model = newModel.(Model)

	if model.searchBar.IsFocused() {
		t.Error("SearchBar should not be focused after ESC")
	}

	if model.mode != "view" {
		t.Errorf("Mode should be 'view' after ESC, got '%s'", model.mode)
	}
}
func TestFilterClearWithESC(t *testing.T) {
	model := NewModel("", "kanagawa", false)

	// Önce bazı entries ekle
	model.entries = []logentry.Entry{
		{Level: logentry.Info, Message: "test1"},
		{Level: logentry.Error, Message: "test2"},
	}
	model.filtered = model.entries

	// Level filter uygula
	model.levelFilter = logentry.Error
	model, _ = model.applyLevelFilter()

	if len(model.filtered) != 1 {
		t.Errorf("Expected 1 filtered entry, got %d", len(model.filtered))
	}

	// ESC tuşu gönder
	escMsg := tea.KeyMsg{Type: tea.KeyEsc}
	newModel, _ := model.Update(escMsg)
	model = newModel.(Model)

	// Filter temizlenmeli
	if model.levelFilter != logentry.Unknown {
		t.Error("Level filter should be cleared after ESC")
	}

	if len(model.filtered) != 2 {
		t.Errorf("Expected all entries after filter clear, got %d", len(model.filtered))
	}
}
