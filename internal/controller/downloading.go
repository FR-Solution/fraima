package controller

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/codeclysm/extract/v3"
	"gopkg.in/yaml.v2"

	"github.com/fraima/fraimactl/internal/config"
)

type downloadItem struct {
	Name       string    `yaml:"name"`
	Src        string    `yaml:"src"`
	CheckSum   *checkSum `yaml:"checkSum"`
	HostPath   string    `yaml:"path"`
	Owner      string    `yaml:"owner"`
	Permission int       `yaml:"permission"`
	Unzip      unzip     `yaml:"unzip"`
}

type unzip struct {
	Status bool     `yaml:"status"`
	Files  []string `yaml:"files"`
}

type checkSum struct {
	Src  string `yaml:"src"`
	Type string `yaml:"type"`
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

		if item.CheckSum != nil {
			fileCheckSum, err := download(item.CheckSum.Src)
			if err != nil {
				return err
			}

			if check(file, fileCheckSum, item.CheckSum.Type) {
				return fmt.Errorf("the file was downloaded incorrectly")
			}
		}

		// var data []byte
		if item.Unzip.Status {
			err = unzipFile(item.Name, file)
			if err != nil {
				return err
			}

			// for _, f := range item.Unzip.Files {
			// 	filePath := getDownloadDir(item.Name, f)
			// 	data, err = os.ReadFile(filePath)
			// 	if err != nil {
			// 		return err
			// 	}

			// 	err = createFile(path.Join(item.HostPath, path.Base(f)), data, item.Permission, item.Owner)
			// 	if err != nil {
			// 		return err
			// 	}
			// }
			continue
		}

		// err = createFile(path.Join(item.HostPath, item.Name), file, item.Permission, item.Owner)
		// if err != nil {
		// 	return err
		// }
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

	yamlData, err := yaml.Marshal(itemMap)
	if err != nil {
		return item, err
	}

	err = yaml.Unmarshal(yamlData, &item)
	return item, err
}

func download(src string) ([]byte, error) {
	resp, err := client.Get(src)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
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
