package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ersanisk/sieve/internal/ui"
)

func TestModelFilePickerKeyHandling(t *testing.T) {
	// Create model without file path (should trigger file picker)
	model := NewModel("", "default", false)

	// Initialize - this returns a batch of commands
	initCmd := model.Init()
	if initCmd == nil {
		t.Fatal("Expected init command")
	}

	// Manually trigger ShowFilePickerMsg
	showPickerMsg := ui.ShowFilePickerMsg{Directory: "."}

	// Process ShowFilePickerMsg
	updatedModel, _ := model.Update(showPickerMsg)
	model = updatedModel.(Model)

	if !model.filePicker.IsVisible() {
		t.Error("File picker should be visible after ShowFilePickerMsg")
	}

	// Simulate receiving log files
	logFilesMsg := ui.LogFilesFoundMsg{
		Files: []string{"/tmp/test1.log", "/tmp/test2.log", "/tmp/test3.log"},
	}
	updatedModel, _ = model.Update(logFilesMsg)
	model = updatedModel.(Model)

	// Test keyboard navigation - down key
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(Model)

	// Test keyboard navigation - up key
	keyMsg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	updatedModel, _ = model.Update(keyMsg)
	model = updatedModel.(Model)

	// Test enter key - should select file
	keyMsg = tea.KeyMsg{Type: tea.KeyEnter}
	updatedModel, enterCmd := model.Update(keyMsg)
	model = updatedModel.(Model)

	if enterCmd == nil {
		t.Fatal("Expected command after enter key")
	}

	// Execute the command
	cmdMsg := enterCmd()
	fileSelectedMsg, ok := cmdMsg.(ui.FileSelectedMsg)
	if !ok {
		t.Fatalf("Expected FileSelectedMsg after enter, got %T", cmdMsg)
	}

	// Process FileSelectedMsg
	updatedModel, _ = model.Update(fileSelectedMsg)
	model = updatedModel.(Model)

	if model.filePicker.IsVisible() {
		t.Error("File picker should be hidden after file selection")
	}
}

func TestModelFilePickerQuitHandling(t *testing.T) {
	model := NewModel("", "default", false)

	// Show file picker
	model.filePicker.Show()
	model.mode = "filepicker"
	model.filePicker.SetFiles([]string{"/tmp/test.log"})

	// Press q to quit
	keyMsg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}}
	updatedModel, cmd := model.Update(keyMsg)
	model = updatedModel.(Model)

	if cmd == nil {
		t.Fatal("Expected command after q key")
	}

	msg := cmd()
	_, ok := msg.(ui.QuitMsg)
	if !ok {
		t.Fatalf("Expected QuitMsg after q, got %T", msg)
	}
}
