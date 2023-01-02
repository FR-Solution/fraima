package controller

import (
	"fmt"
	"io/fs"
	"os"
	"path"
)

func createFile(filepath string, data []byte, perm int) error {
	dir := path.Dir(filepath)
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	err := os.WriteFile(filepath, data, fs.FileMode(perm))
	if err != nil {
		return err
	}

	err = os.Chown(filepath, os.Getuid(), os.Getgid())
	return err
}

func getArgsMap(args map[any]any) map[string]any {
	rArgs := make(map[string]any)
	for k, v := range args {
		if nArgs, ok := v.(map[any]any); ok {
			rArgs[fmt.Sprint(k)] = getArgsMap(nArgs)
			continue
		}
		rArgs[fmt.Sprint(k)] = v
	}
	return rArgs
}
