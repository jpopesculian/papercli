package dropbox

import (
	"github.com/jpopesculian/papercli/pkg/config"
)

type FolderRequest struct {
	DocId Id `json:"doc_id"`
}

type Folder struct {
	Id   Id     `json:"id"`
	Name string `json:"name"`
}

type FolderResult struct {
	FolderSharingPolicyType interface{} `json:"folder_sharing_policy_type"`
	Folders                 []Folder    `json:"folders"`
}

func FolderInfo(docId Id, options *config.CliOptions) (result *FolderResult, err error) {
	params := &FolderRequest{
		DocId: docId,
	}
	result = &FolderResult{}
	request := &Request{
		Url:     "/paper/docs/get_folder_info",
		Params:  params,
		Options: options,
	}
	err = request.EvalStruct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func FolderInfoFuture(docId Id, options *config.CliOptions) (results chan *FolderResult, errors chan error) {
	results = make(chan *FolderResult, 1)
	errors = make(chan error, 1)
	go func() {
		result, err := FolderInfo(docId, options)
		errors <- err
		results <- result
	}()
	return results, errors
}

func BatchFolderInfo(docIds []Id, options *config.CliOptions) (folderInfos []*FolderResult, err error) {
	n := len(docIds)
	folderInfos = make([]*FolderResult, n)
	results := make(chan *FolderResult, n)
	errors := make(chan error, n)
	for _, docId := range docIds {
		go func(docId Id) {
			folders, err := FolderInfo(docId, options)
			errors <- err
			results <- folders
		}(docId)
	}
	for i := 0; i < n; i++ {
		err := <-errors
		if err != nil {
			return nil, err
		}
		folderInfos[i] = <-results
	}
	return folderInfos, nil
}
