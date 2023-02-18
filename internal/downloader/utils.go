package downloader

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"os"
	"path"

	"github.com/codeclysm/extract/v3"
)

// func getDownloadList(spec any) ([]downloadItem, error) {
// 	specItems, ok := spec.([]any)
// 	if !ok {
// 		return nil, fmt.Errorf("downloading spec must be array")
// 	}
// 	downloadList := make([]downloadItem, 0, len(specItems))
// 	for _, item := range specItems {
// 		di, err := getDownloadItem(item)
// 		if err != nil {
// 			return nil, err
// 		}
// 		downloadList = append(downloadList, di)
// 	}
// 	return downloadList, nil
// }

// func getDownloadItem(i any) (downloadItem, error) {
// 	var item downloadItem
// 	itemMap, err := getMap(i)
// 	if err != nil {
// 		return item, err
// 	}

// 	yamlData, err := yaml.Marshal(itemMap)
// 	if err != nil {
// 		return item, err
// 	}

// 	err = yaml.Unmarshal(yamlData, &item)
// 	return item, err
// }

func getMap(i any) (map[string]any, error) {
	rArgs := make(map[string]any)
	err := fmt.Errorf("args converting is not available")
	args, ok := i.(map[any]any)
	if !ok {
		return rArgs, err
	}
	for k, v := range args {
		key := fmt.Sprint(k)
		if nArgs, ok := v.(map[any]any); ok {
			rArgs[key], err = getMap(nArgs)
			if err != nil {
				return rArgs, err
			}
			continue
		}

		rArgs[key] = v
	}
	return rArgs, nil
}

func check(file, fileCheckSum []byte, checkSumType string) bool {
	switch checkSumType {
	case "sha256":
		sum := sha256.Sum256(file)
		bytes.Equal(sum[:], fileCheckSum)
	case "md5":
		sum := md5.Sum(file)
		bytes.Equal(sum[:], fileCheckSum)
	}

	return false
}

func unzipFile(component string, data []byte) error {
	reader := bytes.NewReader(data)
	err := extract.Archive(context.Background(), reader, getDownloadDir(component, ""), nil)
	if err != nil {
		return err
	}
	return nil
}

func getDownloadDir(component, filePath string) string {
	return path.Join(os.TempDir(), "fraima", component, filePath)
}
