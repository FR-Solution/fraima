package controller

import (
	"fmt"
	"strings"

	"go.uber.org/zap"

	"github.com/fraima/fraima/internal/config"
)

const (
	v1       = "v1"
	v1alpha1 = "v1alpha1"
	v1beta1  = "v1beta1"
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
		cfg := configuration{
			apiVersion: getVersion(f.ApiVersion),
			extraArgs:  make(map[string]string),
		}

		if f.ExtraArgs != nil {
			extraArgs, ok := f.ExtraArgs.(map[any]any)
			if !ok {
				zap.L().Error("invalid extaArgs", zap.String("kind", f.Kind))
				continue
			}
			for k, v := range extraArgs {
				cfg.extraArgs[fmt.Sprint(k)] = fmt.Sprint(v)
			}
		}

		err := create(cfg)
		if err != nil {
			zap.L().Error("creating file", zap.String("kind", f.Kind), zap.Error(err))
			continue
		}
		zap.L().Info("file created", zap.String("kind", f.Kind))
	}
	return nil
}

func getVersion(str string) string {
	if strings.Contains(str, v1) {
		return v1
	}
	if strings.Contains(str, v1alpha1) {
		return v1alpha1
	}
	if strings.Contains(str, v1beta1) {
		return v1beta1
	}
	return ""
}
