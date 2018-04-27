package files

import (
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/store"
)

func BuildUpstreamFileTree(document *store.Document, db *store.Store) *DocumentNode {
	documentNode := documentToNode(document)
	var prevNode Node
	var node Node = documentNode
	var entity store.FolderEntity = document
	for entity.InFolder() {
		folder := db.UpstreamFolderById(entity.FolderId())
		prevNode = node
		node = folderToNode(folder)
		node.SetChild(prevNode)
		prevNode.SetParent(node)
		entity = folder
	}
	return documentNode
}
