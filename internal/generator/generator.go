package generator

import (
	"go.uber.org/zap"

	"github.com/fraima/fraimactl/internal/config"
)

var kindCreator map[string]func(config.Generate) error = map[string]func(config.Generate) error{
	"KubeletService":             createKubletService,
	"KubeletConfiguration":       createKubletConfiguration,
	"ContainerdService":          createContainerdService,
	"ContainerdConfiguration":    createContainerdConfiguration,
	"SysctlNetworkConfiguration": createSysctlNetworkConfiguration,
	"ModProbConfiguration":       createModProbConfiguration,
}

func Run(generateList []config.Generate, skippingKinds map[string]struct{}) {
	for _, f := range generateList {
		if _, isSkipping := skippingKinds[f.Kind]; isSkipping {
			continue
		}
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
}
