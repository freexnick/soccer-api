package server

import (
	"context"
	"net/http"
	"soccer-api/internal/infrastructure/observer"
)

type Configuration struct {
	Observer     *observer.Observer
	Port         string
	Handler      http.Handler
	ReadTimeout  uint
	WriteTimeout uint
	IdleTimeout  uint
}

type Lifecycle interface {
	Start(ctx context.Context) error
	Close(ctx context.Context) error
}

type Client struct {
	httpS  *http.Server
	observ *observer.Observer
}
