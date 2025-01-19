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

	// If there's an exclude pattern, check if the file matches it
	if excludePattern != "" {
		matched, err := filepath.Match(excludePattern, filepath.Base(path))
		if err != nil {
			fmt.Println("Error in filepath.Match for exclude pattern:", err)
			return true
		}
		if matched {
			return false
		}
	}

	// If there's an include pattern, the file must match it
	if includePattern != "" {
		matched, err := filepath.Match(includePattern, filepath.Base(path))
		if err != nil {
			fmt.Println("Error in filepath.Match for include pattern:", err)
			return false
		}
		return matched
	}

	return true
}
