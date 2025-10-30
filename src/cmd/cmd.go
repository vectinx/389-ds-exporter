package cmd

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
)

type CmdArguments struct {
	ConfigFile     string
	IsConfigCheck  bool
	PromslogConfig *promslog.Config
}

func ParseCmdArguments(version string) *CmdArguments {
	var (
		configFilePath = kingpin.Flag("config", "Path to configuration file").
				Default("config.yml").
				String()
		checkConfig = kingpin.Flag("check-config", "Validate the current configuration and print it to stdout").Bool()
	)
	args := &CmdArguments{}

	args.PromslogConfig = &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, args.PromslogConfig)

	kingpin.Version(version)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	args.ConfigFile = *configFilePath
	args.IsConfigCheck = *checkConfig

	return args
}
