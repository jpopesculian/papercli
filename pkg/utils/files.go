package utils

import (
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
