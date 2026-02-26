package ui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ersanisk/sieve/internal/theme"
)

func TestFilePickerKeyboardNavigation(t *testing.T) {
	th := theme.Get("default")
	picker := NewFilePicker(th)
	picker.Show()
	picker.SetFiles([]string{"file1.log", "file2.log", "file3.log"})

	// Test down navigation
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedPicker, _ := picker.Update(keyMsg)
	if updatedPicker.selected != 1 {
		t.Errorf("Expected selected = 1, got %d", updatedPicker.selected)
	}

	// Test up navigation
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 0 {
		t.Errorf("Expected selected = 0, got %d", updatedPicker.selected)
	}

	// Test go to bottom
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 2 {
		t.Errorf("Expected selected = 2, got %d", updatedPicker.selected)
	}

	// Test go to top
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 0 {
		t.Errorf("Expected selected = 0, got %d", updatedPicker.selected)
	}
}

func TestFilePickerEnterSelection(t *testing.T) {
	th := theme.Get("default")
	picker := NewFilePicker(th)
	picker.Show()
	picker.SetFiles([]string{"file1.log", "file2.log"})

	// Move to second file
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedPicker, _ := picker.Update(keyMsg)

	// Press enter
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := updatedPicker.Update(keyMsg)

	// Check that command returns FileSelectedMsg
	if cmd == nil {
		t.Fatal("Expected command to be returned")
	}

	msg := cmd()
	fileMsg, ok := msg.(FileSelectedMsg)
	if !ok {
		t.Fatalf("Expected FileSelectedMsg, got %T", msg)
	}

	if fileMsg.Path != "file2.log" {
		t.Errorf("Expected file2.log, got %s", fileMsg.Path)
	}
}

func TestFilePickerSearchModeNavigation(t *testing.T) {
	th := theme.Get("default")
	picker := NewFilePicker(th)
	picker.Show()
	picker.SetFiles([]string{"file1.log", "file2.log", "file3.log", "file4.log"})

	// Enter search mode
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	updatedPicker, _ := picker.Update(keyMsg)

	// Test down navigation with arrow key (empty filter)
	keyMsg = tea.KeyMsg{Type: tea.KeyDown}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 1 {
		t.Errorf("Expected selected = 1, got %d", updatedPicker.selected)
	}

	// Test up navigation with arrow key
	keyMsg = tea.KeyMsg{Type: tea.KeyUp}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 0 {
		t.Errorf("Expected selected = 0, got %d", updatedPicker.selected)
	}

	// Test j key navigation
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 1 {
		t.Errorf("Expected selected = 1 after 'j', got %d", updatedPicker.selected)
	}

	// Test k key navigation
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 0 {
		t.Errorf("Expected selected = 0 after 'k', got %d", updatedPicker.selected)
	}

	// Type a filter character (this should reset selection to 0)
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'1'}}
	updatedPicker, _ = updatedPicker.Update(keyMsg)

	// Verify filter worked
	if len(updatedPicker.filteredFiles) != 1 {
		t.Errorf("Expected 1 filtered file, got %d", len(updatedPicker.filteredFiles))
	}

	// Navigate again after filtering
	keyMsg = tea.KeyMsg{Type: tea.KeyDown}
	updatedPicker, _ = updatedPicker.Update(keyMsg)

	// Verify still in search mode
	if !updatedPicker.searchMode {
		t.Error("Expected to remain in search mode after navigation")
	}
}

func TestFilePickerSearchModeJump(t *testing.T) {
	th := theme.Get("default")
	picker := NewFilePicker(th)
	picker.Show()
	picker.SetFiles([]string{"alpha.log", "beta.log", "gamma.log", "delta.log"})

	// Enter search mode
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	updatedPicker, _ := picker.Update(keyMsg)

	// Test go to bottom with G
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'G'}}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 3 {
		t.Errorf("Expected selected = 3 after 'G', got %d", updatedPicker.selected)
	}

	// Test go to top with g
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 0 {
		t.Errorf("Expected selected = 0 after 'g', got %d", updatedPicker.selected)
	}

	// Verify still in search mode
	if !updatedPicker.searchMode {
		t.Error("Expected to remain in search mode after jump")
	}
}

func TestFilePickerSearchModeSelection(t *testing.T) {
	th := theme.Get("default")
	picker := NewFilePicker(th)
	picker.Show()
	picker.SetFiles([]string{"app1.log", "app2.log", "app3.log"})

	// Enter search mode
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'/'}}
	updatedPicker, _ := picker.Update(keyMsg)

	// Navigate to second file
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedPicker, _ = updatedPicker.Update(keyMsg)
	if updatedPicker.selected != 1 {
		t.Errorf("Expected selected = 1, got %d", updatedPicker.selected)
	}

	// Press enter to confirm and select
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	_, cmd := updatedPicker.Update(keyMsg)

	// Check that command returns FileSelectedMsg with the correct file
	if cmd == nil {
		t.Fatal("Expected command to be returned")
	}

	msg := cmd()
	fileMsg, ok := msg.(FileSelectedMsg)
	if !ok {
		t.Fatalf("Expected FileSelectedMsg, got %T", msg)
	}

	if fileMsg.Path != "app2.log" {
		t.Errorf("Expected app2.log, got %s", fileMsg.Path)
	}

	// Verify search mode was exited
	if updatedPicker.searchMode {
		t.Error("Expected search mode to be exited after enter")
	}
}
