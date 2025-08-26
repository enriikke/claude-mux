# claude-mux 🐙

> Parallel Claude Code execution with isolated git worktrees

[![CI](https://github.com/enriikke/claude-mux/actions/workflows/ci.yml/badge.svg)](https://github.com/enriikke/claude-mux/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/enriikke/claude-mux)](https://github.com/enriikke/claude-mux/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/enriikke/claude-mux)](https://goreportcard.com/report/github.com/enriikke/claude-mux)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## What is claude-mux?

claude-mux enables you to run multiple Claude Code sessions in parallel without conflicts. Each session gets its own isolated git worktree and branch, allowing you to:

- 🎯 **Run multiple AI coding tasks simultaneously** - No more waiting for one task to finish
- 🔀 **Compare different approaches** - Try multiple solutions in parallel
- 🛡️ **Isolate experiments** - Each session is completely isolated from others
- 🔄 **Easy cleanup** - Remove worktrees when done, or keep them for review
- 📦 **Zero dependencies** - Just needs git (which you already have)

## Demo

```bash
# Start a refactoring task
$ claude-mux new refactor-auth
🌳 Creating worktree: refactor-auth-abc123
✅ Worktree created at: .claude-mux/refactor-auth-abc123
🌿 Branch: claude-mux-main-refactor-auth-abc123
🚀 Launching Claude Code...

# In another terminal, start a different task
$ claude-mux new add-tests
🌳 Creating worktree: add-tests-def456
✅ Worktree created at: .claude-mux/add-tests-def456
🌿 Branch: claude-mux-main-add-tests-def456
🚀 Launching Claude Code...

# View active sessions
$ claude-mux list
Active Claude worktrees:

  claude-mux-main-refactor-auth-abc123
    Path:   /project/.claude-mux/refactor-auth-abc123
    Status: active

  claude-mux-main-add-tests-def456
    Path:   /project/.claude-mux/add-tests-def456
    Status: active
```

## Installation

### Using Go

```bash
go install github.com/enriikke/claude-mux/cmd/claude-mux@latest
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/enriikke/claude-mux/releases).

```bash
# Linux/macOS
curl -L https://github.com/enriikke/claude-mux/releases/latest/download/claude-mux_$(uname -s)_$(uname -m).tar.gz | tar xz
sudo mv claude-mux /usr/local/bin/

# Or wget
wget -qO- https://github.com/enriikke/claude-mux/releases/latest/download/claude-mux_$(uname -s)_$(uname -m).tar.gz | tar xz
```

### From Source

```bash
git clone https://github.com/enriikke/claude-mux.git
cd claude-mux
go build -o claude-mux ./cmd/claude-mux
sudo mv claude-mux /usr/local/bin/
```

## Prerequisites

- **Git** - Required for worktree management
- **Claude Code** - Install from [Anthropic](https://docs.anthropic.com/en/docs/claude-code)

## Setup

Add the claude-mux directory to your project's `.gitignore`:

```bash
echo ".claude-mux/" >> .gitignore
```

This prevents worktree directories from being tracked in your repository.

## Usage

### Basic Commands

```bash
# Create a new Claude session with auto-generated name
claude-mux new

# Create a named session
claude-mux new refactor-auth

# List active worktrees
claude-mux list

# Remove a specific worktree
claude-mux remove refactor-auth

# Remove all Claude worktrees
claude-mux prune

# Auto-cleanup after session ends
claude-mux new --cleanup my-task
```

### Options

```bash
claude-mux [command] [flags]

Commands:
  new       Create a new Claude session with isolated worktree
  list      List active Claude worktrees
  remove    Remove a Claude worktree and its branch
  prune     Remove all Claude worktrees

Flags:
  --base-path string    Base path for worktrees (default ".claude-mux")
  --claude-cmd string   Claude Code command (default "claude")
  -v, --verbose         Enable verbose output
  -h, --help           Help for claude-mux
  --version            Version information

New Command Flags:
  -c, --cleanup        Auto-cleanup worktree after Claude exits

Remove Command Flags:
  -f, --force          Force removal even if branch has unmerged changes
```

### Advanced Usage

```bash
# Use custom worktree location
claude-mux new --base-path /tmp/claude-sessions task1

# Use different Claude command
claude-mux new --claude-cmd "claude-dev" experiment

# Verbose output for debugging
claude-mux new -v debug-task
```

## How It Works

1. **Validates** that you're in a git repository
2. **Creates** a new git worktree with a unique branch name
3. **Launches** Claude Code in the isolated worktree directory
4. **Preserves** or cleans up the worktree based on your preference

Each worktree is completely isolated, allowing multiple Claude instances to edit code without conflicts. When you're done, you can merge the best solutions back to your main branch.

## Development

### Building

```bash
# Clone the repository
git clone https://github.com/enriikke/claude-mux.git
cd claude-mux

# Install dependencies
go mod download

# Build
go build -o claude-mux ./cmd/claude-mux

# Run tests
go test ./...

# Run with race detector
go test -race ./...

# Run linter
golangci-lint run

# Build for all platforms
goreleaser build --snapshot --clean
```

### Project Structure

```
claude-mux/
├── cmd/claude-mux/       # Entry point
├── internal/             # Private packages
│   ├── git/             # Git operations
│   ├── worktree/        # Worktree management
│   └── config/          # Configuration
└── pkg/                 # Public packages (future)
```

### Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Roadmap

- [x] Basic worktree creation and management
- [x] Named sessions support
- [x] Auto-cleanup option
- [ ] Session persistence and switching (Phase 1)
- [ ] Process management for attach/detach
- [ ] Container isolation support (Phase 2)
- [ ] Session templates and presets
- [ ] Integration with other AI tools

## FAQ

**Q: What happens to my changes after Claude exits?**
A: By default, worktrees are preserved so you can review and merge changes. Use `--cleanup` to auto-remove.

**Q: Can I run this in a repo with uncommitted changes?**
A: Yes! Worktrees branch from your current HEAD, uncommitted changes stay in your main working directory.

**Q: How do I merge changes from a worktree?**
A: Use standard git commands: `git merge claude-mux-main-task-abc123` or cherry-pick specific commits.

**Q: Does this work with Claude Code's MCP servers?**
A: Yes! Each Claude instance runs normally with full MCP support.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by tmux's session management
- Built for the Claude Code community
- Thanks to all contributors!

---

Made with ❤️ for parallel AI coding
