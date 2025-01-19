package utils

import (
	"fmt"
	"path/filepath"
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
	if excludePattern != "" {
		excludePath := filepath.Join(startDir, excludePattern)
		matches, err := filepath.Glob(excludePath)
		if err != nil {
			fmt.Println("Error in filepath.Glob for exclude pattern:", err)
			return true
		}

		absPath, _ := filepath.Abs(path)
		for _, match := range matches {
			if absPath == match {
				return false
			}
		}
	}

	if includePattern != "" {
		includePath := filepath.Join(startDir, includePattern)
		matches, err := filepath.Glob(includePath)
		if err != nil {
			fmt.Println("Error in filepath.Glob for include pattern:", err)
			return false
		}

		absPath, _ := filepath.Abs(path)
		included := false
		for _, match := range matches {
			if absPath == match {
				included = true
				break
			}
		}
		if !included {
			return false
		}
	}

	return true
}
