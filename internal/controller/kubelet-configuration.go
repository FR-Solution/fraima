package controller

import (
	"encoding/json"
	"fmt"
	"unicode"

	tag "github.com/irbgeo/struct-tag-builder"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	"k8s.io/kubernetes/pkg/kubelet/kubeletconfig/util/codec"

	"github.com/fraima/fraima/internal/config"
)

const (
	kubeletConfigurationFilePath = "/etc/kubernetes/kubelet/config.yaml"
	kubeletConfigurationFilePERM = 0644
)

func createKubletConfiguration(cfg config.File) error {
	data, err := createKubeletConfigurationData(cfg)
	if err != nil {
		return err
	}

	return createFile(kubeletConfigurationFilePath, data, kubeletConfigurationFilePERM)
}

func createKubeletConfigurationData(cfg config.File) ([]byte, error) {
	groupVersion, err := schema.ParseGroupVersion(cfg.APIVersion)
	if err != nil {
		return nil, err
	}

	kubeletConfiguration, err := getKubeletConfiguration(cfg.ExtraArgs)
	if err != nil {
		return nil, err
	}
	kubeletConfiguration.TypeMeta = metav1.TypeMeta{
		Kind:       cfg.Kind,
		APIVersion: cfg.APIVersion,
	}

	return codec.EncodeKubeletConfig(kubeletConfiguration, groupVersion)
}

func getKubeletConfiguration(extraArgs any) (*kubeletconfig.KubeletConfiguration, error) {
	eargs := make(map[string]any)
	if extraArgs != nil {
		args, ok := extraArgs.(map[any]any)
		if !ok {
			return nil, fmt.Errorf("args converting is not available")
		}
		for k, v := range args {
			eargs[fmt.Sprint(k)] = v
		}
	}

	jsonData, err := json.Marshal(eargs)
	if err != nil {
		return nil, err
	}

	cfg := new(kubeletconfig.KubeletConfiguration)
	tb, err := tag.NewTagBuilder(cfg)
	if err != nil {
		return nil, err
	}

	c := tb.Build(getTag)
	err = json.Unmarshal(jsonData, c.Writable())
	if err != nil {
		return nil, err
	}

	err = c.SaveInto(cfg)
	return cfg, err
}

func getTag(fieldName string) string {
	for i, v := range fieldName {
		tagValue := string(unicode.ToLower(v)) + fieldName[i+1:]
		return fmt.Sprintf(`yaml:"%s"`, tagValue)
	}
	return ""
}
