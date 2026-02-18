# AGENTS.md — AI Agent Instructions for Sieve

> This file provides instructions for AI coding agents (Claude Code, OpenCode, Cursor, Copilot, Aider, etc.) working on the Sieve project.

## 🧠 Project Context

Sieve is a **terminal-based JSON log viewer** built in **Go** using the **Bubble Tea** TUI framework. It is designed to be fast, ergonomic, and cross-platform. Think of it as `jq` + `less` + `tail -f` combined into a beautiful TUI.

**Read `CLAUDE.md` first** for full project structure, build commands, and conventions.

---

## 📋 Agent Rules

### General

1. **Always run tests after making changes:** `go test ./...`
2. **Always run the linter before committing:** `golangci-lint run`
3. **Always format code:** `gofmt -w . && goimports -w .`
4. **Never introduce CGo dependencies.** The project must cross-compile cleanly.
5. **Never bypass Bubble Tea's message loop.** No goroutine-to-UI direct writes.
6. **Always wrap errors with context.** Never return bare `err`.
7. **Write tests for every new exported function.** Prefer table-driven tests.
8. **Commit messages must follow conventional commits:** `feat:`, `fix:`, `refactor:`, `test:`, `docs:`, `chore:`

### Code Quality Gates

Before considering any task complete, verify:

- [ ] `go build ./...` — compiles without errors
- [ ] `go test ./...` — all tests pass
- [ ] `go vet ./...` — no vet warnings
- [ ] `golangci-lint run` — no lint issues
- [ ] `gofmt -d .` — no formatting diffs
- [ ] No new dependencies added without justification

---

## 🏛️ Architecture Rules

### Bubble Tea Model Pattern

Every UI component must follow this structure:

```go
package ui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

// Model
type MyComponent struct {
    // state fields
    width  int
    height int
    // ...
}

// Messages (component-specific)
type MyComponentMsg struct { /* ... */ }

// Constructor
func NewMyComponent() MyComponent {
    return MyComponent{}
}

// tea.Model interface
func (m MyComponent) Init() tea.Cmd { return nil }

func (m MyComponent) Update(msg tea.Msg) (MyComponent, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}

func (m MyComponent) View() string {
    // Use lipgloss for styling, NEVER raw ANSI codes
    return ""
}
```

### Package Boundaries

```
cmd/          → CLI parsing only. No business logic.
internal/app/ → Main Bubble Tea model. Orchestrates components.
internal/ui/  → Individual UI components. Each is a tea.Model.
internal/parser/ → Log parsing. Zero UI dependencies.
internal/filter/ → Filter engine. Zero UI dependencies.
internal/search/ → Search logic. Zero UI dependencies.
internal/tail/   → File tailing. Zero UI dependencies.
internal/theme/  → Theming. Only lipgloss dependency.
internal/config/ → Config loading. Only viper dependency.
pkg/logentry/    → Public types. ZERO external dependencies.
```

**Rules:**
- `pkg/` must have zero external dependencies (only stdlib)
- `internal/parser/`, `internal/filter/`, `internal/search/`, `internal/tail/` must NOT import any `internal/ui/` or `internal/app/` packages
- `internal/ui/` components must NOT import each other directly — communicate via messages through `internal/app/`
- `internal/theme/` is the only package that may define `lipgloss.Style` values

### State Management

- All state lives in the main `app.Model` or in individual component models
- Components communicate via Bubble Tea messages, NOT via shared pointers
- Background I/O (file reading, tailing) returns results through `tea.Cmd` → `tea.Msg`
- Never use `sync.Mutex` in UI code — use message passing instead

### File I/O Pattern

```go
// CORRECT: Background I/O via tea.Cmd
func readFileCmd(path string) tea.Cmd {
    return func() tea.Msg {
        data, err := os.ReadFile(path)
        if err != nil {
            return errMsg{err}
        }
        return fileLoadedMsg{data}
    }
}

// WRONG: Blocking I/O in Update()
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    // ❌ NEVER DO THIS — blocks the UI
    data, _ := os.ReadFile("app.log")
}
```

---

## 🧩 Task-Specific Instructions

### When Adding a New Feature

1. Create or modify the appropriate package in `internal/`
2. Define message types for the feature
3. Wire messages into `internal/app/update.go`
4. Add keybindings in `internal/app/keymap.go`
5. Update the help overlay in `internal/ui/help.go`
6. Write tests (unit tests + any integration tests needed)
7. Update `README.md` if the feature is user-facing
8. Update `CLAUDE.md` if the architecture changes

### When Fixing a Bug

1. Write a failing test that reproduces the bug **first**
2. Fix the bug
3. Verify the test passes
4. Check for similar patterns elsewhere in the codebase
5. Run full test suite: `go test -race ./...`

### When Refactoring

1. Ensure full test coverage exists before refactoring
2. Make changes in small, reviewable increments
3. Run tests after each increment
4. Do not change behavior — tests should pass without modification
5. Run benchmarks if the refactoring touches performance-critical paths

### When Adding a New Theme

