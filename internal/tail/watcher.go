package tail

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches a file for changes and sends new lines through a channel.
type Watcher struct {
	path      string
	reader    *Reader
	watcher   *fsnotify.Watcher
	channels  map[context.Context]chan<- []string
	mu        sync.RWMutex
	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	isRunning bool
}

// NewWatcher creates a new Watcher for the specified file.
func NewWatcher(path string) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create fsnotify watcher: %w", err)
	}

	reader, err := NewReader(path)
	if err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	w := &Watcher{
		path:     path,
		reader:   reader,
		watcher:  watcher,
		channels: make(map[context.Context]chan<- []string),
		ctx:      ctx,
		cancel:   cancel,
	}

	if err := w.watcher.Add(path); err != nil {
		reader.Close()
		watcher.Close()
		return nil, fmt.Errorf("failed to add file to watcher: %w", err)
	}

	return w, nil
}

// Start starts the watcher and returns a channel that receives new lines.
// The channel is closed when the context is cancelled.
func (w *Watcher) Start(ctx context.Context) <-chan []string {
	linesCh := make(chan []string, 100)

	w.mu.Lock()
	w.channels[ctx] = linesCh
	w.mu.Unlock()

	if !w.isRunning {
		w.isRunning = true
		w.wg.Add(1)
		go w.watchLoop()
	}

	return linesCh
}

// watchLoop is the main watcher loop.
func (w *Watcher) watchLoop() {
	defer w.wg.Done()

	pollInterval := 250 * time.Millisecond
	lastCheck := time.Now()

	for {
		select {
		case <-w.ctx.Done():
			return

		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			if event.Op&fsnotify.Write == fsnotify.Write {
				w.sendNewLines()
				lastCheck = time.Now()
			}

		case _, ok := <-w.watcher.Errors:
			if !ok {
				return
			}

		case <-time.After(pollInterval):
			if time.Since(lastCheck) > time.Second {
				w.sendNewLines()
				lastCheck = time.Now()
			}
		}
	}
}

// sendNewLines reads new lines from the file and sends them to all channels.
func (w *Watcher) sendNewLines() {
	lines, err := w.reader.ReadNew()
	if err != nil {
		return
	}

	if len(lines) == 0 {
		return
	}

	w.mu.RLock()
	defer w.mu.RUnlock()

	for ctx, ch := range w.channels {
		select {
		case ch <- lines:
		case <-ctx.Done():
		case <-w.ctx.Done():
			return
		}
	}
}

// Stop stops the watcher.
func (w *Watcher) Stop() error {
	w.cancel()
	w.wg.Wait()

	if err := w.watcher.Close(); err != nil {
		return fmt.Errorf("failed to close fsnotify watcher: %w", err)
	}

	if err := w.reader.Close(); err != nil {
		return fmt.Errorf("failed to close reader: %w", err)
	}

	return nil
}

// Path returns the file path being watched.
func (w *Watcher) Path() string {
	return w.path
}

// IsRunning returns true if the watcher is running.
func (w *Watcher) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.isRunning
}

// GetReader returns the reader used by the watcher.
func (w *Watcher) GetReader() *Reader {
	return w.reader
}

// PollImmediately polls the file immediately for new lines.
func (w *Watcher) PollImmediately() ([]string, error) {
	return w.reader.ReadNew()
}

// RemoveChannel removes a channel from the watcher.
func (w *Watcher) RemoveChannel(ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()
	delete(w.channels, ctx)
}

// Watch reads a file and continuously sends new lines through the channel.
// This is a convenience function that creates a watcher and starts it.
func Watch(ctx context.Context, path string) (<-chan []string, error) {
	watcher, err := NewWatcher(path)
	if err != nil {
		return nil, err
	}

	linesCh := watcher.Start(ctx)

	go func() {
		<-ctx.Done()
		_ = watcher.Stop()
	}()

	return linesCh, nil
}

// WatchFromEnd watches a file starting from the end.
func WatchFromEnd(ctx context.Context, path string) (<-chan []string, error) {
	watcher, err := NewWatcher(path)
	if err != nil {
		return nil, err
	}

	if err := watcher.GetReader().ResetPosition(); err != nil {
		_ = watcher.Stop()
		return nil, err
	}

	linesCh := watcher.Start(ctx)

	go func() {
		<-ctx.Done()
		_ = watcher.Stop()
	}()

	return linesCh, nil
}

// WatchFromBeginning watches a file starting from the beginning.
func WatchFromBeginning(ctx context.Context, path string) (<-chan []string, error) {
	watcher, err := NewWatcher(path)
	if err != nil {
		return nil, err
	}

	linesCh := watcher.Start(ctx)

	go func() {
		<-ctx.Done()
		_ = watcher.Stop()
	}()

	return linesCh, nil
}

// Reopen attempts to reopen the file if it was closed or moved.
func (w *Watcher) Reopen() error {
	if err := w.reader.Reopen(); err != nil {
		return fmt.Errorf("failed to reopen reader: %w", err)
	}

	_ = w.watcher.Remove(w.path)
	if err := w.watcher.Add(w.path); err != nil {
		return fmt.Errorf("failed to re-add file to watcher: %w", err)
	}

	return nil
}

// CheckFileExists checks if the watched file still exists.
func (w *Watcher) CheckFileExists() bool {
	_, err := os.Stat(w.path)
	return err == nil
}

// FileSize returns the current size of the watched file.
func (w *Watcher) FileSize() (int64, error) {
	info, err := os.Stat(w.path)
	if err != nil {
		return 0, fmt.Errorf("failed to stat file: %w", err)
	}
	return info.Size(), nil
}
