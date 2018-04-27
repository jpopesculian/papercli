package files

type Node interface {
	Next() Node
	Prev() Node
	IsRoot() bool
	SetChild(Node)
	SetParent(Node)
}
