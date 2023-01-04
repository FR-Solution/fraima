package controller

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/fraima/fraimactl/internal/config"
)

var (
	//go:embed template/k8s.conf.tmpl
	k8sConfigurationTemplateStr string
	k8sConfigurationTemplate    = template.Must(template.New("k8s.conf-service").Parse(k8sConfigurationTemplateStr))
)

const (
	k8sConfigurationServiceFilePath = "/etc/modules-load.d/k8s.conf"
	k8sConfigurationServiceFilePERM = 0644
)

// createModProbConfiguration create k8s.conf file.
func createModProbConfiguration(cfg config.Instruction) error {
	data, err := createModProbConfigurationData(cfg)
	if err != nil {
		return err
	}

	return createFile(k8sConfigurationServiceFilePath, data, k8sConfigurationServiceFilePERM)
}

func createModProbConfigurationData(cfg config.Instruction) ([]byte, error) {
	var (
		args []any
		ok   bool
	)
	if cfg.Spec != nil {
		args, ok = cfg.Spec.([]any)
		if !ok {
			return nil, fmt.Errorf("args converting is not available")
		}
	}

	k8sConfigurationServiceBuffer := new(bytes.Buffer)
	err := k8sConfigurationTemplate.Execute(k8sConfigurationServiceBuffer, args)
	return k8sConfigurationServiceBuffer.Bytes(), err
}
