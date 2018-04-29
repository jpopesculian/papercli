package utils

import (
	"io/ioutil"
	"os"
	"strings"
)

func IsDir(dir string) (bool, error) {
	info, err := os.Stat(dir)
	passed := err == nil && info.IsDir()
	return passed, err
}

func UpDirectory(dir string) string {
	paths := SplitPath(dir)
	return strings.Join(paths[:len(paths)-1], string(os.PathSeparator))
}

func SplitPath(dir string) []string {
	return strings.Split(dir, string(os.PathSeparator))
}

func ReadFileAsync(path string) (chan []byte, chan error) {
	result := make(chan []byte, 1)
	errs := make(chan error, 1)
	go func() {
		content, err := ioutil.ReadFile(path)
		errs <- err
		result <- content
	}()
	return result, errs
}
