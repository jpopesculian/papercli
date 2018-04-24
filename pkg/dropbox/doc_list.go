package dropbox

import (
	"github.com/jpopesculian/papercli/pkg/config"
)

type DocListResult struct {
	DocIds  []Id   `json:"doc_ids"`
	Cursor  Cursor `json:"cursor"`
	HasMore bool   `json:"has_more"`
}

func DocList(options *config.CliOptions) (result *DocListResult, err error) {
	request := &Request{
		Url:     "/paper/docs/list",
		Options: options,
	}
	result = &DocListResult{}
	err = request.EvalStruct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
