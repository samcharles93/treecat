# TreeCat 🌳

TreeCat is a command-line tool that displays directory trees with file contents. It provides a visual representation of your directory structure along with the contents of text files, making it easy to explore and document codebases. The tool uses concurrent processing for improved performance on large directories.

## Features

- 📂 Display directory structure with intuitive icons
- 📄 Show file contents inline
- 🔍 Include/exclude files using glob patterns
- 💾 Save output to file
- 🖥️ Cross-platform support (Linux, macOS, Windows)
- 🎯 Skips hidden files and directories by default

## Installation

### Option 1: Download Binary

Download the latest binary for your platform from the [releases page](https://github.com/samcharles93/treecat/releases).

### Option 2: Build from Source

```bash
# Clone the repository
git clone https://github.com/samcharles93/treecat.git
cd treecat

# Build using Go
go build

# Or use make
make build
```

## Usage

```bash
# Display tree for current directory
treecat

# Display tree for specific directory
treecat /path/to/directory

# Include only specific files
treecat -i "*.go"

# Exclude specific files
treecat -e "*.txt"

# Save output to file
treecat -o output.txt
```

### Command Line Flags

- `-e, --exclude`: Pattern to exclude (glob syntax)
- `-i, --include`: Pattern to include (glob syntax)
- `-o, --output`: Output file path

## Performance

TreeCat uses concurrent processing to improve performance when scanning large directories:
- Parallel processing of directory entries
- Worker pool to limit system resource usage
- 30-second timeout to prevent hanging on large directories
- Safe concurrent operations with proper synchronization

## Known Issues

- On Windows, if the directory is very large, the tool may appear to hang. A 30-second timeout has been implemented to prevent this. If the timeout occurs, try using include/exclude patterns to limit the scope:
  ```bash
  # Example: Only include Go files
  treecat -i "*.go"
  ```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

Copyright (c) 2025 Sam Catlow
