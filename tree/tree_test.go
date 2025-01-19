package tree

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestFiles(t *testing.T, tempDir string) {
	testFiles := map[string]string{
		"file1.txt":            "content1",
		"dir1/file2.txt":       "content2",
		"dir1/dir2/file3.txt":  "content3",
		"dir1/dir2/binary.bin": string([]byte{0x00, 0xFF}),
	}

	for path, content := range testFiles {
		fullPath := filepath.Join(tempDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestBuildTree(t *testing.T) {
	// Test directory structure creation
	tempDir := t.TempDir()
	createTestFiles(t, tempDir)

	// Test cases
	tests := []struct {
		name           string
		excludePattern string
		includePattern string
		wantFiles      int
	}{
		{
			name:           "No filters",
			excludePattern: "",
			includePattern: "",
			wantFiles:      4,
		},
		{
			name:           "Exclude txt files",
			excludePattern: "*.txt",
			includePattern: "",
			wantFiles:      1,
		},
		{
			name:           "Include only txt files",
			excludePattern: "",
			includePattern: "*.txt",
			wantFiles:      3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tree, err := BuildTree(tempDir, tc.excludePattern, tc.includePattern, tempDir, -1, true) // -1 for unlimited depth, force=true for tests
			if err != nil {
				t.Fatal(err)
			}

			// Count files (non-directory nodes)
			var count int
			var countFiles func(*Node)
			countFiles = func(node *Node) {
				if !node.IsDir {
					count++
				}
				for _, child := range node.Children {
					countFiles(child)
				}
			}
			countFiles(tree)

			if count != tc.wantFiles {
				t.Errorf("got %d files, want %d", count, tc.wantFiles)
			}
		})
	}
}
