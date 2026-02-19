package app

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ersanisk/sieve/internal/ui"
)

func TestFindLogFilesCmd(t *testing.T) {
	// Create test directory structure
	tmpDir := t.TempDir()

	// Create test files
	_ = os.WriteFile(filepath.Join(tmpDir, "root.log"), []byte("root"), 0644)
	_ = os.MkdirAll(filepath.Join(tmpDir, "subdir1"), 0755)
	_ = os.WriteFile(filepath.Join(tmpDir, "subdir1", "app.log"), []byte("app"), 0644)
	_ = os.MkdirAll(filepath.Join(tmpDir, "subdir1", "subdir2"), 0755)
	_ = os.WriteFile(filepath.Join(tmpDir, "subdir1", "subdir2", "deep.log"), []byte("deep"), 0644)

	// Also create a non-log file
	_ = os.WriteFile(filepath.Join(tmpDir, "readme.txt"), []byte("readme"), 0644)

	// Execute command
	cmd := findLogFilesCmd(tmpDir)
	msg := cmd()

	// Check result
	logFilesMsg, ok := msg.(ui.LogFilesFoundMsg)
	if !ok {
		t.Fatalf("Expected LogFilesFoundMsg, got %T", msg)
	}

	// Should find 3 log files
	if len(logFilesMsg.Files) != 3 {
		t.Errorf("Expected 3 log files, got %d: %v", len(logFilesMsg.Files), logFilesMsg.Files)
	}

	// Check that files are found
	foundFiles := make(map[string]bool)
	for _, file := range logFilesMsg.Files {
		foundFiles[filepath.Base(file)] = true
	}

	expectedFiles := []string{"root.log", "app.log", "deep.log"}
	for _, expected := range expectedFiles {
		if !foundFiles[expected] {
			t.Errorf("Expected to find %s, but it was not in the list", expected)
		}
	}
}
