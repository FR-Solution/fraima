package controller

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/fraima/fraimactl/internal/config"
)

var (
	//go:embed template/sysctl.d.network.conf.tmpl
	sysctlNetworkConfigurationTemplateStr string
	sysctlNetworkConfigurationTemplate    = template.Must(template.New("sysctl.d-network.conf").Parse(sysctlNetworkConfigurationTemplateStr))
)

const (
	sysctlsFilePath = "/etc/sysctl.d/99-network.conf"
	sysctlsFilePERM = 0644
)

// createSysctlNetworkConfiguration create Sysctls.service file.
func createSysctlNetworkConfiguration(cfg config.Instruction) error {
	data, err := createSysctlsServiceData(cfg)
	if err != nil {
		return err
	}

	return createFile(sysctlsFilePath, data, sysctlsFilePERM)
}

func createSysctlsServiceData(cfg config.Instruction) ([]byte, error) {
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

	sysctlsServiceBuffer := new(bytes.Buffer)
	err := sysctlNetworkConfigurationTemplate.Execute(sysctlsServiceBuffer, extraArgs)
	return sysctlsServiceBuffer.Bytes(), err
}
