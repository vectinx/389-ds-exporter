package cmd

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
)

type CmdArguments struct {
	ConfigFile     string
	IsConfigCheck  bool
	PromslogConfig *promslog.Config
}

func ParseCmdArguments() *CmdArguments {
	var (
		configFilePath = kingpin.Flag("config", "Path to configuration file").
				Default("config.yml").
				String()
		checkConfig = kingpin.Flag("check-config", "Validate the current configuration and print it to stdout").Bool()
	)
	args := &CmdArguments{}

	args.PromslogConfig = &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, args.PromslogConfig)

	kingpin.Version(version.Print("389-ds-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	args.ConfigFile = *configFilePath
	args.IsConfigCheck = *checkConfig

	return args
}
