package repository

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetObjectPath returns the on-disk path for an object hash (.evolution/objects/xx/yyyy...).
func GetObjectPath(hash string) (string, error) {
	if len(hash) < 3 {
		return "", fmt.Errorf("invalid object hash length %d", len(hash))
	}

	dir := hash[:2]
	filename := hash[2:]

	return filepath.Join(RepositoryDir, ObjectsDir, dir, filename), nil
}

// HasBlob checks if a blob with the given hash already exists in objects storage.
func HasBlob(hash string) bool {
	objectPath, err := GetObjectPath(hash)
	if err != nil {
		return false
	}

	info, err := os.Stat(objectPath)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

// WriteBlob stores file content as a blob object in .evolution/objects/xx/yyyy...
// If the content already exists (identical SHA-256 hash), writing is skipped (deduplication).
func WriteBlob(data []byte) (string, error) {
	hash := HashContent(ObjectTypeBlob, data)

	objectPath, err := GetObjectPath(hash)
	if err != nil {
		return "", err
	}

	// Automatic Deduplication: if object file already exists, return hash without re-writing
	if HasBlob(hash) {
		return hash, nil
	}

	// Ensure subdirectory (first 2 chars) exists
	dir := filepath.Dir(objectPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating object directory %s: %w", dir, err)
	}

	// Write payload with header prefix
	header := fmt.Sprintf("%s %d\x00", ObjectTypeBlob, len(data))
	payload := append([]byte(header), data...)

	if err := os.WriteFile(objectPath, payload, 0444); err != nil { // 0444 = read-only for immutability
		return "", fmt.Errorf("writing blob object: %w", err)
	}

	return hash, nil
}

// ReadBlob reads a blob object from .evolution/objects/xx/yyyy... and strips its header prefix.
func ReadBlob(hash string) ([]byte, error) {
	objectPath, err := GetObjectPath(hash)
	if err != nil {
		return nil, err
	}

	rawPayload, err := os.ReadFile(objectPath)
	if err != nil {
		return nil, fmt.Errorf("reading blob %s: %w", hash[:8], err)
	}

	// Strip header prefix ("blob <size>\0<content>")
	headerEnd := -1
	for i := 0; i < len(rawPayload); i++ {
		if rawPayload[i] == 0 { // Null byte delimiter
			headerEnd = i
			break
		}
	}

	if headerEnd == -1 || headerEnd >= len(rawPayload) {
		return nil, fmt.Errorf("corrupted blob header for object %s", hash[:8])
	}

	content := rawPayload[headerEnd+1:]
	return content, nil
}
