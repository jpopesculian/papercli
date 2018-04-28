package api

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/files"
	"github.com/jpopesculian/papercli/pkg/store"
	// "log"
)

func Update(options *config.CliOptions) {
	store := store.NewStore(options)
	id := store.FetchFirstId()
	document := store.UpstreamDocumentById(id)
	tree := files.BuildUpstreamFileTree(document, store)
	spew.Dump(tree)
}
