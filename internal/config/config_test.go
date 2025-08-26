package config

import (
	"reflect"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	got := DefaultConfig()
	want := Config{
		WorktreeBasePath: ".claude-mux",
		ClaudeCommand:    "claude",
		AutoCleanup:      false,
		Verbose:          false,
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("DefaultConfig() = %v, want %v", got, want)
	}
}

func TestConfig_Fields(t *testing.T) {
	cfg := Config{
		WorktreeBasePath: "/custom/path",
		ClaudeCommand:    "claude-dev",
		AutoCleanup:      true,
		Verbose:          true,
	}

	tests := []struct {
		name  string
		field string
		want  interface{}
	}{
		{"WorktreeBasePath", "WorktreeBasePath", "/custom/path"},
		{"ClaudeCommand", "ClaudeCommand", "claude-dev"},
		{"AutoCleanup", "AutoCleanup", true},
		{"Verbose", "Verbose", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.field {
			case "WorktreeBasePath":
				if cfg.WorktreeBasePath != tt.want {
					t.Errorf("Config.%s = %v, want %v", tt.field, cfg.WorktreeBasePath, tt.want)
				}
			case "ClaudeCommand":
				if cfg.ClaudeCommand != tt.want {
					t.Errorf("Config.%s = %v, want %v", tt.field, cfg.ClaudeCommand, tt.want)
				}
			case "AutoCleanup":
				if cfg.AutoCleanup != tt.want {
					t.Errorf("Config.%s = %v, want %v", tt.field, cfg.AutoCleanup, tt.want)
				}
			case "Verbose":
				if cfg.Verbose != tt.want {
					t.Errorf("Config.%s = %v, want %v", tt.field, cfg.Verbose, tt.want)
				}
			}
		})
	}
}
