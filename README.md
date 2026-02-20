# Sieve

**A blazing-fast, multi-platform JSON log viewer for your terminal.**

Sieve is a modern TUI (Terminal User Interface) application built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea) that transforms the way you interact with JSON logs. Filter, search, follow, and analyze your logs — all from the comfort of your terminal.

[![CI](https://github.com/ersisk/sieve/actions/workflows/ci.yml/badge.svg)](https://github.com/ersisk/sieve/actions/workflows/ci.yml)
[![Release](https://github.com/ersisk/sieve/actions/workflows/release.yml/badge.svg)](https://github.com/ersisk/sieve/actions/workflows/release.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/ersisk/sieve?style=flat-square&logo=go)](https://go.dev/)
[![Go Report Card](https://goreportcard.com/badge/github.com/ersisk/sieve?style=flat-square)](https://goreportcard.com/report/github.com/ersisk/sieve)
[![License](https://img.shields.io/github/license/ersisk/sieve?style=flat-square)](LICENSE)
[![Release](https://img.shields.io/github/v/release/ersisk/sieve?style=flat-square)](https://github.com/ersisk/sieve/releases/latest)
[![Homebrew](https://img.shields.io/badge/Homebrew-available-orange?style=flat-square&logo=homebrew)](https://github.com/ersisk/homebrew-tap)

---

## ✨ Features

### 🎨 Syntax Highlighting & Colorization
- Automatic colorization of JSON log levels (`DEBUG`, `INFO`, `WARN`, `ERROR`, `FATAL`)
- Customizable color themes with full 256-color and true-color support
- Semantic highlighting for timestamps, keys, values, and nested objects
- Distinct visual indicators for different log sources

### 🔎 Fuzzy Finder & Search
- Built-in fuzzy finder for lightning-fast log file discovery
- Real-time incremental search across log entries
- Regex-powered pattern matching
- Multi-field search — filter by message, level, service, or any JSON key
- Search history with recall support

### 📡 Live Tail Mode (`-f`)
- Follow log files in real-time, just like `tail -f` but supercharged
- Auto-scroll with smart pause on manual scroll-up
- Support for multiple simultaneous file tailing
- Stdin pipe support for streaming log pipelines

### 🧹 Advanced Filtering
- JQ-style expression filtering on any JSON field
- Compound filters with AND / OR / NOT logic
- Time-range based filtering with human-friendly inputs (`--since "2h ago"`)
- Log level filtering with threshold support (`--level warn` shows WARN and above)
- Bookmarkable filter presets for repeated use

### 📂 Multi-Source Input
- Read from local files, directories, or glob patterns
- Stdin pipe support (`cat app.log | sieve`)
- Watch entire directories for new log files
- Automatic detection of JSON, JSONL, and mixed-format logs

### 🧭 Navigation & Interaction
- Vim-style keybindings (`j/k`, `g/G`, `/`, `n/N`)
- Expandable/collapsible JSON tree view for deep inspection
- Line bookmarking and quick-jump
- Context view — see surrounding log entries around a match
- Copy selected log entry or field to clipboard

### 📊 Dashboard & Analytics
- Real-time log level distribution histogram
- Logs-per-second throughput indicator
- Error rate sparkline
- Top recurring error messages summary

### ⚡ Performance
- Lazy loading and virtual scrolling for multi-GB log files
- Memory-mapped file I/O for near-instant startup
- Concurrent log parsing with goroutine pool
- Intelligent caching and pagination

---

## Installation

### Homebrew (macOS / Linux) - Recommended

```bash
# Add the tap
brew tap ersisk/tap

# Install sieve
brew install sieve

# Or install directly
brew install ersisk/tap/sieve
```

### Using Go

```bash
go install github.com/ersisk/sieve@latest
```

### Debian / Ubuntu (.deb)

```bash
# Download the latest .deb package
curl -LO https://github.com/ersisk/sieve/releases/latest/download/sieve_VERSION_linux_amd64.deb

# Install
sudo dpkg -i sieve_VERSION_linux_amd64.deb
```

### Fedora / RHEL / CentOS (.rpm)

```bash
# Download the latest .rpm package
curl -LO https://github.com/ersisk/sieve/releases/latest/download/sieve_VERSION_linux_amd64.rpm

# Install
sudo rpm -i sieve_VERSION_linux_amd64.rpm
```

### Alpine Linux (.apk)

```bash
# Download the latest .apk package
curl -LO https://github.com/ersisk/sieve/releases/latest/download/sieve_VERSION_linux_amd64.apk

# Install
sudo apk add --allow-untrusted sieve_VERSION_linux_amd64.apk
```

### From Source

```bash
git clone https://github.com/ersisk/sieve.git
cd sieve
make build
```

### Pre-built Binaries

Download the latest release for your platform from the [Releases](https://github.com/ersisk/sieve/releases) page.

| Platform | Architecture | Download |
|----------|--------------|----------|
| Linux    | amd64        | [sieve_VERSION_linux_amd64.tar.gz](https://github.com/ersisk/sieve/releases/latest) |
| Linux    | arm64        | [sieve_VERSION_linux_arm64.tar.gz](https://github.com/ersisk/sieve/releases/latest) |
| macOS    | Intel        | [sieve_VERSION_darwin_amd64.tar.gz](https://github.com/ersisk/sieve/releases/latest) |
| macOS    | Apple Silicon| [sieve_VERSION_darwin_arm64.tar.gz](https://github.com/ersisk/sieve/releases/latest) |
| Windows  | amd64        | [sieve_VERSION_windows_amd64.zip](https://github.com/ersisk/sieve/releases/latest) |

#### Manual Installation

```bash
# Example for macOS Apple Silicon
curl -LO https://github.com/ersisk/sieve/releases/latest/download/sieve_VERSION_darwin_arm64.tar.gz
tar -xzf sieve_VERSION_darwin_arm64.tar.gz
sudo mv sieve /usr/local/bin/

# Verify installation
sieve --version
```

---

## 🚀 Quick Start

```bash
# View a JSON log file
sieve app.log

# Follow a log file in real-time
sieve -f /var/log/myservice/app.log

# Pipe from stdin
kubectl logs my-pod | sieve

# Open with a filter preset
sieve --level error app.log

# Fuzzy find and open a log file
sieve --find /var/log/
```

---

## ⌨️ Keybindings

| Key | Action |
|---|---|
| `j` / `↓` | Scroll down |
| `k` / `↑` | Scroll up |
| `g` | Jump to top |
| `G` | Jump to bottom |
| `Enter` | Expand / collapse JSON entry |
| `/` | Open search |
| `n` / `N` | Next / previous search result |
| `f` | Open filter panel |
| `F` | Toggle live tail (follow) mode |
| `Tab` | Switch between panels |
| `b` | Bookmark current line |
| `'` | Jump to next bookmark |
| `y` | Copy current entry to clipboard |
| `t` | Toggle theme (light / dark) |
| `d` | Toggle dashboard panel |
| `?` | Show help |
| `q` / `Ctrl+C` | Quit |

---

## 🛠️ Usage & Examples

### Basic Viewing

```bash
# Single file
sieve server.log

# Multiple files
sieve api.log worker.log scheduler.log

# Glob patterns
sieve /var/log/myapp/*.log
```

### Filtering

```bash
# Show only errors and above
sieve --level error app.log

# Filter by a specific JSON field
sieve --filter '.status == 500' app.log

# Combine filters
sieve --level warn --filter '.service == "auth"' --since "1h ago" app.log

# Exclude patterns
sieve --exclude '.path == "/healthz"' app.log
```

### Live Tail

```bash
# Follow a single file
sieve -f app.log

# Follow with a filter
sieve -f --level error --filter '.service == "payment"' app.log

# Follow multiple files
sieve -f /var/log/myapp/*.log
```

### Fuzzy File Finder

```bash
# Interactive fuzzy find in a directory
sieve --find /var/log/

# Recursive search with depth limit
sieve --find --depth 3 /var/log/
```

### Piping & Integration

```bash
# Docker logs
docker logs -f my-container 2>&1 | sieve

# Kubernetes logs
kubectl logs -f deployment/api | sieve

# journalctl
journalctl -u myservice -o json | sieve

# Combined pipeline
cat app.log | jq -c 'select(.level == "error")' | sieve
```

---

## ⚙️ Configuration

Sieve looks for a configuration file at `~/.config/sieve/config.yaml`.

```yaml
# ~/.config/sieve/config.yaml

theme: "kanagawa"          # kanagawa (default) | monokai | dracula | gruvbox | nord

colors:
  debug: "#6272A4"
  info: "#50FA7B"
  warn: "#F1FA8C"
  error: "#FF5555"
  fatal: "#FF79C6"
  timestamp: "#8BE9FD"
  key: "#BD93F9"

keybindings:
  scroll_down: "j"
  scroll_up: "k"
  search: "/"
  quit: "q"

display:
  timestamp_format: "2006-01-02 15:04:05"
  max_line_width: 0          # 0 = auto
  show_line_numbers: true
  wrap_lines: false
  json_indent: 2

filters:
  presets:
    errors-only:
      level: error
    auth-issues:
      level: warn
      filter: '.service == "auth"'
    slow-requests:
      filter: '.duration_ms > 1000'

performance:
  max_buffer_size: 100000     # max lines in memory
  worker_count: 4             # parsing goroutines
```

---

## 🎨 Themes

Sieve ships with several built-in themes:

| Theme | Description |
|---|---|
| `kanagawa` | Inspired by famous painting, balanced and calm (default) |
| `monokai` | Classic dark theme |
| `dracula` | Popular dark purple theme |
| `gruvbox` | Retro warm dark theme |
| `nord` | Arctic, north-bluish palette |

You can also define fully custom themes in your config file.

---

## 🏗️ Architecture

```
sieve/
├── cmd/                    # CLI entrypoint & flag parsing
│   └── root.go
├── internal/
│   ├── app/                # Bubble Tea application model
│   │   ├── model.go
│   │   ├── update.go
│   │   └── view.go
│   ├── parser/             # JSON log parsing engine
│   │   ├── json.go
│   │   └── detector.go
│   ├── filter/             # Filter expression engine
│   │   ├── engine.go
│   │   ├── expr.go
│   │   └── preset.go
│   ├── search/             # Fuzzy finder & regex search
│   │   ├── fuzzy.go
│   │   └── regex.go
│   ├── tail/               # Live file tailing
│   │   ├── watcher.go
│   │   └── reader.go
│   ├── ui/                 # TUI components
│   │   ├── logview.go
│   │   ├── sidebar.go
│   │   ├── searchbar.go
│   │   ├── statusbar.go
│   │   ├── dashboard.go
│   │   └── help.go
│   ├── theme/              # Color themes & styling
│   │   ├── theme.go
│   │   └── builtin.go
│   └── config/             # Configuration loader
│       └── config.go
├── pkg/
│   └── logentry/           # Public log entry types
│       └── entry.go
├── Makefile
├── go.mod
├── go.sum
└── README.md
```

### Tech Stack

| Component | Library |
|---|---|
| TUI Framework | [Bubble Tea](https://github.com/charmbracelet/bubbletea) |
| Layout & Styling | [Lip Gloss](https://github.com/charmbracelet/lipgloss) |
| Text Input | [Bubbles](https://github.com/charmbracelet/bubbles) |
| Fuzzy Matching | [go-fuzzyfinder](https://github.com/ktr0731/go-fuzzyfinder) |
| File Watching | [fsnotify](https://github.com/fsnotify/fsnotify) |
| Config Parsing | [Viper](https://github.com/spf13/viper) |
| CLI Flags | [Cobra](https://github.com/spf13/cobra) |

---

## Development

### Prerequisites

- Go 1.23 or higher
- [golangci-lint](https://golangci-lint.run/usage/install/)
- [GoReleaser](https://goreleaser.com/install/) (for releases)

### Getting Started

```bash
# Clone the repo
git clone https://github.com/ersisk/sieve.git
cd sieve

# Install dependencies
go mod download

# Run in development
go run main.go testdata/sample.log

# Run tests
make test

# Run linter
make lint

# Build
make build

# Run full CI pipeline locally
make ci
```

### Creating a Release

Releases are automated via GitHub Actions. To create a new release:

```bash
# Tag the release
git tag -a v1.0.0 -m "Release v1.0.0"

# Push the tag
git push origin v1.0.0
```

The release workflow will automatically:
- Build binaries for all platforms (Linux, macOS, Windows)
- Create packages (.deb, .rpm, .apk)
- Generate checksums and changelog
- Update the Homebrew tap
- Update the Scoop bucket
- Publish the release to GitHub

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) before submitting a pull request.

---

## 🗺️ Roadmap

- [X] Core JSON log viewing with colorization
- [X] Live tail mode (`-f`)
- [X] Fuzzy file finder
- [X] Advanced filtering engine
- [X] Vim-style keybindings
- [ ] Remote log source support (SSH, S3)
- [ ] Log format auto-detection (logfmt, CLF, CSV)
- [ ] Export filtered results to file
- [ ] Plugin system for custom parsers
- [ ] Log diffing between two files
- [ ] Session save & restore
- [ ] Built-in log aggregation from multiple hosts
- [ ] Lua scripting for custom transformations

---

## 📄 License

This project is licensed under the Apache 2.0 License — see the [LICENSE](LICENSE) file for details.

---

<p align="center">
  Built with ❤️ and <a href="https://github.com/charmbracelet/bubbletea">Bubble Tea</a>
</p>
