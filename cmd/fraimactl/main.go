package main

import (
	"flag"
	"strings"

	"go.uber.org/zap"

	"github.com/fraima/fraimactl/internal/config"
	"github.com/fraima/fraimactl/internal/controller"
	"github.com/fraima/fraimactl/internal/downloader"
	"github.com/fraima/fraimactl/internal/generator"
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
	// TODO: clear
	flag.StringVar(&configFile, "config", "/home/irbgeo/projects/fraima/fraima/config-debug.yaml", "path to dir with configs")
	flag.StringVar(&skipKindList, "skip-phases", "", `list of skipped phases`)
	flag.Parse()

	if configFile == "" {
		zap.L().Fatal("the path to the config file is not set")
	}

	skippingPhasesList := strings.Split(skipKindList, " ")
	skippingPhasesMap := make(map[string]struct{})
	for _, p := range skippingPhasesList {
		skippingPhasesMap[p] = struct{}{}
	}

	zap.L().Debug("configuration", zap.String("version", Version), zap.Strings("skipping phases", skippingPhasesList))

	instructionList, err := config.GetInstructionList(configFile)
	if err != nil {
		zap.L().Fatal("read config", zap.Error(err))
	}

	generator := generator.New()

	downloader := downloader.New()

	cntl := controller.New(
		generator,
		downloader,
	)

	zap.L().Info("started")

	cntl.Run(instructionList, skippingPhasesMap)

	zap.L().Info("goodbye")
}
