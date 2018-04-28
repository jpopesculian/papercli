package files

import (
	"github.com/jpopesculian/papercli/pkg/store"
	"io/ioutil"
	"path/filepath"
)

type DocumentNode struct {
	Id      store.Id
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

func (node *DocumentNode) IsLeaf() bool {
	return node.Next() == nil
}

func (node *DocumentNode) SetParent(parent Node) {
	node.Parent = parent
}

func (node *DocumentNode) SetChild(child Node) {
}

func (node *DocumentNode) Create(dir string) error {
	path := filepath.Join(dir, node.FsName())
	return ioutil.WriteFile(path, node.Content, 0644)
}

func (node *DocumentNode) FsName() string {
	return node.Name + ".md"
}

func documentToNode(document *store.Document) *DocumentNode {
	return &DocumentNode{
		Id:      document.Id,
		Name:    document.Title,
		Content: document.Content,
	}
}
