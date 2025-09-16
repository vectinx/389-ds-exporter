package cmd

import (
	"github.com/alecthomas/kingpin/v2"
)

type CmdArguments struct {
	ConfigFile    string
	IsConfigCheck bool
}

func ParseCmdArguments(version string) *CmdArguments {
	var (
		configFilePath = kingpin.Flag("config", "Path to configuration file").
				Default("config.yml").
				String()
		checkConfig = kingpin.Flag("check-config", "Check current configuration and print it to stdout").Bool()
	)
	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	return &CmdArguments{
		ConfigFile:    *configFilePath,
		IsConfigCheck: *checkConfig,
	}
}
