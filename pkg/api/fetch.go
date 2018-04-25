package api

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	dp "github.com/jpopesculian/papercli/pkg/dropbox"
	"github.com/jpopesculian/papercli/pkg/store"
	"log"
)

func Fetch(options *config.CliOptions) {
	store := store.NewStore()
	list, err := dp.DocList(options)
	if err != nil {
		log.Fatal(err)
	}
	folders, err := fetchFolders(list.DocIds, options)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(folders)
	err = store.SaveUpstreamFolders(folders)
	if err != nil {
		log.Fatal(err)
	}
}

type uniqueFolderList struct {
	folderPresence map[store.Id]bool
	folders        []store.Folder
}

func newUniqueFolderList(n int) *uniqueFolderList {
	return &uniqueFolderList{
		folderPresence: map[store.Id]bool{},
		folders:        make([]store.Folder, 0, n),
	}
}
func (list *uniqueFolderList) add(folderResult *dp.FolderResult) {
	folders := folderResultToFolderList(folderResult)
	for _, folder := range folders {
		if !list.folderPresence[folder.Id] {
			list.folderPresence[folder.Id] = true
			list.folders = append(list.folders, folder)
		}
	}
}

func folderResultToFolderList(folders *dp.FolderResult) (result []store.Folder) {
	result = make([]store.Folder, len(folders.Folders))
	for index, folder := range folders.Folders {
		result[index] = store.Folder{
			Id:   store.Id(folder.Id),
			Name: folder.Name,
		}
		if index != 0 {
			result[index].Parent = result[index-1].Id
		}
	}
	return result
}

func fetchFolders(docIds []dp.Id, options *config.CliOptions) (result []store.Folder, err error) {
	n := len(docIds)
	results := make(chan *dp.FolderResult, n)
	errors := make(chan error, n)
	for _, docId := range docIds {
		go func(docId dp.Id) {
			folders, err := dp.FolderInfo(docId, options)
			errors <- err
			results <- folders
		}(docId)
	}
	folderList := newUniqueFolderList((n + 1) * 2)
	for i := 0; i < n; i++ {
		err := <-errors
		if err != nil {
			return nil, err
		}
		folders := <-results
		folderList.add(folders)
	}
	return folderList.folders, nil
}
