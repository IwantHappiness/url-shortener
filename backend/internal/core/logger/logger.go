package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerContextKey struct{}

var key = loggerContextKey{}

type Logger struct {
	*zap.Logger

	file *os.File
}

// TODO: add logger interface

func ToContext(ctx context.Context, logger *Logger) context.Context {
	return context.WithValue(ctx, key, logger)
}

func FromContext(ctx context.Context) *Logger {
	logger, ok := ctx.Value(key).(*Logger)
	if !ok {
		panic("no logger in context")
	}

	return logger
}

func NewLogger(config Config) (*Logger, error) {
	zapLvL := zap.NewAtomicLevel()
	if err := zapLvL.UnmarshalText([]byte(config.Level)); err != nil {
		return nil, fmt.Errorf("unmarshal log level: %w", err)
	}

	if err := os.MkdirAll(config.Folder, 0755); err != nil {
		return nil, fmt.Errorf("create log folder: %w", err)
	}

	timestamp := time.Now().UTC().Format("2006-01-02T15-04-05.00000")
	logFilePath := filepath.Join(config.Folder, fmt.Sprintf("%s.log", timestamp))

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY /* |os.O_APPEND */, 0644)
	if err != nil {
		return nil, fmt.Errorf("open log file: %w", err)
	}

	zapConfig := zap.NewDevelopmentEncoderConfig()
	zapConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02T15:04:05.00000")

	zapEncoder := zapcore.NewConsoleEncoder(zapConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(zapEncoder, zapcore.AddSync(os.Stdout), zapLvL),
		zapcore.NewCore(zapEncoder, zapcore.AddSync(logFile), zapLvL),
	)

	zapLogger := zap.New(core, zap.AddCaller())

	return &Logger{Logger: zapLogger, file: logFile}, nil
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{
		Logger: l.Logger.With(fields...),
		file:   l.file,
	}
}

func (l *Logger) Close() {
	if err := l.file.Close(); err != nil {
		fmt.Println("failed to close application logger:", err)
	}
}
