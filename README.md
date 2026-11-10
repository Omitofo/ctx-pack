# ctx-pack

**CLI tool that packages a folder into clean, structured context for LLMs.**

Turn any project directory into high-quality, ready-to-use context. Respects `.gitignore`, supports custom ignore rules, generates clean file trees, and provides useful metadata.

> Status: Early development · Written in Go

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

## Planned Features

### MVP
- Recursive directory walking
- Full `.gitignore` support
- Custom ignore patterns
- Clean file tree generation
- File content inclusion with size limits
- Markdown output

### Later
- Approximate token counting
- JSON output mode
- Configurable depth and file size limits
- Smart binary file detection
- Optional file summaries

---

## Development

### Requirements
- Go 1.22+

### Build
```bash
go build -o ctx-pack ./cmd/ctx-pack
```

### Run
```bash
./ctx-pack --version
```

---

## Project Structure

```
ctx-pack/
├── cmd/
│   └── ctx-pack/       # Entry point
├── internal/           # Private packages (coming soon)
├── go.mod
└── README.md
```

---

## License

MIT (planned)

---

This project lives under the [engineering-sketches](https://github.com/Omitofo/engineering-sketches) collection.
