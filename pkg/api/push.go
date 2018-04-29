package api

import (
	"errors"
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	dp "github.com/jpopesculian/papercli/pkg/dropbox"
	"github.com/jpopesculian/papercli/pkg/files"
	"github.com/jpopesculian/papercli/pkg/store"
	difftool "github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func pushCreate(content []byte, folderId store.Id, db *store.Store, options *config.CliOptions) (*dp.UpdateResult, error) {
	log.Println("create")
	return dp.Create(content, dp.Id(folderId), options)
}

func pushUpdate(content []byte, docId store.Id, folderId store.Id, db *store.Store, options *config.CliOptions) (*dp.UpdateResult, error) {
	revision := db.UpstreamRevisionById(docId)
	return dp.Update(content, dp.Id(docId), revision, options)
}

func doPush(path string, db *store.Store, options *config.CliOptions) error {
	rel, err := options.RelToRoot(path)
	if err != nil {
		return err
	}
	folder := filepath.Dir(rel)
	folderId := store.Id("")
	if folder != "." {
		folderId := <-db.FolderIdByPath(folder)
		if len(folderId) < 1 {
			err = errors.New("Folder '" + folder + "' doesn't exist on server! Need to `pull` after creating it remotely.")
			return err
		}
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	docId := <-db.DocIdByPath(rel)
	var result *dp.UpdateResult
	if len(docId) > 0 {
		result, err = pushUpdate(content, docId, folderId, db, options)
	} else {
		// possibly find with difftool if renamed?
		result, err = pushCreate(content, folderId, db, options)
	}
	if err != nil {
		log.Println("result fail")
		return err
	}
	if err = fetchOne(result.DocId, result.Title, db, options); err != nil {
		return err
	}
	// TODO: hangs for some reason
	// err = updateOne(store.Id(result.DocId), db, options)
	// if err != nil {
	// 	return err
	// }
	return os.Remove(path)
}

func Push(options *config.CliOptions) {
	Fetch(options)
	db := store.NewStore(options)
	files.ChangedFiles(db, options, func(path string, dmp *difftool.DiffMatchPatch, diffs []difftool.Diff, err error) {
		if err := doPush(path, db, options); err != nil {
			log.Println(err)
		}
	})
	db.Close()
	Update(options)
}
