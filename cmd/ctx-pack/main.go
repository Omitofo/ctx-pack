package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Omitofo/ctx-pack/internal/walker"
)

const version = "0.1.0"

func main() {
	var (
		showVersion    = flag.Bool("version", false, "Print version and exit")
		showHelp       = flag.Bool("help", false, "Show help")
		maxDepth       = flag.Int("max-depth", 0, "Maximum directory depth (0 = unlimited)")
		maxFiles       = flag.Int("max-files", 5000, "Maximum number of files to include (safety limit)")
		followSymlinks = flag.Bool("follow-symlinks", false, "Follow symbolic links (disabled by default for safety)")
	)

	flag.Parse()

	if *showVersion {
		fmt.Printf("ctx-pack %s\n", version)
		os.Exit(0)
	}

	if *showHelp || flag.NArg() == 0 {
		printHelp()
		os.Exit(0)
	}

	root := flag.Arg(0)

	// Safety: resolve to absolute path early
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot resolve path: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Scanning: %s\n", absRoot)
	if *maxDepth > 0 {
		fmt.Printf("Max depth: %d\n", *maxDepth)
	}
	fmt.Printf("Max files: %d\n", *maxFiles)
	fmt.Println()

	result, err := walker.Walk(walker.Options{
		Root:           absRoot,
		MaxDepth:       *maxDepth,
		MaxFiles:       *maxFiles,
		FollowSymlinks: *followSymlinks,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Directories: %d\n", len(result.Dirs))
	fmt.Printf("Files:       %d\n", len(result.Files))
	if result.Skipped > 0 {
		fmt.Printf("Skipped:     %d\n", result.Skipped)
	}
	if result.Truncated {
		fmt.Println("\n⚠️  Result was truncated due to safety limits.")
	}

	fmt.Println("\nFile tree (relative):")
	for _, f := range result.Files {
		fmt.Printf("  %s\n", f)
	}
}

func printHelp() {
	fmt.Println(`ctx-pack — package folders into clean context for LLMs

Usage:
  ctx-pack [flags] <directory>

Examples:
  ctx-pack .
  ctx-pack ./my-project
  ctx-pack --max-depth 3 --max-files 1000 ./src

Flags:`)
	flag.PrintDefaults()
	fmt.Println(`
Safety notes:
  - Read-only: never modifies any files
  - Does not follow symlinks by default
  - Has file and depth limits to prevent runaway scans`)
}
