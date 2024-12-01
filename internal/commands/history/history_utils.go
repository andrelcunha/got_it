package history

import (
	"fmt"
	"got_it/internal/commands/config"
	"got_it/internal/logger"
	"got_it/internal/models"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// ReadRefFromHEAD reads the current commit from the HEAD file
func ReadRefFromHEAD(conf *config.Config, logger *logger.Logger) (string, error) {
	// Get the current commit from the HEAD
	headPath := filepath.Join(conf.GotDir, "HEAD")
	headRefBytes, err := os.ReadFile(headPath)
	if err != nil {
		logger.Debug("Error reading HEAD file: %s", err)
		return "", err
	}
	headRef := string(headRefBytes)

	//find the prefix "refs: " in  the headRef
	if !strings.HasPrefix(string(headRef), "ref: ") {
		logger.Debug("unespected HEAD format: %s", headRef)
		return "", err
	}

	//remove the prefix "ref: "
	headRef = strings.TrimSpace(headRef)
	headRef = headRef[5:]
	headRef = filepath.Join(conf.GotDir, headRef)

	return headRef, nil
}

// getFirstCommitHash returns the hash of the parent commit (the HEAD commit)
// and the name of the file pointed to by the HEAD
func GetFirstCommitHash(conf *config.Config, logger *logger.Logger) (string, string, error) {
	/// Read the current commit from the HEAD file
	headRef, err := ReadRefFromHEAD(conf, logger)
	if err != nil {
		return "", "", err
	}
	// get HEAD name (branch) from headRef
	headBranch := filepath.Base(headRef)

	// Verrify if the file exists
	if _, err := os.Stat(headRef); os.IsNotExist(err) {
		logger.Debug("File does not exist: %s", headRef)
		return "", headBranch, err
	}

	// Read the content of the file pointed to by the HEAD reference
	commitHashBytes, err := os.ReadFile(headRef)
	if err != nil {
		logger.Debug("Error reading commit file: %s", err)
		return "", headBranch, err
	}

	return string(commitHashBytes), headBranch, nil
}

// reconstruct file content from deltas
func reconstructFileContent(conf *config.Config, logger *logger.Logger, filename string) (string, error) {
	// create a stack of treesEntries
	deltasHashes := []string{}

	// Get the first parent commit hash
	HEAD_Hash, _, err := GetFirstCommitHash(conf, logger)
	if err != nil {
		return "", err
	}
	// stop loop if 'parent' field in the commit metadata does not exist or type is 'blob' in the tree
	content, err := getBlobContentOrKeepDiving(conf, logger, &deltasHashes, filename, HEAD_Hash)

	if deltasHashes != nil {
		// reconstruct the file content from the deltas
		content, err = aplyDeltas(conf, logger, content, deltasHashes)
		if err != nil {
			return "", err
		}
	}

	return content, nil
}

func getBlobContentOrKeepDiving(conf *config.Config, logger *logger.Logger, stack *[]string, wantedFilename, commitHash string) (string, error) {
	// get metadata from commitHash
	parser := models.NewCommitDataParser(logger)
	rawMetadata, err := getContentFromHash(conf, logger, commitHash)
	if err != nil {
		return "", err
	}
	parsedMetadata, err := parser.Parse(rawMetadata)
	if err != nil {
		return "", err
	}

	treeContent, err := getContentFromHash(conf, logger, parsedMetadata.Tree)
	if err != nil {
		return "", err
	}
	// get the file fileHash from the tree content and push it to the stack
	// if the type is "blob", there is nothing else to do, fust return the content
	fileHash, fileType, err := findFileInTree(treeContent, wantedFilename)
	if err != nil {
		return "", err
	}
	if fileType == string(models.TT_BLOB) {
		logger.Debug("File is a	blob")
		content, err := getContentFromHash(conf, logger, fileHash)
		if err != nil {
			return "", err
		}
		return content, nil
	} else if fileType == string(models.TT_DELTA) {
		// stack the  delta hash
		*stack = append(*stack, fileHash)
		// keep diving:
		// - read the metadata of the parent commit and start again
		parentHash := parsedMetadata.Parent
		getBlobContentOrKeepDiving(conf, logger, stack, wantedFilename, parentHash)
	}
	return "", fmt.Errorf("file type not supported")
}

// getContentFromHash returns the content of the commit (or tree) object file
func getContentFromHash(conf *config.Config, logger *logger.Logger, commitHash string) (string, error) {
	file := filepath.Join(conf.GotDir, "objects", commitHash[:2], commitHash[2:])
	file, err := filepath.Abs(file)
	if err != nil {
		logger.Debug("Error getting absolute path: %s", err)
		return "", err
	}
	obj, err := os.Open(file)
	if err != nil {
		logger.Debug("Error opening object file: %s", err)
		return "", err
	}
	defer obj.Close()
	content, err := os.ReadFile(file)
	if err != nil {
		logger.Debug("Error reading object file: %s", err)
		return "", err
	}
	return string(content), nil
}

// findFileInTree returns the hash of the file in the tree and the type (blob or delta)
func findFileInTree(treeContent string, fileName string) (string, string, error) {
	// break path into parts
	pathParts := strings.Split(fileName, "/")
	pathIndex := 0

	lines := strings.Split(treeContent, "\n")
	for _, line := range lines {
		parts := strings.Split(line, " ")
		if len(parts) < 4 {
			continue
		}
		treeEntry := models.TreeEntry{
			// mode: parts[models.TreeFormatMap[models.TK_MODE]],
			Hash: parts[models.TreeFormatMap[models.TK_HASH]],
			Type: parts[models.TreeFormatMap[models.TK_TYPE]],
			Name: parts[models.TreeFormatMap[models.TK_NAME]],
		}
		if pathIndex < len(pathParts) && treeEntry.Name == pathParts[pathIndex] {
			if treeEntry.Type == string(models.TT_TREE) {
				pathIndex++
				continue
			}
			if pathIndex == len(pathParts)-1 {
				if treeEntry.Type == string(models.TT_BLOB) || treeEntry.Type == string(models.TT_DELTA) {
					return treeEntry.Hash, treeEntry.Type, nil
				}
			}
		}
	}

	return "", "", fmt.Errorf("file %s not found in tree", fileName)
}

// ReconstructFileContent reconstructs the file content from the deltas
func aplyDeltas(conf *config.Config, logger *logger.Logger, baseContent string, deltas []string) (string, error) {

	currentContent := baseContent
	dmp := diffmatchpatch.New()
	for _, deltaHash := range deltas {
		deltaContent, err := getContentFromHash(conf, logger, deltaHash)
		if err != nil {
			return "", err
		}
		patches, err := dmp.PatchFromText(deltaContent)
		if err != nil {
			return "", err
		}
		currentContent, _ = dmp.PatchApply(patches, currentContent)
	}
	return currentContent, nil
}
