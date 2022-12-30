package controller

import (
	_ "embed"
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
	extraArgs := make(map[string]any)
	if cfg.ExtraArgs != nil {
		args, ok := cfg.ExtraArgs.(map[any]any)
		if !ok {
			return nil, fmt.Errorf("args converting is not available")
		}
		for k, v := range args {
			extraArgs[fmt.Sprint(k)] = v
		}
	}

	st, err := structure.New(new(containerd.Config))
	if err != nil {
		return nil, err
	}

	err = st.AssignFrom(extraArgs)
	if err != nil {
		return nil, err
	}

	data, err := toml.Marshal(st.Struct())
	return data, err
}
