package files

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const ROOT_FILE_NAME = ".paperroot"

func CreateRootFile() {
	os.OpenFile(ROOT_FILE_NAME, os.O_RDONLY|os.O_CREATE, 0666)
}

func FindRootFile() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for len(dir) > 0 {
		rootFile := filepath.Join(dir, ROOT_FILE_NAME)
		dir = upDirectory(dir)
		if _, err := os.Stat(rootFile); err == nil {
			return rootFile
		}
	}
	return ""
}

func RootDir() string {
	return filepath.Dir(FindRootFile())
}

func upDirectory(dir string) string {
	separator := string(os.PathSeparator)
	paths := strings.Split(dir, separator)
	return strings.Join(paths[:len(paths)-1], separator)
}
