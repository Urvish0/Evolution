package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// IndexEntry represents a single file staged in the index.
type IndexEntry struct {
	Path       string `json:"path"`
	Hash       string `json:"hash"`
	Mode       string `json:"mode"`
	ModifiedAt string `json:"modified_at"`
}

// Index represents the staging area holding all currently staged files.
type Index struct {
	Entries map[string]IndexEntry `json:"entries"`
}

// NewIndex initializes an empty Index struct.
func NewIndex() Index {
	return Index{
		Entries: make(map[string]IndexEntry),
	}
}

// LoadIndex loads the index from .evolution/index. Returns an empty Index if the file doesn't exist yet.
func LoadIndex() (Index, error) {
	idx := NewIndex()
	indexPath := filepath.Join(RepositoryDir, IndexFile)

	data, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return idx, nil
		}
		return idx, fmt.Errorf("reading index file: %w", err)
	}

	if err := json.Unmarshal(data, &idx); err != nil {
		return idx, fmt.Errorf("parsing index file: %w", err)
	}

	if idx.Entries == nil {
		idx.Entries = make(map[string]IndexEntry)
	}

	return idx, nil
}

// SaveIndex saves the Index to .evolution/index.
func SaveIndex(idx Index) error {
	indexPath := filepath.Join(RepositoryDir, IndexFile)

	data, err := json.MarshalIndent(idx, "", "  ")
	if err != nil {
		return fmt.Errorf("serializing index: %w", err)
	}

	if err := os.WriteFile(indexPath, data, 0644); err != nil {
		return fmt.Errorf("writing index file: %w", err)
	}

	return nil
}

// StagePath stages a file or directory at targetPath into the Index.
func StagePath(targetPath string) error {
	idx, err := LoadIndex()
	if err != nil {
		return err
	}

	cleanPath := filepath.Clean(targetPath)
	info, err := os.Stat(cleanPath)
	if err != nil {
		return fmt.Errorf("staging path %s: %w", cleanPath, err)
	}

	if info.IsDir() {
		// Recursively walk directory and stage all files
		err := filepath.Walk(cleanPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			relPath, err := filepath.Rel(".", path)
			if err != nil {
				relPath = path
			}
			relPath = filepath.ToSlash(relPath)

			// Skip .evolution repository directory and hidden files
			if strings.HasPrefix(relPath, RepositoryDir) || strings.HasPrefix(filepath.Base(relPath), ".") {
				if info.IsDir() && relPath != "." {
					return filepath.SkipDir
				}
				return nil
			}

			if !info.IsDir() {
				if err := stageSingleFile(&idx, path, relPath); err != nil {
					return err
				}
			}

			return nil
		})

		if err != nil {
			return fmt.Errorf("staging directory %s: %w", cleanPath, err)
		}
	} else {
		relPath, err := filepath.Rel(".", cleanPath)
		if err != nil {
			relPath = cleanPath
		}
		relPath = filepath.ToSlash(relPath)

		if err := stageSingleFile(&idx, cleanPath, relPath); err != nil {
			return err
		}
	}

	return SaveIndex(idx)
}

func stageSingleFile(idx *Index, fullPath, relPath string) error {
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("reading file %s: %w", fullPath, err)
	}

	blobHash, err := WriteBlob(content)
	if err != nil {
		return fmt.Errorf("creating blob for %s: %w", relPath, err)
	}

	idx.Entries[relPath] = IndexEntry{
		Path:       relPath,
		Hash:       blobHash,
		Mode:       ModeFile,
		ModifiedAt: time.Now().Format(time.RFC3339),
	}

	return nil
}

// BuildTreeFromIndex constructs a Merkle Tree from the staged index entries and returns the Root Tree Hash.
func BuildTreeFromIndex(idx Index) (string, error) {
	if len(idx.Entries) == 0 {
		return "", nil
	}

	// Group entries by directory hierarchy
	// e.g. "tools/search.json" -> dir "tools", file "search.json"
	type node struct {
		isDir    bool
		blobHash string
		mode     string
		children map[string]*node
	}

	root := &node{isDir: true, children: make(map[string]*node)}

	for relPath, entry := range idx.Entries {
		parts := strings.Split(filepath.ToSlash(relPath), "/")
		curr := root

		for i, part := range parts {
			if i == len(parts)-1 {
				// Leaf file
				curr.children[part] = &node{
					isDir:    false,
					blobHash: entry.Hash,
					mode:     entry.Mode,
				}
			} else {
				// Subdirectory
				if _, exists := curr.children[part]; !exists {
					curr.children[part] = &node{isDir: true, children: make(map[string]*node)}
				}
				curr = curr.children[part]
			}
		}
	}

	var buildNodeTree func(n *node) (string, error)
	buildNodeTree = func(n *node) (string, error) {
		var tree Tree

		// Sort child names for deterministic tree hashing
		var names []string
		for name := range n.children {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			child := n.children[name]
			if child.isDir {
				childHash, err := buildNodeTree(child)
				if err != nil {
					return "", err
				}
				tree.Entries = append(tree.Entries, TreeEntry{
					Mode: ModeTree,
					Type: ObjectTypeTree,
					Hash: childHash,
					Name: name,
				})
			} else {
				tree.Entries = append(tree.Entries, TreeEntry{
					Mode: child.mode,
					Type: ObjectTypeBlob,
					Hash: child.blobHash,
					Name: name,
				})
			}
		}

		return WriteTree(tree)
	}

	return buildNodeTree(root)
}
