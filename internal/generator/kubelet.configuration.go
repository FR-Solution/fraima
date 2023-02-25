package generator

import (
	"k8s.io/kubelet/config/v1beta1"
	"sigs.k8s.io/yaml"

	"github.com/fraima/fraimactl/internal/config"
	"github.com/fraima/fraimactl/internal/utils"
)

const (
	kubeletConfigurationFilePath = "./config.yaml"
	kubeletConfigurationFilePERM = 0644
)

func createKubeletConfiguration(i config.Instruction) error {
	data, err := getKubeletConfigurationData(i)
	if err != nil {
		return err
	}

	return utils.CreateFile(kubeletConfigurationFilePath, data, kubeletConfigurationFilePERM, "root:root")
}

func getKubeletConfigurationData(i config.Instruction) ([]byte, error) {
	eargs, err := getMap(i.Spec.Configuration.ExtraArgs)
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

	cfg.APIVersion = i.APIVersion
	cfg.Kind = "KubeletConfiguration"

	data, err := yaml.Marshal(cfg)

	return data, err
}
