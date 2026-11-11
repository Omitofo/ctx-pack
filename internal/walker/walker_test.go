package walker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWalkBasic(t *testing.T) {
	dir := t.TempDir()

	// Create a simple structure
	os.MkdirAll(filepath.Join(dir, "src"), 0755)
	os.WriteFile(filepath.Join(dir, "src", "main.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test"), 0644)
	os.WriteFile(filepath.Join(dir, ".gitignore"), []byte("*.log\n"), 0644)
	os.WriteFile(filepath.Join(dir, "app.log"), []byte("log data"), 0644)

	result, err := Walk(Options{
		Root:     dir,
		MaxFiles: 100,
	})
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	// Should find main.go and README.md, but not app.log (gitignored)
	found := map[string]bool{}
	for _, f := range result.Files {
		found[f] = true
	}

	if !found["src/main.go"] && !found[filepath.Join("src", "main.go")] {
		// Accept either slash style depending on OS
		foundMain := false
		for _, f := range result.Files {
			if filepath.Base(f) == "main.go" {
				foundMain = true
				break
			}
		}
		if !foundMain {
			t.Error("expected to find main.go")
		}
	}

	foundReadme := false
	for _, f := range result.Files {
		if filepath.Base(f) == "README.md" {
			foundReadme = true
			break
		}
	}
	if !foundReadme {
		t.Error("expected to find README.md")
	}

	// app.log should be ignored
	for _, f := range result.Files {
		if filepath.Base(f) == "app.log" {
			t.Error("app.log should have been ignored by .gitignore")
		}
	}
}

func TestWalkMaxFiles(t *testing.T) {
	dir := t.TempDir()

	for i := 0; i < 10; i++ {
		name := filepath.Join(dir, "file"+string(rune('a'+i))+".txt")
		os.WriteFile(name, []byte("content"), 0644)
	}

	result, err := Walk(Options{
		Root:     dir,
		MaxFiles: 3,
	})
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	if len(result.Files) > 3 {
		t.Errorf("expected at most 3 files, got %d", len(result.Files))
	}
	if !result.Truncated {
		t.Error("expected Truncated=true when hitting MaxFiles")
	}
}

func TestWalkMaxDepth(t *testing.T) {
	dir := t.TempDir()

	// Create nested structure: a/b/c/file.txt
	nested := filepath.Join(dir, "a", "b", "c")
	os.MkdirAll(nested, 0755)
	os.WriteFile(filepath.Join(nested, "deep.txt"), []byte("deep"), 0644)
	os.WriteFile(filepath.Join(dir, "top.txt"), []byte("top"), 0644)

	result, err := Walk(Options{
		Root:     dir,
		MaxDepth: 1,
		MaxFiles: 100,
	})
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	for _, f := range result.Files {
		if filepath.Base(f) == "deep.txt" {
			t.Error("deep.txt should have been excluded by MaxDepth=1")
		}
	}
}

func TestWalkCtxIgnore(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "keep.txt"), []byte("keep"), 0644)
	os.WriteFile(filepath.Join(dir, "secret.env"), []byte("SECRET=1"), 0644)
	os.WriteFile(filepath.Join(dir, ".ctxignore"), []byte("*.env\n"), 0644)

	result, err := Walk(Options{
		Root:     dir,
		MaxFiles: 100,
	})
	if err != nil {
		t.Fatalf("Walk failed: %v", err)
	}

	for _, f := range result.Files {
		if filepath.Base(f) == "secret.env" {
			t.Error("secret.env should have been ignored by .ctxignore")
		}
	}
}

func TestWalkNonExistent(t *testing.T) {
	_, err := Walk(Options{
		Root: "/this/path/should/not/exist/at/all",
	})
	if err == nil {
		t.Error("expected error for non-existent path")
	}
}
