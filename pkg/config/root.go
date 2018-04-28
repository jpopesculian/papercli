package config

import (
	"github.com/jpopesculian/papercli/pkg/utils"
	"log"
	"os"
	"path/filepath"
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
		if passed, _ := utils.IsDir(configDir); passed {
			return configDir
		}
		dir = utils.UpDirectory(dir)
	}
	return ""
}

func findRootDir() string {
	configDir := findConfigDir()
	if len(configDir) > 0 {
		return utils.UpDirectory(configDir)
	}
	return ""
}

func configDirFromRoot(dir string) string {
	return filepath.Join(dir, CONFIG_DIR_NAME)
}
