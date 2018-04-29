package dropbox

import (
	"github.com/jpopesculian/papercli/pkg/config"
)

const UPDATE_POLICY = "overwrite_all"

type UpdateRequest struct {
	DocId    Id     `json:"doc_id"`
	Revision int    `json:"revision"`
	Format   string `json:"import_format"`
	Policy   string `json:"doc_update_policy"`
}

type UpdateResult struct {
	DocId    Id     `json:"doc_id"`
	Revision int    `json:"revision"`
	Title    string `json:"title"`
}

func Update(content []byte, id Id, revision int, options *config.CliOptions) (result *UpdateResult, err error) {
	params := &UpdateRequest{
		DocId:    id,
		Revision: revision,
		Policy:   UPDATE_POLICY,
		Format:   MARKDOWN,
	}
	result = &UpdateResult{}
	request := &Request{
		Url:            "/paper/docs/update",
		Params:         params,
		Options:        options,
		Data:           content,
		ParamsInHeader: true,
	}
	err = request.EvalStruct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
