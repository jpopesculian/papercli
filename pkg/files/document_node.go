package files

import (
	"github.com/jpopesculian/papercli/pkg/store"
)

type DocumentNode struct {
	Name    string
	Parent  Node
	Content []byte
}

func (node *DocumentNode) Next() Node {
	return nil
}

func (node *DocumentNode) Prev() Node {
	return node.Parent
}

func (node *DocumentNode) IsRoot() bool {
	return node.Prev() == nil
}

func (node *DocumentNode) SetParent(parent Node) {
	node.Parent = parent
}

func (node *DocumentNode) SetChild(child Node) {
}

func documentToNode(document *store.Document) *DocumentNode {
	return &DocumentNode{
		Name:    document.Title,
		Content: document.Content,
	}
}
