package controller

import (
	"bytes"
	_ "embed"
	"os"
	"text/template"
)

var (
	//go:embed template/kubelet.service.tmpl
	kubeletTemplateStr string
	kubeletTemplate    = template.Must(template.New("kubelet-service").Parse(kubeletTemplateStr))
)

const (
	kubeletServiceFilePath = "/etc/systemd/system/kubelet.service"
	kubeletServiceFilePERM = 0644
)

// createKubletService create kubelet.service file.
func createKubletService(cfg configuration) error {
	kubletServiceBuffer := new(bytes.Buffer)
	err := kubeletTemplate.Execute(kubletServiceBuffer, cfg.extraArgs)
	if err != nil {
		return err
	}

	err = os.WriteFile(kubeletServiceFilePath, kubletServiceBuffer.Bytes(), kubeletServiceFilePERM)
	if err != nil {
		return err
	}

	err = os.Chown(kubeletServiceFilePath, os.Getuid(), os.Getgid())
	return err
}
