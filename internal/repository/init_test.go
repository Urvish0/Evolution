package repository

import (
	"testing"
)

func TestInitAndExists(t *testing.T) {
	setupTestRepo(t)

	// 1. Repository should not exist initially
	if Exists() {
		t.Errorf("expected Exists() to be false before Init()")
	}

	// 2. Initialize repository
	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 3. Repository should now exist
	if !Exists() {
		t.Errorf("expected Exists() to be true after Init()")
	}

	// 4. Double initialization should fail
	if err := Init(); err == nil {
		t.Errorf("expected second Init() to return error, but got nil")
	}
}
