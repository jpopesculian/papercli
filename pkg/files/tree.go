package files

import (
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/store"
	"path/filepath"
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

func CreateFile(node Node, options *config.CliOptions) error {
	dir, err := options.RootDir()
	if err != nil {
		return err
	}
	for !node.IsRoot() {
		node = node.Prev()
	}
	for node != nil {
		err = node.Create(dir)
		if err != nil {
			return err
		}
		dir = filepath.Join(dir, node.FsName())
		node = node.Next()
	}
	return nil
}

func RelativePath(node Node) string {
	path := ""
	for node != nil {
		path = filepath.Join(node.FsName(), path)
		node = node.Prev()
	}
	return path
}

func TreeToFolderList(node Node) (folders []store.Folder) {
	folders = make([]store.Folder, 0)
	for node != nil {
		folderNode, ok := node.(*FolderNode)
		if ok {
			folders = append(folders, *nodeToFolder(folderNode))
		}
		node = node.Prev()
	}
	return folders
}
