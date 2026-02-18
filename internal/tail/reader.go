package tail

import (
	"bufio"
	"fmt"
	"os"
	"sync"
)

// Reader reads a file incrementally, tracking the last read position.
type Reader struct {
	path   string
	file   *os.File
	pos    int64
	mu     sync.Mutex
	isOpen bool
}

// NewReader creates a new Reader for the specified file.
func NewReader(path string) (*Reader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return &Reader{
		path:   path,
		file:   file,
		pos:    0,
		isOpen: true,
	}, nil
}

// ReadNew reads all new lines since the last read.
// Updates the internal position to the end of the file.
func (r *Reader) ReadNew() ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isOpen {
		return nil, fmt.Errorf("reader is closed")
	}

	if _, err := r.file.Seek(r.pos, 0); err != nil {
		return nil, fmt.Errorf("failed to seek to position %d: %w", r.pos, err)
	}

	var lines []string
	scanner := bufio.NewScanner(r.file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	newPos, _ := r.file.Seek(0, 2)
	r.pos = newPos

	return lines, nil
}

// ReadAll reads the entire file from the beginning.
// Does not update the internal position.
func (r *Reader) ReadAll() ([]string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isOpen {
		return nil, fmt.Errorf("reader is closed")
	}

	if _, err := r.file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("failed to seek to beginning: %w", err)
	}

	var lines []string
	scanner := bufio.NewScanner(r.file)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file: %w", err)
	}

	return lines, nil
}

// Position returns the current read position.
func (r *Reader) Position() int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.pos
}

// SetPosition sets the read position.
func (r *Reader) SetPosition(pos int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isOpen {
		return fmt.Errorf("reader is closed")
	}

	_, err := r.file.Seek(pos, 0)
	if err != nil {
		return fmt.Errorf("failed to seek to position %d: %w", pos, err)
	}

	r.pos = pos
	return nil
}

// ResetPosition resets the read position to the beginning of the file.
func (r *Reader) ResetPosition() error {
	return r.SetPosition(0)
}

// Close closes the reader and the underlying file.
func (r *Reader) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if !r.isOpen {
		return nil
	}

	if err := r.file.Close(); err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	r.isOpen = false
	return nil
}

// IsClosed returns true if the reader is closed.
func (r *Reader) IsClosed() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return !r.isOpen
}

// Reopen reopens the file if it was closed.
func (r *Reader) Reopen() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isOpen {
		return nil
	}

	file, err := os.Open(r.path)
	if err != nil {
		return fmt.Errorf("failed to reopen file: %w", err)
	}

	r.file = file
	r.isOpen = true
	r.pos = 0

	return nil
}
