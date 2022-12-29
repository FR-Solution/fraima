package main

import (
	"flag"

	"go.uber.org/zap"

	"github.com/fraima/fraima/internal/config"
	"github.com/fraima/fraima/internal/controller"
)

var (
	Version = "undefined"
)

func main() {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.Level.SetLevel(zap.DebugLevel)
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	var configFile string
	flag.StringVar(&configFile, "config", "/home/geo/projects/fraima/fraima/config-example.yaml", "path to dir with configs")
	flag.Parse()

	if configFile == "" {
		zap.L().Fatal("the path to the config file is not set")
	}

	zap.L().Debug("configuration", zap.String("version", Version))

	files, err := config.ReadConfig(configFile)
	if err != nil {
		zap.L().Fatal("read config", zap.Error(err))
	}

	zap.L().Info("started")

	err = controller.Generation(files)
	if err != nil {
		zap.L().Fatal("generation", zap.Error(err))
	}

	zap.L().Info("goodbye")
}
