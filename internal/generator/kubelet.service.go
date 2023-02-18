package generator

import (
	"bytes"
	_ "embed"
	"text/template"

	"github.com/fraima/fraimactl/internal/config"
	"github.com/fraima/fraimactl/internal/utils"
)

var (
	//go:embed template/kubelet.service.tmpl
	kubeletTemplateStr string
	kubeletTemplate    = template.Must(template.New("kubelet-service").Parse(kubeletTemplateStr))
)

const (
	kubeletServiceName      = "kubelet"
	kubeletServiceFilePath  = "/etc/systemd/system/kubelet.service"
	kubeletServiceFilePERM  = 0644
	kubeletServiceFileOwner = "root:root"
)

// createKubletService create kubelet.service file.
func createKubeletService(i config.Instruction) error {
	data, err := createKubleteServiceData(i.Spec.Service)
	if err != nil {
		return err
	}

	err = utils.CreateFile(kubeletServiceFilePath, data, kubeletServiceFilePERM, kubeletServiceFileOwner)
	if err != nil {
		return err
	}

	return err
}

func createKubleteServiceData(cfg *config.Config) ([]byte, error) {
	eargs, err := getMap(cfg.ExtraArgs)
	if err != nil {
		return nil, err
	}

	kubletServiceBuffer := new(bytes.Buffer)
	err = kubeletTemplate.Execute(kubletServiceBuffer, eargs)
	return kubletServiceBuffer.Bytes(), err
}
