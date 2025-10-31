package cmd

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

type CmdArguments struct {
	ConfigFile           string
	IsConfigCheck        bool
	PromslogConfig       *promslog.Config
	ExporterToolkitFlags *web.FlagConfig
	MetricsPath          string
}

func ParseCmdArguments() *CmdArguments {
	var (
		configFilePath = kingpin.Flag("config.file", "Path to configuration file").
				Default("config.yml").
				String()
		checkConfig = kingpin.Flag("config.check", "Validate the current configuration and print it to stdout").Bool()
		metricsPath = kingpin.Flag(
			"web.metrics.path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		toolkitFlags = kingpinflag.AddFlags(kingpin.CommandLine, ":9389")
	)
	args := &CmdArguments{}

	args.PromslogConfig = &promslog.Config{}
	flag.AddFlags(kingpin.CommandLine, args.PromslogConfig)

	kingpin.Version(version.Print("389-ds-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	args.ConfigFile = *configFilePath
	args.IsConfigCheck = *checkConfig
	args.ExporterToolkitFlags = toolkitFlags
	args.MetricsPath = *metricsPath

	return args
}
