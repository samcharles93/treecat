package tree

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	maxSafeFileCount = 1000     // Maximum number of files to process without force flag
	maxSafeFileSize  = 10485760 // 10MB per file limit without force flag
)

// checkDirectorySize performs safety checks on directory size and file count
func checkDirectorySize(path string) (int, int64, error) {
	var fileCount int
	var totalSize int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileCount++
			totalSize += info.Size()
		}
		return nil
	})

	return fileCount, totalSize, err
}

// validateDirectorySize checks if directory meets safety limits
func validateDirectorySize(path string, force bool) error {
	if force {
		return nil
	}

	fileCount, totalSize, err := checkDirectorySize(path)
	if err != nil {
		return fmt.Errorf("error checking directory size: %w", err)
	}

	if fileCount == 0 {
		return nil
	}

	if fileCount > maxSafeFileCount {
		return fmt.Errorf("directory contains %d files (max safe limit is %d). Use --force to override",
			fileCount, maxSafeFileCount)
	}

	avgFileSize := float64(totalSize) / float64(fileCount)
	if avgFileSize > float64(maxSafeFileSize) {
		return fmt.Errorf("average file size %.2fMB exceeds safe limit of %.2fMB. Use --force to override",
			avgFileSize/1024/1024, float64(maxSafeFileSize)/1024/1024)
	}

	return nil
}
