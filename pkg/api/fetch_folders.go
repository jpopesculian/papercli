package api

import (
	dp "github.com/jpopesculian/papercli/pkg/dropbox"
	"github.com/jpopesculian/papercli/pkg/store"
)

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
func (list *uniqueFolderList) add(folders []store.Folder) {
	for _, folder := range folders {
		if !list.folderPresence[folder.Id] {
			list.folderPresence[folder.Id] = true
			list.folders = append(list.folders, folder)
		}
	}
}

func folderInfoToFolderList(folderInfo *dp.FolderResult) (result []store.Folder) {
	result = make([]store.Folder, len(folderInfo.Folders))
	for index, folder := range folderInfo.Folders {
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
