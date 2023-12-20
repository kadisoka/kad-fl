package service

import (
	"context"
)

// Server abstracts all service servers
type Server interface {
	// ServerInfo returns basic information about the server.
	ServerInfo() ServerInfo

	// Serve starts the server. This method is blocking and won't return
	// until the server is stopped, e.g., through Shutdown.
	Serve(context.Context) error

	// Shutdown gracefully stops the server.
	Shutdown(context.Context) error

	// IsAcceptingClients returns true if the service is ready to serve clients.
	IsAcceptingClients() bool

	// IsHealthy returns true if the service is considerably healthy.
	IsHealthy() bool
}

type ServerInfo struct {
	Name string

	ServiceInfo Info
}
