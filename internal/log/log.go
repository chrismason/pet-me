package log

import (
	"fmt"
	"log"
	"time"

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

func NewLogger(lvl Level, cfg *config.ServerConfig) *Logger {
	client := appinsights.NewTelemetryClient(cfg.InstrumentationKey)

	if lvl == DisableLogging || cfg.InstrumentationKey == "" {
		client.SetIsEnabled(false)
	}

	return &Logger{
		Level:  lvl,
		client: client,
	}
}

func (l *Logger) appInsightsLog(lvl Level, msg string) {
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

func (l *Logger) logLog(lvl Level, msg string) {
	if l.Level > lvl {
		return
	}

	log.Println(msg)
}

func (l *Logger) Debug(msg string) {
	if l.client.IsEnabled() {
		l.appInsightsLog(DebugLevel, msg)
	} else {
		l.logLog(DebugLevel, msg)
	}
}

func (l *Logger) Info(msg string) {
	if l.client.IsEnabled() {
		l.appInsightsLog(InfoLevel, msg)
	} else {
		l.logLog(InfoLevel, msg)
	}
}

func (l *Logger) Error(msg string) {
	if l.client.IsEnabled() {
		l.appInsightsLog(ErrorLevel, msg)
	} else {
		l.logLog(ErrorLevel, msg)
	}
}

func (l *Logger) Request(method string, url string, duration time.Duration, statusCode string) {
	if l.client.IsEnabled() {
		l.client.TrackRequest(method, url, duration, statusCode)
	} else {
		l.logLog(InfoLevel, fmt.Sprintf("Method '%s' at '%s' took '%v' with status code '%s'", method, url, duration, statusCode))
	}
}

func (l *Logger) Dependency(name string, target string, success bool) {
	if l.client.IsEnabled() {
		l.client.TrackRemoteDependency(name, "HTTP", target, success)
	} else {
		l.logLog(InfoLevel, fmt.Sprintf("Calling HTTP endpoint '%s' with target '%s' with a success status of %t", name, target, success))
	}
}
