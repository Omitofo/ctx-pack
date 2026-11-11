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

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "error: missing directory argument")
		fmt.Fprintln(os.Stderr, "\nUsage: ctx-pack [flags] <directory>")
		fmt.Fprintln(os.Stderr, "Run 'ctx-pack --help' for more information.")
		os.Exit(1)
	}

	if flag.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "error: too many arguments (only one directory is supported)")
		os.Exit(1)
	}

	root := flag.Arg(0)
	formatName := strings.ToLower(strings.TrimSpace(*format))

	if formatName != "markdown" && formatName != "json" {
		fmt.Fprintf(os.Stderr, "error: unsupported format %q (supported: markdown, json)\n", *format)
		os.Exit(1)
	}

	if *maxFiles < 0 {
		fmt.Fprintln(os.Stderr, "error: --max-files cannot be negative")
		os.Exit(1)
	}
	if *maxDepth < 0 {
		fmt.Fprintln(os.Stderr, "error: --max-depth cannot be negative")
		os.Exit(1)
	}
	if *maxFileSize < 0 {
		fmt.Fprintln(os.Stderr, "error: --max-file-size cannot be negative")
		os.Exit(1)
	}

	// Resolve to absolute path early
	absRoot, err := filepath.Abs(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot resolve path %q: %v\n", root, err)
		os.Exit(1)
	}

	info, err := os.Stat(absRoot)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "error: path does not exist: %s\n", absRoot)
		} else {
			fmt.Fprintf(os.Stderr, "error: cannot access %s: %v\n", absRoot, err)
		}
		os.Exit(1)
	}
	if !info.IsDir() {
		fmt.Fprintf(os.Stderr, "error: path is not a directory: %s\n", absRoot)
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

	fmt.Fprintf(os.Stderr, "→ Scanning %s\n", absRoot)

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
		fmt.Fprintf(os.Stderr, "error: failed to pack directory: %v\n", err)
		os.Exit(1)
	}

	if len(result.Files) == 0 {
		fmt.Fprintln(os.Stderr, "warning: no files found (check ignore rules or path)")
	}

	// Summary on stderr so it doesn't pollute the output
	fmt.Fprintf(os.Stderr, "→ Found %d files", len(result.Files))
	if result.TotalTokens > 0 {
		fmt.Fprintf(os.Stderr, " (~%d tokens)", result.TotalTokens)
	}
	if result.Skipped > 0 {
		fmt.Fprintf(os.Stderr, ", skipped %d", result.Skipped)
	}
	if result.Truncated {
		fmt.Fprint(os.Stderr, " [truncated]")
	}
	fmt.Fprintln(os.Stderr)

	var output string

	switch formatName {
	case "json":
		output, err = pack.ToJSON(result)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to generate JSON: %v\n", err)
			os.Exit(1)
		}
	default:
		output = pack.ToMarkdown(result)
	}

	if *outputFile != "" {
		err := os.WriteFile(*outputFile, []byte(output), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to write %s: %v\n", *outputFile, err)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, "→ Wrote %s\n", *outputFile)
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
  - Supports markdown (default) and json output formats
  - Progress and summary messages go to stderr; content goes to stdout (or --output)`)
}
