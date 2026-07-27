package repository

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// Object Types supported by Evolution storage engine.
const (
	ObjectTypeBlob   = "blob"
	ObjectTypeTree   = "tree"
	ObjectTypeCommit = "commit"
)

// HashContent computes the SHA-256 hash of an object prefixed with its type and byte size.
// Format: "<type> <len>\0<content>"
func HashContent(objectType string, data []byte) string {
	header := fmt.Sprintf("%s %d\x00", objectType, len(data))

	hasher := sha256.New()
	hasher.Write([]byte(header))
	hasher.Write(data)

	return hex.EncodeToString(hasher.Sum(nil))
}

// HashRaw computes the raw SHA-256 hash of a byte slice without a header prefix.
func HashRaw(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
