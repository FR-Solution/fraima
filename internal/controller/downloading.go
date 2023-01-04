package controller

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/fraima/fraimactl/internal/config"
)

type downloadItem struct {
	Name       string `json:"name"`
	Src        string `json:"src"`
	HostPath   string `json:"hostpath"`
	Permission int    `json:"permission"`
	Unzip      bool   `json:"unzip"`
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
	downloadFile := path.Base(i.Src)
	hostPath := path.Join(i.HostPath, downloadFile)

	downloadPath := hostPath
	if i.Unzip {
		downloadPath = path.Join(i.HostPath, downloadFile)
	}

	err := wget(i.Src, downloadPath)
	if err != nil {
		return err
	}

	if i.Unzip {
		err = unzip(downloadPath, hostPath)
		if err != nil {
			return err
		}
	}

	err = os.Chmod(hostPath, fs.FileMode(i.Permission))
	if err != nil {
		return err
	}

	err = os.Chown(hostPath, os.Getuid(), os.Getgid())
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

func unzip(src, dst string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(dst, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil

}
