package controller

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/fraima/fraimactl/internal/config"
)

var (
	//go:embed template/sysctl.conf.tmpl
	sysctlConfigurationTemplateStr string
	sysctlConfigurationTemplate    = template.Must(template.New("sysctl.conf").Parse(sysctlConfigurationTemplateStr))
)

const (
	sysctlFilePath = "/etc/sysctl.d/99-fraima.conf"
	sysctlFilePERM = 0644
)

// createSysctlNetworkConfiguration create Sysctl.service file.
func createSysctlConfiguration(cfg config.Instruction) error {
	data, err := createSysctlServiceData(cfg)
	if err != nil {
		return err
	}

	return createFile(sysctlFilePath, data, sysctlFilePERM, "root:root")
}

func createSysctlServiceData(cfg config.Instruction) ([]byte, error) {
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

	sysctlServiceBuffer := new(bytes.Buffer)
	err := sysctlConfigurationTemplate.Execute(sysctlServiceBuffer, extraArgs)
	return sysctlServiceBuffer.Bytes(), err
}
