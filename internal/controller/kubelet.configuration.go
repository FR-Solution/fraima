package controller

import (
	"k8s.io/kubelet/config/v1beta1"
	"sigs.k8s.io/yaml"

	"github.com/fraima/fraimactl/internal/config"
)

const (
	kubeletConfigurationFilePath = "./config.yaml"
	kubeletConfigurationFilePERM = 0644
)

func createKubletConfiguration(cfg config.Instruction) error {
	data, err := getKubeletConfigurationData(cfg.APIVersion, cfg.Spec)
	if err != nil {
		return err
	}

	return createFile(kubeletConfigurationFilePath, data, kubeletConfigurationFilePERM, "root:root")
}

func getKubeletConfigurationData(apiVersion string, spec any) ([]byte, error) {
	eargs, err := getMap(spec)
	if err != nil {
		return nil, err
	}

	yamlData, err := yaml.Marshal(eargs)
	if err != nil {
		return nil, err
	}

	cfg := new(v1beta1.KubeletConfiguration)
	err = yaml.Unmarshal(yamlData, cfg)
	if err != nil {
		return nil, err
	}

	cfg.APIVersion = apiVersion
	cfg.Kind = "KubeletConfiguration"

	data, err := yaml.Marshal(cfg)
	return data, err
}
