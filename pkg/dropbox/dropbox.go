package dropbox

import (
	"github.com/jpopesculian/papercli/pkg/config"
	"io/ioutil"
	"log"
	"net/http"
)

func Test(options *config.CliOptions) {
	req, err := http.NewRequest(
		"POST",
		"https://api.dropboxapi.com/2/users/get_current_account",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+*options.AccessKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	resString := string(resData)

	log.Printf(resString)
}
