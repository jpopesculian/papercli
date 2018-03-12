package config

import (
	"flag"
	"fmt"
	"os"
)

type CliOptions struct {
	AccessKey *string
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
