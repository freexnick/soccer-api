package server

import (
	"context"
	"errors"

	"fmt"
	"net/http"
	"time"
)

func New(conf Configuration) (*Client, error) {
	if conf.Handler == nil || conf.Port == "" {
		return nil, errors.New("missing required options")
	}

	return &Client{
		observ: conf.Observ,
		httpS: &http.Server{
			Addr:         conf.Port,
			Handler:      conf.Handler,
			ReadTimeout:  time.Duration(conf.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(conf.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(conf.IdleTimeout) * time.Second,
		},
	}, nil
}

func (s *Client) Start(ctx context.Context) error {
	s.observ.Info(ctx, fmt.Sprintf("HTTP server is listening on %s", s.httpS.Addr))

	err := s.httpS.ListenAndServe()

	if err != nil && err != http.ErrServerClosed {
		s.observ.Error(ctx, fmt.Errorf("HTTP server failed: %v", err))
		return err
	}

	return nil
}

func (s *Client) Close(ctx context.Context) error {
	err := s.httpS.Shutdown(ctx)

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		s.observ.Error(ctx, fmt.Errorf("error shutting down HTTP server: %v", err))
		return err
	}

	s.observ.Info(ctx, "HTTP server shut down gracefully")
	return nil
}
