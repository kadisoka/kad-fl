package webservice

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"sync"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/kadisoka/kad-fl/service"
)

type ServerConfig struct {
	ServePort int `env:"SERVE_PORT"`

	// ServiceInfo is used to override server's service info.
	ServiceInfo *service.Info `env:"-"`

	// Logger is used to pass the logger to the server instance.
	Logger *slog.Logger `env:"-"`
}

func NewServer(
	serviceIdentifier string,
	handler http.Handler,
	config ServerConfig,
) (*Server, error) {
	if config.ServePort <= 0 {
		return nil, errors.Arg("config.ServePort invalid")
	}
	if handler == nil {
		return nil, errors.Arg("handler is unspecified")
	}

	log := config.Logger
	if log == nil {
		log = slog.Default()
	}
	log = log.With("svc", serviceIdentifier)

	return &Server{
		log:     log,
		config:  config,
		handler: handler,
	}, nil
}

// Server is a wrapper for http.Server.
type Server struct {
	mutex sync.RWMutex
	log   *slog.Logger

	config       ServerConfig
	shuttingDown bool
	httpServer   *http.Server
	handler      http.Handler
}

var _ service.Service = &Server{}

var serviceInfo = service.Info{
	Name:        "Generic HTTP Service",
	Description: "A generic service for aggregating web servers",
}

func (srv *Server) ServiceInfo() service.Info {
	var svcInfo *service.Info

	func() {
		srv.mutex.RLock()
		defer srv.mutex.RUnlock()

		svcInfo = srv.config.ServiceInfo
	}()

	if svcInfo != nil {
		return *svcInfo
	}
	return serviceInfo
}

func (srv *Server) Serve(_ context.Context) error {
	var httpServer *http.Server

	err := func() error {
		srv.mutex.Lock()
		defer srv.mutex.Unlock()

		if srv.httpServer != nil {
			return errors.Msg("server is already running")
		}
		if srv.shuttingDown {
			return errors.Msg("server is shutting down")
		}

		httpServer = &http.Server{
			Addr:    fmt.Sprintf(":%d", srv.config.ServePort),
			Handler: srv.handler}
		srv.httpServer = httpServer

		return nil
	}()
	if err != nil {
		return err
	}

	err = httpServer.ListenAndServe()

	srv.mutex.RLock()
	defer srv.mutex.RUnlock()

	if err == nil {
		if !srv.shuttingDown {
			return errors.Msg("server stopped unexpectedly")
		}
		srv.log.Info("Done.")
		return nil
	}
	if err == http.ErrServerClosed && srv.shuttingDown {
		srv.log.Info("Done.")
		return nil
	}
	return err
}

func (srv *Server) ShutdownService(ctx context.Context) error {
	srv.mutex.Lock()
	srv.shuttingDown = true
	httpServer := srv.httpServer
	srv.mutex.Unlock()

	return httpServer.Shutdown(ctx)
}

func (srv *Server) ServiceStatus() service.Status {
	status := service.Status{}

	func() {
		srv.mutex.RLock()
		defer srv.mutex.RUnlock()

		if srv.httpServer != nil && !srv.shuttingDown {
			status.Ready = true
		}
	}()

	return status
}
