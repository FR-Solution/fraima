package generator

import (
	_ "embed"
	"fmt"

	"github.com/pelletier/go-toml"

	"github.com/fraima/fraimactl/internal/config"
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

	tomlData, err := toml.Marshal(eargs)
	if err != nil {
		return nil, err
	}

	// TODO:
	// for "github.com/pelletier/go-toml/v2"
	// When will https://github.com/pelletier/go-toml/issues/836 close
	// cc, err := structure.New(new(containerd.Config))
	// if err != nil {
	// 	return nil, err
	// }

	// cc.AddTags(getContainerdTag)

	// err = toml.Unmarshal(tomlData, cc.Struct())
	// if err != nil {
	// 	return nil, err
	// }

	// return toml.Marshal(cc.Struct())

	return tomlData, err
}

// TODO:
// for "github.com/pelletier/go-toml/v2"
// When will https://github.com/pelletier/go-toml/issues/836 close
// var regexpContainerdTag = regexp.MustCompile(`"$`)

// func getContainerdTag(fieldName, fieldTag string) string {
// 	if fieldTag != "" {
// 		fieldTag = regexpContainerdTag.ReplaceAllString(fieldTag, `,omitempty"`)
// 	}
// 	return fieldTag
// }
