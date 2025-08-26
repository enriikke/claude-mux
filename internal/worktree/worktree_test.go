package worktree

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/enriikke/claude-mux/internal/config"
)

// restoreDirectory returns a function that restores the working directory
// and reports any error via t.Errorf
func restoreDirectory(t *testing.T, dir string) func() {
	return func() {
		if err := os.Chdir(dir); err != nil {
			t.Errorf("Failed to restore directory: %v", err)
		}
	}
}

func setupTestRepo(t *testing.T) string {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}

	// Initialize git repo
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	commands := [][]string{
		{"git", "init"},
		{"git", "config", "user.email", "test@example.com"},
		{"git", "config", "user.name", "Test User"},
	}

	for _, cmd := range commands {
		if err := exec.Command(cmd[0], cmd[1:]...).Run(); err != nil {
			t.Fatalf("Failed to run %v: %v", cmd, err)
		}
	}

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	if err := exec.Command("git", "add", ".").Run(); err != nil {
		t.Fatalf("Failed to git add: %v", err)
	}
	if err := exec.Command("git", "commit", "-m", "initial").Run(); err != nil {
		t.Fatalf("Failed to git commit: %v", err)
	}

	// Return to original directory
	if err := os.Chdir(originalDir); err != nil {
		t.Fatalf("Failed to restore directory: %v", err)
	}

	return tmpDir
}

func TestManager_generateWorktreeDetails(t *testing.T) {
	cfg := config.Config{
		WorktreeBasePath: ".claude-mux-test",
		ClaudeCommand:    "echo",
		AutoCleanup:      false,
		Verbose:          false,
	}

	manager := NewManager(cfg)

	// Setup test repo
	repoDir := setupTestRepo(t)
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer restoreDirectory(t, originalDir)()
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change to repo directory: %v", err)
	}

	tests := []struct {
		name     string
		input    string
		wantName bool
	}{
		{"with name", "my-task", true},
		{"empty name", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			details, err := manager.generateWorktreeDetails(tt.input)
			if err != nil {
				t.Fatalf("generateWorktreeDetails() error = %v", err)
			}

			// Check if name is included in the result
			if tt.wantName && !strings.Contains(details.Name, tt.input) {
				t.Errorf("Expected name to contain %q, got %q", tt.input, details.Name)
			}

			// Check branch format
			if !strings.HasPrefix(details.Branch, "claude-mux-") {
				t.Errorf("Expected branch to start with 'claude-mux-', got %q", details.Branch)
			}

			// Check path includes base path
			if !strings.Contains(details.Path, cfg.WorktreeBasePath) {
				t.Errorf("Expected path to contain base path %q, got %q", cfg.WorktreeBasePath, details.Path)
			}
		})
	}
}

func TestManager_List(t *testing.T) {
	cfg := config.Config{
		WorktreeBasePath: ".claude-mux-test",
		ClaudeCommand:    "echo",
		AutoCleanup:      false,
		Verbose:          false,
	}

	manager := NewManager(cfg)

	// Setup test repo
	repoDir := setupTestRepo(t)
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer restoreDirectory(t, originalDir)()
	if err := os.Chdir(repoDir); err != nil {
		t.Fatalf("Failed to change to repo directory: %v", err)
	}

	// Test listing with no worktrees (should not error)
	err = manager.List()
	if err != nil {
		t.Errorf("List() with no worktrees should not error: %v", err)
	}

	// Create a worktree manually
	worktreePath := filepath.Join(repoDir, ".claude-mux-test", "test-worktree")
	if err := exec.Command("git", "worktree", "add", "-b", "claude-mux-test", worktreePath).Run(); err != nil {
		t.Fatalf("Failed to create worktree: %v", err)
	}

	// Test listing with a worktree
	err = manager.List()
	if err != nil {
		t.Errorf("List() with worktrees should not error: %v", err)
	}
}

func TestWorktreeDetails(t *testing.T) {
	details := WorktreeDetails{
		Name:   "test-task",
		Branch: "claude-mux-main-test-task",
		Path:   "/tmp/claude-mux/test-task",
	}

	if details.Name != "test-task" {
		t.Errorf("Expected Name to be 'test-task', got %q", details.Name)
	}

	if details.Branch != "claude-mux-main-test-task" {
		t.Errorf("Expected Branch to be 'claude-mux-main-test-task', got %q", details.Branch)
	}

	if details.Path != "/tmp/claude-mux/test-task" {
		t.Errorf("Expected Path to be '/tmp/claude-mux/test-task', got %q", details.Path)
	}
}
