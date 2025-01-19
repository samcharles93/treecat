package tree

import (
	"path/filepath"
	"sync"
)

// TreeNode represents a file or directory in the tree
type TreeNode struct {
	Path     string      // Absolute path
	Name     string      // Base name
	IsDir    bool        // Whether this is a directory
	Children []*TreeNode // Child nodes (for directories)
	Content  string      // File content (for files)
	mu       sync.Mutex  // Protects Children
}

// ResolveAbsolutePath converts a relative path to absolute
func ResolveAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}
