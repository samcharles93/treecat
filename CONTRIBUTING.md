# Quick Guide

## Making Changes

```bash
# Start from develop branch
git checkout develop

# Make your changes
vim main.go

# Run tests
go test ./...
# Fix any failures

# Check what changed
git status
git diff

# Commit changes
git add main.go
git commit -m "type: brief description"
# Types: fix, feat, docs, chore, etc.
```

## Testing

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./tree/...
go test ./utils/...

# Run tests with output
go test -v ./...

# Run specific test
go test -v -run TestBuildTree ./tree/...
```

## Versioning

```bash
# Check current version
git describe --tags

# Update version in main.go
Version = "0.1.1" # Match git tag

# Test everything again
go test ./...

# Commit version change
git add main.go
git commit -m "chore: bump version to 0.1.1"

# Create tag
git tag -a v0.1.1 -m "Release v0.1.1 - Brief description"
```

## Building

```bash
# Build Windows binary
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -X main.Version=0.1.1" -trimpath -o dist/treecat-windows-amd64.exe

# Build all platforms
make build-all

# Test the binary
./dist/treecat-windows-amd64.exe --version
```

## Pushing

```bash
# Push changes
git push origin develop

# Push tags
git push origin v0.1.1
```

## Releases

1. Go to GitHub releases
2. Create new release from tag
3. Upload binary from dist/
4. Add brief notes about changes
5. Test downloaded binary works
