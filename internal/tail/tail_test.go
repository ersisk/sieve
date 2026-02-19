package tail

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewReader(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	reader, err := NewReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewReader() failed: %v", err)
	}
	defer reader.Close()

	if reader.path != tmpFile.Name() {
		t.Errorf("Reader.path = %v, want %v", reader.path, tmpFile.Name())
	}

	if reader.IsClosed() {
		t.Error("Reader is closed after creation")
	}
}

func TestReader_ReadNew(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testLines := []string{"line 1", "line 2", "line 3"}
	for _, line := range testLines {
		_, err := tmpFile.WriteString(line + "\n")
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}
	tmpFile.Close()

	reader, err := NewReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewReader() failed: %v", err)
	}
	defer reader.Close()

	lines, err := reader.ReadNew()
	if err != nil {
		t.Fatalf("ReadNew() failed: %v", err)
	}

	if len(lines) != len(testLines) {
		t.Errorf("ReadNew() got %d lines, want %d", len(lines), len(testLines))
	}

	for i, line := range lines {
		if line != testLines[i] {
			t.Errorf("ReadNew()[%d] = %v, want %v", i, line, testLines[i])
		}
	}

	lines, err = reader.ReadNew()
	if err != nil {
		t.Fatalf("ReadNew() failed: %v", err)
	}

	if len(lines) != 0 {
		t.Errorf("ReadNew() got %d new lines, want 0", len(lines))
	}
}

func TestReader_ReadAll(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testLines := []string{"line 1", "line 2", "line 3"}
	for _, line := range testLines {
		_, err := tmpFile.WriteString(line + "\n")
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}
	tmpFile.Close()

	reader, err := NewReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewReader() failed: %v", err)
	}
	defer reader.Close()

	lines, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll() failed: %v", err)
	}

	if len(lines) != len(testLines) {
		t.Errorf("ReadAll() got %d lines, want %d", len(lines), len(testLines))
	}

	lines2, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("ReadAll() failed: %v", err)
	}

	if len(lines2) != len(testLines) {
		t.Errorf("ReadAll() got %d lines on second call, want %d", len(lines2), len(testLines))
	}
}

func TestReader_Position(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("line 1\nline 2\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	reader, err := NewReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewReader() failed: %v", err)
	}
	defer reader.Close()

	if reader.Position() != 0 {
		t.Errorf("Position() = %d, want 0", reader.Position())
	}

	_, err = reader.ReadNew()
	if err != nil {
		t.Fatalf("ReadNew() failed: %v", err)
	}

	if reader.Position() == 0 {
		t.Error("Position() = 0 after ReadNew, want > 0")
	}
}

func TestReader_SetPosition(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	testLines := []string{"line 1", "line 2", "line 3"}
	for _, line := range testLines {
		_, err := tmpFile.WriteString(line + "\n")
		if err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
	}
	tmpFile.Close()

	reader, err := NewReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewReader() failed: %v", err)
	}
	defer reader.Close()

	_, err = reader.ReadNew()
	if err != nil {
		t.Fatalf("ReadNew() failed: %v", err)
	}

	initialPos := reader.Position()
	if err := reader.SetPosition(0); err != nil {
		t.Fatalf("SetPosition() failed: %v", err)
	}

	if reader.Position() != 0 {
		t.Errorf("Position() = %d after SetPosition(0), want 0", reader.Position())
	}

	if err := reader.SetPosition(initialPos); err != nil {
		t.Fatalf("SetPosition() failed: %v", err)
	}

	if reader.Position() != initialPos {
		t.Errorf("Position() = %d, want %d", reader.Position(), initialPos)
	}
}

func TestReader_ResetPosition(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("line 1\nline 2\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	reader, err := NewReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewReader() failed: %v", err)
	}
	defer reader.Close()

	_, err = reader.ReadNew()
	if err != nil {
		t.Fatalf("ReadNew() failed: %v", err)
	}

	if err := reader.ResetPosition(); err != nil {
		t.Fatalf("ResetPosition() failed: %v", err)
	}

	if reader.Position() != 0 {
		t.Errorf("Position() = %d after ResetPosition(), want 0", reader.Position())
	}
}

