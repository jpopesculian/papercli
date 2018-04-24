package dropbox

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jpopesculian/papercli/pkg/config"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type Id string

type Cursor struct {
	Value      string    `json:"value"`
	Expiration time.Time `json:"expiration"`
}

type Request struct {
	Url            string
	Params         interface{}
	Options        *config.CliOptions
	ParamsInHeader bool
	httpResponse   *http.Response
}

func (request *Request) newHttpReq(body io.Reader) (req *http.Request, err error) {
	method := "POST"
	reqUrl := "https://api.dropboxapi.com/2" + request.Url
	return http.NewRequest(method, reqUrl, body)
}

func (request *Request) newJsonHttpReq(data []byte) (req *http.Request, err error) {
	req, err = request.newHttpReq(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (request *Request) newHeaderApiHttpReq(data []byte) (req *http.Request, err error) {
	req, err = request.newHttpReq(nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Dropbox-API-Arg", string(data))
	return req, nil
}

func (request *Request) buildHttpReq() (req *http.Request, err error) {
	if request.Params != nil {
		data, err := json.Marshal(request.Params)
		if err != nil {
			return nil, err
		}
		if request.ParamsInHeader {
			req, err = request.newHeaderApiHttpReq(data)
		} else {
			req, err = request.newJsonHttpReq(data)
		}
	} else {
		req, err = request.newHttpReq(nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+*request.Options.AccessKey)
	return req, nil
}

func (request *Request) doHttpReq() (res *http.Response, err error) {
	if request.httpResponse != nil {
		return request.httpResponse, nil
	}
	req, err := request.buildHttpReq()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	request.httpResponse = res
	if res.StatusCode != 200 {
		message, err := readHttpResBody(res)
		if err != nil {
			return res, err
		} else {
			return res, errors.New(message)
		}
	}
	return res, nil
}

func readHttpResBody(res *http.Response) (result string, err error) {
	defer res.Body.Close()
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(resData), nil
}

func (request *Request) EvalString() (result string, err error) {
	res, err := request.doHttpReq()
	if err != nil {
		return "", err
	}
	if request.ParamsInHeader {
		return res.Header.Get("Dropbox-Api-Result"), nil
	} else {
		return readHttpResBody(res)
	}
}

func (request *Request) EvalStruct(object interface{}) (err error) {
	result, err := request.EvalString()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(result), object)
}

func (request *Request) EvalFile(path string) (err error) {
	res, err := request.doHttpReq()
	if err != nil {
		return err
	}
	defer res.Body.Close()
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, res.Body)
	return nil
}
