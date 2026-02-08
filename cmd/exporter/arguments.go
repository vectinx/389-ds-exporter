package main

import (
	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/common/promslog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

// Arguments provides a structure for storing the cli parameters of the exporter.
type Arguments struct {
	ConfigFile           string
	IsConfigCheck        bool
	PromslogConfig       *promslog.Config
	ExporterToolkitFlags *web.FlagConfig
	MetricsPath          string
}

// ParseArguments parses the arguments of the conmad string into the CmdArguments structure.
func ParseArguments() (*kingpin.Application, *Arguments) {
	app := kingpin.New("389-ds-exporter", "389 Directory Server Prometheus exporter")
	args := &Arguments{}

	configFilePath := new(string)
	checkConfig := app.Flag("config.check", "Validate the current configuration and print it to stdout").Bool()
	metricsPath := app.Flag("web.metrics.path", "Path under which to expose metrics.").Default("/metrics").String()
	toolkitFlags := kingpinflag.AddFlags(app, ":9389")

	app.Flag("config.file", "Path to configuration file").
		Default("config.yml").
		StringVar(configFilePath)

	app.Flag("config", "[DEPRECATED]. Use --config.file instead.").
		Hidden().Default("config.yml").
		StringVar(configFilePath)

	args.PromslogConfig = &promslog.Config{}
	flag.AddFlags(app, args.PromslogConfig)

	app.Version(version.Print("389-ds-exporter"))
	app.HelpFlag.Short('h')

	app.Action(func(*kingpin.ParseContext) error {
		args.ConfigFile = *configFilePath
		args.IsConfigCheck = *checkConfig
		args.ExporterToolkitFlags = toolkitFlags
		args.MetricsPath = *metricsPath
		return nil
	})

	return app, args
}
