package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LogrusLogger struct {
	log *logrus.Logger
}

func NewLogrusLogger() *LogrusLogger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
		PrettyPrint:     true,
	})
	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(os.Stdout)
	return &LogrusLogger{
		log: log,
	}
}

func (l *LogrusLogger) WithError(err error) LoggerItf {
	return &LogrusEntry{
		entry: l.log.WithError(err),
	}
}

func (l *LogrusLogger) Debug(args ...any) {
	l.log.Debug(args...)
}

func (l *LogrusLogger) Debugf(format string, args ...any) {
	l.log.Debugf(format, args...)
}

func (l *LogrusLogger) Info(args ...any) {
	l.log.Info(args...)
}

func (l *LogrusLogger) Infof(format string, args ...any) {
	l.log.Infof(format, args...)
}

func (l *LogrusLogger) Warn(args ...any) {
	l.log.Warn(args...)
}

func (l *LogrusLogger) Warnf(format string, args ...any) {
	l.log.Warnf(format, args...)
}

func (l *LogrusLogger) Error(args ...any) {
	l.log.Error(args...)
}

func (l *LogrusLogger) Errorf(format string, args ...any) {
	l.log.Errorf(format, args...)
}

func (l *LogrusLogger) Fatal(args ...any) {
	l.log.Fatal(args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...any) {
	l.log.Fatalf(format, args...)
}

func (l *LogrusLogger) WithField(key string, value any) LoggerItf {
	return &LogrusEntry{
		entry: l.log.WithField(key, value),
	}
}

func (l *LogrusLogger) WithFields(fields map[string]any) LoggerItf {
	return &LogrusEntry{
		entry: l.log.WithFields(fields),
	}
}
