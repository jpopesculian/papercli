package api

import (
	// "github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	dp "github.com/jpopesculian/papercli/pkg/dropbox"
	"github.com/jpopesculian/papercli/pkg/store"
	"log"
)

func Fetch(options *config.CliOptions) {
	store := store.NewStore(options)
	defer store.Close()
	list, err := dp.DocList(options)
	if err != nil {
		log.Fatal(err)
	}
	folders, documents, err := fetchDocInfos(list.DocIds, options)
	if err != nil {
		log.Fatal(err)
	}
	if err = store.SaveUpstreamFolders(folders); err != nil {
		log.Fatal(err)
	}
	if err = store.SaveUpstreamDocuments(documents); err != nil {
		log.Fatal(err)
	}
}

func fetchDocInfo(docId dp.Id, options *config.CliOptions) (folderList []store.Folder, document *store.Document, err error) {
	folderInfoC, folderErrC := dp.FolderInfoFuture(docId, options)
	downloadC, downloadErrC := dp.DownloadFuture(docId, options)
	err = <-downloadErrC
	if err != nil {
		return nil, nil, err
	}
	err = <-folderErrC
	if err != nil {
		return nil, nil, err
	}
	folderInfo := <-folderInfoC
	download := <-downloadC
	folderList = folderInfoToFolderList(folderInfo)
	document = downloadResultToDocument(docId, download, folderList)
	return folderList, document, nil
}

func fetchDocInfos(docIds []dp.Id, options *config.CliOptions) (folderResult []store.Folder, documentResult []store.Document, err error) {
	n := len(docIds)
	folderLists := make(chan []store.Folder, n)
	documents := make(chan *store.Document, n)
	errs := make(chan error, n*2)
	for _, docId := range docIds {
		go func(docId dp.Id) {
			folderList, document, err := fetchDocInfo(docId, options)
			errs <- err
			folderLists <- folderList
			documents <- document
		}(docId)
	}
	uniqueFolderList := newUniqueFolderList((n + 1) * 2)
	documentResult = make([]store.Document, n)
	for i := 0; i < n; i++ {
		err := <-errs
		if err != nil {
			return nil, nil, err
		}
		folderList := <-folderLists
		document := <-documents
		uniqueFolderList.add(folderList)
		documentResult[i] = *document
	}
	return uniqueFolderList.folders, documentResult, nil
}
