package pack

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/Omitofo/ctx-pack/internal/walker"
)

// Options for packing.
type Options struct {
	Root           string
	MaxDepth       int
	MaxFiles       int
	MaxFileSize    int64 // bytes. 0 = default (64KB)
	FollowSymlinks bool
	IgnorePatterns []string
	IncludeContent bool
}

// FileEntry represents one file in the pack result.
type FileEntry struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	Content string `json:"content,omitempty"`
	Binary  bool   `json:"binary,omitempty"`
	Skipped bool   `json:"skipped,omitempty"`
	Tokens  int    `json:"tokens,omitempty"`
}

// Result is the full pack output.
type Result struct {
	Root        string      `json:"root"`
	Files       []FileEntry `json:"files"`
	Dirs        []string    `json:"dirs,omitempty"`
	Skipped     int         `json:"skipped"`
	Truncated   bool        `json:"truncated,omitempty"`
	TotalTokens int         `json:"total_tokens"`
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

		// Simple binary detection
		if isBinary(data) {
			entry.Binary = true
			entry.Skipped = true
			result.Files = append(result.Files, entry)
			continue
		}

		entry.Content = string(data)
		entry.Tokens = estimateTokens(entry.Content)
		result.TotalTokens += entry.Tokens
		result.Files = append(result.Files, entry)
	}

	return result, nil
}

// estimateTokens provides a reasonable approximation.
// Real tokenizers vary, but chars/4 is a widely used rough estimate for English/code.
func estimateTokens(text string) int {
	if text == "" {
		return 0
	}
	// +3/4 so it rounds up
	return (len(text) + 3) / 4
}

// isBinary uses a simple heuristic.
func isBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	if strings.Contains(string(data), "\x00") {
		return true
	}
	return !utf8.Valid(data)
}

// ToMarkdown generates a clean Markdown representation.
func ToMarkdown(r *Result) string {
	var b strings.Builder

	b.WriteString("# Context Pack\n\n")
	b.WriteString(fmt.Sprintf("**Root:** `%s`\n\n", r.Root))
	b.WriteString(fmt.Sprintf("- Directories: %d\n", len(r.Dirs)))
	b.WriteString(fmt.Sprintf("- Files: %d\n", len(r.Files)))
	b.WriteString(fmt.Sprintf("- Estimated tokens: **%d**\n", r.TotalTokens))
	if r.Skipped > 0 {
		b.WriteString(fmt.Sprintf("- Skipped: %d\n", r.Skipped))
	}
	if r.Truncated {
		b.WriteString("- **Truncated** due to safety limits\n")
	}
	b.WriteString("\n> Token count is an approximation (≈ characters / 4).\n\n")
	b.WriteString("---\n\n")

	// File tree
	b.WriteString("## File Tree\n\n```\n")
	for _, f := range r.Files {
		line := f.Path
		if f.Binary {
			line += "  (binary)"
		} else if f.Skipped {
			line += "  (skipped - too large)"
		} else if f.Tokens > 0 {
			line += fmt.Sprintf("  (~%d tokens)", f.Tokens)
		}
		b.WriteString(line + "\n")
	}
	b.WriteString("```\n\n")

	// File contents
	b.WriteString("## File Contents\n\n")
	for _, f := range r.Files {
		if f.Content == "" {
			continue
		}

		b.WriteString(fmt.Sprintf("### `%s`\n\n", f.Path))
		b.WriteString(fmt.Sprintf("_≈ %d tokens_\n\n", f.Tokens))
		b.WriteString("```" + detectLang(f.Path) + "\n")
		b.WriteString(f.Content)
		if !strings.HasSuffix(f.Content, "\n") {
			b.WriteString("\n")
		}
		b.WriteString("```\n\n")
	}

	return b.String()
}

// ToJSON generates a pretty-printed JSON representation.
func ToJSON(r *Result) (string, error) {
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return "", fmt.Errorf("marshal json: %w", err)
	}
	return string(data) + "\n", nil
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
