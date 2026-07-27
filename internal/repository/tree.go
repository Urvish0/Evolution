package repository

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Standard file modes in Git / Evolution.
const (
	ModeFile = "100644"
	ModeExec = "100755"
	ModeTree = "040000"
)

// TreeEntry represents a single item inside a directory (file or subdirectory).
type TreeEntry struct {
	Mode string `json:"mode"`
	Type string `json:"type"` // "blob" or "tree"
	Hash string `json:"hash"` // SHA-256 hash of object
	Name string `json:"name"` // filename or directory name
}

// Tree represents a directory snapshot containing entries.
type Tree struct {
	Entries []TreeEntry `json:"entries"`
}

// SerializeTree converts a Tree into a canonical binary byte slice format:
// Each entry: "<mode> <type> <hash> <name>\n"
func SerializeTree(tree Tree) []byte {
	// Sort entries alphabetically by name for deterministic hashing
	sort.Slice(tree.Entries, func(i, j int) bool {
		return tree.Entries[i].Name < tree.Entries[j].Name
	})

	var buf bytes.Buffer
	for _, entry := range tree.Entries {
		line := fmt.Sprintf("%s %s %s %s\n", entry.Mode, entry.Type, entry.Hash, entry.Name)
		buf.WriteString(line)
	}

	return buf.Bytes()
}

// ParseTree parses a canonical binary byte slice back into a Tree struct.
func ParseTree(data []byte) (Tree, error) {
	var tree Tree
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, " ", 4)
		if len(parts) < 4 {
			return tree, fmt.Errorf("corrupted tree entry line %q", line)
		}

		tree.Entries = append(tree.Entries, TreeEntry{
			Mode: parts[0],
			Type: parts[1],
			Hash: parts[2],
			Name: parts[3],
		})
	}

	return tree, nil
}

// WriteTree stores a Tree object in .evolution/objects/xx/yyyy...
func WriteTree(tree Tree) (string, error) {
	serialized := SerializeTree(tree)
	hash := HashContent(ObjectTypeTree, serialized)

	objectPath, err := GetObjectPath(hash)
	if err != nil {
		return "", err
	}

	// Automatic Deduplication
	if HasBlob(hash) {
		return hash, nil
	}

	dir := filepath.Dir(objectPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("creating tree object directory: %w", err)
	}

	header := fmt.Sprintf("%s %d\x00", ObjectTypeTree, len(serialized))
	payload := append([]byte(header), serialized...)

	if err := os.WriteFile(objectPath, payload, 0444); err != nil {
		return "", fmt.Errorf("writing tree object: %w", err)
	}

	return hash, nil
}

// ReadTree loads a Tree object from .evolution/objects/xx/yyyy...
func ReadTree(hash string) (Tree, error) {
	var tree Tree
	objectPath, err := GetObjectPath(hash)
	if err != nil {
		return tree, err
	}

	rawPayload, err := os.ReadFile(objectPath)
	if err != nil {
		return tree, fmt.Errorf("reading tree %s: %w", hash[:8], err)
	}

	// Strip header prefix
	headerEnd := -1
	for i := 0; i < len(rawPayload); i++ {
		if rawPayload[i] == 0 {
			headerEnd = i
			break
		}
	}

	if headerEnd == -1 || headerEnd >= len(rawPayload) {
		return tree, fmt.Errorf("corrupted tree header for object %s", hash[:8])
	}

	serialized := rawPayload[headerEnd+1:]
	return ParseTree(serialized)
}

// BuildTreeFromDirectory recursively walks dirPath, writes Blobs for files and Trees for subdirectories,
// and returns the root Tree's SHA-256 hash.
func BuildTreeFromDirectory(dirPath string) (string, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("reading directory %s: %w", dirPath, err)
	}

	var tree Tree
	for _, entry := range entries {
		name := entry.Name()

		// Skip hidden files and .evolution repository directory
		if name == RepositoryDir || strings.HasPrefix(name, ".") {
			continue
		}

		fullPath := filepath.Join(dirPath, name)

		if entry.IsDir() {
			// Recursive step: Build child tree for subdirectory
			childTreeHash, err := BuildTreeFromDirectory(fullPath)
			if err != nil {
				return "", err
			}

			tree.Entries = append(tree.Entries, TreeEntry{
				Mode: ModeTree,
				Type: ObjectTypeTree,
				Hash: childTreeHash,
				Name: name,
			})
		} else {
			// File step: Read content and store as Blob
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return "", fmt.Errorf("reading file %s: %w", fullPath, err)
			}

			blobHash, err := WriteBlob(content)
			if err != nil {
				return "", err
			}

			tree.Entries = append(tree.Entries, TreeEntry{
				Mode: ModeFile,
				Type: ObjectTypeBlob,
				Hash: blobHash,
				Name: name,
			})
		}
	}

	return WriteTree(tree)
}
