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
