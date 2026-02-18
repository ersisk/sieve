# CLAUDE.md — Sieve Project Intelligence

## Project Overview

Sieve is an open-source, multi-platform JSON log viewer TUI (Terminal User Interface) written in Go using the Bubble Tea framework. It provides advanced features like syntax highlighting, fuzzy finding, live tailing, filtering, and a real-time dashboard.

## Tech Stack

- **Language:** Go 1.22+
- **TUI Framework:** Bubble Tea (github.com/charmbracelet/bubbletea)
- **Styling:** Lip Gloss (github.com/charmbracelet/lipgloss)
- **Components:** Bubbles (github.com/charmbracelet/bubbles)
- **CLI:** Cobra (github.com/spf13/cobra)
- **Config:** Viper (github.com/spf13/viper)
- **File Watching:** fsnotify (github.com/fsnotify/fsnotify)
- **Testing:** Go standard testing + testify

## Project Structure

```
sieve/
├── cmd/                        # CLI entrypoint
│   └── root.go                 # Cobra root command & flags
├── internal/                   # Private application code
│   ├── app/                    # Bubble Tea main model
│   │   ├── model.go            # App model, Init, Update, View
│   │   ├── update.go           # Message handlers
│   │   ├── view.go             # Rendering logic
│   │   └── keymap.go           # Keybinding definitions
│   ├── parser/                 # Log parsing engine
│   │   ├── json.go             # JSON/JSONL parser
│   │   ├── detector.go         # Auto-detect log format
│   │   └── parser_test.go
│   ├── filter/                 # Filter expression engine
│   │   ├── engine.go           # Filter evaluation
│   │   ├── expr.go             # Expression AST & parser
│   │   ├── preset.go           # Named filter presets
│   │   └── filter_test.go
│   ├── search/                 # Search subsystem
│   │   ├── fuzzy.go            # Fuzzy file finder
│   │   ├── regex.go            # Regex search engine
│   │   └── search_test.go
│   ├── tail/                   # Live file tailing
│   │   ├── watcher.go          # fsnotify-based watcher
│   │   ├── reader.go           # Incremental file reader
│   │   └── tail_test.go
│   ├── ui/                     # TUI component library
│   │   ├── logview.go          # Main log viewport
│   │   ├── sidebar.go          # Sidebar panel
│   │   ├── searchbar.go        # Search input bar
│   │   ├── statusbar.go        # Bottom status bar
│   │   ├── filterbar.go        # Filter input panel
│   │   ├── dashboard.go        # Analytics dashboard
│   │   ├── help.go             # Help overlay
│   │   └── treeview.go         # JSON tree expand/collapse
│   ├── theme/                  # Color themes & styling
│   │   ├── theme.go            # Theme interface & loader
│   │   ├── builtin.go          # Built-in themes
│   │   └── theme_test.go
│   └── config/                 # Configuration management
│       ├── config.go           # Viper-based config loader
│       └── defaults.go         # Default configuration values
├── pkg/                        # Public reusable packages
│   └── logentry/               # Log entry types
│       ├── entry.go            # LogEntry struct & methods
│       └── level.go            # Log level enum & parsing
├── testdata/                   # Test fixtures
│   ├── sample.log
│   ├── multiline.log
│   └── mixed-format.log
├── .goreleaser.yaml            # GoReleaser config
├── Makefile
├── go.mod
├── go.sum
├── CLAUDE.md                   # This file
├── AGENTS.md                   # AI agent instructions
└── README.md
```

## Build & Run Commands

```bash
# Build
go build -o sieve .
make build

# Run development
go run main.go testdata/sample.log
go run main.go -f /var/log/syslog

# Test
go test ./...
make test

# Test with verbose output
go test -v ./...

# Test a specific package
go test -v ./internal/parser/...
go test -v ./internal/filter/...

# Test with race detector
go test -race ./...

# Lint
golangci-lint run
make lint

# Format
gofmt -w .
goimports -w .

# Generate (if needed)
go generate ./...

# Install locally
go install .

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o sieve-linux-amd64 .
GOOS=darwin GOARCH=arm64 go build -o sieve-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o sieve-windows-amd64.exe .

# Release (via GoReleaser)
goreleaser release --snapshot --clean
```

