package walker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Options controls how the walker behaves.
// All limits exist for safety.
type Options struct {
	Root       string // directory to walk
	MaxDepth   int    // 0 = unlimited
	MaxFiles   int    // 0 = unlimited
	FollowSymlinks bool // default false (safer)
}

// Result holds the outcome of a walk.
type Result struct {
	Root       string
	Files      []string
	Dirs       []string
	Skipped    int
	Truncated  bool // true if we hit MaxFiles or MaxDepth
}

// Walk performs a safe, read-only directory walk.
// It never follows symlinks unless explicitly allowed.
// It never writes anything.
func Walk(opts Options) (*Result, error) {
	root, err := filepath.Abs(opts.Root)
	if err != nil {
		return nil, fmt.Errorf("resolve root: %w", err)
	}

	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("stat root: %w", err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("root is not a directory: %s", root)
	}

	result := &Result{
		Root: root,
	}

	fileCount := 0

	err = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			// Don't fail the whole walk on permission errors etc.
			result.Skipped++
			return nil
		}

		// Safety: never follow symlinks unless explicitly enabled
		if d.Type()&fs.ModeSymlink != 0 && !opts.FollowSymlinks {
			if d.IsDir() {
				return fs.SkipDir
			}
			result.Skipped++
			return nil
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			result.Skipped++
			return nil
		}

		// Skip the root itself in the lists
		if rel == "." {
			return nil
		}

		// Depth check
		depth := strings.Count(rel, string(os.PathSeparator))
		if opts.MaxDepth > 0 && depth > opts.MaxDepth {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			result.Dirs = append(result.Dirs, rel)
			return nil
		}

		// File limit safety
		if opts.MaxFiles > 0 && fileCount >= opts.MaxFiles {
			result.Truncated = true
			return fs.SkipAll
		}

		result.Files = append(result.Files, rel)
		fileCount++
		return nil
	})

	if err != nil {
		return result, fmt.Errorf("walk: %w", err)
	}

	return result, nil
}
