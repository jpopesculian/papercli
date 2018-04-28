package api

import (
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/files"
	"github.com/jpopesculian/papercli/pkg/store"
	// "log"
)

func Update(options *config.CliOptions) {
	db := store.NewStore(options)
	defer db.Close()
	documents, cont := db.UpstreamDocuments()
	for <-cont {
		go func(document *store.Document) {
			tree := files.BuildUpstreamFileTree(document, db)
			files.CreateFile(tree, options)
		}(<-documents)
	}
}
