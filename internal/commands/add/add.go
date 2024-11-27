package add

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"got_it/internal/commands/config"
	_init "got_it/internal/commands/init"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Add struct {
	verbose bool
	config  *config.Config
}

type IndexEntry struct {
	FilePath string
	FileHash string
}

func NewAdd(verbose bool, config *config.Config) *Add {
	return &Add{
		verbose: verbose,
		config:  config,
	}
}

// Exetue is a shortcut for running the add command
func Execute(files []string, verbose bool) {
	c := config.NewConfig()
	a := NewAdd(verbose, c)
	a.runAdd(files, a.verbose)
}

// Execute adds files to the staging area
func (a *Add) runAdd(files []string, verbose bool) {
	i := _init.NewInit()
	// Get the absolute path of the index file
	// indexPath := filepath.Join(a.config.GotDir, a.config.IndexFile)
	// indexFile, err := filepath.Abs(indexPath)
	indexFile := a.config.GetIndexPath()

	// Ensure the .got directory exists
	if !i.IsInitialized() {
		return
	}

	// Get staged files
	stagedFiles, err := readIndex(indexFile)

	// Get the absolute path of the repository root
	repoRoot, err := filepath.Abs(".")
	if err != nil {
		fmt.Println("Error getting repository root:", err)
		return
	}

	// Add files to the staging area
	for _, file := range files {
		// Get the absolute path of the file
		absFile, err := filepath.Abs(file)
		if err != nil {
			fmt.Printf("Error getting absolute path of %v\n", err)
			continue
		}

		// Check if the file is within the repository
		if repoRoot != absFile {
			if mached, err := filepath.Match(repoRoot+"/*", absFile); err != nil || !mached {
				fmt.Printf("Error: %s is outside the repository\n", absFile)
				fmt.Println("root:", repoRoot)
				continue
			}
		}

		// Get file information
		fileInfo, err := os.Stat(file)
		if err != nil {
			fmt.Printf("Error %v\n", err)
			continue
		}

		if fileInfo.IsDir() {
			filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if !info.IsDir() {
					a.stageFile(path, stagedFiles, verbose)
				}
				return nil
			})
		} else {
			a.stageFile(file, stagedFiles, verbose)
		}
	}
}

func (a *Add) isGotDir(file string) bool {
	// if filepath has prefix = gotDir, skip it
	fileAbs, _ := filepath.Abs(file)
	gotDirAbs, _ := filepath.Abs(a.config.GetGotDir())
	gotDirAbs += string(filepath.Separator)
	return strings.HasPrefix(fileAbs, gotDirAbs)

}

func (a *Add) stageFile(file string, stagedFiles map[string]string, verbose bool) {
	indexFile := a.config.GetIndexPath()
	//check if file is already staged
	if _, alreadyStaged := stagedFiles[file]; alreadyStaged {
		if verbose {
			fmt.Printf("File %s is already staged\n", file)
		}
		return
	}

	if a.ignoreFile(file) {
		return
	}
	hash, err := hashFile(file)
	if err != nil {
		fmt.Printf("Error hashing file %v\n", err)
		return
	}

	err = a.storeFileContet(file, hash)
	if err != nil {
		fmt.Printf("Error storing file content for %s: %v\n", file, err)
		return
	}

	err = addToIndex(indexFile, file, hash)
	if err != nil {
		fmt.Printf("Error adding file %s to index: %v\n", file, err)
		return
	}
	if verbose {
		fmt.Printf("add '%s'\n", file)
	}
}

// ignoreFile tests if file matches the ignore patterns on .gotignore
func (a *Add) ignoreFile(file string) bool {
	shallIgnore := false

	if isEssentialFile(file, config.GetEssentilFiles()) {
		return false
	}

	// Read .gotignore file
	ignoreFile, err := os.Open(".gotignore")
	if err != nil {
		return false
	}
	defer ignoreFile.Close()
	scanner := bufio.NewScanner(ignoreFile)

	var ignorePatterns, negatePatterns []string

	for scanner.Scan() {
		line := scanner.Text()
		pattern := strings.TrimSpace(strings.Split(line, "#")[0])

		if pattern == "" {
			continue
		}

		if strings.HasPrefix(pattern, "!") {
			negatePatterns = append(negatePatterns, pattern)
		} else {
			ignorePatterns = append(ignorePatterns, pattern)
		}
	}

	// append the gotDir to the ignorePatterns
	ignorePatterns = append(ignorePatterns, a.config.GetGotDir())

	shallIgnore = matchPatterns(ignorePatterns, file)

	// If the file matches any negate patterns, do NOT ignore it
	if matchPatterns(negatePatterns, file) {
		shallIgnore = false
	}

	return shallIgnore
}

func matchPatterns(ignorePatterns []string, file string) bool {
	isMatch := false
	for _, pattern := range ignorePatterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		for _, match := range matches {
			if file == match || isSubPath(match, file) {
				isMatch = true
			}
		}
	}
	return isMatch
}

func isSubPath(dir, file string) bool {
	rel, err := filepath.Rel(dir, file)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "..") && !strings.Contains(rel, "../")
}

// Essential files that should never be ignored
func isEssentialFile(file string, essentials []string) bool {

	// Check if the file is in the list of essential files
	for _, essential := range essentials {
		if filepath.Base(file) == essential {
			return true
		}
	}
	return false
}

// hashFile returns the SHA1 hash of the file
func hashFile(filePath string) (string, error) {
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

// storeFileContent saves the file content in the .got/objects directory
func (a *Add) storeFileContet(filePath, hash string) error {
	// take the firt two characters of the hash as the directory name
	objDir := filepath.Join(a.config.GetGotDir(), "objects", hash[:2])

	if err := os.MkdirAll(objDir, 0755); err != nil {
		return err
	}

	// take the rest of the hash as the file name
	objFile := filepath.Join(objDir, hash[2:])

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	return os.WriteFile(objFile, fileContent, 0644)
}

func addToIndex(indexFile, filePath, hash string) error {
	entry := IndexEntry{FilePath: filePath, FileHash: hash}

	file, err := os.OpenFile(indexFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "%s %s\n", entry.FilePath, entry.FileHash)
	return err
}

// Read the index file and return a list of file paths
func readIndex(indexFile string) (map[string]string, error) {
	stagedFiles := make(map[string]string)
	file, err := os.Open(indexFile)
	if err != nil {
		return stagedFiles, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry IndexEntry
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) == 2 {
			entry.FilePath = parts[0]
			entry.FileHash = parts[1]
			stagedFiles[entry.FilePath] = entry.FileHash
		}
	}
	if err := scanner.Err(); err != nil {
		return stagedFiles, err
	}
	return stagedFiles, nil
}
