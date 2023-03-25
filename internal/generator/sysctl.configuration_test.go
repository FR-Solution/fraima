package generator

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fraima/fraimactl/internal/config"
)

func TestCreateSysctlConfigurationData(t *testing.T) {
	testCfg := &config.Config{
		ExtraArgs: map[any]any{
			"net.ipv4.ip_forward": 1,
		},
	}

	expectedData := "net.ipv4.ip_forward=1"
	actualData, err := createSysctlConfigurationData(testCfg)
	require.NoError(t, err)
	require.Equal(t, expectedData, string(actualData))
}
