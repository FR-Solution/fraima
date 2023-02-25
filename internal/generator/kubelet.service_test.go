package generator

import (
	"testing"
)

func TestCreateKubleteServiceData(t *testing.T) {
	// testCfg := &config.Config{
	// 	ExtraArgs: map[any]any{
	// 		"config": "/etc/kubernetes/containerd/config.toml",
	// 	},
	// }

	// expectedData := "[Unit]\nDescription=containerd container runtime\nDocumentation=https://containerd.io\nAfter=network.target\n\n[Service]\nExecStartPre=/sbin/modprobe overlay\nExecStart=/usr/bin/containerd --config=/etc/kubernetes/containerd/config.toml\nRestart=always\nRestartSec=5\nDelegate=yes\nKillMode=process\nOOMScoreAdjust=-999\nLimitNOFILE=1048576\nLimitNPROC=infinity\nLimitCORE=infinity\n\n[Install]\nWantedBy=multi-user.target"
	// actualData, err := createContainerdServiceData(testCfg)
	// require.NoError(t, err)
	// require.Equal(t, expectedData, string(actualData))
}
