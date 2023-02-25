package generator

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fraima/fraimactl/internal/config"
)

func TestContainerdConfigurationData(t *testing.T) {
	testCfg := &config.Config{
		ExtraArgs: map[any]any{
			"version": 2,
			"plugins": map[any]any{
				"io.containerd.grpc.v1.cri": map[any]any{
					"containerd": map[any]any{
						"runtimes": map[any]any{
							"runc": map[any]any{
								"runtime_type": "io.containerd.runc.v2",
								"options": map[any]any{
									"SystemdCgroup": true,
								},
							},
						},
					},
				},
			},
		},
	}

	expectedData := "version = 2\n\n[plugins]\n  [plugins.'io.containerd.grpc.v1.cri']\n    [plugins.'io.containerd.grpc.v1.cri'.containerd]\n      [plugins.'io.containerd.grpc.v1.cri'.containerd.runtimes]\n        [plugins.'io.containerd.grpc.v1.cri'.containerd.runtimes.runc]\n          runtime_type = 'io.containerd.runc.v2'\n\n          [plugins.'io.containerd.grpc.v1.cri'.containerd.runtimes.runc.options]\n            SystemdCgroup = true\n"
	actualData, err := createContainerdConfigurationData(testCfg)
	require.NoError(t, err)
	require.Equal(t, expectedData, string(actualData))
}