## Code Style & Conventions

### Go Conventions
- Follow standard Go conventions: `gofmt`, `goimports`, effective Go
- Use `internal/` for private packages, `pkg/` for public reusable packages
- Error handling: always wrap errors with `fmt.Errorf("context: %w", err)`
- Naming: use idiomatic Go names (e.g., `LogEntry` not `Log_Entry`)
- Interfaces should be small and focused (1-3 methods)
- Prefer table-driven tests

### Bubble Tea Patterns
- Each UI component is a `tea.Model` with its own `Init()`, `Update()`, `View()`
- Use message passing between components, never share mutable state directly
- Keep `Update()` pure — side effects go into `tea.Cmd` functions
- Use `tea.Batch()` for concurrent commands
- Styles belong in `theme/` package, not inline in view functions
- Use Lip Gloss for all styling — no raw ANSI escape codes

### Naming Conventions
- Files: `snake_case.go`
- Packages: short, lowercase, single-word preferred
- Exported types: `PascalCase`
- Unexported: `camelCase`
- Test files: `*_test.go` in same package
- Interfaces: verb-based (e.g., `Parser`, `Watcher`, `Renderer`)

### Error Handling
- Return `error` as the last return value
- Wrap errors with context: `fmt.Errorf("parsing log entry at line %d: %w", line, err)`
- Use sentinel errors for expected conditions: `var ErrInvalidJSON = errors.New("invalid JSON")`
- Log unexpected errors, return expected errors

### Commit Messages
- Use conventional commits: `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`
- Examples:
  - `feat: add fuzzy file finder with directory traversal`
  - `fix: prevent panic on malformed JSON input`
  - `refactor: extract filter expression parser into separate package`
  - `test: add table-driven tests for log level parsing`

## Key Design Decisions

1. **Virtual Scrolling:** The log viewport uses virtual scrolling — only visible lines are rendered. This is critical for handling multi-GB files.

2. **Lazy Parsing:** Log entries are parsed on-demand as they enter the viewport, not eagerly on file load. Parsed entries are cached.

3. **Message-Driven Architecture:** All state changes flow through Bubble Tea's `Update()` cycle. Components communicate via typed messages, never direct method calls.

4. **Theme Abstraction:** All colors and styles are defined through the `theme.Theme` interface. Components never hardcode colors.

5. **Filter Expressions:** Filters use a simple expression language (JQ-inspired) that compiles to an evaluator function for performance.

6. **Concurrent I/O:** File reading and parsing happen in background goroutines. Results are sent to the TUI via `tea.Cmd` channels.

## Important Patterns

### Adding a New UI Component
1. Create a new file in `internal/ui/`
2. Define a model struct implementing `tea.Model`
3. Define component-specific messages
4. Wire into the main app model in `internal/app/model.go`
5. Add keybindings in `internal/app/keymap.go`

### Adding a New Theme
1. Add theme definition in `internal/theme/builtin.go`
2. Register in the theme map
3. Add to config schema in `internal/config/defaults.go`

### Adding a New Filter Operation
1. Define the operator in `internal/filter/expr.go`
2. Implement evaluation in `internal/filter/engine.go`
3. Add parser support in the expression parser
4. Write table-driven tests in `internal/filter/filter_test.go`

## Performance Targets

- Startup: < 100ms for files under 100MB
- Scrolling: 60fps smooth rendering
- Search: < 500ms for 1M lines
- Memory: < 200MB for a 1GB log file (virtual scrolling)
- Tail latency: < 50ms from file write to screen update

## Dependencies Policy

- Prefer stdlib over external dependencies
- Charmbracelet ecosystem is the only "large" dependency group (bubbletea, lipgloss, bubbles)
- Avoid CGo dependencies — pure Go only for cross-compilation
- Pin all dependency versions in `go.mod`
- Review dependency security with `govulncheck`
