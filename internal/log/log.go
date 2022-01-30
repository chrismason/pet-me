package log

import (
	"github.com/chrismason/pet-me/internal/config"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/microsoft/ApplicationInsights-Go/appinsights/contracts"
)

type Level int

const (
	DebugLevel Level = 0
	InfoLevel
	ErrorLevel
	DisableLogging
)

type Logger struct {
	Level  Level
	client appinsights.TelemetryClient
}

func New(lvl Level, cfg *config.ServerConfig) *Logger {
	client := appinsights.NewTelemetryClient("")

	if lvl == DisableLogging || cfg.InstrumentationKey == "" {
		client.SetIsEnabled(false)
	}

	return &Logger{
		Level:  lvl,
		client: client,
	}
}

func (l *Logger) log(lvl Level, msg string) {
	if l.Level > lvl {
		return
	}

	var sev contracts.SeverityLevel
	switch int(l.Level) {
	case 0:
		sev = appinsights.Verbose
	case 1:
		sev = appinsights.Information
	case 2:
		sev = appinsights.Error
	}

	l.client.TrackTrace(msg, sev)
}

func (l *Logger) Debug(msg string) {
	l.log(DebugLevel, msg)
}

func (l *Logger) Info(msg string) {
	l.log(InfoLevel, msg)
}

func (l *Logger) Error(msg string) {
	l.log(ErrorLevel, msg)
}
