# TreeCat 🌳

Display directory trees with file contents.

## Install

```bash
# Download binary from releases or build from source:
git clone https://github.com/samcharles93/treecat.git
cd treecat
go build
```

## Use

```bash
# Basic usage
treecat                    # Show current directory
treecat /some/path        # Show specific directory

# Options
-d N     # Depth (default: 1, -1 for unlimited)
-i "*.go" # Include only Go files
-e "*.txt" # Exclude text files
-f       # Process large directories
-o file  # Save to file
```

MIT License - Copyright (c) 2025 Sam Catlow
