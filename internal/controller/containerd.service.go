package controller

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/fraima/fraimactl/internal/config"
)

var (
	//go:embed template/containerd.service.tmpl
	containerdServiceTemplateStr string
	containerdServiceTemplate    = template.Must(template.New("containerd-service").Parse(containerdServiceTemplateStr))
)

const (
	containerdServiceFilePath = "/etc/systemd/system/containerd.service"
	containerdServiceFilePERM = 0644
)

// createContainerdService create containerd.service file.
func createContainerdService(cfg config.Instruction) error {
	data, err := createContainerdServiceData(cfg)
	if err != nil {
		return err
	}

	return createFile(containerdServiceFilePath, data, containerdServiceFilePERM, "root:root")
}

func createContainerdServiceData(cfg config.Instruction) ([]byte, error) {
	extraArgs := make(map[string]string)
	if cfg.Spec != nil {
		args, ok := cfg.Spec.(map[any]any)
		if !ok {
			return nil, fmt.Errorf("args converting is not available")
		}
		for k, v := range args {
			extraArgs[fmt.Sprint(k)] = fmt.Sprint(v)
		}
	}

	containerdServiceBuffer := new(bytes.Buffer)
	err := containerdServiceTemplate.Execute(containerdServiceBuffer, extraArgs)
	return containerdServiceBuffer.Bytes(), err
}
