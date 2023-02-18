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
	Run(fileType string, instruction config.Instruction) error
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
	wgRun := &sync.WaitGroup{}
	for _, i := range instructions {
		if _, isSkipping := skippingPhases[getPhaseName(i.Kind)]; isSkipping {
			continue
		}

		wgRun.Add(1)
		go func(wgRun *sync.WaitGroup, i config.Instruction) {
			defer wgRun.Done()

			wg := &sync.WaitGroup{}

			if i.Spec.Service != nil {
				wg.Add(1)
				go s.generation(wg, serviceFileType, i)
			}
			if i.Spec.Configuration != nil {
				wg.Add(1)
				go s.generation(wg, configurationFileType, i)
			}
			wg.Add(1)
			go s.downloading(wg, i.Metadata, i.Spec.Download)
			wg.Wait()

			s.starting(i.Kind, i.Spec.Starting)
		}(wgRun, i)
	}
	wgRun.Wait()
}

func (s *controller) generation(wg *sync.WaitGroup, fileType string, instruction config.Instruction) {
	defer wg.Done()
	zap.L().Info("start_generation", zap.Any("apiVersion", instruction.Metadata.APIVersion), zap.Any("kind", instruction.Metadata.Kind), zap.String("type", fileType))
	if err := s.generator.Run(fileType, instruction); err != nil {
		zap.L().Error("generation", zap.Any("apiVersion", instruction.Metadata.APIVersion), zap.Any("kind", instruction.Metadata.Kind), zap.String("type", fileType), zap.Error(err))
		return
	}
	zap.L().Info("finish_generation", zap.Any("apiVersion", instruction.Metadata.APIVersion), zap.Any("kind", instruction.Metadata.Kind), zap.String("type", fileType))
}

func (s *controller) downloading(wg *sync.WaitGroup, meta config.Metadata, instructions []config.DownloadInstruction) {
	defer wg.Done()
	for _, instruction := range instructions {
		zap.L().Info("start_downloading", zap.Any("apiVersion", meta.APIVersion), zap.String("kind", meta.Kind), zap.Any("instruction", instruction))
		if err := s.downloader.Run(instruction); err != nil {
			zap.L().Error("downloading", zap.Any("apiVersion", meta.APIVersion), zap.String("kind", meta.Kind), zap.Any("instruction", instruction), zap.Error(err))
		}
		zap.L().Info("finish_downloading", zap.Any("apiVersion", meta.APIVersion), zap.String("kind", meta.Kind), zap.Any("instruction", instruction))
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

func getPhaseName(kind string) string {
	for _, n := range phaseNames {
		if strings.Contains(strings.ToLower(kind), n) {
			return n
		}
	}
	return ""
}
