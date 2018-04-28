package api

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/files"
	"github.com/jpopesculian/papercli/pkg/store"
	"log"
)

func saveUpstreamDocuments(db *store.Store, options *config.CliOptions) (int, chan *store.Document, chan *files.DocumentNode, chan error) {
	documents, cont := db.UpstreamDocuments()
	count := 0
	savedDocuments := make(chan *store.Document)
	savedTrees := make(chan *files.DocumentNode)
	errors := make(chan error)
	for <-cont {
		go func(document *store.Document) {
			tree := files.BuildUpstreamFileTree(document, db)
			err := files.CreateFile(tree, options)
			errors <- err
			if err == nil {
				savedDocuments <- document
				savedTrees <- tree
			} else {
				savedDocuments <- nil
				savedTrees <- nil
			}
		}(<-documents)
		count++
	}
	return count, savedDocuments, savedTrees, errors
}

func Update(options *config.CliOptions) {
	db := store.NewStore(options)
	defer db.Close()
	count, documents, trees, errors := saveUpstreamDocuments(db, options)
	folderList := newUniqueFolderList(0)
	for i := 0; i < count; i++ {
		err := <-errors
		if err != nil {
			log.Println(err)
		}
		document := <-documents
		if document != nil {
			db.SaveLocalDocument(document)
		}
		tree := <-trees
		if tree != nil {
			folderList.add(files.TreeToFolderList(tree))
		}
	}
	if err := db.SaveLocalFolders(folderList.folders); err != nil {
		log.Println(err)
	}
}
