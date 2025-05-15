package backend

import (
	"os"
	"path/filepath"
	"strings"
	"tofi/internal/utils"
	"github.com/charmbracelet/log"
)


func Log(level log.Level, silent bool, msg any, keyvals ...any) {
	if !silent {
		log.Log(level, msg, keyvals)
	}
}

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

func List(PATH string, silent bool) ([]string, error) {
	var pathArray []string
	PATHS := strings.SplitSeq(PATH, ":")

	PATHS(func(v string) bool {
		execs, err := walkDir(v)

		if err != nil {
			utils.Log(log.ErrorLevel, silent, err)
		}

		for _, exe := range execs {
			pathArray = append(pathArray, exe)
		}

		return true
	})

	return pathArray, nil
}
