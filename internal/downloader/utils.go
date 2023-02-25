package downloader

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/sha256"
	"os"
	"path"

	"github.com/codeclysm/extract/v3"
)

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
