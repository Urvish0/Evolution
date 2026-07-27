package repository

import (
	"bytes"
	"testing"
)

func TestBlobWriteAndRead(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	content := []byte("System Prompt: You are a legal research AI.")

	// 1. Write blob
	hash, err := WriteBlob(content)
	if err != nil {
		t.Fatalf("WriteBlob() failed: %v", err)
	}

	if len(hash) != 64 {
		t.Fatalf("expected 64-character hash, got %d", len(hash))
	}

	// 2. Verify HasBlob is true
	if !HasBlob(hash) {
		t.Errorf("expected HasBlob(%s) to be true", hash[:8])
	}

	// 3. Read blob back
	readContent, err := ReadBlob(hash)
	if err != nil {
		t.Fatalf("ReadBlob() failed: %v", err)
	}

	if !bytes.Equal(readContent, content) {
		t.Errorf("expected read content %q, got %q", string(content), string(readContent))
	}
}

func TestBlobDeduplication(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	content := []byte("Duplicate content test")

	hash1, err := WriteBlob(content)
	if err != nil {
		t.Fatalf("first WriteBlob() failed: %v", err)
	}

	hash2, err := WriteBlob(content)
	if err != nil {
		t.Fatalf("second WriteBlob() failed: %v", err)
	}

	if hash1 != hash2 {
		t.Errorf("expected identical hashes for identical content, got %q and %q", hash1, hash2)
	}
}

func TestReadNonExistentBlob(t *testing.T) {
	setupTestRepo(t)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	fakeHash := "0000000000000000000000000000000000000000000000000000000000000000"
	_, err := ReadBlob(fakeHash)
	if err == nil {
		t.Errorf("expected error reading non-existent blob, got nil")
	}
}
