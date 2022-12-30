package controller

import (
	"encoding/json"
	"fmt"

	"github.com/irbgeo/go-structure"
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
	groupVersion, err := schema.ParseGroupVersion(cfg.APIVersion)
	if err != nil {
		return err
	}

	kubeletConfiguration, err := getKubeletConfiguration(cfg.ExtraArgs)
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

	return createFile(kubeletConfigurationFilePath, data, kubeletConfigurationFilePERM)
}

func getKubeletConfiguration(extraArgs any) (*kubeletconfig.KubeletConfiguration, error) {
	var eargs map[string]any
	if extraArgs != nil {
		args, ok := extraArgs.(map[any]any)
		if !ok {
			return nil, fmt.Errorf("args converting is not available")
		}
		eargs = getArgsMap(args)
	}

	jsonData, err := json.Marshal(eargs)
	if err != nil {
		return nil, err
	}

	kc, err := structure.New(new(kubeletconfig.KubeletConfiguration))
	if err != nil {
		return nil, err
	}

	kc.AddTags(getTag)

	err = json.Unmarshal(jsonData, kc.Struct())
	if err != nil {
		return nil, err
	}

	cfg := new(kubeletconfig.KubeletConfiguration)
	err = kc.SaveInto(cfg)
	return cfg, err
}
