package main

import (
	"os"
)

// Cache represents a file-based cache system for storing byte arrays.
type Cache struct {
	FilePath string
}

// NewCache creates a new cache with a specified file path.
func NewCache(filePath string) *Cache {
	return &Cache{FilePath: filePath}
}

func DeleteCache(filePath string) error {
	return os.Remove(filePath)
}

// Set writes a byte array to the cache file.
func (c *Cache) Set(data []byte) error {
	// Write the data to the file, setting the permissions to 0644.
	err := os.WriteFile(c.FilePath, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

// Get reads the byte array from the cache file.
func (c *Cache) Get() ([]byte, error) {
	// Read the data from the file.
	data, err := os.ReadFile(c.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	return data, nil
}
