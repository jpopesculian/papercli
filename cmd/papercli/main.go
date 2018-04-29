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
	"fmt"
	"github.com/jpopesculian/papercli/pkg/api"
	"github.com/jpopesculian/papercli/pkg/config"
	"github.com/jpopesculian/papercli/pkg/dropbox"
	"github.com/jpopesculian/papercli/pkg/version"
	"os"
)

func main() {
	command, options := config.ParseArgs()
	switch command {
	case "test":
		dropbox.Test(options)
	case "list":
		dropbox.DocList(options)
	case "folder":
		dropbox.FolderTest(options)
	case "download":
		dropbox.DownloadTest(options)
	case "fetch":
		api.Fetch(options)
	case "update":
		api.Update(options)
	case "pull":
		api.Fetch(options)
		api.Update(options)
	case "push":
		api.Push(options)
	case "init":
		config.CreateConfigDir(options)
	case "version":
		fmt.Println(version.VERSION)
	default:
		config.PrintUsage()
		os.Exit(0)
	}
}
