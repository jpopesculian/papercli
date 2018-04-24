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
