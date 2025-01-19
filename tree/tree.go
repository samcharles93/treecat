package tree

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/samcharles93/treecat/utils"
)

func ResolveAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

func BuildTree(path string, excludePattern, includePattern, startDir string) (*TreeNode, error) {
	// Add timeout context to prevent hanging
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return buildTreeWithContext(ctx, path, excludePattern, includePattern, startDir)
}

func buildTreeWithContext(ctx context.Context, path string, excludePattern, includePattern, startDir string) (*TreeNode, error) {
	// Check for context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	node := &TreeNode{
		Path:  absPath,
		Name:  filepath.Base(absPath),
		IsDir: info.IsDir(),
	}

	if !info.IsDir() {
		if utils.ShouldIncludeFile(path, excludePattern, includePattern, startDir) {
			content, err := os.ReadFile(path)
			if err == nil {
				if utils.IsLikelyBinary(content) {
					node.Content = "\n[binary data omitted]\n"
				} else {
					node.Content = string(content)
				}
			}
		}
		return node, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	// Use a worker pool for processing entries
	const maxWorkers = 4
	sem := make(chan struct{}, maxWorkers)
	var wg sync.WaitGroup

	// Process entries concurrently
	errChan := make(chan error, len(entries))
	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		if utils.ShouldIncludeFile(childPath, excludePattern, includePattern, startDir) {
			wg.Add(1)
			go func(entryPath string) {
				defer wg.Done()

				// Acquire semaphore
				sem <- struct{}{}
				defer func() { <-sem }()

				child, err := buildTreeWithContext(ctx, entryPath, excludePattern, includePattern, startDir)
				if err != nil {
					errChan <- fmt.Errorf("error processing %s: %w", entryPath, err)
					return
				}
				if child != nil {
					node.mu.Lock()
					node.Children = append(node.Children, child)
					node.mu.Unlock()
				}
			}(childPath)
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(errChan)

	// Check for any errors
	for err := range errChan {
		// Log the error but continue processing
		fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
	}

	return node, nil
}

func PrintTreeWithOutput(node *TreeNode, prefix string, isLast bool, outputFile *os.File, startDir string) {
	// Calculate relative path
	relativePath, _ := filepath.Rel(startDir, node.Path)

	icon := "📄"
	if node.IsDir {
		if len(node.Children) == 0 {
			icon = "📁"
		} else {
			icon = "📂"
		}
	}

	connector := "├── "
	if isLast {
		connector = "└── "
	}
	fmt.Fprintf(outputFile, "%s%s%s %s\n", prefix, connector, icon, relativePath)

	if node.Content != "" {
		contentPrefix := prefix
		if isLast {
			contentPrefix += "    "
		} else {
			contentPrefix += "│   "
		}

		contentLines := strings.Split(node.Content, "\n")
		for _, line := range contentLines {
			fmt.Fprintf(outputFile, "%s%s\n", contentPrefix, line)
		}

		fmt.Fprintf(outputFile, "%s\n", contentPrefix)
	}

	for i, child := range node.Children {
		var newPrefix string
		if isLast {
			newPrefix = prefix + "    "
		} else {
			newPrefix = prefix + "│   "
		}
		isLastChild := i == len(node.Children)-1
		PrintTreeWithOutput(child, newPrefix, isLastChild, outputFile, startDir)
	}
}
