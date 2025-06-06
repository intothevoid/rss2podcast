package fileutil

import (
	"os"
	"path/filepath"
	"strings"
)

// CleanupFolder accepts a list of file types and a folder path. It will delete
// all files in the folder that have the specified file types.
func CleanupFolder(folderPath string, fileTypes []string) error {
	// Get file paths of all files in the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	// Iterate over all files in the folder
	for _, file := range files {
		// Check if the file is of the specified file types
		for _, fileType := range fileTypes {
			if strings.HasSuffix(file.Name(), fileType) {
				// Delete the file
				err := os.Remove(filepath.Join(folderPath, file.Name()))
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// MoveFilesToArchiveFolder moves files of specified types to an archive folder
func MoveFilesToArchiveFolder(folderPath string, fileTypes []string) error {
	// Create archive folder if it doesn't exist
	archivePath := filepath.Join(folderPath, "archive")
	if err := os.MkdirAll(archivePath, 0755); err != nil {
		return err
	}

	// Get file paths of all files in the folder
	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	// Iterate over all files in the folder
	for _, file := range files {
		// Skip if it's a directory
		if file.IsDir() {
			continue
		}

		// Check if the file is of the specified file types
		for _, fileType := range fileTypes {
			if strings.HasSuffix(file.Name(), fileType) {
				// Move the file to archive folder
				oldPath := filepath.Join(folderPath, file.Name())
				newPath := filepath.Join(archivePath, file.Name())
				if err := os.Rename(oldPath, newPath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
