package controller

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/fraimactl/fraimactl/internal/config"
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
func createKubletService(cfg config.File) error {
	data, err := createKubleteServiceData(cfg)
	if err != nil {
		return err
	}

	return createFile(kubeletServiceFilePath, data, kubeletServiceFilePERM)
}

func createKubleteServiceData(cfg config.File) ([]byte, error) {
	extraArgs := make(map[string]string)
	if cfg.ExtraArgs != nil {
		args, ok := cfg.ExtraArgs.(map[any]any)
		if !ok {
			return nil, fmt.Errorf("args converting is not available")
		}
		for k, v := range args {
			extraArgs[fmt.Sprint(k)] = fmt.Sprint(v)
		}
	}

	kubletServiceBuffer := new(bytes.Buffer)
	err := kubeletTemplate.Execute(kubletServiceBuffer, extraArgs)
	return kubletServiceBuffer.Bytes(), err
}
