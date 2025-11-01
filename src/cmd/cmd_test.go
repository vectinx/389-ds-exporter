package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDefaultValues(t *testing.T) {
	app, args := ParseCmdArguments()
	_, err := app.Parse([]string{})
	require.NoError(t, err, "Empty cmd args should be parsed successfully")
	require.Equal(t, "config.yml", args.ConfigFile, "--config.file default value should be config.yml")
	require.Equal(t, "/metrics", args.MetricsPath, "--web.metrics.path default value should be /metrics")
}

func TestDeprecatedConfigFlag(t *testing.T) {
	app, args := ParseCmdArguments()
	_, err := app.Parse([]string{"--config", "old.yml"})
	require.NoError(t, err, "Parsing args with deprecated backward-compatibility '--config' flag should not fail")
	require.Equal(t, "old.yml", args.ConfigFile, "--config value should be correctly parsed to args")
}
