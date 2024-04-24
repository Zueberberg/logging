package logging

import (
	"os"
	"time"
)

type LogLevel uint

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

var configLevelNames = map[string]LogLevel{
	"DEBUG":   DEBUG,
	"INFO":    INFO,
	"WARNING": WARNING,
	"ERROR":   ERROR,
}

type Logger struct {
	Name     string
	LogLevel LogLevel
	Handlers []Handlerable
}

type LogValues struct {
	LoggerName string
	Time       time.Time
	Level      LogLevel
	LevelName  string
	Msg        []any
}

func (l Logger) log(vals LogValues) {
	if l.LogLevel <= vals.Level {
		now := time.Now()
		for _, handler := range l.Handlers {
			vals.Time = now
			handler.WriteLog(vals)
		}
	}
}

func (l Logger) Fatal(a ...any) {
	l.log(LogValues{
		LoggerName: l.Name,
		Level:      ERROR,
		LevelName:  logMsgFatal,
		Msg:        a,
	})
	os.Exit(1)
}

func (l Logger) Error(a ...any) {
	l.log(LogValues{
		LoggerName: l.Name,
		Level:      ERROR,
		LevelName:  logMsgError,
		Msg:        a,
	})
}

func (l Logger) Fail(a ...any) {
	l.log(LogValues{
		LoggerName: l.Name,
		Level:      DEBUG,
		LevelName:  logMsgFail,
		Msg:        a,
	})
}

func (l Logger) Trace(a ...any) {
	l.log(LogValues{
		LoggerName: l.Name,
		Level:      ERROR,
		LevelName:  logMsgTrace,
		Msg:        a,
	})
}

func (l Logger) Success(a ...any) {
	l.log(LogValues{
		LoggerName: l.Name,
		Level:      DEBUG,
		LevelName:  logMsgSuccess,
		Msg:        a,
	})
}

func (l Logger) Warning(a ...any) {
	l.log(LogValues{
		LoggerName: l.Name,
		Level:      WARNING,
		LevelName:  logMsgWarning,
		Msg:        a,
	})
}

func (l Logger) Info(a ...any) {
	l.log(LogValues{
		LoggerName: l.Name,
		Level:      INFO,
		LevelName:  logMsgInfo,
		Msg:        a,
	})
}

func (l Logger) Debug(a ...any) {
	l.log(LogValues{
		LoggerName: l.Name,
		Level:      DEBUG,
		LevelName:  logMsgDebug,
		Msg:        a,
	})
}

func (l Logger) WithHandler(h Handlerable) Logger {
	l.Handlers = append(l.Handlers, h)
	return l
}

func (l Logger) WithName(n string) Logger {
	l.Name = n
	return l
}

func (l Logger) WithLevel(ll LogLevel) Logger {
	l.LogLevel = ll
	return l
}
