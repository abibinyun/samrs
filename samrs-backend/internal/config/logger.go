package config

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLogFile       = "./logs/samrs.log"
	defaultLogMaxSizeMB  = 100
	defaultLogMaxBackups = 7
	defaultLogMaxAgeDays = 7
	defaultLogCompress   = true
)

// NewLogger builds a production logger that writes JSON logs to stdout
// and to a rotating file (lumberjack) for local persistence.
func NewLogger() (*zap.Logger, func()) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoder := zapcore.NewJSONEncoder(encoderCfg)
	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level),
	}

	logFile := strings.TrimSpace(os.Getenv("LOG_FILE"))
	if logFile == "" {
		logFile = defaultLogFile
	}

	if logFile != "-" {
		_ = os.MkdirAll(filepath.Dir(logFile), 0o755)
		fileSink := &lumberjack.Logger{
			Filename:   logFile,
			MaxSize:    getEnvInt("LOG_MAX_SIZE_MB", defaultLogMaxSizeMB),
			MaxBackups: getEnvInt("LOG_MAX_BACKUPS", defaultLogMaxBackups),
			MaxAge:     getEnvInt("LOG_MAX_AGE_DAYS", defaultLogMaxAgeDays),
			Compress:   getEnvBool("LOG_COMPRESS", defaultLogCompress),
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(fileSink), level))
	}

	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	cleanup := func() {
		_ = logger.Sync()
	}
	return logger, cleanup
}

func getEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		return fallback
	}
	return val
}

func getEnvBool(key string, fallback bool) bool {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	val, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return val
}
