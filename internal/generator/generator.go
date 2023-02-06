package generator

import (
	"strings"

	"github.com/fraima/fraimactl/internal/config"
)

type generator struct {
	kindHandlers map[string]map[string]func(apiVersion string, instruction config.Instruction) error
}

func New() *generator {
	return &generator{
		kindHandlers: map[string]map[string]func(apiVersion string, instruction config.Instruction) error{
			"kubelet": {
				"service":       createKubeletService,
				"configuration": createKubeletConfiguration,
			},
			"containerd": {
				"service":       createContainerdService,
				"configuration": createContainerdConfiguration,
			},
			"sysctl": {
				"configuration": createSysctlConfiguration,
			},
			"modprob": {
				"configuration": createModProbeConfiguration,
			},
		},
	}
}

func (s *generator) Run(apiVersion, fileType string, instruction config.Instruction) error {
	handlers, isExist := s.kindHandlers[strings.ToLower(instruction.Kind)]
	if !isExist {
		return errUnknownKind
	}
	handler, isExist := handlers[fileType]
	if !isExist {
		return errUnknownFileType
	}

	return handler(instruction)
}
