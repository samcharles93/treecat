package utils

import (
	"fmt"
	"path/filepath"
	"strings"
)

func IsLikelyBinary(data []byte) bool {
	if len(data) == 0 {
		return false
	}
	nonPrintableCount := 0
	for _, b := range data {
		if b < 32 || b > 126 {
			nonPrintableCount++
		}
	}
	return float64(nonPrintableCount)/float64(len(data)) > 0.2
}

func ShouldIncludeFile(path, excludePattern, includePattern, startDir string) bool {
	// Skip hidden files and directories by default
	parts := strings.Split(path, string(filepath.Separator))
	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return false
		}
	}

	// Get absolute paths for both startDir and path
	absStartDir, err := filepath.Abs(startDir)
	if err != nil {
		fmt.Println("Error getting absolute start dir path:", err)
		return false
	}
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("Error getting absolute file path:", err)
		return false
	}

	// Get relative path for pattern matching
	relPath, err := filepath.Rel(absStartDir, absPath)
	if err != nil {
		fmt.Println("Error getting relative path:", err)
		return false
	}

	// Convert Windows paths to forward slashes for consistent matching
	relPath = filepath.ToSlash(relPath)

	// If there's an exclude pattern, check if the file matches it
	if excludePattern != "" {
		// Convert Windows-style pattern to forward slashes
		excludePattern = filepath.ToSlash(excludePattern)
		matched, err := filepath.Match(excludePattern, relPath)
		if err != nil {
			fmt.Println("Error in filepath.Match for exclude pattern:", err)
			return false
		}
		if matched {
			return false
		}
	}

	// If there's an include pattern, the file must match it
	if includePattern != "" {
		// Convert Windows-style pattern to forward slashes
		includePattern = filepath.ToSlash(includePattern)
		matched, err := filepath.Match(includePattern, relPath)
		if err != nil {
			fmt.Println("Error in filepath.Match for include pattern:", err)
			return false
		}
		return matched
	}

	return true
}
