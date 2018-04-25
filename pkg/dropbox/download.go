package dropbox

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
)

const MARKDOWN = "markdown"

type DownloadRequest struct {
	DocId  Id     `json:"doc_id"`
	Format string `json:"export_format"`
}

type DownloadResult struct {
	Owner    string `json:"owner"`
	Title    string `json:"title"`
	Revision int    `json:"revision"`
	MimeType string `json:"mime_type"`
	Content  []byte `json:"content"`
}

func Download(docId Id, options *config.CliOptions) (result *DownloadResult, err error) {
	params := &DownloadRequest{
		DocId:  docId,
		Format: MARKDOWN,
	}
	result = &DownloadResult{}
	request := &Request{
		Url:            "/paper/docs/download",
		Params:         params,
		Options:        options,
		ParamsInHeader: true,
	}
	err = request.EvalStruct(result)
	if err != nil {
		return nil, err
	}
	result.Content, err = request.EvalFile()
	if err != nil {
		return nil, err
	}
	spew.Dump(err)
	return result, nil
}
