package tree

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PrintTreeWithOutput prints the tree structure to the specified output file
func PrintTreeWithOutput(node *Node, prefix string, isLast bool, outputFile *os.File, startDir string) {
	relativePath, _ := filepath.Rel(startDir, node.Path)

	printNodeLine(node, prefix, isLast, relativePath, outputFile)

	if node.Content != "" {
		printNodeContent(node, prefix, isLast, outputFile)
	}

	if node.IsDir && len(node.Children) > 0 {
		fmt.Fprintln(outputFile)
	}

	for i, child := range node.Children {
		newPrefix := getContentPrefix(prefix, isLast)
		isLastChild := i == len(node.Children)-1
		PrintTreeWithOutput(child, newPrefix, isLastChild, outputFile, startDir)
	}
}

// getConnector determines the appropriate connector string based on whether it's the last node.
func getConnector(isLast bool) string {
	if isLast {
		return "└── "
	}
	return "├── "
}

// printNodeLine prints the line representing the current node (file or folder).
func printNodeLine(node *Node, prefix string, isLast bool, relativePath string, outputFile *os.File) {
	icon := getNodeIcon(node)
	connector := getConnector(isLast)
	fmt.Fprintf(outputFile, "%s%s%s %s\n", prefix, connector, icon, relativePath)
}

// getNodeIcon determines the appropriate icon for the node.
func getNodeIcon(node *Node) string {
	if node.IsDir {
		if len(node.Children) == 0 {
			return "📁"
		}
		return "📂"
	}
	return "📄"
}

// printNodeContent prints the content of the node, if any.
func printNodeContent(node *Node, prefix string, isLast bool, outputFile *os.File) {
	contentPrefix := getContentPrefix(prefix, isLast)
	contentLines := strings.Split(node.Content, "\n")
	for _, line := range contentLines {
		fmt.Fprintf(outputFile, "%s%s\n", contentPrefix, line)
	}
	fmt.Fprintf(outputFile, "%s\n", contentPrefix)
}

// getContentPrefix determines the prefix for content lines.
func getContentPrefix(prefix string, isLast bool) string {
	if isLast {
		return prefix + "    "
	}
	return prefix + "│   "
}
