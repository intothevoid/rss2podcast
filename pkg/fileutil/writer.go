package fileutil

import (
	"os"
	"sync"
)

var (
	fileMutexes = make(map[string]*sync.Mutex)
	mutexLock   sync.Mutex
)

func getFileMutex(filename string) *sync.Mutex {
	mutexLock.Lock()
	defer mutexLock.Unlock()

	if mutex, exists := fileMutexes[filename]; exists {
		return mutex
	}

	mutex := &sync.Mutex{}
	fileMutexes[filename] = mutex
	return mutex
}

func FlushMapToFile(filename string, textMap map[string]string) error {
	// Get the mutex for this file
	mutex := getFileMutex(filename)
	mutex.Lock()
	defer mutex.Unlock()

	// Open file in append mode, create if it doesn't exist
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Create a temporary buffer to hold all content
	var buffer []byte
	for title, text := range textMap {
		buffer = append(buffer, []byte(title+"\n"+text+"\n")...)
	}

	// Write the entire buffer in one operation
	if _, err = f.Write(buffer); err != nil {
		return err
	}

	// Ensure the write is flushed to disk
	return f.Sync()
}

func FlushStringToFile(filename string, content string) error {
	// Get the mutex for this file
	mutex := getFileMutex(filename)
	mutex.Lock()
	defer mutex.Unlock()

	// Open file in append mode, create if it doesn't exist
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Convert the string content to a byte slice
	buffer := []byte(content + "\n")

	// Write the entire buffer in one operation
	if _, err = f.Write(buffer); err != nil {
		return err
	}

	// Ensure the write is flushed to disk
	return f.Sync()
}

func ReadFileContent(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
