package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/fraima/fraimactl/internal/config"
	"github.com/fraima/fraimactl/internal/utils"
)

type downloader struct {
	client *http.Client
}

func New() *downloader {
	return &downloader{
		client: &http.Client{},
	}
}

func (s *downloader) Run(instruction config.DownloadInstruction) error {
	file, err := s.download(instruction.Src)
	if err != nil {
		return err
	}

	if instruction.CheckSum != nil {
		fileCheckSum, err := s.download(instruction.CheckSum.Src)
		if err != nil {
			return err
		}

		if check(file, fileCheckSum, instruction.CheckSum.Type) {
			return fmt.Errorf("the file was downloaded incorrectly")
		}
	}

	if instruction.Unzip.Status {
		err = unzipFile(instruction.Name, file)
		if err != nil {
			return err
		}

		for _, f := range instruction.Unzip.Files {
			filePath := getDownloadDir(instruction.Name, f)
			data, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			err = utils.CreateFile(path.Join(instruction.HostPath, path.Base(f)), data, instruction.Permission, instruction.Owner)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err = utils.CreateFile(path.Join(instruction.HostPath, instruction.Name), file, instruction.Permission, instruction.Owner)
	if err != nil {
		return err
	}
	return nil
}

func (s *downloader) download(src string) ([]byte, error) {
	resp, err := s.client.Get(src)
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
