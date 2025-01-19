package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PrintTreeWithOutput prints the tree structure to the specified output file
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

	// Add spacing between folders
	if node.IsDir && len(node.Children) > 0 {
		fmt.Fprintln(outputFile)
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
