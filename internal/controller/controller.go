package controller

import (
	"go.uber.org/zap"

	"github.com/fraima/fraima/internal/config"
)

var kindCreator map[string]func(configuration) error = map[string]func(configuration) error{
	"KubeletServiceConfiguration": createKubletService,
}

func Generation(files []config.File) error {
	for _, f := range files {
		create, isExist := kindCreator[f.Kind]
		if !isExist {
			zap.L().Warn("unknown kind", zap.String("kind", f.Kind))
			continue
		}
		var (
			cfg configuration
			ok  bool
		)

		if f.ExtraArgs != nil {
			cfg.extraArgs, ok = f.ExtraArgs.(map[string]any)
			if !ok {
				zap.L().Error("invalid extaArgs", zap.String("kind", f.Kind))
				continue
			}
		}

		err := create(cfg)
		if err != nil {
			zap.L().Error("creating file", zap.String("kind", f.Kind), zap.Error(err))
			continue
		}
	}
	return nil
}
