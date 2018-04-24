package dropbox

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"log"
)

func DownloadTest(options *config.CliOptions) {
	list, err := DocList(options)
	if err != nil {
		log.Fatal(err)
	}
	docId := list.DocIds[0]
	Download(docId, options)
}

func Test(options *config.CliOptions) {
	request := &Request{
		Url:     "/users/get_current_account",
		Options: options,
	}
	result, err := request.EvalString()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(result)
}

func FolderTest(options *config.CliOptions) {
	list, err := DocList(options)
	if err != nil {
		log.Fatal(err)
	}
	n := len(list.DocIds)
	results := make(chan *FolderResult, n)
	errors := make(chan error, n)
	for _, docId := range list.DocIds {
		go func(docId Id) {
			folders, err := FolderInfo(docId, options)
			errors <- err
			results <- folders
		}(docId)
	}
	for i := 0; i < n; i++ {
		err := <-errors
		if err != nil {
			log.Fatal(err)
		}
		folders := <-results
		spew.Dump(folders)
	}
}
