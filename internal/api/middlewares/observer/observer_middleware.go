package middleware

import (
	"net/http"
	"time"

	"soccer-api/internal/infrastructure/observer"
)

func New(c Configuration) *ObserverMiddleware {
	return &ObserverMiddleware{
		Observer: c.Observer,
	}
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
func (om *ObserverMiddleware) Observe(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		start := time.Now()

		defer func() {
			om.Observer.Info(r.Context(), "Request completed",
				observer.KV{Key: "status", Value: rw.statusCode},
				observer.KV{Key: "duration", Value: time.Since(start).Seconds()},
			)
		}()

		om.Observer.Info(r.Context(), "Request started",
			observer.KV{Key: "method", Value: r.Method},
			observer.KV{Key: "url", Value: r.URL.String()},
		)

		next.ServeHTTP(rw, r)
	})
}
