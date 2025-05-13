package backend

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
)

func walkDir(PATH string) ([]string, error) {
	paths := []string{}
	err := filepath.WalkDir(PATH, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		// Check if it's a regular file and executable (any exec bit set)
		if info.Mode().IsRegular() && info.Mode().Perm()&0111 != 0 {
			paths = append(paths, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return paths, nil
}

func List(PATH string) ([]string, error) {
	var pathArray []string
	PATHS := strings.SplitSeq(PATH, ":")

	PATHS(func(v string) bool {
		execs, err := walkDir(v)

		if err != nil {
			log.Error(err)
		}

		for _, exe := range execs {
			pathArray = append(pathArray, exe)
		}

		return true
	})

	return pathArray, nil
}
