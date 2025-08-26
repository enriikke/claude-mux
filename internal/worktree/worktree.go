package worktree

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/enriikke/claude-mux/internal/config"
	"github.com/enriikke/claude-mux/internal/git"
)

// Manager handles git worktree operations for Claude sessions
type Manager struct {
	config config.Config
	git    *git.Client
}

// NewManager creates a new worktree manager
func NewManager(cfg config.Config) *Manager {
	return &Manager{
		config: cfg,
		git:    git.NewClient(cfg.Verbose),
	}
}

// CreateAndLaunch creates a new worktree and launches Claude Code
func (m *Manager) CreateAndLaunch(name string) error {
	// Validate we're in a git repository
	if err := m.git.ValidateRepo(); err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	// Generate worktree details
	details, err := m.generateWorktreeDetails(name)
	if err != nil {
		return fmt.Errorf("failed to generate worktree details: %w", err)
	}

	// Create the worktree
	fmt.Printf("🌳 Creating worktree: %s\n", details.Name)
	if err := m.createWorktree(details); err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	fmt.Printf("✅ Worktree created at: %s\n", details.Path)
	fmt.Printf("🌿 Branch: %s\n", details.Branch)

	// Launch Claude Code
	fmt.Printf("\n🚀 Launching Claude Code...\n")
	if err := m.launchClaude(details.Path); err != nil {
		if m.config.AutoCleanup {
			_ = m.cleanup(details)
		}
		return fmt.Errorf("failed to launch Claude: %w", err)
	}

	// Cleanup if requested
	if m.config.AutoCleanup {
		fmt.Printf("\n🧹 Cleaning up worktree...\n")
		return m.cleanup(details)
	}

	fmt.Printf("\n✨ Session completed. Worktree preserved at: %s\n", details.Path)
	fmt.Printf("💡 To remove: claude-mux remove %s\n", details.Name)
	return nil
}

// List shows all active Claude worktrees
func (m *Manager) List() error {
	worktrees, err := m.git.ListWorktrees()
	if err != nil {
		return err
	}

	// Filter for claude-mux worktrees
	var claudeWorktrees []git.Worktree
	for _, wt := range worktrees {
		if strings.Contains(wt.Branch, "claude-mux-") ||
			strings.Contains(wt.Path, m.config.WorktreeBasePath) {
			claudeWorktrees = append(claudeWorktrees, wt)
		}
	}

	if len(claudeWorktrees) == 0 {
		fmt.Println("No active Claude worktrees found.")
		return nil
	}

	fmt.Println("Active Claude worktrees:")
	fmt.Println()
	for _, wt := range claudeWorktrees {
		status := "active"
		if wt.Locked {
			status = "locked"
		}
		fmt.Printf("  %s\n", wt.Branch)
		fmt.Printf("    Path:   %s\n", wt.Path)
		fmt.Printf("    Status: %s\n", status)
		fmt.Println()
	}

	return nil
}

// Remove deletes a specific worktree and its branch
func (m *Manager) Remove(name string, force bool) error {
	// Find the worktree
	worktrees, err := m.git.ListWorktrees()
	if err != nil {
		return err
	}

	var target *git.Worktree
	for _, wt := range worktrees {
		if strings.Contains(wt.Branch, name) || strings.HasSuffix(wt.Path, name) {
			target = &wt
			break
		}
	}

	if target == nil {
		return fmt.Errorf("worktree '%s' not found", name)
	}

	details := WorktreeDetails{
		Name:   name,
		Branch: target.Branch,
		Path:   target.Path,
	}

	return m.cleanup(details)
}

