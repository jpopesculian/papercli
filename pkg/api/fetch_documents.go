package api

import (
	dp "github.com/jpopesculian/papercli/pkg/dropbox"
	"github.com/jpopesculian/papercli/pkg/store"
)

func downloadResultToDocument(id dp.Id, download *dp.DownloadResult, folders []store.Folder) *store.Document {
	folderId := store.Id("")
	if len(folders) > 0 {
		folderId = folders[len(folders)-1].Id
	}
	return &store.Document{
		Id:       store.Id(id),
		Title:    download.Title,
		Revision: download.Revision,
		Content:  download.Content,
		Folder:   folderId,
	}
}
