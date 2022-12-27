package controller

import (
	"bytes"
	_ "embed"
	"fmt"
	"text/template"

	"github.com/fraima/fraima/internal/config"
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

// createK8sConfiguration create k8s.conf file.
func createK8sConfiguration(cfg config.File) error {
	data, err := createK8sConfigurationData(cfg)
	if err != nil {
		return err
	}

	return createFile(k8sConfigurationServiceFilePath, data, k8sConfigurationServiceFilePERM)
}

func createK8sConfigurationData(cfg config.File) ([]byte, error) {
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

	k8sConfigurationServiceBuffer := new(bytes.Buffer)
	err := k8sConfigurationTemplate.Execute(k8sConfigurationServiceBuffer, extraArgs)
	return k8sConfigurationServiceBuffer.Bytes(), err
}
