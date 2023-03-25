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
					"apiVersion":   "kubelet.config.k8s.io/v1beta1",
					"kind":         "KubeletConfiguration",
					"registerNode": true,

					"tlsCertFile":       "kubelet-server.pem",
					"tlsPrivateKeyFile": "kubelet-server-key.pem",

					"authentication": map[any]any{
						"anonymous": map[any]any{
							"enabled": false,
						},
						"webhook": map[any]any{
							"cacheTTL": "0s",
							"enabled":  true,
						},
						"x509": map[any]any{
							"clientCAFile": "kubernetes-ca.pem",
						},
					},
					"authorization": map[any]any{
						"mode": "Webhook",
						"webhook": map[any]any{
							"cacheAuthorizedTTL":   "0s",
							"cacheUnauthorizedTTL": "0s",
						},
					},
					"cgroupDriver":              "systemd",
					"clusterDNS":                []string{"29.64.0.10"},
					"clusterDomain":             "cluster.local",
					"cpuManagerReconcilePeriod": "0s",
					"fileCheckFrequency":        "0s",
					"healthzBindAddress":        "127.0.0.1",
					"healthzPort":               10248,
					"httpCheckFrequency":        "0s",
					"imageMinimumGCAge":         "0s",
					"logging": map[any]any{
						"flushFrequency": 0,
						"options": map[any]any{
							"json": map[any]any{
								"infoBufferSize": "0",
							},
						},
						"verbosity": 0,
					},
					"memorySwap":                      map[any]any{},
					"nodeStatusReportFrequency":       "1s",
					"nodeStatusUpdateFrequency":       "1s",
					"resolvConf":                      "/run/systemd/resolve/resolv.conf",
					"runtimeRequestTimeout":           "0s",
					"shutdownGracePeriod":             "15s",
					"shutdownGracePeriodCriticalPods": "5s",
				},
			},
		},
	}

	expectedData := "authentication:\n  anonymous: {}\n  webhook:\n    cacheTTL: 0s\n  x509: {}\nauthorization:\n  webhook:\n    cacheAuthorizedTTL: 0s\n    cacheUnauthorizedTTL: 0s\ncpuManagerReconcilePeriod: 0s\nevictionPressureTransitionPeriod: 0s\nfileCheckFrequency: 0s\nhttpCheckFrequency: 0s\nimageMinimumGCAge: 0s\nnodeStatusReportFrequency: 0s\nnodeStatusUpdateFrequency: 0s\nruntimeRequestTimeout: 0s\nstreamingConnectionIdleTimeout: 0s\nsyncFrequency: 0s\nvolumeStatsAggPeriod: 0s\n"
	actualData, err := getKubeletConfigurationData(testInstruction.Spec.Configuration)
	require.NoError(t, err)
	require.Equal(t, expectedData, string(actualData))
}