1. Define the theme colors in `internal/theme/builtin.go`:
```go
var MyTheme = Theme{
    Name: "mytheme",
    Colors: ThemeColors{
        Debug:     lipgloss.Color("#..."),
        Info:      lipgloss.Color("#..."),
        Warn:      lipgloss.Color("#..."),
        Error:     lipgloss.Color("#..."),
        Fatal:     lipgloss.Color("#..."),
        Timestamp: lipgloss.Color("#..."),
        Key:       lipgloss.Color("#..."),
        Value:     lipgloss.Color("#..."),
        Background: lipgloss.Color("#..."),
        Foreground: lipgloss.Color("#..."),
    },
}
```
2. Register it in the theme registry map
3. Add it to `README.md` theme table
4. Test it: `go run main.go --theme mytheme testdata/sample.log`

### When Working with the Filter Engine

- Filters follow a JQ-inspired syntax: `.field == "value"`, `.status > 400`
- The filter pipeline: `string → parse → AST → compile → evaluator func`
- The evaluator is a `func(logentry.Entry) bool` for performance
- Always add test cases to the expression parser tests
- Supported operators: `==`, `!=`, `>`, `<`, `>=`, `<=`, `contains`, `matches`
- Compound: `and`, `or`, `not`

### When Working with the Log Parser

- The parser must handle: valid JSON, JSONL, mixed text+JSON, malformed lines
- Malformed lines should be displayed as plain text, never crash
- Common JSON log fields to detect: `level`, `msg`/`message`, `time`/`timestamp`/`ts`, `caller`, `error`
- Support multiple level formats: `"INFO"`, `"info"`, `"I"`, `30` (numeric Bunyan-style)
- Parser must be safe for concurrent use (called from goroutine pool)

---

## 🚫 Anti-Patterns to Avoid

| ❌ Don't | ✅ Do Instead |
|---|---|
| Raw ANSI escape codes | Use `lipgloss.Style` |
| Global mutable state | Component-local state + messages |
| `log.Fatal()` in library code | Return `error` |
| Blocking I/O in `Update()` | Background `tea.Cmd` |
| `interface{}` / `any` for log fields | Typed `logentry.Entry` with field accessors |
| Large switch statements in `Update()` | Delegate to component `Update()` methods |
| Hardcoded colors | Use `theme.Theme` interface |
| `os.Exit()` anywhere except `main()` | Return errors up the call stack |
| `panic()` for expected errors | Return `error` with context |
| Importing `internal/ui` from `internal/parser` | Keep packages decoupled |

---

## 📐 Testing Guidelines

### Test Structure

```go
func TestParseLogEntry(t *testing.T) {
    tests := []struct {
        name     string
        input    string
        want     logentry.Entry
        wantErr  bool
    }{
        {
            name:  "valid JSON log",
            input: `{"level":"info","msg":"hello","ts":"2024-01-01T00:00:00Z"}`,
            want:  logentry.Entry{Level: logentry.Info, Message: "hello"},
        },
        {
            name:    "malformed JSON",
            input:   `{invalid`,
            wantErr: true,
        },
        // ... more cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := parser.Parse(tt.input)
            if tt.wantErr {
                require.Error(t, err)
                return
            }
            require.NoError(t, err)
            assert.Equal(t, tt.want.Level, got.Level)
            assert.Equal(t, tt.want.Message, got.Message)
        })
    }
}
```

### Test Coverage Expectations

| Package | Minimum Coverage |
|---|---|
| `pkg/logentry` | 90% |
| `internal/parser` | 85% |
| `internal/filter` | 90% |
| `internal/search` | 80% |
| `internal/tail` | 75% |
| `internal/config` | 80% |
| `internal/ui` | 60% (TUI testing is harder) |
| `internal/app` | 50% |

### Benchmark Tests

Performance-critical packages must include benchmarks:

```go
func BenchmarkParseJSON(b *testing.B) {
    line := `{"level":"info","msg":"test","ts":"2024-01-01T00:00:00Z"}`
    for i := 0; i < b.N; i++ {
        parser.Parse(line)
    }
}
```

Run benchmarks: `go test -bench=. -benchmem ./internal/parser/`

---

## 🔧 Useful Commands for Agents

```bash
# Quick validation cycle
go build ./... && go test ./... && go vet ./...

# Full CI check
make ci    # build + test + lint + vet

# Run with sample data
go run main.go testdata/sample.log

# Run in follow mode for testing
echo '{"level":"info","msg":"test"}' >> /tmp/test.log && go run main.go -f /tmp/test.log

# Profile CPU
go test -cpuprofile cpu.prof -bench=. ./internal/parser/
go tool pprof cpu.prof

# Check for race conditions
go test -race ./...

# Dependency audit
go mod tidy
govulncheck ./...
```

---

## 📚 Reference Links

- [Bubble Tea Docs](https://github.com/charmbracelet/bubbletea)
- [Bubble Tea Examples](https://github.com/charmbracelet/bubbletea/tree/master/examples)
- [Lip Gloss Docs](https://github.com/charmbracelet/lipgloss)
- [Bubbles Components](https://github.com/charmbracelet/bubbles)
- [Cobra CLI Docs](https://cobra.dev/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
