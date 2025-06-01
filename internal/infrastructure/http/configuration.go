package server

import (
	"net/http"
	"soccer-api/internal/infrastructure/observer"
)

type Configuration struct {
	Observ       *observer.Observer
	Port         string
	Handler      http.Handler
	ReadTimeout  uint
	WriteTimeout uint
	IdleTimeout  uint
}

type Client struct {
	httpS  *http.Server
	observ *observer.Observer
}
