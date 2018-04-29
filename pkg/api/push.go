package api

import (
	"github.com/jpopesculian/papercli/pkg/config"
	"log"
	"os"
	"path/filepath"
)

func documentPaths(options *config.CliOptions, fn func(string) error) error {
	root, err := options.RootDir()
	if err != nil {
		return err
	}
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		return fn(path)
	})
}

func Push(options *config.CliOptions) {
	err := documentPaths(options, func(path string) error {
		log.Println(path)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
