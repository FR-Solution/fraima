package generator

import (
	"bytes"
	_ "embed"
	"os/exec"
	"text/template"

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
func createSysctlConfiguration(extraArgs any) error {
	data, err := createSysctlServiceData(extraArgs)
	if err != nil {
		return err
	}

	if err = utils.CreateFile(sysctlFilePath, data, sysctlFilePERM, "root:root"); err != nil {
		return err
	}

	err = exec.Command("sysctl", "--system").Run()
	return err
}

func createSysctlServiceData(extraArgs any) ([]byte, error) {
	eargs, err := getMap(extraArgs)
	if err != nil {
		return nil, err
	}

	sysctlServiceBuffer := new(bytes.Buffer)
	err = sysctlConfigurationTemplate.Execute(sysctlServiceBuffer, eargs)
	return sysctlServiceBuffer.Bytes(), err
}
