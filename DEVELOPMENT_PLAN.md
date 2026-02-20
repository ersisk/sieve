# Sieve Development Plan

## Context

Sieve is a terminal-based JSON log viewer TUI built with Go + Bubble Tea. The project is architecturally complete (README, CLAUDE.md, AGENTS.md) but has **zero implementation code**. This plan builds the project from scratch in 7 phases.

---

## Phase 1: Project Skeleton & Core Types (`pkg/logentry`)

**Goal:** Initialize Go module, create `pkg/logentry` package with zero external dependencies.

**Files:**
- `go.mod` — `module github.com/ersanisk/sieve`, Go 1.22+
- `main.go` — Minimal entry point
- `pkg/logentry/level.go` — `Level` type (iota enum: Debug..Fatal), `ParseLevel(string) Level`, `String()`
- `pkg/logentry/entry.go` — `Entry` struct (Level, Message, Timestamp, Caller, Fields, Raw, Line, IsJSON), `GetField()`
- `pkg/logentry/level_test.go` — Table-driven tests (upper/lowercase, single char, numeric Bunyan)
- `pkg/logentry/entry_test.go` — GetField tests
- `testdata/sample.log`, `multiline.log`, `mixed-format.log` — Test fixtures

**Output:** `go build ./...` + `go test ./pkg/logentry/...` pass

---

## Phase 2: Configuration, Themes & CLI Skeleton

**Goal:** Viper-based config, Lip Gloss theme interface, Cobra CLI skeleton.

**Files:**
- `internal/config/defaults.go` — Default values
- `internal/config/config.go` — `Config` struct, `Load() (*Config, error)` (via Viper)
- `internal/config/config_test.go`
- `internal/theme/theme.go` — `Theme` interface: `LevelStyle()`, `TimestampStyle()`, `KeyStyle()`, etc.
- `internal/theme/builtin.go` — Kanagawa (default), Monokai, Dracula, Gruvbox, Nord themes; `Registry`, `Get(name)`
- `internal/theme/theme_test.go`
- `cmd/root.go` — Cobra root command: `--theme`, `--level`, `--filter`, `-f`, `--config` flags

**New dependencies:** cobra, viper, lipgloss

**Output:** `go run main.go --help` shows help, flags are parsed

---

## Phase 3: JSON Parser & Format Detector

**Goal:** Parser that converts log lines into `logentry.Entry`. Zero UI dependency.

**Files:**
- `internal/parser/json.go` — `Parser` struct, `ParseLine(raw, lineNum)`, `ParseLines(reader)`. Field detection: level/lvl/severity, msg/message, time/timestamp/ts/@timestamp, caller/source
- `internal/parser/detector.go` — `Format` type, `DetectFormat(reader)` — detects format by sampling first N lines
- `internal/parser/parser_test.go` — Comprehensive table-driven tests (valid JSON, malformed JSON, plain text, nested JSON, different timestamp formats)
- `internal/parser/bench_test.go` — `BenchmarkParseLine` (target <1μs/line)

**Output:** Parser tests + benchmarks pass

---

## Phase 4: Filter Engine & Search Infrastructure

**Goal:** JQ-inspired filter expression engine and fuzzy/regex search. Zero UI dependency.

**Files:**
- `internal/filter/expr.go` — AST definitions: `Expr` interface, `FieldAccess`, `Literal`, `BinaryOp`, `UnaryOp`, `CompoundExpr`; `Parse(input) (Expr, error)`
- `internal/filter/engine.go` — `Compile(Expr) (Evaluator, error)`, `ByLevel(min) Evaluator`
- `internal/filter/preset.go` — Preset definitions and conversion
- `internal/filter/filter_test.go`
- `internal/search/fuzzy.go` — `FuzzyMatch(entries, query) []SearchResult`
- `internal/search/regex.go` — `RegexMatch(entries, pattern) ([]SearchResult, error)`
- `internal/search/search_test.go`

**Supported operators:** `==`, `!=`, `>`, `<`, `>=`, `<=`, `contains`, `matches`, `and`, `or`, `not`

**Output:** Filter + search tests pass, benchmarks run

---

## Phase 5: File Tailing & Incremental Reader

**Goal:** fsnotify-based file watcher and incremental reader.

