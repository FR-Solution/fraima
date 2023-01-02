package downloader

import (
	"io/fs"
	"os"
	"os/exec"

	"github.com/fraima/fraimactl/internal/config"
	"go.uber.org/zap"
)

func Run(downloadList []config.Download) {
	for _, d := range downloadList {
		err := download(d)
		if err != nil {
			zap.L().Error("downloading", zap.Any("download", d), zap.Error(err))
		}
	}
}

func download(d config.Download) error {
	err := wget(d.URL, d.Filepath)
	if err != nil {
		return err
	}

	err = os.Chmod(d.Filepath, fs.FileMode(d.Permission))
	return err
}

func wget(url, filepath string) error {
	// run shell `wget URL -O filepath`
	cmd := exec.Command("wget", url, "-O", filepath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
