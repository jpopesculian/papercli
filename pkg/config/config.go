package config

import (
	"errors"
	"flag"
	"fmt"
	"github.com/jpopesculian/papercli/pkg/utils"
	"os"
)

type CliOptions struct {
	AccessKey *string
	Dir       *string
	RestArgs  []string
}

func PrintUsage() {
	fmt.Printf("Usage: %s [OPTIONS] command ...\n", os.Args[0])
	flag.PrintDefaults()
}

func ParseArgs() (string, *CliOptions) {
	options := CliOptions{
		AccessKey: flag.String(
			"accessKey",
			os.Getenv("PAPER_ACCESS_KEY"),
			"Access Key for self authorized testing",
		),
		Dir: flag.String(
			"dir",
			findRootDir(),
			"Root directory for PaperCLI",
		),
	}
	flag.Parse()
	if len(flag.Args()) == 0 {
		PrintUsage()
		os.Exit(1)
	}
	args := flag.Args()
	command := args[0]
	options.RestArgs = args[1:]
	return command, &options
}

func (options *CliOptions) RootDir() (string, error) {
	if len(*options.Dir) < 1 {
		return "", errors.New("PaperCLI can't find root directory!")
	}
	if passed, err := utils.IsDir(*options.Dir); !passed {
		return "", err
	}
	return *options.Dir, nil
}

func (options *CliOptions) ConfigDir() (string, error) {
	rootDir, err := options.RootDir()
	if err != nil {
		return "", err
	}
	configDir := configDirFromRoot(rootDir)
	if passed, err := utils.IsDir(*options.Dir); !passed {
		return "", err
	}
	return configDir, nil
}
