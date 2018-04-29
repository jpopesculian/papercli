package files

import (
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/store"
	"github.com/jpopesculian/papercli/pkg/utils"
	difftool "github.com/sergi/go-diff/diffmatchpatch"
	"log"
	"os"
	"path/filepath"
	"sync"
)

func documentPaths(options *config.CliOptions, fn func(string)) {
	root, err := options.RootDir()
	if err != nil {
		log.Fatal(err)
	}
	var wg sync.WaitGroup
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err != nil || info.IsDir() || filepath.Ext(path) != ".md" {
				return
			}
			fn(path)
		}()
		return nil
	})
	wg.Wait()
}

func doDiff(path string, db *store.Store, options *config.CliOptions) (dmp *difftool.DiffMatchPatch, diffs []difftool.Diff, err error) {
	rel, err := options.RelToRoot(path)
	if err != nil {
		return nil, nil, err
	}
	fileC, errC := utils.ReadFileAsync(path)
	lastPushC := db.LastPushByPath(rel)
	if err = <-errC; err != nil {
		return nil, nil, err
	}
	lastPush := <-lastPushC
	file := <-fileC
	dmp = difftool.New()
	diffs = dmp.DiffMain(string(lastPush), string(file), false)
	return dmp, diffs, err
}

func isEqual(diffs []difftool.Diff) bool {
	for _, diff := range diffs {
		if diff.Type != difftool.DiffEqual {
			return false
		}
	}
	return true
}

func ChangedFiles(db *store.Store, options *config.CliOptions, fn func(string, *difftool.DiffMatchPatch, []difftool.Diff, error)) {
	documentPaths(options, func(path string) {
		dmp, diffs, err := doDiff(path, db, options)
		if err != nil {
			log.Println(err)
			return
		}
		if !isEqual(diffs) {
			fn(path, dmp, diffs, err)
		}
	})
}
