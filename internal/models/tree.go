package models

// TreeFormatMap is a map of the columns in the tree file
var TreeFormatMap = map[TreeKey]int{
	"mode": int(treeModeValue),
	"hash": int(treeHashValue),
	"type": int(treeTypeValue),
	"name": int(treeNameValue),
}

type treeValue int

const (
	treeModeValue treeValue = iota
	treeTypeValue
	treeHashValue
	treeNameValue
)

// TreeKey is a typed for the collumns' names in the tree file
type TreeKey string

const (
	TK_MODE TreeKey = "mode"
	TK_HASH TreeKey = "hash"
	TK_TYPE TreeKey = "type"
	TK_NAME TreeKey = "name"
)

// TreeEntryType is a type for the tree entres in the tree file
type TreeEntryType string

const (
	TT_BLOB  TreeEntryType = "blob"
	TT_TREE  TreeEntryType = "tree"
	TT_DELTA TreeEntryType = "delta"
)

// TreeEntry is a struct for each line in the tree file
type TreeEntry struct {
	Mode string
	Hash string
	Type string
	Name string
}
