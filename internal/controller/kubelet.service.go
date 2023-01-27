package controller

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/fraima/fraimactl/internal/config"
)

var (
	//go:embed template/kubelet.service.tmpl
	kubeletTemplateStr string
	kubeletTemplate    = template.Must(template.New("kubelet-service").Parse(kubeletTemplateStr))
)

const (
	kubeletServiceName     = "kubelet"
	kubeletServiceFilePath = "/etc/systemd/system/kubelet.service"
	kubeletServiceFilePERM = 0644
)

// createKubletService create kubelet.service file.
func createKubletService(cfg config.Instruction) error {
	data, err := createKubleteServiceData(cfg)
	if err != nil {
		return err
	}

	err = createFile(kubeletServiceFilePath, data, kubeletServiceFilePERM, "root:root")
	if err != nil {
		return err
	}

	err = startService(kubeletServiceName)
	return err
}

func createKubleteServiceData(cfg config.Instruction) ([]byte, error) {
	eargs, err := getMap(cfg.Spec)
	if err != nil {
		return nil, err
	}

	kubletServiceBuffer := new(bytes.Buffer)
	err = kubeletTemplate.Execute(kubletServiceBuffer, eargs)
	return kubletServiceBuffer.Bytes(), err
}
