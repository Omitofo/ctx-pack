# ctx-pack — Simple Guide

This guide is for people who are **not programmers** but want to use `ctx-pack`.

---

## What is this tool?

`ctx-pack` takes a folder on your computer and turns it into a clean, organized text summary.

That summary is useful when you want to:

- Give an AI (like ChatGPT, Claude, or Grok) a clear picture of a project
- Share the structure and content of a folder in one single file
- Understand what’s inside a project without opening every file

It reads the files, skips junk and big/binary files, and produces a nice Markdown (or JSON) document.

**Important:** It only *reads* files. It never changes, deletes, or moves anything.

---

## What you need first

You need **Go** installed on your computer (version 1.22 or newer).

### How to check if Go is installed

Open a terminal (Command Prompt / Terminal / PowerShell) and type:

```bash
go version
```

If you see something like `go version go1.23.0 ...`, you’re good.

If you get an error, install Go from the official site:  
→ [https://go.dev/dl/](https://go.dev/dl/)

---

## How to get the tool

1. Download or clone this project.
2. Open a terminal **inside** the project folder.
3. Build the tool with this command:

```bash
go build -o ctx-pack ./cmd/ctx-pack
```

After this, you will have a file called `ctx-pack` (or `ctx-pack.exe` on Windows).

---

## Basic usage

The simplest way to use it:

```bash
./ctx-pack .
```

(The `.` means “the current folder”.)

This will print a clean summary of the folder to the screen.

---

## Most useful commands

### 1. Scan the current folder
```bash
./ctx-pack .
```
Shows a summary of the current folder.

### 2. Scan a specific folder
```bash
./ctx-pack ./my-project
```
or
```bash
./ctx-pack /path/to/some/folder
```

### 3. Save the result to a file
```bash
./ctx-pack --output context.md .
```
Creates a file called `context.md` with the summary.  
You can then open it or paste it into an AI chat.

### 4. Get JSON instead of Markdown
```bash
./ctx-pack --format json --output context.json .
```
Useful if another program needs to read the result.

### 5. Limit how deep it goes
```bash
./ctx-pack --max-depth 3 .
```
Only looks 3 folders deep. Good for big projects.

### 6. Limit the number of files
```bash
./ctx-pack --max-files 200 .
```
Stops after 200 files (safety feature).

### 7. Ignore extra things
```bash
./ctx-pack --ignore "dist,*.log,temp" .
```
Tells the tool to skip folders/files that match those patterns.

### 8. See the version
```bash
./ctx-pack --version
```

### 9. See help
```bash
./ctx-pack --help
```

---

## Tips for beginners

- Always run the command **from inside** the project folder, or give it the full path.
- The tool respects `.gitignore` and `.ctxignore` files automatically.  
  If something is being ignored and you don’t want that, check those files.
- Big files and binary files (images, videos, executables, etc.) are skipped on purpose.
- Progress messages appear in the terminal. The actual summary goes either to the screen or to the file you specified with `--output`.

---

## Example workflow

1. Open terminal in your project folder.
2. Run:

```bash
./ctx-pack --output context.md .
```

3. Open the file `context.md`.
4. Copy its content and paste it into an AI chat if you want the AI to understand your project.

---

## Safety

- The tool **only reads** files.
- It does **not** follow shortcuts/symlinks by default.
- It has limits so it won’t try to process millions of files by accident.

You’re safe to try it on any normal project folder.
