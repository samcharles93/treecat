package tree

type TreeNode struct {
	Path     string
	Name     string
	IsDir    bool
	Children []*TreeNode
	Content  string
}
