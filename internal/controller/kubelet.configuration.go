package controller

import (
	"strings"

	"github.com/irbgeo/go-structure"
	"gopkg.in/yaml.v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"k8s.io/kubernetes/pkg/kubelet/kubeletconfig/util/codec"

	"github.com/fraima/fraimactl/internal/config"
)

const (
	kubeletConfigurationFilePath = "/etc/kubernetes/kubelet/config.yaml"
	kubeletConfigurationFilePERM = 0644
)

func createKubletConfiguration(cfg config.Instruction) error {
	groupVersion, err := schema.ParseGroupVersion(cfg.APIVersion)
	if err != nil {
		return err
	}

	kubeletConfiguration, err := getKubeletConfiguration(cfg.Spec)
	if err != nil {
		return err
	}
	kubeletConfiguration.TypeMeta = metav1.TypeMeta{
		Kind:       cfg.Kind,
		APIVersion: cfg.APIVersion,
	}

	data, err := codec.EncodeKubeletConfig(kubeletConfiguration, groupVersion)
	if err != nil {
		return err
	}

	return createFile(kubeletConfigurationFilePath, data, kubeletConfigurationFilePERM, "root:root")
}

func getKubeletConfiguration(spec any) (*kubeletconfig.KubeletConfiguration, error) {
	eargs, err := getMap(spec)
	if err != nil {
		return nil, err
	}

	yamlData, err := yaml.Marshal(eargs)
	if err != nil {
		return nil, err
	}

	kc, err := structure.New(new(kubeletconfig.KubeletConfiguration))
	if err != nil {
		return nil, err
	}

	kc.ChangeTags(getTag)

	err = yaml.Unmarshal(yamlData, kc.Struct())
	if err != nil {
		return nil, err
	}

	cfg := new(kubeletconfig.KubeletConfiguration)
	err = kc.SaveInto(cfg)
	return cfg, err
}

var prevFieldName string

func getTag(fieldName, fieldTag, fieldType string) string {
	if strings.ToLower(fieldType) == "duration" && strings.ToLower(fieldName) != "duration" {
		prevFieldName = fieldName
		return `yaml:",inline"`
	}
	if strings.Contains(fieldTag, "name=duration") {
		return strings.ToLower(prevFieldName)
	}
	if fieldTag == "" {
		return strings.ToLower(fieldName)
	}
	return ""
}
