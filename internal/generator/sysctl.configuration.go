package generator

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/fraima/fraimactl/internal/config"
	"github.com/fraima/fraimactl/internal/utils"
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
func createSysctlConfiguration(i config.Instruction) error {
	data, err := createSysctlConfigurationData(i.Spec.Configuration)
	if err != nil {
		return err
	}

	if err = utils.CreateFile(sysctlFilePath, data, sysctlFilePERM, "root:root"); err != nil {
		return err
	}

	return nil
}

func createSysctlConfigurationData(cfg *config.Config) ([]byte, error) {
	eargs, err := getMap(cfg.ExtraArgs)
	if err != nil {
		return nil, err
	}

	sysctlServiceBuffer := new(bytes.Buffer)
	err = sysctlConfigurationTemplate.Execute(sysctlServiceBuffer, eargs)
	return sysctlServiceBuffer.Bytes(), err
}
