package repository

import (
	"testing"
)

func TestConfigLifecycle(t *testing.T) {
	setupTestRepo(t)

	// Initialize repo to create config
	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// Load config
	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() failed: %v", err)
	}

	if cfg.Version != 1 {
		t.Errorf("expected Version to be 1, got %d", cfg.Version)
	}

	if cfg.DefaultBranch != DefaultBranch {
		t.Errorf("expected DefaultBranch to be %q, got %q", DefaultBranch, cfg.DefaultBranch)
	}

	if cfg.RepositoryID == "" {
		t.Errorf("expected RepositoryID to be non-empty UUID")
	}
}
