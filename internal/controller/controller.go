package controller

import (
	"go.uber.org/zap"

	"github.com/fraima/fraimactl/internal/config"
)

var kindCreator map[string]func(config.Instruction) error = map[string]func(config.Instruction) error{
	"KubeletService":          createKubletService,
	"KubeletConfiguration":    createKubletConfiguration,
	"ContainerdService":       createContainerdService,
	"ContainerdConfiguration": createContainerdConfiguration,
	"SysctlConfiguration":     createSysctlConfiguration,
	"ModProbeConfiguration":   createModProbeConfiguration,
	"DownloadConfiguration":   downloading,
}

func Run(instructions []config.Instruction, skippingPhases []string) {
	for _, i := range instructions {
		handler, isExist := kindCreator[i.Kind]
		if !isExist {
			zap.L().Warn("unknown kind", zap.String("kind", i.Kind))
			continue
		}
		err := handler(i)
		if err != nil {
			zap.L().Error(i.Kind, zap.String("api_version", i.APIVersion), zap.Error(err))
			continue
		}
		zap.L().Info(i.Kind, zap.String("api_version", i.APIVersion))
	}
}
