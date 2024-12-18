package add

import (
	"bufio"
	"fmt"
	"got_it/internal/commands/config"
	_init "got_it/internal/commands/init"
	"got_it/internal/logger"
	"got_it/internal/models"
	"got_it/internal/utils"
	"os"
	"path/filepath"
	"strings"
)

type Add struct {
	config *config.Config
	logger *logger.Logger
}

func NewAdd(config *config.Config, logger *logger.Logger) *Add {
	return &Add{
		config: config,
		logger: logger,
	}
}

// Exetue is a shortcut for running the add command
func Execute(files []string, verbose bool) {
	// Fetch debug flag from environment
	debug := os.Getenv("GOT_DEBUG") == "true"
	c := config.NewConfig()
	l := logger.NewLogger(verbose, debug)
	a := NewAdd(c, l)
	a.runAdd(files)
}

// Execute adds files to the staging area
func (a *Add) runAdd(files []string) {
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
	stagedFiles, err := utils.ReadIndex(indexFile)

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
				if !info.IsDir() && !a.isGotDir(path) {
					a.stageFile(path, stagedFiles)
				}
				return nil
			})
		} else {
			a.stageFile(file, stagedFiles)
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

func (a *Add) stageFile(file string, stagedFiles map[string]string) {
	indexFile := a.config.GetIndexPath()
	//check if file is already staged
	isStaged, isChanged := a.checkStagedAndChanged(stagedFiles, file)
	if isStaged {
		return
	}

	if a.ignoreFile(file) {
		return
	}
	hash, err := utils.HashFile(file)
	if err != nil {
		fmt.Printf("Error hashing file %v\n", err)
		return
	}

	err = a.storeFileContet(file, hash)
	if err != nil {
		fmt.Printf("Error storing file content for %s: %v\n", file, err)
		return
	}

	if isChanged {

		// Create a delta between the old and new file contents
		stagedFilePath := stagedFiles[file]
		stagedFile, err := os.ReadFile(stagedFilePath)
		fileBytes, err := os.ReadFile(file)

		delta := utils.CreateDelta(fileBytes, stagedFile)
		if err != nil {
			fmt.Printf("Error creating delta for %s: %v\n", file, err)
			return
		}
		deltaHash := utils.HashContent(string(delta))
		// TODO: Store the delta in the .got/deltas directory
		// avaliate if we need to store the delta in a separate folder

		// Udate the index file with the new hash
		err = a.updateHashChangedFileInIndex(file, deltaHash)
		a.logger.Log("add '%s' (modified)\n", file)
	} else {
		err = addToIndex(indexFile, file, hash)
		if err != nil {
			fmt.Printf("Error adding file %s to index: %v\n", file, err)
			return
		}
		a.logger.Log("add '%s'\n", file)
	}
}

// checkStagedAndChanged checks if the file is already staged and if it has changed
// if it is already staged, it returns true, false
// if it has changed, it returns false, true
// if it is not staged (new file), it returns false, false
func (a *Add) checkStagedAndChanged(stagedFiles map[string]string, file string) (bool, bool) {
	hashStaged, alreadyStaged := stagedFiles[file]
	if alreadyStaged {
		// Get file content and calculate hash
		hashFromFile, err := utils.HashFile(file)
		if err != nil {
			a.logger.Debug("Error hashing file %v\n", err)
			return true, false
		}
		// Check if the hash matches the one in the index
		if hashStaged == hashFromFile {
			a.logger.Log("File %s is already staged\n", file)
			return true, false
		} else {
			a.logger.Debug("File %s has changed\n", file)
			return false, true
		}
	}
	return false, false
}

// updateHashChangedFileOnIndex updates the hash of a file in the index
func (a *Add) updateHashChangedFileInIndex(file string, hash string) error {
	// Open the index file for reading
	indexFile := a.config.GetIndexPath()
	indexFilePath, _ := filepath.Abs(indexFile)
	a.logger.Debug("Index: %s \n", file)
	index, err := os.Open(indexFilePath)
	if err != nil {
		return err
	}
	defer index.Close()
	// Open a temporary file for writing
	gotDir := a.config.GetGotDir()
	gotDirPath, _ := filepath.Abs(gotDir)
	tempFile, err := os.CreateTemp(gotDirPath, "index_*")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())
	// Read the index file line by line
	scanner := bufio.NewScanner(index)
	writer := bufio.NewWriter(tempFile)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		// Field 0 is the path	& field 1 is the hash
		if len(fields) >= 2 {
			a.logger.Debug("\nComparing line '%s'\n", line)
			a.logger.Debug("field '%s' <=> %s", fields[models.IndexKeyValue[models.PATH_KEY]], file)
			if fields[models.IndexKeyValue[models.PATH_KEY]] == file {
				// Update the hash in the line
				fields[models.IndexKeyValue[models.HASH_KEY]] = hash
				line = strings.Join(fields, " ")

			}
		}
		// Write the modified line to the temporary file
		writer.WriteString(line + "\n")
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	writer.Flush()
	// pring the temp file content

	// Replace the original index file with the temporary file
	if err := os.Rename(tempFile.Name(), indexFile); err != nil {
		a.logger.Debug("Error renaming temporary file: %v\n", err)
		return err
	}
	return nil
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

	gotDirPattern, _ := filepath.Abs(a.config.GetGotDir())
	gotDirPattern += fmt.Sprintf("%s*", string(filepath.Separator))
	ignorePatterns = append(ignorePatterns, gotDirPattern)
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
	entry := models.IndexEntry{Path: filePath, Hash: hash}

	file, err := os.OpenFile(indexFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	var entryLine [2]string
	entryLine[models.IndexKeyValue[models.PATH_KEY]] = entry.Path
	entryLine[models.IndexKeyValue[models.HASH_KEY]] = entry.Hash
	entryStr := strings.Join(entryLine[:], " ")

	_, err = fmt.Fprintf(file, "%s \n", entryStr)
	return err
}
