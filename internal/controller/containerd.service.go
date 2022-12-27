package controller

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/fraima/fraima/internal/config"
)

var (
	//go:embed template/containerd.service.tmpl
	containerdTemplateStr string
	containerdTemplate    = template.Must(template.New("containerd-service").Parse(containerdTemplateStr))
)

const (
	containerdServiceFilePath = "/etc/systemd/system/containerd.service"
	containerdServiceFilePERM = 0644
)

// createContainerdService create containerd.service file.
func createContainerdService(cfg config.File) error {
	data, err := createContainerdServiceData(cfg)
	if err != nil {
		return err
	}

	return createFile(containerdServiceFilePath, data, containerdServiceFilePERM)
}

func createContainerdServiceData(cfg config.File) ([]byte, error) {
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

	containerdServiceBuffer := new(bytes.Buffer)
	err := containerdTemplate.Execute(containerdServiceBuffer, extraArgs)
	return containerdServiceBuffer.Bytes(), err
}
