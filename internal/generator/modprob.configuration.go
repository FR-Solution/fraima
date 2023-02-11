package generator

import (
	"bytes"
	_ "embed"
	"os/exec"
	"text/template"

	"github.com/fraima/fraimactl/internal/config"
	"github.com/fraima/fraimactl/internal/utils"
)

var (
	//go:embed template/modprobe.conf.tmpl
	k8sConfigurationTemplateStr string
	k8sConfigurationTemplate    = template.Must(template.New("k8s.conf-service").Parse(k8sConfigurationTemplateStr))
)

const (
	k8sConfigurationServiceFilePath  = "/etc/modules-load.d/k8s.conf"
	k8sConfigurationServiceFilePERM  = 0644
	k8sConfigurationServiceFileOwner = "root:root"
)

// createModProbeConfiguration create k8s.conf file.
func createModProbeConfiguration(i config.Instruction) error {
	eargs, ok := i.Spec.Configuration.ExtraArgs.([]string)
	if !ok {
		return errArgsUnavailable
	}

	data, err := createModProbeConfigurationData(eargs)
	if err != nil {
		return err
	}

	err = utils.CreateFile(k8sConfigurationServiceFilePath, data, k8sConfigurationServiceFilePERM, k8sConfigurationServiceFileOwner)
	if err != nil {
		return err
	}

	for _, a := range eargs {
		err = exec.Command("modprobe", a).Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func createModProbeConfigurationData(eargs []string) ([]byte, error) {
	k8sConfigurationServiceBuffer := new(bytes.Buffer)
	err := k8sConfigurationTemplate.Execute(k8sConfigurationServiceBuffer, eargs)
	return k8sConfigurationServiceBuffer.Bytes(), err
}
