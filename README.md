# lnkr

A Link helper for managing hard links and symbolic links with configuration files.

## Overview

`lnkr` is a CLI tool that helps you manage links (hard links and symbolic links) using a configuration file. It allows you to define links in a `.lnkr.toml` file and create/remove them as needed.

## Installation

### From Source

```bash
go build -o lnkr .
```

### From Release

You can download pre-built binaries from the [releases page](https://github.com/longkey1/lnkr/releases).

### Using go install

```bash
go install github.com/longkey1/lnkr@latest
```

## Usage

### Initialization

First, initialize your project with `lnkr`:

```bash
# Basic initialization
lnkr init

# Initialize with remote directory
lnkr init --remote /path/to/remote/directory

# Initialize with remote directory and create it if it doesn't exist
lnkr init --remote /path/to/remote/directory --create-remote
```

This will:
- Create a `.lnkr.toml` configuration file
- Add `.lnkr.toml` to `.git/info/exclude` to prevent it from being tracked

**Remote Path Behavior:**

The `--remote` flag accepts both absolute and relative paths, and they behave differently:

**Absolute Paths:**
- When you specify an absolute path (e.g., `/backup/data`), it is used as-is
- The path must exist unless `--create-remote` is specified
- Example: `lnkr init --remote /backup/project` → remote = `/backup/project`

**Relative Paths:**
- When you specify a relative path, it is resolved relative to the base directory
- The base directory is determined by the `LNK_REMOTE_ROOT` environment variable
- If `LNK_REMOTE_ROOT` is not set, it defaults to `$HOME/.config/lnkr`
- Example: `lnkr init --remote myproject` with `LNK_REMOTE_ROOT=/backup` → remote = `/backup/myproject`

**Default Remote Path (when --remote is not specified):**
- Uses the `LNK_REMOTE_ROOT` environment variable as the base directory
- Appends the project name (current directory name) to the base
- Example: If current directory is `/workspace/myproject` and `LNK_REMOTE_ROOT=/backup`, the remote will be `/backup/myproject`

**Environment Variables:**
- `LNK_REMOTE_DEPTH`: Controls how many directory levels to include in the default remote path (default: 2)
  - Example: `/a/b/c` with `LNK_REMOTE_DEPTH=2` → `b/c`
  - Example: `/a/b/c` with `LNK_REMOTE_DEPTH=1` → `c`
  - Example: `/a/b/c/d` with `LNK_REMOTE_DEPTH=3` → `b/c/d`
- `LNK_REMOTE_ROOT`: Base directory for remote paths (if set, uses `LNK_REMOTE_ROOT/project-name`)

### Adding Links

Add files or directories to your link configuration:

```bash
# Add a single file (hard link by default)
lnkr add /path/to/file.txt

# Add a directory (requires --recursive for hard links)
lnkr add /path/to/directory --recursive

# Add with symbolic link
lnkr add /path/to/file.txt --symbolic

# Add using source directory as base (relative paths)
lnkr add file.txt

# Add using remote directory as base (relative paths)
lnkr add file.txt --from-remote
```

**Options:**
- `--recursive, -r`: Add all subdirectories recursively
- `--symbolic, -s`: Create symbolic link instead of hard link
- `--from-remote`: Use remote directory as base for relative paths

**Note:** 
- `--recursive` and `--symbolic` cannot be used together
- When adding directories, `--recursive` is required for hard links
- Symbolic links can be used with directories

### Creating Links

Create the actual links based on your configuration:

```bash
# Create links using source directory as base
lnkr link

# Create links using remote directory as base
lnkr link --from-remote
```

**Options:**
- `--from-remote`: Use remote directory as base for link source paths

### Removing Links

Remove links from the filesystem:

```bash
# Remove all links defined in .lnkr.toml
lnkr unlink
```

### Removing from Configuration

Remove entries from the configuration file:

```bash
# Remove a specific path and its subdirectories
lnkr remove /path/to/remove
```

### Checking Status

Check the status of your links:

```bash
# Show status of all configured links
lnkr status
```

### Cleaning

Clean up the configuration:

```bash
# Remove .lnkr.toml and clean up git exclusions
lnkr clean
```

## Configuration File (.lnkr.toml)

The configuration file uses TOML format:

```toml
# Base directories
source = "/workspace"
remote = "/path/to/remote/directory"

# Link definitions
[[links]]
path = "relative/path/to/file.txt"
type = "hard"

[[links]]
path = "another/file.txt"
type = "symbolic"
```

**Fields:**
- `source`: Source directory (where links will be created)
- `remote`: Remote directory (alternative base for relative paths)
- `links`: Array of link definitions
  - `path`: Relative path from base directory
  - `type`: Link type (`hard` or `symbolic`)

## Examples

### Basic Workflow

```bash
# 1. Initialize project
lnkr init --remote /backup/data

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

### Using Remote Directory

```bash
# Initialize with remote directory
lnkr init --remote /backup/project

# Add files using remote as base
lnkr add /backup/project/file1.txt --from-remote
lnkr add /backup/project/subdir/ --recursive --from-remote

# Create links using remote as base
lnkr link --from-remote
```

### Using Environment Variables

```bash
# Set environment variables
export LNK_REMOTE_ROOT="/backup"
export LNK_REMOTE_DEPTH=2

# Initialize (will use /backup/project-name)
lnkr init

# Or initialize with custom depth
LNK_REMOTE_DEPTH=1 lnkr init
```

## Link Types

### Hard Links
- Default link type
- Share the same inode as the original file
- Cannot cross filesystem boundaries
- Cannot link directories (creates directories instead)

### Symbolic Links
- Use `--symbolic` flag
- Point to the original file/directory
- Can cross filesystem boundaries
- Can link directories
- Cannot be used with `--recursive` option

## Environment Variables

### LNK_REMOTE_DEPTH
Controls how many directory levels to include in the default remote path when initializing without specifying `--remote`.

- **Default**: 2 (parent directory + current directory)
- **Examples**:
  - `/a/b/c` with `LNK_REMOTE_DEPTH=2` → `b/c`
  - `/a/b/c` with `LNK_REMOTE_DEPTH=1` → `c`
  - `/a/b/c/d` with `LNK_REMOTE_DEPTH=3` → `b/c/d`

### LNK_REMOTE_ROOT
Base directory for remote paths. If set, the tool will use `LNK_REMOTE_ROOT/project-name` as the remote directory.

- **Example**: If `LNK_REMOTE_ROOT=/backup` and current directory is `/workspace/myproject`, the remote will be `/backup/myproject`

## Notes

- The tool automatically sorts paths in the configuration file
- Links are created relative to the source or remote directory
- The configuration file is automatically excluded from git tracking
- All paths in the configuration are stored as relative paths from the base directory
