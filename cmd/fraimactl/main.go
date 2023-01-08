package main

import (
	"flag"
	"strings"

	"go.uber.org/zap"

	"github.com/fraima/fraimactl/internal/config"
	"github.com/fraima/fraimactl/internal/controller"
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
	var skipKindList string
	flag.StringVar(&configFile, "config", "/home/geo/projects/fraima/fraima/config-example.yaml", "path to dir with configs")
	flag.StringVar(&skipKindList, "skip-phases", "", `list of skipped phases`)
	flag.Parse()

	if configFile == "" {
		zap.L().Fatal("the path to the config file is not set")
	}

	skippingPhases := strings.Split(skipKindList, " ")

	zap.L().Debug("configuration", zap.String("version", Version), zap.Strings("skipping phases", skippingPhases))

	instructionList, err := config.GetInstructionList(configFile)
	if err != nil {
		zap.L().Fatal("read config", zap.Error(err))
	}

	zap.L().Info("started")

	controller.Run(instructionList, skippingPhases)

	zap.L().Info("goodbye")
}
