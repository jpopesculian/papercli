package api

import (
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/files"
	"github.com/jpopesculian/papercli/pkg/store"
	difftool "github.com/sergi/go-diff/diffmatchpatch"
	"log"
)

func Push(options *config.CliOptions) {
	db := store.NewStore(options)
	defer db.Close()
	files.ChangedPaths(db, options, func(path string, dmp *difftool.DiffMatchPatch, diffs []difftool.Diff, err error) {
		log.Println(path)
		log.Println(dmp.DiffPrettyText(diffs))
	})
}