func TestReader_Close(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	reader, err := NewReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewReader() failed: %v", err)
	}

	if err := reader.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	if !reader.IsClosed() {
		t.Error("Reader is not closed after Close()")
	}

	if err := reader.Close(); err != nil {
		t.Errorf("Close() failed on already closed reader: %v", err)
	}
}

func TestReader_Reopen(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("line 1\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	reader, err := NewReader(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewReader() failed: %v", err)
	}

	if err := reader.Close(); err != nil {
		t.Fatalf("Close() failed: %v", err)
	}

	if !reader.IsClosed() {
		t.Fatal("Reader should be closed")
	}

	if err := reader.Reopen(); err != nil {
		t.Fatalf("Reopen() failed: %v", err)
	}

	if reader.IsClosed() {
		t.Error("Reader is still closed after Reopen()")
	}

	if err := reader.Close(); err != nil {
		t.Fatalf("Close() failed after Reopen(): %v", err)
	}
}

func TestNewWatcher(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)
	tmpFile.Close()

	watcher, err := NewWatcher(tmpPath)
	if err != nil {
		t.Fatalf("NewWatcher() failed: %v", err)
	}
	defer func() {
		_ = watcher.Stop()
	}()

	if watcher.path != tmpPath {
		t.Errorf("Watcher.path = %v, want %v", watcher.path, tmpPath)
	}

	if watcher.IsRunning() {
		t.Error("Watcher is running after creation, should not be")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_ = watcher.Start(ctx)

	if !watcher.IsRunning() {
		t.Error("Watcher is not running after Start()")
	}
}

func TestWatcher_Start(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	watcher, err := NewWatcher(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewWatcher() failed: %v", err)
	}
	defer watcher.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	linesCh := watcher.Start(ctx)

	if linesCh == nil {
		t.Fatal("Start() returned nil channel")
	}

	if !watcher.IsRunning() {
		t.Error("Watcher is not running after Start()")
	}
}

