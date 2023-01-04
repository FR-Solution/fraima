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

func getMap(i any) (map[string]any, error) {
	rArgs := make(map[string]any)
	err := fmt.Errorf("args converting is not available")
	args, ok := i.(map[any]any)
	if !ok {
		return rArgs, err
	}
	for k, v := range args {
		if nArgs, ok := v.(map[any]any); ok {
			rArgs[fmt.Sprint(k)], err = getMap(nArgs)
			if err != nil {
				return rArgs, err
			}
			continue
		}
		rArgs[fmt.Sprint(k)] = v
	}
	return rArgs, nil
}
