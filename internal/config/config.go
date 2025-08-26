package config

// Config holds the configuration for claude-mux
type Config struct {
	// WorktreeBasePath is the base directory for all worktrees
	WorktreeBasePath string

	// ClaudeCommand is the command to launch Claude Code
	ClaudeCommand string

	// AutoCleanup determines if worktrees are removed after Claude exits
	AutoCleanup bool

	// Verbose enables detailed output
	Verbose bool
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		WorktreeBasePath: ".claude-mux",
		ClaudeCommand:    "claude",
		AutoCleanup:      false,
		Verbose:          false,
	}
}
