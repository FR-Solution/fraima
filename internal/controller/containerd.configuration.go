package controller

import (
	_ "embed"
	"encoding/json"
	"fmt"

	containerd "github.com/containerd/containerd/cmd/containerd/command"
	"github.com/irbgeo/go-structure"
	"github.com/pelletier/go-toml"

	"github.com/fraima/fraima/internal/config"
)

const (
	containerdConfigurationFilePath = "/etc/kubernetes/containerd/config.toml"
	containerdConfigurationFilePERM = 0644
)

// createContainerdConfiguration create containerd.service file.
func createContainerdConfiguration(cfg config.File) error {
	data, err := createContainerdConfigurationData(cfg)
	if err != nil {
		return err
	}

	return createFile(containerdConfigurationFilePath, data, containerdConfigurationFilePERM)
}

func createContainerdConfigurationData(cfg config.File) ([]byte, error) {
	var eargs map[string]any
	if cfg.ExtraArgs != nil {
		args, ok := cfg.ExtraArgs.(map[any]any)
		if !ok {
			return nil, fmt.Errorf("args converting is not available")
		}
		eargs = getArgsMap(args)
	}

	jsonData, err := json.Marshal(eargs)
	if err != nil {
		return nil, err
	}

	cc, err := structure.New(new(containerd.Config))
	if err != nil {
		return nil, err
	}

	cc.AddTags(getTag)

	err = json.Unmarshal(jsonData, cc.Struct())
	if err != nil {
		return nil, err
	}

	containerdCfg := new(containerd.Config)
	err = cc.SaveInto(containerdCfg)
	if err != nil {
		return nil, err
	}

	data, err := toml.Marshal(containerdCfg)
	return data, err
}
