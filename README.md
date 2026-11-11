# ctx-pack

**CLI tool that packages a folder into clean, structured context for LLMs.**

Turn any project directory into high-quality, ready-to-use context. Respects `.gitignore`, supports custom ignore rules, generates clean file trees, and provides useful metadata.

> Status: Early development · Written in Go

---

## Current Capabilities

- Safe recursive directory walking
- **Respects `.gitignore`** automatically
- Extra ignore patterns via `--ignore`
- Read-only (never modifies anything)
- Does **not** follow symlinks by default
- Configurable depth and file count limits (safety)

```bash
# Build
go build -o ctx-pack ./cmd/ctx-pack

# Examples
./ctx-pack .
./ctx-pack --max-depth 3 --max-files 1000 ./my-project
./ctx-pack --ignore "dist,*.log,tmp/*" .
./ctx-pack --version
./ctx-pack --help
```

---

## Vision

`ctx-pack` aims to be a small, reliable, and delightful tool for preparing context for AI systems. It should feel intentional, fast, and well-engineered.

### Core Goals
- Clean and predictable output
- Strong respect for ignore rules
- Useful metadata (token estimates, structure, etc.)
- Excellent CLI experience
- Easy to install and use

---

## Roadmap

### Done
- [x] Basic CLI
- [x] Safe directory walker (no symlink following by default)
- [x] Safety limits (max depth, max files)
- [x] `.gitignore` support
- [x] Custom ignore patterns (`--ignore`)

### Next
- [ ] Clean file tree + content output (Markdown)
- [ ] Max file size limit
- [ ] Binary file detection

### Later
- Approximate token counting
- JSON output mode
- Better formatting / tree view

---

## Safety Principles

- **Read-only**: Never writes or modifies any files
- **No symlink following** by default
- Hard limits on depth and number of files
- Clear error messages

---

## Project Structure

```
ctx-pack/
├── cmd/ctx-pack/         # CLI entry point
├── internal/
│   └── walker/           # Safe directory walking + gitignore
├── go.mod
├── Makefile
└── README.md
```

---

## License

MIT (planned)
