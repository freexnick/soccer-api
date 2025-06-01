package middleware

import (
	"net/http"

	"soccer-api/internal/infrastructure/observer"
)

type Configuration struct {
	Observer *observer.Observer
}

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

type ObserverMiddleware struct {
	Observer *observer.Observer
}
