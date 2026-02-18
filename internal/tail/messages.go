package tail

import "time"

// NewLinesMsg is sent when new lines are read from the file.
type NewLinesMsg struct {
	Lines     []string
	Timestamp time.Time
}

// TailErrorMsg is sent when an error occurs during tailing.
type TailErrorMsg struct {
	Error error
	Path  string
}

// TailStartedMsg is sent when tailing is started.
type TailStartedMsg struct {
	Path string
}

// TailStoppedMsg is sent when tailing is stopped.
type TailStoppedMsg struct {
	Path string
}

// TailCmd represents a command to control the tail functionality.
type TailCmd struct {
	Action  string
	Path    string
	Options map[string]any
}

// TailActions defines the available tail actions.
const (
	TailActionStart   = "start"
	TailActionStop    = "stop"
	TailActionRestart = "restart"
	TailActionPause   = "pause"
	TailActionResume  = "resume"
)

// FileRotatedMsg is sent when the watched file is rotated.
type FileRotatedMsg struct {
	OldPath string
	NewPath string
}

// FileTruncatedMsg is sent when the watched file is truncated.
type FileTruncatedMsg struct {
	Path    string
	OldSize int64
	NewSize int64
}

// FileCreatedMsg is sent when the watched file is created.
type FileCreatedMsg struct {
	Path string
}

// FileDeletedMsg is sent when the watched file is deleted.
type FileDeletedMsg struct {
	Path string
}

// TailStatus represents the current status of the tail operation.
type TailStatus int

const (
	TailStatusStopped TailStatus = iota
	TailStatusStarting
	TailStatusRunning
	TailStatusPaused
	TailStatusError
)

func (s TailStatus) String() string {
	switch s {
	case TailStatusStopped:
		return "stopped"
	case TailStatusStarting:
		return "starting"
	case TailStatusRunning:
		return "running"
	case TailStatusPaused:
		return "paused"
	case TailStatusError:
		return "error"
	default:
		return "unknown"
	}
}

// TailStatusMsg is sent when the tail status changes.
type TailStatusMsg struct {
	Status   TailStatus
	Path     string
	Error    error
	Timeline []TailStatusEvent
}

// TailStatusEvent represents a status change event.
type TailStatusEvent struct {
	Status    TailStatus
	Timestamp time.Time
	Error     error
}
