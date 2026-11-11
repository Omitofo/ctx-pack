# ctx-pack

**CLI tool that packages a folder into clean, structured context for LLMs.**

Turn any project directory into high-quality, ready-to-use context. Respects `.gitignore`, supports custom ignore rules, skips binaries, estimates tokens, and outputs clean Markdown.

> Status: Early development · Written in Go

---

## Current Capabilities

- Safe recursive directory walking
- Respects `.gitignore` automatically
- Extra ignore patterns via `--ignore`
- Skips binary files
- Skips files larger than a configurable size (default 64KB)
- **Approximate token counting** (total + per file)
- Generates clean **Markdown** output (tree + file contents)
- Read-only (never modifies the scanned directory)
- Does **not** follow symlinks by default
- Configurable safety limits

```bash
# Build
go build -o ctx-pack ./cmd/ctx-pack

# Basic usage (prints Markdown to stdout)
./ctx-pack .

# Save to file
./ctx-pack --output context.md .

# With limits
./ctx-pack --max-depth 4 --max-files 300 --max-file-size 32768 ./src

# Extra ignores
./ctx-pack --ignore "dist,*.log,tmp/*" .
```

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
- [x] Custom ignore patterns
- [x] Binary file detection
- [x] Markdown output (tree + contents)
- [x] Approximate token counting

### Next
- [ ] Better tree formatting
- [ ] JSON output mode
- [ ] `.ctxignore` support (in addition to `.gitignore`)

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
│   ├── walker/            # Safe walking + gitignore
│   └── pack/              # Content reading + Markdown + tokens
├── go.mod
├── Makefile
└── README.md
```

---

## License

MIT (planned)
