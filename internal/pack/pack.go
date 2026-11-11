package pack

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/Omitofo/ctx-pack/internal/walker"
)

// Options for packing.
type Options struct {
	Root            string
	MaxDepth        int
	MaxFiles        int
	MaxFileSize     int64 // bytes. 0 = default (64KB)
	FollowSymlinks  bool
	IgnorePatterns  []string
	IncludeContent  bool
}

// FileEntry represents one file in the pack result.
type FileEntry struct {
	Path    string
	Size    int64
	Content string // empty if binary or too large
	Binary  bool
	Skipped bool // true if content was skipped (too large / binary)
}

// Result is the full pack output.
type Result struct {
	Root      string
	Files     []FileEntry
	Dirs      []string
	Skipped   int
	Truncated bool
}

// Pack walks the directory and optionally reads file contents.
func Pack(opts Options) (*Result, error) {
	if opts.MaxFileSize <= 0 {
		opts.MaxFileSize = 64 * 1024 // 64 KB default
	}

	wResult, err := walker.Walk(walker.Options{
		Root:           opts.Root,
		MaxDepth:       opts.MaxDepth,
		MaxFiles:       opts.MaxFiles,
		FollowSymlinks: opts.FollowSymlinks,
		IgnorePatterns: opts.IgnorePatterns,
	})
	if err != nil {
		return nil, err
	}

	result := &Result{
		Root:      wResult.Root,
		Dirs:      wResult.Dirs,
		Skipped:   wResult.Skipped,
		Truncated: wResult.Truncated,
	}

	for _, rel := range wResult.Files {
		fullPath := filepath.Join(wResult.Root, rel)

		info, err := os.Stat(fullPath)
		if err != nil {
			result.Skipped++
			continue
		}

		entry := FileEntry{
			Path: rel,
			Size: info.Size(),
		}

		if !opts.IncludeContent {
			result.Files = append(result.Files, entry)
			continue
		}

		// Safety: skip files that are too large
		if info.Size() > opts.MaxFileSize {
			entry.Skipped = true
			result.Files = append(result.Files, entry)
			continue
		}

		data, err := os.ReadFile(fullPath)
		if err != nil {
			result.Skipped++
			continue
		}

		// Simple binary detection: contains null byte or invalid UTF-8
		if isBinary(data) {
			entry.Binary = true
			entry.Skipped = true
			result.Files = append(result.Files, entry)
			continue
		}

		entry.Content = string(data)
		result.Files = append(result.Files, entry)
	}

	return result, nil
}

// isBinary uses a simple heuristic.
func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	// Null byte is a strong indicator
	if strings.Contains(string(data), "\x00") {
		return true
	}
	// Invalid UTF-8
	return !utf8.Valid(data)
}

// ToMarkdown generates a clean Markdown representation.
func ToMarkdown(r *Result) string {
	var b strings.Builder

	b.WriteString("# Context Pack\n\n")
	b.WriteString(fmt.Sprintf("**Root:** `%s`\n\n", r.Root))
	b.WriteString(fmt.Sprintf("- Directories: %d\n", len(r.Dirs)))
	b.WriteString(fmt.Sprintf("- Files: %d\n", len(r.Files)))
	if r.Skipped > 0 {
		b.WriteString(fmt.Sprintf("- Skipped: %d\n", r.Skipped))
	}
	if r.Truncated {
		b.WriteString("- **Truncated** due to safety limits\n")
	}
	b.WriteString("\n---\n\n")

	// File tree
	b.WriteString("## File Tree\n\n```\n")
	for _, f := range r.Files {
		b.WriteString(f.Path)
		if f.Binary {
			b.WriteString("  (binary)")
		} else if f.Skipped {
			b.WriteString("  (skipped - too large)")
		}
		b.WriteString("\n")
	}
	b.WriteString("```\n\n")

	// File contents
	b.WriteString("## File Contents\n\n")
	for _, f := range r.Files {
		if f.Content == "" {
			continue
		}

		b.WriteString(fmt.Sprintf("### `%s`\n\n", f.Path))
		b.WriteString("```" + detectLang(f.Path) + "\n")
		b.WriteString(f.Content)
		if !strings.HasSuffix(f.Content, "\n") {
			b.WriteString("\n")
		}
		b.WriteString("```\n\n")
	}

	return b.String()
}

// detectLang returns a markdown language hint based on extension.
func detectLang(path string) string {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return "go"
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx":
		return "javascript"
	case ".py":
		return "python"
	case ".rs":
		return "rust"
	case ".md":
		return "markdown"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".html", ".htm":
		return "html"
	case ".css":
		return "css"
	case ".sh":
		return "bash"
	case ".sql":
		return "sql"
	default:
		return ""
	}
}
