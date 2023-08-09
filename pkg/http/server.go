package http

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

const (
	defaultAddr            = ":8080"
	defaultReadTimeout     = 500 * time.Millisecond
	defaultWriteTimeout    = 500 * time.Millisecond
	defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server          *http.Server
	shutdownTimeout time.Duration
}

func NewServer(handler http.Handler, options ...Option) *Server {
	s := &Server{
		server: &http.Server{
			Addr:         defaultAddr,
			Handler:      handler,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, apply := range options {
		apply(s)
	}

	return s
}

func (s *Server) Run() {
	go func() {
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve error: %v\n", err)
		}
	}()
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
