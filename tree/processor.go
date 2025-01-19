package tree

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/samcharles93/treecat/utils"
)

// ProcessOptions contains all options for tree processing
type ProcessOptions struct {
	ExcludePattern string
	IncludePattern string
	StartDir       string
	MaxDepth       int
	CurrentDepth   int
	Force          bool
}

// BuildTree builds a tree structure starting from the given path
func BuildTree(path string, excludePattern, includePattern, startDir string, maxDepth int, force bool) (*Node, error) {
	opts := &ProcessOptions{
		ExcludePattern: excludePattern,
		IncludePattern: includePattern,
		StartDir:       startDir,
		MaxDepth:       maxDepth,
		Force:          force,
	}

	// Add timeout context to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Perform safety checks
	if err := validateDirectorySize(path, opts.Force); err != nil {
		return nil, err
	}

	startTime := time.Now()
	tree, err := buildTreeWithContext(ctx, path, opts)
	if err == nil {
		elapsed := time.Since(startTime)
		fmt.Fprintf(os.Stderr, "\nTree built in %v\n", elapsed)
	}
	return tree, err
}

// processFile handles processing of a single file
func processFile(path string, opts *ProcessOptions) (*Node, error) {
	if !utils.ShouldIncludeFile(path, opts.ExcludePattern, opts.IncludePattern, opts.StartDir) {
		return nil, nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Path:  absPath,
		Name:  filepath.Base(absPath),
		IsDir: false,
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		fmt.Fprintf(os.Stderr, "Warning: cannot read %s: %v\n", path, err)
		return node, nil
	}

	if utils.IsLikelyBinary(content) {
		node.Content = "\n[binary data omitted]\n"
	} else {
		node.Content = string(content)
	}
	return node, nil
}

// processDirectory handles processing of a directory
func processDirectory(ctx context.Context, path string, opts *ProcessOptions) (*Node, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	node := &Node{
		Path:  absPath,
		Name:  filepath.Base(absPath),
		IsDir: true,
	}

	// Skip processing children if we've reached max depth
	if opts.MaxDepth != -1 && opts.CurrentDepth >= opts.MaxDepth {
		return node, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	childOpts := *opts
	childOpts.CurrentDepth++

	err = processChildren(ctx, node, entries, path, &childOpts)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// processEntry handles processing of a single directory entry
func processEntry(ctx context.Context, entryPath string, node *Node, opts *ProcessOptions, errChan chan<- error) {
	info, err := os.Stat(entryPath)
	if err != nil {
		if !os.IsNotExist(err) {
			errChan <- fmt.Errorf("error processing %s: %v", entryPath, err)
		}
		return
	}

	var child *Node
	if info.IsDir() {
		child, err = processDirectory(ctx, entryPath, opts)
	} else {
		child, err = processFile(entryPath, opts)
	}

	if err != nil {
		if !os.IsNotExist(err) {
			errChan <- fmt.Errorf("error processing %s: %v", entryPath, err)
		}
		return
	}

	if child != nil {
		node.mu.Lock()
		node.Children = append(node.Children, child)
		node.mu.Unlock()
	}
}

// processChildren handles concurrent processing of directory entries
func processChildren(ctx context.Context, node *Node, entries []os.DirEntry, path string, opts *ProcessOptions) error {
	const maxWorkers = 4
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup
	errChan := make(chan error, len(entries))

	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		wg.Add(1)
		go func(entryPath string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			processEntry(ctx, entryPath, node, opts, errChan)
		}(childPath)
	}

	wg.Wait()
	close(errChan)

	// Log any errors that occurred
	for err := range errChan {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	return nil
}

// buildTreeWithContext is the main tree building function
func buildTreeWithContext(ctx context.Context, path string, opts *ProcessOptions) (*Node, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	if info.IsDir() {
		return processDirectory(ctx, path, opts)
	}
	return processFile(path, opts)
}
