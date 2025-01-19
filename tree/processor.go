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

// BuildTree builds a tree structure starting from the given path
func BuildTree(path string, excludePattern, includePattern, startDir string, maxDepth int, force bool) (*TreeNode, error) {
	// Add timeout context to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Perform safety checks
	if err := validateDirectorySize(path, force); err != nil {
		return nil, err
	}

	startTime := time.Now()
	tree, err := buildTreeWithContext(ctx, path, excludePattern, includePattern, startDir, maxDepth, 0)
	if err == nil {
		elapsed := time.Since(startTime)
		fmt.Fprintf(os.Stderr, "\nTree built in %v\n", elapsed)
	}
	return tree, err
}

// processFile handles processing of a single file
func processFile(path string, excludePattern, includePattern, startDir string) (*TreeNode, error) {
	if !utils.ShouldIncludeFile(path, excludePattern, includePattern, startDir) {
		return nil, nil
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	node := &TreeNode{
		Path:  absPath,
		Name:  filepath.Base(absPath),
		IsDir: false,
	}

	content, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		// Log error but continue processing
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
func processDirectory(ctx context.Context, path string, excludePattern, includePattern, startDir string, maxDepth, currentDepth int) (*TreeNode, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	node := &TreeNode{
		Path:  absPath,
		Name:  filepath.Base(absPath),
		IsDir: true,
	}

	// Skip processing children if we've reached max depth
	if maxDepth != -1 && currentDepth >= maxDepth {
		return node, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	err = processChildren(ctx, node, entries, path, excludePattern, includePattern, startDir, maxDepth, currentDepth)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// processChildren handles concurrent processing of directory entries
func processChildren(ctx context.Context, node *TreeNode, entries []os.DirEntry, path, excludePattern, includePattern, startDir string, maxDepth, currentDepth int) error {
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

			var child *TreeNode
			var err error

			info, err := os.Stat(entryPath)
			if err != nil {
				if !os.IsNotExist(err) {
					errChan <- fmt.Errorf("error processing %s: %v", entryPath, err)
				}
				return
			}

			if info.IsDir() {
				child, err = processDirectory(ctx, entryPath, excludePattern, includePattern, startDir, maxDepth, currentDepth+1)
			} else {
				child, err = processFile(entryPath, excludePattern, includePattern, startDir)
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
func buildTreeWithContext(ctx context.Context, path string, excludePattern, includePattern, startDir string, maxDepth, currentDepth int) (*TreeNode, error) {
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
		return processDirectory(ctx, path, excludePattern, includePattern, startDir, maxDepth, currentDepth)
	}
	return processFile(path, excludePattern, includePattern, startDir)
}
