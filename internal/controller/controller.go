package controller

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/fraima/fraimactl/internal/config"
	"go.uber.org/zap"
)

type generator interface {
	Run(metadata config.Metadata, fileType string, extraArgs any) error
}

type downloader interface {
	Run(instruction config.DownloadInstruction) error
}

var (
	phaseNames []string = []string{
		"containerd",
		"kubelet",
		"modprob",
		"sysctl",
	}

	configurationFileType = "configuration"
	serviceFileType       = "service"
)

type controller struct {
	generator  generator
	downloader downloader
}

func New(
	generator generator,
	downloader downloader,
) *controller {
	return &controller{
		generator:  generator,
		downloader: downloader,
	}
}

func (s *controller) Run(instructions []config.Instruction, skippingPhases map[string]struct{}) {
	for _, i := range instructions {
		if _, isSkipping := skippingPhases[strings.ToLower(i.Kind)]; isSkipping {
			continue
		}

		go func(i config.Instruction) {
			wg := &sync.WaitGroup{}

			wg.Add(3)
			if i.Spec.Service != nil {
				go s.generation(wg, i.Metadata, configurationFileType, i.Spec.Service.ExtraArgs)
			}
			if i.Spec.Configuration != nil {
				go s.generation(wg, i.Metadata, configurationFileType, i.Spec.Configuration.ExtraArgs)
			}
			go s.downloading(wg, i.Kind, i.Spec.Download)
			wg.Wait()
		}(i)
	}
}

func (s *controller) generation(wg *sync.WaitGroup, metadata config.Metadata, fileType string, extraArgs any) {
	defer wg.Done()
	if err := s.generator.Run(metadata, fileType, extraArgs); err != nil {
		zap.L().Error("generation", zap.Any("metadata", metadata), zap.String("type", fileType), zap.Error(err))
	}
}

func (s *controller) downloading(wg *sync.WaitGroup, kind string, instructions []config.DownloadInstruction) {
	defer wg.Done()
	for _, instruction := range instructions {
		if err := s.downloader.Run(instruction); err != nil {
			zap.L().Error("downloading", zap.String("kind", kind), zap.Any("instruction", instruction), zap.Error(err))
		}
	}
}

func (s *controller) starting(kind string, instructions []string) error {
	for _, i := range instructions {
		commands := strings.Split(i, " ")
		var err error
		if len(commands) == 1 {
			err = exec.Command(commands[0]).Run()
		} else {
			err = exec.Command(commands[0], commands[1:]...).Run()
		}

		if err != nil {
			return fmt.Errorf("trying execution command: %s", i)
		}
	}
	return nil
}