**Files:**
- `internal/tail/reader.go` — `Reader` struct, `NewReader(path)`, `ReadNew() ([]string, error)` — reads new lines from last position
- `internal/tail/watcher.go` — `Watcher` struct, `NewWatcher(path)`, `Start(ctx) <-chan []string`, `Stop()`
- `internal/tail/messages.go` — Bubble Tea message types: `NewLinesMsg`, `TailErrorMsg`, `TailCmd`
- `internal/tail/tail_test.go` — Tests with temp files for reading/watching

**New dependencies:** fsnotify

**Output:** `-f` mode works in CLI as plain text output (no TUI yet)

---

## Phase 6: TUI Components (`internal/ui/`)

**Goal:** Build UI components as Bubble Tea models. Each component is independently testable.

**Files:**
- `internal/ui/logview.go` — Main log viewport: **virtual scrolling** (only render visible lines), themed coloring, line numbers
- `internal/ui/statusbar.go` — Bottom status bar: file name, line count, active filter, mode
- `internal/ui/searchbar.go` — Search input bar (bubbles/textinput)
- `internal/ui/filterbar.go` — Filter input panel
- `internal/ui/sidebar.go` — JSON detail view
- `internal/ui/help.go` — Help overlay (keybinding list)
- `internal/ui/treeview.go` — JSON tree view (expandable/collapsible)
- `internal/ui/dashboard.go` — Dashboard panel (level distribution, lines/sec)
- `internal/ui/messages.go` — All shared UI message types

**Critical:** `logview.go` must handle 100K+ entries smoothly via virtual scrolling. Only lines between `offset` and `offset+height` are rendered.

**New dependencies:** bubbletea, bubbles

**Output:** Each component passes isolated tests

---

## Phase 7: Main Application Model & Integration

**Goal:** Wire everything together into the main Bubble Tea model. Produce a working TUI.

**Files:**
- `internal/app/model.go` — `Model` struct (all sub-components + data + state), `NewModel()`, `Init()`, `Update()`, `View()`
- `internal/app/keymap.go` — `KeyMap` struct: j/k, g/G, /, n/N, f, F, Tab, b, d, ?, q
- `internal/app/update.go` — Message routing, mode transitions, file/tail/filter/search messages
- `internal/app/view.go` — Layout management (via lipgloss)
- `internal/app/commands.go` — `tea.Cmd` functions: `loadFileCmd`, `tailCmd`, `searchCmd`, `filterCmd`
- `cmd/root.go` — Update: `config.Load()` → `theme.Get()` → `app.NewModel()` → `tea.NewProgram().Run()`
- `main.go` — Update: `cmd.Execute()`, `version`/`buildTime` variables

**Output (fully working application):**
- `./sieve testdata/sample.log` — TUI opens, JSON logs visible, scrolling works
- `./sieve -f /tmp/test.log` — Live tail mode
- `./sieve --level error testdata/sample.log` — Level filter
- `./sieve --filter '.service == "auth"' testdata/sample.log` — Expression filter
- `/` search, `f` filter, `?` help, `d` dashboard, `Tab` sidebar

---

## Dependency Graph

```
Phase 1: pkg/logentry ──────────────────────────────── (stdlib only)
  ↓
Phase 2: internal/config + theme + cmd/ ─────────────── (cobra, viper, lipgloss)
  ↓
Phase 3: internal/parser ───────────────────────────── (logentry + stdlib)
  ↓
Phase 4: internal/filter + search ──────────────────── (logentry + stdlib)
  ↓
Phase 5: internal/tail ─────────────────────────────── (fsnotify + stdlib)
  ↓
Phase 6: internal/ui/ ──────────────────────────────── (logentry, theme, bubbletea, bubbles)
  ↓
Phase 7: internal/app/ + main.go ───────────────────── (wires everything together)
```

## Key Rules

1. **`go build ./...` and `go test ./...` must pass after each phase** — no breakage allowed
2. **Package boundaries must be respected:** parser/filter/search/tail → zero UI dependency; pkg/logentry → zero external dependency
3. **UI components must not import each other** — communication via messages only
4. **Rename `.golangci` to `.golangci.yml`**
5. **Tests must be written alongside implementation**, benchmarks added from Phase 3 onwards

## Verification

After each phase:
```bash
go build ./...          # No compilation errors
go test -race ./...     # All tests pass
go vet ./...            # No vet warnings
gofmt -d .              # No format diffs
```

After Phase 7 (full validation):
```bash
make ci                       # build + test + lint + vet
./sieve testdata/sample.log   # TUI works
./sieve -f /tmp/test.log      # Tail mode works
```
