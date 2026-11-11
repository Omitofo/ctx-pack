package walker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	ignore "github.com/sabhiram/go-gitignore"
)

// Options controls how the walker behaves.
// All limits exist for safety.
type Options struct {
	Root           string   // directory to walk
	MaxDepth       int      // 0 = unlimited
	MaxFiles       int      // 0 = unlimited
	FollowSymlinks bool     // default false (safer)
	IgnorePatterns []string // additional gitignore-style patterns
}

// Result holds the outcome of a walk.
type Result struct {
	Root      string
	Files     []string
	Dirs      []string
	Skipped   int
	Truncated bool // true if we hit MaxFiles or MaxDepth
}

// Walk performs a safe, read-only directory walk.
// It never follows symlinks unless explicitly allowed.
// It never writes anything.
// It respects .gitignore and .ctxignore in the root + any extra patterns.
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

	// Build ignore matcher from multiple sources
	var lines []string

	// 1. Load .gitignore from the root if it exists
	loadIgnoreFile(filepath.Join(root, ".gitignore"), &lines)

	// 2. Load .ctxignore from the root if it exists (project-specific for ctx-pack)
	loadIgnoreFile(filepath.Join(root, ".ctxignore"), &lines)

	// 3. Add any extra patterns passed via CLI
	lines = append(lines, opts.IgnorePatterns...)

	var ign *ignore.GitIgnore
	if len(lines) > 0 {
		ign = ignore.CompileIgnoreLines(lines...)
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

		// Skip the root itself
		if rel == "." {
			return nil
		}

		// Respect all ignore sources
		// go-gitignore works best with forward slashes
		relSlash := filepath.ToSlash(rel)
		if ign != nil && ign.MatchesPath(relSlash) {
			if d.IsDir() {
				return fs.SkipDir
			}
			result.Skipped++
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

// loadIgnoreFile appends non-empty lines from a file into the lines slice.
func loadIgnoreFile(path string, lines *[]string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return // file doesn't exist or can't be read — that's fine
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			*lines = append(*lines, line)
		}
	}
}
