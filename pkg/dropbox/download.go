package dropbox

import (
	"github.com/jpopesculian/papercli/pkg/config"
)

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
		Url:             "/paper/docs/download",
		Params:          params,
		Options:         options,
		ParamsInHeader:  true,
		ResultsInHeader: true,
	}
	err = request.EvalStruct(result)
	if err != nil {
		return nil, err
	}
	result.Content, err = request.EvalFile()
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DownloadFuture(docId Id, options *config.CliOptions) (results chan *DownloadResult, errors chan error) {
	results = make(chan *DownloadResult, 1)
	errors = make(chan error, 1)
	go func() {
		result, err := Download(docId, options)
		errors <- err
		results <- result
	}()
	return results, errors
}

func BatchDownload(docIds []Id, options *config.CliOptions) (downloads []*DownloadResult, err error) {
	n := len(docIds)
	downloads = make([]*DownloadResult, n)
	results := make(chan *DownloadResult, n)
	errors := make(chan error, n)
	for _, docId := range docIds {
		go func(docId Id) {
			download, err := Download(docId, options)
			errors <- err
			results <- download
		}(docId)
	}
	for i := 0; i < n; i++ {
		err := <-errors
		if err != nil {
			return nil, err
		}
		downloads[i] = <-results
	}
	return downloads, nil
}
