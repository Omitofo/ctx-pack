# ctx-pack

**CLI tool that packages a folder into clean, structured context for LLMs.**

Turn any project directory into high-quality, ready-to-use context. Respects `.gitignore` and `.ctxignore`, supports custom ignore rules, skips binaries, estimates tokens, and outputs clean Markdown or JSON.

> Status: Early development · Written in Go

**New to this?** → Read the simple [Beginner Guide](GUIDE.md)

---

## Current Capabilities

- Safe recursive directory walking
- Respects `.gitignore` and **`.ctxignore`** automatically
- Extra ignore patterns via `--ignore`
- Skips binary files
- Skips files larger than a configurable size (default 64KB)
- Approximate token counting (total + per file)
- Hierarchical file tree in Markdown output
- **Markdown** and **JSON** output formats
- Read-only (never modifies the scanned directory)
- Does **not** follow symlinks by default
- Configurable safety limits

```bash
# Build
go build -o ctx-pack ./cmd/ctx-pack

# Run tests
make test
# or: go test ./...

# Markdown (default)
./ctx-pack .
./ctx-pack --output context.md .

# JSON
./ctx-pack --format json .
./ctx-pack --format json --output context.json .

# With limits
./ctx-pack --max-depth 4 --max-files 300 --max-file-size 32768 ./src

# Extra ignores
./ctx-pack --ignore "dist,*.log,tmp/*" .
```

### Ignore files

`ctx-pack` automatically loads ignore patterns from:

1. `.gitignore` (if present)
2. `.ctxignore` (if present) — use this for context-specific exclusions without touching your main `.gitignore`
3. Patterns passed via `--ignore`

---

## Vision

`ctx-pack` aims to be a small, reliable, and delightful tool for preparing context for AI systems. It should feel intentional, fast, and well-engineered.

---

## Roadmap

### Done
- [x] Basic CLI
- [x] Safe directory walker
- [x] Safety limits (depth, files, file size)
- [x] `.gitignore` support
- [x] `.ctxignore` support
- [x] Custom ignore patterns
- [x] Binary file detection
- [x] Markdown output (tree + contents)
- [x] Approximate token counting
- [x] JSON output mode
- [x] Hierarchical tree formatting
- [x] Basic tests
- [x] Improved error messages & UX feedback
- [x] Beginner-friendly guide

### Next
- LICENSE (MIT)
- Version bump / release readiness

---

## Safety Principles

- **Read-only** on the scanned directory
- No symlink following by default
- Hard limits on depth, number of files, and file size
- Binary files are detected and skipped
- Clear feedback on what was skipped / truncated

---

## Project Structure

```
ctx-pack/
├── cmd/ctx-pack/          # CLI entry point
├── internal/
│   ├── walker/            # Safe walking + gitignore + ctxignore
│   └── pack/              # Content reading + Markdown/JSON + tokens + tree
├── go.mod
├── Makefile
├── GUIDE.md               # Beginner guide
└── README.md
```

---

## License

MIT (planned)
