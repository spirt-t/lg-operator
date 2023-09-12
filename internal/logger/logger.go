package logger

import (
	"fmt"

	"github.com/spirt-t/lg-operator/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logLevelKey = "log.level"
)

// NewLogger with level defined in config.
func NewLogger(cfg config.Manager) (*zap.Logger, error) {
	var levelVal string

	if err := cfg.UnmarshalKey(logLevelKey, &levelVal); err != nil {
		return nil, err
	}

	level, err := zapcore.ParseLevel(levelVal)
	if err != nil {
		return nil, fmt.Errorf("fail to define log_level: %w", err)
	}

	encoderConf := zap.NewDevelopmentEncoderConfig()
	encoderConf.EncodeLevel = zapcore.CapitalColorLevelEncoder

	zapCfg := zap.Config{
		Level:         zap.NewAtomicLevelAt(level),
		Encoding:      "console",
		EncoderConfig: encoderConf,
		OutputPaths:   []string{"stdout"},
	}

	return zapCfg.Build()
}
