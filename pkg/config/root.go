package config

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

const CONFIG_DIR_NAME = ".paperconfig"

func CreateConfigDir(options *CliOptions) {
	root := *options.Dir
	dir := CONFIG_DIR_NAME
	if len(root) > 0 {
		dir = configDirFromRoot(root)
	}
	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func findConfigDir() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	for len(dir) > 0 {
		configDir := configDirFromRoot(dir)
		if passed, _ := isDir(configDir); passed {
			return configDir
		}
		dir = upDirectory(dir)
	}
	return ""
}

func findRootDir() string {
	configDir := findConfigDir()
	if len(configDir) > 0 {
		return upDirectory(configDir)
	}
	return ""
}

func configDirFromRoot(dir string) string {
	return filepath.Join(dir, CONFIG_DIR_NAME)
}

func isDir(dir string) (bool, error) {
	info, err := os.Stat(dir)
	passed := err == nil && info.IsDir()
	return passed, err
}

func upDirectory(dir string) string {
	paths := splitPath(dir)
	return strings.Join(paths[:len(paths)-1], string(os.PathSeparator))
}

func splitPath(dir string) []string {
	return strings.Split(dir, string(os.PathSeparator))
}