// Prune removes all Claude worktrees
func (m *Manager) Prune() error {
	worktrees, err := m.git.ListWorktrees()
	if err != nil {
		return err
	}

	count := 0
	for _, wt := range worktrees {
		if strings.Contains(wt.Branch, "claude-mux-") ||
			strings.Contains(wt.Path, m.config.WorktreeBasePath) {
			details := WorktreeDetails{
				Name:   filepath.Base(wt.Path),
				Branch: wt.Branch,
				Path:   wt.Path,
			}
			if err := m.cleanup(details); err != nil {
				fmt.Printf("⚠️  Failed to remove %s: %v\n", wt.Path, err)
			} else {
				count++
			}
		}
	}

	fmt.Printf("🧹 Removed %d worktree(s)\n", count)
	return nil
}

// WorktreeDetails contains information about a worktree
type WorktreeDetails struct {
	Name   string
	Branch string
	Path   string
}

// generateWorktreeDetails creates unique names for a new worktree
func (m *Manager) generateWorktreeDetails(name string) (WorktreeDetails, error) {
	// Get current branch for context
	currentBranch, err := m.git.CurrentBranch()
	if err != nil {
		return WorktreeDetails{}, err
	}

	// Generate unique identifier
	timestamp := time.Now().Format("20060102-150405")
	randomBytes := make([]byte, 3)
	if _, err := rand.Read(randomBytes); err != nil {
		// Fallback to timestamp only if random fails
		randomBytes = []byte{0, 0, 0}
	}
	randomHex := hex.EncodeToString(randomBytes)

	// Build names
	var sessionName string
	if name != "" {
		sessionName = fmt.Sprintf("%s-%s", name, randomHex)
	} else {
		sessionName = fmt.Sprintf("%s-%s", timestamp, randomHex)
	}

	branch := fmt.Sprintf("claude-mux-%s-%s", currentBranch, sessionName)

	// Clean branch name (git branch naming rules)
	branch = strings.ReplaceAll(branch, "/", "-")
	branch = strings.ReplaceAll(branch, " ", "-")

	// Get absolute path for worktree
	basePath := m.config.WorktreeBasePath
	if !filepath.IsAbs(basePath) {
		cwd, err := os.Getwd()
		if err != nil {
			return WorktreeDetails{}, err
		}
		basePath = filepath.Join(cwd, basePath)
	}

	return WorktreeDetails{
		Name:   sessionName,
		Branch: branch,
		Path:   filepath.Join(basePath, sessionName),
	}, nil
}

// createWorktree creates a new git worktree
func (m *Manager) createWorktree(details WorktreeDetails) error {
	// Create parent directory with secure permissions
	parentDir := filepath.Dir(details.Path)
	if err := os.MkdirAll(parentDir, 0750); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create worktree with new branch
	return m.git.CreateWorktree(details.Path, details.Branch)
}

// launchClaude starts Claude Code in the specified directory
func (m *Manager) launchClaude(worktreePath string) error {
	// Change to worktree directory
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Chdir(worktreePath); err != nil {
		return fmt.Errorf("failed to change directory: %w", err)
	}

	defer func() {
		err := os.Chdir(originalDir)
		if err != nil {
			fmt.Printf("⚠️  Failed to restore original directory: %v\n", err)
		}
	}()

	// Launch Claude Code
	// #nosec G204 -- ClaudeCommand comes from user config, not untrusted input
	cmd := exec.Command(m.config.ClaudeCommand)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// cleanup removes a worktree and its branch
func (m *Manager) cleanup(details WorktreeDetails) error {
	// Remove worktree
	if err := m.git.RemoveWorktree(details.Path); err != nil {
		fmt.Printf("⚠️  Failed to remove worktree: %v\n", err)
	} else {
		fmt.Printf("✅ Removed worktree: %s\n", details.Path)
	}

	// Try to delete branch
	if err := m.git.DeleteBranch(details.Branch, false); err != nil {
		// Try force delete
		if err := m.git.DeleteBranch(details.Branch, true); err != nil {
			fmt.Printf("ℹ️  Branch preserved (has unmerged changes): %s\n", details.Branch)
		} else {
			fmt.Printf("✅ Deleted branch: %s\n", details.Branch)
		}
	} else {
		fmt.Printf("✅ Deleted branch: %s\n", details.Branch)
	}

	return nil
}
