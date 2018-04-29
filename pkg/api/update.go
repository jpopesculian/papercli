package api

import (
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/files"
	"github.com/jpopesculian/papercli/pkg/store"
	"log"
)

func saveUpstreamDocuments(db *store.Store, options *config.CliOptions, fn func(*store.Document, *files.DocumentNode, error)) {
	db.UpstreamDocuments(func(document *store.Document) {
		tree := files.BuildUpstreamFileTree(document, db)
		err := files.CreateFile(tree, options)
		if err == nil {
			document.Path = files.RelativePath(tree)
		}
		fn(document, tree, err)
	})
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
