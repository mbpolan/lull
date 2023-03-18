package logger

import (
	"fmt"
	"github.com/mbpolan/lull/internal/system"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"path/filepath"
	"time"
)

var sugar *zap.SugaredLogger

// Setup prepares global logging for the app with a verbosity level.
func Setup(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return errors.Wrap(err, "invalid log level")
	}

	cfgDir, err := system.GetConfigDir()
	if err != nil {
		return errors.Wrap(err, "cannot determine config directory")
	}

	// initialize log directory
	logDir := filepath.Join(cfgDir, "logs")
	err = system.CreateDir(logDir)
	if err != nil {
		return errors.Wrap(err, "could not create log directory")
	}

	// name the log file after the current date/time
	logFile := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02-15:04:05"))

	// configure logging to the log file
	cfg := zap.NewProductionConfig()
	cfg.Encoding = "console"
	cfg.Level = lvl
	cfg.OutputPaths = []string{
		filepath.Join(cfgDir, "logs", logFile),
	}

	logger, err := cfg.Build()
	if err != nil {
		return err
	}

	sugar = logger.Sugar()
	return nil
}

// Infof logs an info-level message with a format and template args.
func Infof(format string, args ...any) {
	sugar.Infof(format, args)
}

// Errorf logs an error-level message with a format and template args.
func Errorf(format string, args ...any) {
	sugar.Errorf(format, args)
}
