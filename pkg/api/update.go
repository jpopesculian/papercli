package api

import (
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/files"
	"github.com/jpopesculian/papercli/pkg/store"
	"log"
)

func updateFile(document *store.Document, db *store.Store, options *config.CliOptions) (*store.Document, *files.DocumentNode, error) {
	tree := files.BuildUpstreamFileTree(document, db)
	err := files.CreateFile(tree, options)
	if err == nil {
		document.Path = files.RelativePath(tree)
	}
	return document, tree, err
}

func saveUpstreamDocuments(db *store.Store, options *config.CliOptions, fn func(*store.Document, *files.DocumentNode, error)) {
	db.UpstreamDocuments(func(document *store.Document) {
		document, tree, err := updateFile(document, db, options)
		fn(document, tree, err)
	})
}

func updateOne(docId store.Id, db *store.Store, options *config.CliOptions) error {
	document := db.UpstreamDocumentById(docId)
	document, tree, err := updateFile(document, db, options)
	if err != nil {
		return err
	}
	if err = db.SaveLocalDocument(document); err != nil {
		return err
	}
	if err = db.SaveLocalFolders(files.TreeToFolderList(tree)); err != nil {
		return err
	}
	return nil
}

func Update(options *config.CliOptions) {
	db := store.NewStore(options)
	defer db.Close()
	folderList := newUniqueFolderList(0)
	saveUpstreamDocuments(db, options, func(document *store.Document, tree *files.DocumentNode, err error) {
		if err != nil {
			log.Println(err)
			return
		}
		if err = db.SaveLocalDocument(document); err != nil {
			log.Println(err)
		}
		folderList.add(files.TreeToFolderList(tree))
	})
	if err := db.SaveLocalFolders(folderList.folders); err != nil {
		log.Println(err)
	}
}