func TestWatcher_Watch(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.log")

	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	linesCh, err := Watch(ctx, tmpFile)
	if err != nil {
		t.Fatalf("Watch() failed: %v", err)
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		f, err := os.OpenFile(tmpFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			t.Errorf("Failed to open file: %v", err)
			return
		}
		defer f.Close()

		for i := 0; i < 3; i++ {
			_, err := f.WriteString("test line\n")
			if err != nil {
				t.Errorf("Failed to write line: %v", err)
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	totalLines := 0
	timeout := time.After(2 * time.Second)

loop:
	for {
		select {
		case lines := <-linesCh:
			totalLines += len(lines)
			if totalLines >= 3 {
				break loop
			}
		case <-timeout:
			t.Fatalf("Timeout waiting for lines (got %d, want 3)", totalLines)
		case <-ctx.Done():
			t.Fatalf("Context cancelled")
		}
	}

	if totalLines < 3 {
		t.Errorf("Got %d lines, want at least 3", totalLines)
	}
}

func TestWatcher_WatchFromEnd(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.log")

	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	for i := 0; i < 5; i++ {
		_, err := f.WriteString("initial line\n")
		if err != nil {
			t.Fatalf("Failed to write initial lines: %v", err)
		}
	}
	f.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	linesCh, err := WatchFromEnd(ctx, tmpFile)
	if err != nil {
		t.Fatalf("WatchFromEnd() failed: %v", err)
	}

	go func() {
		time.Sleep(200 * time.Millisecond)
		f, err := os.OpenFile(tmpFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			t.Errorf("Failed to open file: %v", err)
			return
		}
		defer f.Close()

		_, err = f.WriteString("new line\n")
		if err != nil {
			t.Errorf("Failed to write line: %v", err)
		}
	}()

	var allLines []string
	timeout := time.After(2 * time.Second)

	for {
		select {
		case lines := <-linesCh:
			allLines = append(allLines, lines...)
			foundNewLine := false
			for _, line := range allLines {
				if line == "new line" {
					foundNewLine = true
					break
				}
			}
			if foundNewLine {
				return
			}
		case <-timeout:
			t.Errorf("Timeout waiting for 'new line'. Got lines: %v", allLines)
			return
		case <-ctx.Done():
			return
		}
	}
}

func TestWatcher_Stop(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	watcher, err := NewWatcher(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewWatcher() failed: %v", err)
	}

	if err := watcher.Stop(); err != nil {
		t.Fatalf("Stop() failed: %v", err)
	}

	if watcher.IsRunning() {
		t.Error("Watcher is still running after Stop()")
	}
}

func TestWatcher_PollImmediately(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.log")

	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	f.Close()

	watcher, err := NewWatcher(tmpFile)
	if err != nil {
		t.Fatalf("NewWatcher() failed: %v", err)
	}
	defer func() {
		_ = watcher.Stop()
	}()

	lines, err := watcher.PollImmediately()
	if err != nil {
		t.Fatalf("PollImmediately() failed: %v", err)
	}

	if len(lines) != 0 {
		t.Errorf("PollImmediately() got %d lines, want 0", len(lines))
	}

	f, err = os.OpenFile(tmpFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}
	if _, err := f.WriteString("test line\n"); err != nil {
		t.Fatalf("Failed to write to file: %v", err)
	}
	f.Close()

	lines, err = watcher.PollImmediately()
	if err != nil {
		t.Fatalf("PollImmediately() failed: %v", err)
	}

	if len(lines) != 1 {
		t.Errorf("PollImmediately() got %d lines, want 1", len(lines))
	}

	if lines[0] != "test line" {
		t.Errorf("PollImmediately()[0] = %v, want 'test line'", lines[0])
	}
}

func TestNewReader_NonExistentFile(t *testing.T) {
	_, err := NewReader("/nonexistent/file.log")
	if err == nil {
		t.Error("NewReader() succeeded for non-existent file, expected error")
	}
}

func TestNewWatcher_NonExistentFile(t *testing.T) {
	_, err := NewWatcher("/nonexistent/file.log")
	if err == nil {
		t.Error("NewWatcher() succeeded for non-existent file, expected error")
	}
}

func TestWatcher_RemoveChannel(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	watcher, err := NewWatcher(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewWatcher() failed: %v", err)
	}
	defer watcher.Stop()

	ctx1, cancel1 := context.WithCancel(context.Background())
	defer cancel1()

	ctx2, cancel2 := context.WithCancel(context.Background())
	defer cancel2()

	_ = watcher.Start(ctx1)
	_ = watcher.Start(ctx2)

	watcher.RemoveChannel(ctx1)

	if err := watcher.Stop(); err != nil {
		t.Fatalf("Stop() failed: %v", err)
	}
}

func TestWatcher_FileSize(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString("test content\n")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	watcher, err := NewWatcher(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewWatcher() failed: %v", err)
	}
	defer watcher.Stop()

	size, err := watcher.FileSize()
	if err != nil {
		t.Fatalf("FileSize() failed: %v", err)
	}

	if size == 0 {
		t.Error("FileSize() returned 0 for non-empty file")
	}
}

func TestWatcher_CheckFileExists(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "test-*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tmpFile.Close()

	watcher, err := NewWatcher(tmpFile.Name())
	if err != nil {
		t.Fatalf("NewWatcher() failed: %v", err)
	}
	defer watcher.Stop()

	if !watcher.CheckFileExists() {
		t.Error("CheckFileExists() returned false for existing file")
	}

	os.Remove(tmpFile.Name())

	if watcher.CheckFileExists() {
		t.Error("CheckFileExists() returned true for deleted file")
	}
}
