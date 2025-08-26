package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Client handles git operations
type Client struct {
	verbose bool
}

// NewClient creates a new git client
func NewClient(verbose bool) *Client {
	return &Client{verbose: verbose}
}

// Worktree represents a git worktree
type Worktree struct {
	Path   string
	Branch string
	Commit string
	Locked bool
}

// ValidateRepo checks if we're in a git repository
func (c *Client) ValidateRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("not in a git repository")
	}
	return nil
}

// CurrentBranch returns the current git branch name
func (c *Client) CurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		// Handle detached HEAD
		cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
		output, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("failed to get current branch: %w", err)
		}
		return fmt.Sprintf("detached-%s", strings.TrimSpace(string(output))), nil
	}
	return strings.TrimSpace(string(output)), nil
}

// CreateWorktree creates a new worktree with a new branch
func (c *Client) CreateWorktree(path, branch string) error {
	cmd := exec.Command("git", "worktree", "add", "-b", branch, path)
	if c.verbose {
		cmd.Stdout = &bytes.Buffer{}
		cmd.Stderr = &bytes.Buffer{}
	}

	if err := cmd.Run(); err != nil {
		if c.verbose {
			return fmt.Errorf("failed to create worktree: %w\nOutput: %s\nError: %s",
				err, cmd.Stdout, cmd.Stderr)
		}
		return fmt.Errorf("failed to create worktree: %w", err)
	}
	return nil
}

// ListWorktrees returns all git worktrees
func (c *Client) ListWorktrees() ([]Worktree, error) {
	cmd := exec.Command("git", "worktree", "list", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list worktrees: %w", err)
	}

	var worktrees []Worktree
	var current Worktree

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			if current.Path != "" {
				worktrees = append(worktrees, current)
				current = Worktree{}
			}
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 1 {
			continue
		}

		switch parts[0] {
		case "worktree":
			if len(parts) > 1 {
				current.Path = parts[1]
			}
		case "branch":
			if len(parts) > 1 {
				// Remove refs/heads/ prefix
				branch := strings.TrimPrefix(parts[1], "refs/heads/")
				current.Branch = branch
			}
		case "HEAD":
			if len(parts) > 1 {
				current.Commit = parts[1]
			}
		case "locked":
			current.Locked = true
		}
	}

	// Add last worktree if exists
	if current.Path != "" {
		worktrees = append(worktrees, current)
	}

	return worktrees, nil
}

// RemoveWorktree removes a git worktree
func (c *Client) RemoveWorktree(path string) error {
	cmd := exec.Command("git", "worktree", "remove", path, "--force")
	return cmd.Run()
}

// DeleteBranch deletes a git branch
func (c *Client) DeleteBranch(branch string, force bool) error {
	flag := "-d"
	if force {
		flag = "-D"
	}
	cmd := exec.Command("git", "branch", flag, branch)
	return cmd.Run()
}
