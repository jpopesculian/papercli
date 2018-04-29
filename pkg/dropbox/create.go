package dropbox

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
)

type CreateRequest struct {
	Format   string `json:"import_format"`
	FolderId Id     `json:"parent_folder_id"`
}

func Create(content []byte, folderId Id, options *config.CliOptions) (result *UpdateResult, err error) {
	params := &CreateRequest{
		FolderId: folderId,
		Format:   MARKDOWN,
	}
	result = &UpdateResult{}
	request := &Request{
		Url:            "/paper/docs/create",
		Params:         params,
		Options:        options,
		Data:           content,
		ParamsInHeader: true,
	}
	res, err := request.doHttpReq()
	spew.Dump(res.Body)
	spew.Dump(res.Header)
	spew.Dump(err.Error())
	err = request.EvalStruct(result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
