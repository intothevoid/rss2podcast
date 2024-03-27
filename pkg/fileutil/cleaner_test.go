package fileutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCleanupFolder(t *testing.T) {
	// Create a temporary folder for testing
	tempFolder := t.TempDir()

	// Create some test files in the temporary folder
	testFiles := []struct {
		name     string
		fileType string
	}{
		{name: "file1.txt", fileType: ".txt"},
		{name: "file2.txt", fileType: ".txt"},
		{name: "file3.jpg", fileType: ".jpg"},
		{name: "file4.jpg", fileType: ".jpg"},
		{name: "file5.png", fileType: ".png"},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tempFolder, tf.name)
		file, err := os.Create(filePath)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		defer file.Close()
	}

	// Define the expected files to be deleted
	expectedDeletedFiles := []string{
		filepath.Join(tempFolder, "file1.txt"),
		filepath.Join(tempFolder, "file2.txt"),
		filepath.Join(tempFolder, "file3.jpg"),
		filepath.Join(tempFolder, "file4.jpg"),
	}

	// Call the CleanupFolder function
	err := CleanupFolder(tempFolder, []string{".txt", ".jpg"})
	if err != nil {
		t.Fatalf("CleanupFolder returned an error: %v", err)
	}

	// Check if the expected files have been deleted
	for _, filePath := range expectedDeletedFiles {
		_, err := os.Stat(filePath)
		if !os.IsNotExist(err) {
			t.Errorf("File %s was not deleted", filePath)
		}
	}
}
