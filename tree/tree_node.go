package tree

import "sync"

type TreeNode struct {
	Path     string
	Name     string
	IsDir    bool
	Children []*TreeNode
	Content  string
	mu       sync.Mutex // Protects Children
}
