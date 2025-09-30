package utils

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	slogmulti "github.com/samber/slog-multi"

	"389-ds-exporter/src/config"
)

const LogFileMode os.FileMode = 0o644

// BuildLogHandler creates a log handler with the specified type.
func BuildLogHandler(format string, w io.Writer, level slog.Level) slog.Handler {
	addSource := level == slog.LevelDebug
	switch format {
	case "text":
		return slog.NewTextHandler(w, &slog.HandlerOptions{AddSource: addSource, Level: level})
	case "json":
		return slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: addSource, Level: level})
	}

	return slog.Default().Handler()
}

// SetupLogger creates a logger based on the provided configuration.
func SetupLogger(cfg *config.ExporterConfig) (*slog.Logger, *os.File, error) {
	var logLevel slog.Level
	handlers := []slog.Handler{}
	var logFile *os.File

	strLogLevel := cfg.LogLevel
	levelMap := map[string]slog.Level{
		"DEBUG":   slog.LevelDebug,
		"INFO":    slog.LevelInfo,
		"WARNING": slog.LevelWarn,
		"ERROR":   slog.LevelError,
	}

	logLevel, ok := levelMap[strLogLevel]
	if !ok {
		return nil, nil, fmt.Errorf("unknown logging level: '%s'", strLogLevel)
	}

	if cfg.LogHandler == "stdout" || cfg.LogHandler == "both" {
		handlers = append(handlers, BuildLogHandler(cfg.LogStdoutFormat, os.Stdout, logLevel))
	}
	if cfg.LogHandler == "file" || cfg.LogHandler == "both" {
		var err error
		logFile, err = os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, LogFileMode)
		if err != nil {
			return nil, nil, fmt.Errorf("error opening log file: %w", err)
		}
		handlers = append(handlers, BuildLogHandler(cfg.LogFileFormat, logFile, logLevel))
	}

	if len(handlers) == 0 {
		return nil, nil, errors.New("unable to create logger - logging handlers not specified")
	}

	logger := slog.New(slogmulti.Fanout(handlers...))

	return logger, logFile, nil
}

// ReopenLogFile reopens the log file. This function is needed to handle log rotation.
func ReopenLogFile(cfg *config.ExporterConfig, old_file *os.File) (*os.File, error) {
	if old_file != nil {
		_ = old_file.Close()
	}

	newLogFile, err := os.OpenFile(cfg.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, LogFileMode)
	if err != nil {
		return nil, fmt.Errorf("failed to open new log file: %w", err)
	}

	newLogger, _, err := SetupLogger(cfg)
	if err != nil {
		_ = newLogFile.Close()

		return nil, fmt.Errorf("failed to set up new logger: %w", err)
	}

	slog.SetDefault(newLogger)
	return newLogFile, nil
}
