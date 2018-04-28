package files

type Node interface {
	Next() Node
	Prev() Node
	IsRoot() bool
	IsLeaf() bool
	SetChild(Node)
	SetParent(Node)
	Create(string) error
	FsName() string
}
