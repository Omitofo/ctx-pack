package pack

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		name string
		text string
		want int
	}{
		{"empty", "", 0},
		{"short", "hello", 2},         // 5 chars → (5+3)/4 = 2
		{"exact", "abcd", 1},          // 4 chars → 1
		{"longer", "abcdefgh", 2},     // 8 chars → 2
		{"code-ish", "func main() {}", 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := estimateTokens(tt.text)
			if got != tt.want {
				t.Errorf("estimateTokens(%q) = %d, want %d", tt.text, got, tt.want)
			}
		})
	}
}

func TestIsBinary(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want bool
	}{
		{"empty", []byte{}, false},
		{"text", []byte("hello world"), false},
		{"utf8", []byte("café 你好"), false},
		{"null-byte", []byte("hello\x00world"), true},
		{"invalid-utf8", []byte{0xff, 0xfe, 0xfd}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isBinary(tt.data)
			if got != tt.want {
				t.Errorf("isBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDetectLang(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"main.go", "go"},
		{"app.tsx", "typescript"},
		{"script.js", "javascript"},
		{"utils.py", "python"},
		{"README.md", "markdown"},
		{"data.json", "json"},
		{"config.yaml", "yaml"},
		{"style.css", "css"},
		{"run.sh", "bash"},
		{"unknown.xyz", ""},
		{"noext", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := detectLang(tt.path)
			if got != tt.want {
				t.Errorf("detectLang(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestRenderTree(t *testing.T) {
	files := []FileEntry{
		{Path: "cmd/ctx-pack/main.go", Tokens: 100},
		{Path: "internal/pack/pack.go", Tokens: 200},
		{Path: "README.md", Tokens: 50},
	}

	tree := renderTree(files)

	// Basic structure checks
	if !strings.Contains(tree, "cmd") {
		t.Error("tree should contain 'cmd'")
	}
	if !strings.Contains(tree, "main.go") {
		t.Error("tree should contain 'main.go'")
	}
	if !strings.Contains(tree, "internal") {
		t.Error("tree should contain 'internal'")
	}
	if !strings.Contains(tree, "README.md") {
		t.Error("tree should contain 'README.md'")
	}
	if !strings.Contains(tree, "├──") && !strings.Contains(tree, "└──") {
		t.Error("tree should use box-drawing characters")
	}
}

func TestPackBasic(t *testing.T) {
	// Create a temporary directory with a couple of files
	dir := t.TempDir()

	// Write a small text file
	err := os.WriteFile(filepath.Join(dir, "hello.txt"), []byte("hello world"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	// Write a "binary" file (contains null byte)
	err = os.WriteFile(filepath.Join(dir, "binary.bin"), []byte("hello\x00world"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	result, err := Pack(Options{
		Root:           dir,
		MaxFiles:       100,
		MaxFileSize:    64 * 1024,
		IncludeContent: true,
	})
	if err != nil {
		t.Fatalf("Pack failed: %v", err)
	}

	if len(result.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(result.Files))
	}

	var foundText, foundBinary bool
	for _, f := range result.Files {
		switch f.Path {
		case "hello.txt":
			foundText = true
			if f.Content != "hello world" {
				t.Errorf("unexpected content for hello.txt: %q", f.Content)
			}
			if f.Tokens == 0 {
				t.Error("expected non-zero tokens for hello.txt")
			}
		case "binary.bin":
			foundBinary = true
			if !f.Binary {
				t.Error("expected binary.bin to be marked as binary")
			}
			if f.Content != "" {
				t.Error("binary file should have empty content")
			}
		}
	}

	if !foundText || !foundBinary {
		t.Errorf("missing expected files (text=%v, binary=%v)", foundText, foundBinary)
	}
}

func TestToMarkdownContainsTreeAndContent(t *testing.T) {
	r := &Result{
		Root:        "/tmp/project",
		TotalTokens: 150,
		Files: []FileEntry{
			{Path: "main.go", Content: "package main\n", Tokens: 100},
			{Path: "README.md", Content: "# Hello\n", Tokens: 50},
		},
	}

	md := ToMarkdown(r)

	if !strings.Contains(md, "# Context Pack") {
		t.Error("missing title")
	}
	if !strings.Contains(md, "Estimated tokens") {
		t.Error("missing token summary")
	}
	if !strings.Contains(md, "## File Tree") {
		t.Error("missing file tree section")
	}
	if !strings.Contains(md, "## File Contents") {
		t.Error("missing file contents section")
	}
	if !strings.Contains(md, "main.go") {
		t.Error("missing main.go in output")
	}
}
