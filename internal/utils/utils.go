package utils

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"got_it/internal/models"
	"io"
	"os"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
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

func GeneratePatch(oldContent, newContent []byte) []byte {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(oldContent), string(newContent), false)
	patches := dmp.PatchMake(diffs)
	return []byte(dmp.PatchToText(patches))
}

func CreateDelta(oldContent, newContent []byte) []byte {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(oldContent), string(newContent), false)
	return []byte(dmp.DiffPrettyText(diffs))
}

// Read the index file and return a list of file paths
func ReadIndex(indexFile string) (map[string]string, error) {
	stagedFiles := make(map[string]string)
	file, err := os.Open(indexFile)
	if err != nil {
		return stagedFiles, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry models.IndexEntry
		line := scanner.Text()
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			entry.Path = parts[models.IndexKeyValue[models.PathKey]]
			entry.Hash = parts[models.IndexKeyValue[models.HashKey]]
			stagedFiles[entry.Path] = entry.Hash
		}
	}
	if err := scanner.Err(); err != nil {
		return stagedFiles, err
	}
	return stagedFiles, nil
}
