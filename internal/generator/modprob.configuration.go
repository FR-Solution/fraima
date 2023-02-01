package generator

import (
	"bytes"
	_ "embed"
	"fmt"
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
	k8sConfigurationServiceFilePath = "/etc/modules-load.d/k8s.conf"
	k8sConfigurationServiceFilePERM = 0644
)

// createModProbeConfiguration create k8s.conf file.
func createModProbeConfiguration(cfg config.Instruction) error {
	var (
		eargs []string
		ok    bool
	)
	if cfg.Spec != nil {
		eargs, ok = cfg.Spec.([]string)
		if !ok {
			return fmt.Errorf("args converting is not available")
		}
	}

	data, err := createModProbeConfigurationData(eargs)
	if err != nil {
		return err
	}

	err = utils.CreateFile(k8sConfigurationServiceFilePath, data, k8sConfigurationServiceFilePERM, "root:root")
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
