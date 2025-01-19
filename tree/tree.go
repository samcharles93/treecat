package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/samcharles93/treecat/utils"
)

func ResolveAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

func BuildTree(path string, excludePattern, includePattern, startDir string) (*TreeNode, error) {
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

	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		if utils.ShouldIncludeFile(childPath, excludePattern, includePattern, startDir) {
			child, err := BuildTree(childPath, excludePattern, includePattern, startDir)
			if err != nil {
				// Handle error (e.g., log it)
			}
			if child != nil {
				node.Children = append(node.Children, child)
			}
		}
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
