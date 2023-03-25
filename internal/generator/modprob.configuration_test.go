package generator

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/fraima/fraimactl/internal/config"
)

func TestCreateModProbeConfigurationData(t *testing.T) {
	testCfg := &config.Config{
		ExtraArgs: []any{
			"br_netfilter",
			"overlay",
		},
	}

	expectedData := "br_netfilter\noverlay\n"
	actualData, err := createModProbeConfigurationData(testCfg)
	require.NoError(t, err)
	require.Equal(t, expectedData, string(actualData))
}
