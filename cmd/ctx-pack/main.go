package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Omitofo/ctx-pack/internal/pack"
)

const version = "0.1.0"

func main() {
	var (
		showVersion    = flag.Bool("version", false, "Print version and exit")
		showHelp       = flag.Bool("help", false, "Show help")
		maxDepth       = flag.Int("max-depth", 0, "Maximum directory depth (0 = unlimited)")
		maxFiles       = flag.Int("max-files", 2000, "Maximum number of files to include (safety limit)")
		maxFileSize    = flag.Int64("max-file-size", 64*1024, "Maximum file size in bytes to include content (default 64KB)")
		followSymlinks = flag.Bool("follow-symlinks", false, "Follow symbolic links (disabled by default for safety)")
		ignoreList     = flag.String("ignore", "", "Additional gitignore-style patterns (comma-separated)")
		outputFile     = flag.String("output", "", "Write output to this file instead of stdout")
		format         = flag.String("format", "markdown", "Output format: markdown | json")
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
	formatName := strings.ToLower(strings.TrimSpace(*format))

	if formatName != "markdown" && formatName != "json" {
		fmt.Fprintf(os.Stderr, "error: unsupported format %q (use markdown or json)\n", *format)
		os.Exit(1)
	}

	// Safety: resolve to absolute path early
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot resolve path: %v\n", err)
		os.Exit(1)
	}

	var extraIgnores []string
	if *ignoreList != "" {
		for _, p := range strings.Split(*ignoreList, ",") {
			p = strings.TrimSpace(p)
			if p != "" {
				extraIgnores = append(extraIgnores, p)
			}
		}
	}

	fmt.Fprintf(os.Stderr, "Scanning: %s\n", absRoot)

	result, err := pack.Pack(pack.Options{
		Root:           absRoot,
		MaxDepth:       *maxDepth,
		MaxFiles:       *maxFiles,
		MaxFileSize:    *maxFileSize,
		FollowSymlinks: *followSymlinks,
		IgnorePatterns: extraIgnores,
		IncludeContent: true,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var output string

	switch formatName {
	case "json":
		output, err = pack.ToJSON(result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		output = pack.ToMarkdown(result)
	}

	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error writing output: %v\n", err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "Wrote %s (%d files, ~%d tokens)\n", *outputFile, len(result.Files), result.TotalTokens)
	} else {
		fmt.Print(output)
	}
}

func printHelp() {
	fmt.Println(`ctx-pack — package folders into clean context for LLMs

Usage:
  ctx-pack [flags] <directory>

Examples:
  ctx-pack .
  ctx-pack --format json .
  ctx-pack --format json --output context.json .
  ctx-pack --max-depth 3 --max-files 500 ./src
  ctx-pack --ignore "dist,*.log" --output context.md .

Flags:`)
	flag.PrintDefaults()
	fmt.Println(`
Notes:
  - Automatically respects .gitignore and .ctxignore in the target directory
  - Use .ctxignore for context-specific exclusions without touching .gitignore
  - Read-only: never modifies any files in the scanned directory
  - Does not follow symlinks by default
  - Skips binary files and files larger than --max-file-size
  - Supports markdown (default) and json output formats`)
}
