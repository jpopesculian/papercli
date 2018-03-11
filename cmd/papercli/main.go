/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	accessKey string
}

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
}

func pull(config apiConfig) {
	req, err := http.NewRequest(
		"POST",
		"https://api.dropboxapi.com/2/users/get_current_account",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+config.accessKey)
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

func main() {
	var (
		accessKey = flag.String(
			"accessKey",
			os.Getenv("PAPER_ACCESS_KEY"),
			"Access Key for self authorized testing",
		)
	)
	flag.Usage = myUsage
	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	args := flag.Args()
	command := args[0]
	if command == "help" {
		flag.Usage()
		os.Exit(0)
	}
	config := apiConfig{
		accessKey: *accessKey,
	}
	switch command := args[0]; command {
	case "pull":
		pull(config)
	default:
		flag.Usage()
		os.Exit(0)
	}
}
