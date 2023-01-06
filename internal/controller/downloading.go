package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/codeclysm/extract/v3"

	"github.com/fraima/fraimactl/internal/config"
)

type downloadItem struct {
	Name       string `json:"name"`
	Src        string `json:"src"`
	HostPath   string `json:"path"`
	Permission int    `json:"permission"`
	Unzip      unzip  `json:"unzip"`
}

type unzip struct {
	Status bool     `json:"status"`
	Files  []string `json:"files"`
}

var client http.Client

func downloading(d config.Instruction) error {
	downloadList, err := getDownloadList(d.Spec)
	if err != nil {
		return fmt.Errorf("get download list: %w", err)
	}
	for _, item := range downloadList {
		file, err := download(item.Src)
		if err != nil {
			return err
		}
		defer file.Close()

		var data []byte
		if item.Unzip.Status {
			err = unzipFile(item.Name, file)
			if err != nil {
				return err
			}

			for _, f := range item.Unzip.Files {
				filePath := getDownloadDir(item.Name, f)
				data, err = os.ReadFile(filePath)
				if err != nil {
					return err
				}

				err = createFile(path.Join(item.HostPath, path.Base(f)), data, item.Permission)
				if err != nil {
					return err
				}
			}
		} else {
			data, err = ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			err = createFile(path.Join(item.HostPath, item.Name), data, item.Permission)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getDownloadList(spec any) ([]downloadItem, error) {
	specItems, ok := spec.([]any)
	if !ok {
		return nil, fmt.Errorf("downloading spec must be array")
	}
	downloadList := make([]downloadItem, 0, len(specItems))
	for _, item := range specItems {
		di, err := getDownloadItem(item)
		if err != nil {
			return nil, err
		}
		downloadList = append(downloadList, di)
	}
	return downloadList, nil
}

func getDownloadItem(i any) (downloadItem, error) {
	var item downloadItem
	itemMap, err := getMap(i)
	if err != nil {
		return item, err
	}

	jsonData, err := json.Marshal(itemMap)
	if err != nil {
		return item, err
	}

	err = json.Unmarshal(jsonData, &item)
	return item, err
}

func download(src string) (io.ReadCloser, error) {
	resp, err := client.Get(src)
	if err != nil {
		return nil, err
	}
	return resp.Body, err
}

func unzipFile(component string, file io.Reader) error {
	err := extract.Archive(context.Background(), file, getDownloadDir(component, ""), nil)
	if err != nil {
		return err
	}
	return nil
}

func getDownloadDir(component, filePath string) string {
	return path.Join(os.TempDir(), "fraima", component, filePath)
}
