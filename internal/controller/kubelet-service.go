package controller

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"text/template"
)

var (
	//go:embed template/kubelet.service.tmpl
	kubeletTemplateStr string
	kubeletTemplate    = template.Must(template.New("kubelet-service").Parse(kubeletTemplateStr))

	args map[string]string = map[string]string{
		"node-labels":                  "true",
		"register-node":                "true",
		"cloud-provider":               "external",
		"image-pull-progress-deadline": "2m",
		"feature-gates":                "RotateKubeletServerCertificate: true",
		"cert-dir":                     "/etc/kubernetes/pki/certs/kubelet",
		"authorization-mode":           "Webhook",
		"anonymous-auth":               "false",
		"cni-bin-dir":                  "/opt/cni/bin",
		"cni-conf-dir":                 "/etc/cni/net.d",
		"network-plugin":               "cni",
		"config":                       "/etc/kubernetes/kubelet/config.yaml",
		"root-dir":                     "/var/lib/kubelet",
		"v":                            "2",
		"kubeconfig":                   "/etc/kubernetes/kubelet/kubeconfig",
		"bootstrap-kubeconfig":         "/etc/kubernetes/kubelet/bootstrap-kubeconfig",
		"container-runtime":            "remote",
		"container-runtime-endpoint":   "/run/containerd/containerd.sock",
		"pod-infra-container-image":    "k8s.gcr.io/pause:3.6",
	}
)

const (
	kubeletServiceFilePath = "/etc/systemd/system/kubelet.service"
	kubeletServiceFilePERM = 0644
)

// KubeletService create kubelet.service file.
func createKubletService(cfg configuration) error {
	for k, v := range cfg.extraArgs {
		args[k] = fmt.Sprint(v)
	}

	kubletServiceBuffer := new(bytes.Buffer)
	err := kubeletTemplate.Execute(kubletServiceBuffer, args)
	if err != nil {
		return err
	}

	err = os.WriteFile(kubeletServiceFilePath, kubletServiceBuffer.Bytes(), kubeletServiceFilePERM)
	if err != nil {
		return err
	}

	err = os.Chown(kubeletServiceFilePath, os.Getuid(), os.Getgid())
	return err
}
