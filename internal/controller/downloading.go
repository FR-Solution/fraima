package controller

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path"

	"github.com/fraima/fraimactl/internal/config"
)

type downloadItem struct {
	Name       string `json:"name"`
	Src        string `json:"src"`
	HostPath   string `json:"hostpath"`
	Permission int    `json:"permission"`
}

func downloading(d config.Instruction) error {
	downloadList, err := getDownloadList(d.Spec)
	if err != nil {
		return fmt.Errorf("get download list: %w", err)
	}
	for _, item := range downloadList {
		err := download(item)
		if err != nil {
			return err
		}
	}
	return nil
}

func download(i downloadItem) error {
	download := path.Base(i.Src)

	err := wget(i.Src, path.Join(i.HostPath, download))
	if err != nil {
		return err
	}

	err = os.Chmod(i.HostPath, fs.FileMode(i.Permission))
	if err != nil {
		return err
	}

	err = os.Chown(i.HostPath, os.Getuid(), os.Getgid())
	return err
}

func wget(url, filepath string) error {
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Run()
	return cmd.Err
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
