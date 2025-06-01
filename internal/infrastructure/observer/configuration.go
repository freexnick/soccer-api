package observer

import "context"

type KV = struct {
	Key   string
	Value any
}

type Observer struct {
	logger Logger
}

type Logger interface {
	Error(context.Context, error, ...KV)
	Warn(context.Context, string, ...KV)
	Info(context.Context, string, ...KV)
	Debug(context.Context, string, ...KV)
	Close(context.Context) error
}

type Configuration struct {
	Logger Logger
}
