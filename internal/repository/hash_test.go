package repository

import (
	"testing"
)

func TestHashContentDeterministic(t *testing.T) {
	data := []byte("You are a helpful legal assistant.")

	hash1 := HashContent(ObjectTypeBlob, data)
	hash2 := HashContent(ObjectTypeBlob, data)

	if hash1 != hash2 {
		t.Errorf("expected deterministic hashing, got %q and %q", hash1, hash2)
	}

	if len(hash1) != 64 {
		t.Errorf("expected 64-character SHA-256 hex string, got length %d", len(hash1))
	}
}

func TestHashContentKnownVector(t *testing.T) {
	expected := "8aec4e4876f854f688d0ebfc8f37598f38e5fd6903cccc850ca36591175aeb60"
	actual := HashContent(ObjectTypeBlob, []byte("hello"))

	if actual != expected {
		t.Errorf("expected hash %q, got %q", expected, actual)
	}
}

func TestHashContentDifferentTypes(t *testing.T) {
	data := []byte("same content")

	blobHash := HashContent(ObjectTypeBlob, data)
	treeHash := HashContent(ObjectTypeTree, data)

	if blobHash == treeHash {
		t.Errorf("expected different hashes for different object types, got same: %q", blobHash)
	}
}
