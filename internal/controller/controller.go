package controller

import (
	"go.uber.org/zap"

	"github.com/fraima/fraima/internal/config"
)

var kindCreator map[string]func(config.File) error = map[string]func(config.File) error{
	"KubeletService":             createKubletService,
	"KubeletConfiguration":       createKubletConfiguration,
	"ContainerdService":          createContainerdService,
	"ContainerdConfiguration":    createContainerdConfiguration,
	"SysctlNetworkConfiguration": createSysctlNetworkConfiguration,
	"K8sConfiguration":           createK8sConfiguration,
}

func Generation(files []config.File) error {
	for _, f := range files {
		create, isExist := kindCreator[f.Kind]
		if !isExist {
			zap.L().Warn("unknown kind", zap.String("kind", f.Kind))
			continue
		}
		err := create(f)
		if err != nil {
			zap.L().Error("creating file", zap.String("kind", f.Kind), zap.Error(err))
			continue
		}
		zap.L().Info("file created", zap.String("kind", f.Kind))
	}
	return nil
}
