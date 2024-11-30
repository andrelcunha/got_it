package models

type IndexEntry struct {
	Path string
	Hash string
}

type indexValue int

const (
	hashValue indexValue = iota // o
	pathValue                   // 1
)

type IndexKey string

const (
	PathKey IndexKey = "path"
	HashKey IndexKey = "hash"
)

var IndexKeyValue = map[IndexKey]int{
	"path": int(pathValue),
	"hash": int(hashValue),
}
