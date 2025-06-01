package logger

import "github.com/rs/zerolog"

type LogLevel int8

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
)

type LogFormat int8

const (
	LogFormatPlain LogFormat = iota
	LogFormatJSON
)

type Logger struct {
	stdoutLogger zerolog.Logger
	stderrLogger zerolog.Logger
}

type Configuration struct {
	LogFormat      LogFormat
	LogLevel       LogLevel
	SkipFrameCount int
	AppVersion     string
	GitCommit      string
}

type KV = struct {
	Key   string
	Value any
}
