package utils

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

// hashFile returns the SHA1 hash of the file
func HashFile(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha1.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func HashContent(content string) string {
	hasher := sha1.New()
	hasher.Write([]byte(content))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
