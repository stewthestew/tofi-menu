package backend

import (
	"os"
	"path/filepath"
	"strings"
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
	PATHS := strings.Split(PATH, ":")

	for _, v := range PATHS {
		if v == "" {
			continue
		}

		execs, err := walkDir(v)
		if err != nil {
			continue
		}

		for _, exe := range execs {
			pathArray = append(pathArray, exe)
		}
	}
	return pathArray, nil
}
