package api

import (
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/files"
	"github.com/jpopesculian/papercli/pkg/store"
	"log"
)

func Update(options *config.CliOptions) {
	db := store.NewStore(options)
	defer db.Close()
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
	folderList := newUniqueFolderList(0)
	for i := 0; i < count; i++ {
		err := <-errors
		if err != nil {
			log.Println(err)
		}
		document := <-savedDocuments
		if document != nil {
			db.SaveLocalDocument(document)
		}
		tree := <-savedTrees
		if tree != nil {
			folderList.add(files.TreeToFolderList(tree))
		}
	}
	if err := db.SaveLocalFolders(folderList.folders); err != nil {
		log.Println(err)
	}
}
