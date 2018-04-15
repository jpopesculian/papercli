package dropbox

import (
	"bytes"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/jpopesculian/papercli/pkg/config"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Cursor struct {
	Value      string    `json:"value"`
	Expiration time.Time `json:"expiration"`
}

type ListResult struct {
	DocIds  []string `json:"doc_ids"`
	Cursor  Cursor   `json:"cursor"`
	HasMore bool     `json:"has_more"`
}

type FolderRequest struct {
	DocId string `json:"doc_id"`
}

type Folder struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type FolderResult struct {
	FolderSharingPolicyType interface{} `json:"folder_sharing_policy_type"`
	Folders                 []Folder    `json:"folders"`
}

func buildReq(url string, params interface{}, options *config.CliOptions) (result *http.Request, err error) {
	var req *http.Request
	method := "POST"
	reqUrl := "https://api.dropboxapi.com/2" + url
	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, reqUrl, bytes.NewBuffer(data))
	} else {
		req, err = http.NewRequest(method, reqUrl, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+*options.AccessKey)
	if params != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return req, nil
}

func request(url string, params interface{}, options *config.CliOptions) (result string, err error) {
	req, err := buildReq(url, params, options)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(resData), nil
}

func reqStruct(url string, params interface{}, object interface{}, options *config.CliOptions) (err error) {
	result, err := request(url, params, options)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), object)
}

func Test(options *config.CliOptions) {
	result, err := request("/users/get_current_account", nil, options)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(result)
}

func List(options *config.CliOptions) {
	var list ListResult
	err := reqStruct("/paper/docs/list", nil, &list, options)
	if err != nil {
		log.Fatal(err)
	}
	spew.Dump(list)
}

func FolderInfo(options *config.CliOptions) {
	var list ListResult
	var folders FolderResult
	err := reqStruct("/paper/docs/list", nil, &list, options)
	if err != nil {
		log.Fatal(err)
	}
	params := FolderRequest{
		DocId: list.DocIds[0],
	}
	for i := 1; i < 100; i++ {
		err = reqStruct("/paper/docs/get_folder_info", params, &folders, options)
		if err != nil {
			log.Fatal(err)
		}
		spew.Dump(folders)
		params = FolderRequest{
			DocId: list.DocIds[i],
		}
	}
}
