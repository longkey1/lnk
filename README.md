# lnkr

A simple CLI tool for managing hard links and symbolic links with configuration files.

## Overview

`lnkr` helps you manage links between local and remote directories using a `.lnkr.toml` configuration file. It automatically handles git exclusions and supports both hard links and symbolic links.

## Installation

```bash
# From source
go build -o lnkr .

# Using go install
go install github.com/longkey1/lnkr@latest
```

## Quick Start

```bash
# 1. Initialize project
lnkr init --remote /backup/project

# 2. Add files to link
lnkr add important.txt
lnkr add config/ --recursive

# 3. Create the links
lnkr link

# 4. Check status
lnkr status

# 5. Remove links when done
lnkr unlink
```

## Commands

### init
Initialize a new lnkr project.

```bash
# Basic initialization
lnkr init

# With remote directory
lnkr init --remote /path/to/remote

# Create remote directory if it doesn't exist
lnkr init --remote /path/to/remote --with-create-remote

# Custom git exclude path
lnkr init --git-exclude-path .gitignore
```

### add
Add files or directories to the link configuration.

```bash
# Add single file (hard link by default)
lnkr add file.txt

# Add directory recursively
lnkr add directory/ --recursive

# Add with symbolic link
lnkr add file.txt --symbolic

# Add from remote directory
lnkr add file.txt --from-remote
```

### link
Create the actual links based on configuration.

```bash
# Create links (local -> remote)
lnkr link

# Create links (remote -> local)
lnkr link --from-remote
```

### unlink
Remove all links from the filesystem.

```bash
lnkr unlink
```

### status
Check the status of configured links.

```bash
lnkr status
```

### remove
Remove entries from the configuration.

```bash
lnkr remove path/to/remove
```

### clean
Remove configuration file and clean up git exclusions.

```bash
lnkr clean
```

## Configuration (.lnkr.toml)

```toml
local = "/workspace"
remote = "/backup/project"
git_exclude_path = ".git/info/exclude"

[[links]]
path = "file.txt"
type = "hard"

[[links]]
path = "config/"
type = "symbolic"
```

## Environment Variables

- `LNKR_REMOTE_ROOT`: Base directory for remote paths (default: `$HOME/.config/lnkr`)
- `LNKR_REMOTE_DEPTH`: Directory levels to include in default remote path (default: 2)

## Link Types

- **Hard Links**: Share the same inode as the original file (default)
- **Symbolic Links**: Point to the original file/directory (use `--symbolic` flag)

## Platform Support

- Linux (AMD64, ARM64, ARMv6, ARMv7)
- macOS (AMD64, ARM64)

Windows is not supported due to filesystem differences.
