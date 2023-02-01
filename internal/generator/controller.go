package generator

import (
	"github.com/fraima/fraimactl/internal/config"
)

type generator struct {
	kindHandler map[string]func(config.Instruction) error
}

func New() *generator {
	return &generator{
		kindHandler: map[string]func(config.Instruction) error{
			"KubeletService":          createKubletService,
			"KubeletConfiguration":    createKubletConfiguration,
			"ContainerdService":       createContainerdService,
			"ContainerdConfiguration": createContainerdConfiguration,
			"SysctlConfiguration":     createSysctlConfiguration,
			"ModProbeConfiguration":   createModProbeConfiguration,
		},
	}
}

func (s *generator) Run(instruction config.Instruction) error {
	handler, isExist := s.kindHandler[instruction.Kind]
	if !isExist {
		return errUnknownKind
	}
	err := handler(instruction)
	if err != nil {
		return err
	}
	return nil
}
