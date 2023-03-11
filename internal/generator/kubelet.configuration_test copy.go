package generator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fraima/fraimactl/internal/config"
)

func TestCreateKubleteConfigurationData(t *testing.T) {
	testInstruction := config.Instruction{
		Spec: config.Spec{
			Configuration: &config.Config{
				ExtraArgs: map[any]any{
					"v":                          2,
					"cert-dir":                   "/etc/kubernetes/pki/certs/kubelet",
					"config":                     "/etc/kubernetes/kubelet/config.yaml",
					"kubeconfig":                 "/etc/kubernetes/kubelet/kubeconfig",
					"bootstrap-kubeconfig":       "/etc/kubernetes/kubelet/bootstrap-kubeconfig",
					"container-runtime-endpoint": "/run/containerd/containerd.sock",
					"container-runtime":          "remote",
					"pod-infra-container-image":  "k8s.gcr.io/pause:3.6",
					"cloud-provider":             "external",
				},
			},
		},
	}

	expectedData := "[Unit]\nDescription=kubelet: The Kubernetes Node Agent\nDocumentation=https://kubernetes.io/docs/home/\nWants=network-online.target\nAfter=network-online.target\n\n[Service]\nExecStart=/usr/bin/kubelet \\\n--bootstrap-kubeconfig=/etc/kubernetes/kubelet/bootstrap-kubeconfig \\\n--cert-dir=/etc/kubernetes/pki/certs/kubelet \\\n--cloud-provider=external \\\n--config=/etc/kubernetes/kubelet/config.yaml \\\n--container-runtime=remote \\\n--container-runtime-endpoint=/run/containerd/containerd.sock \\\n--kubeconfig=/etc/kubernetes/kubelet/kubeconfig \\\n--pod-infra-container-image=k8s.gcr.io/pause:3.6 \\\n--v=2\nRestart=always\nStartLimitInterval=0\nRestartSec=10\n\n[Install]\nWantedBy=multi-user.target\n"
	actualData, err := getKubeletConfigurationData(testInstruction)
	require.NoError(t, err)
	require.Equal(t, expectedData, string(actualData))
}
