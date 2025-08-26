package main

import (
	"fmt"
	"os"

	"github.com/enriikke/claude-mux/internal/config"
	"github.com/enriikke/claude-mux/internal/worktree"
	"github.com/spf13/cobra"
)

var (
	// Version info, populated by goreleaser
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func execute() error {
	var cfg config.Config

	rootCmd := &cobra.Command{
		Use:   "claude-mux",
		Short: "Parallel Claude Code execution with git worktrees",
		Long: `claude-mux enables parallel Claude Code sessions by creating isolated git worktrees.
Each session runs in its own branch and directory, preventing conflicts when running
multiple AI coding tasks simultaneously.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfg.WorktreeBasePath, "base-path", ".claude-mux", "Base path for worktrees")
	rootCmd.PersistentFlags().StringVar(&cfg.ClaudeCommand, "claude-cmd", "claude", "Claude Code command")
	rootCmd.PersistentFlags().BoolVarP(&cfg.Verbose, "verbose", "v", false, "Enable verbose output")

	// New command - creates worktree and launches Claude
	newCmd := &cobra.Command{
		Use:     "new [name]",
		Short:   "Create a new Claude session with isolated worktree",
		Aliases: []string{"create", "start"},
		Args:    cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var name string
			if len(args) > 0 {
				name = args[0]
			}

			autoCleanup, _ := cmd.Flags().GetBool("cleanup")
			cfg.AutoCleanup = autoCleanup

			manager := worktree.NewManager(cfg)
			return manager.CreateAndLaunch(name)
		},
	}
	newCmd.Flags().BoolP("cleanup", "c", false, "Auto-cleanup worktree after Claude exits")

	// List command - shows active worktrees
	listCmd := &cobra.Command{
		Use:     "list",
		Short:   "List active Claude worktrees",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := worktree.NewManager(cfg)
			return manager.List()
		},
	}

	// Remove command - cleanup specific worktree
	removeCmd := &cobra.Command{
		Use:     "remove <name>",
		Short:   "Remove a Claude worktree and its branch",
		Aliases: []string{"rm", "delete"},
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			force, _ := cmd.Flags().GetBool("force")
			manager := worktree.NewManager(cfg)
			return manager.Remove(args[0], force)
		},
	}
	removeCmd.Flags().BoolP("force", "f", false, "Force removal even if branch has unmerged changes")

	// Prune command - cleanup all claude-mux worktrees
	pruneCmd := &cobra.Command{
		Use:   "prune",
		Short: "Remove all Claude worktrees",
		RunE: func(cmd *cobra.Command, args []string) error {
			manager := worktree.NewManager(cfg)
			return manager.Prune()
		},
	}

	rootCmd.AddCommand(newCmd, listCmd, removeCmd, pruneCmd)
	return rootCmd.Execute()
}
