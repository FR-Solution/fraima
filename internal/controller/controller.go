package controller

import "github.com/fraima/fraimactl/internal/config"

type generator interface {
	Run(kind string, instruction GenerateInstruction) error
}

type downloader interface {
	Run(kind string, instruction DownloadInstruction) error
}

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

func (s *controller) Run(instructions []config.Instruction) {

}

func (s *controller) phasesSplit(instructions []config.Instruction) []any

func (s *controller) getPhaseName() (string, error)
