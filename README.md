# ctx-pack

**CLI tool that packages a folder into clean, structured context for LLMs.**

Turn any project directory into high-quality, ready-to-use context. Respects `.gitignore`, supports custom ignore rules, generates clean file trees, and provides useful metadata.

> Status: Early development

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

## Tech Direction

- **Language**: TypeScript (Node.js) or Go — final decision pending
- Clean architecture with clear separation of concerns
- Solid test coverage
- Good documentation and examples

---

## License

MIT (planned)

---

This project lives under the [engineering-sketches](https://github.com/Omitofo/engineering-sketches) collection.
