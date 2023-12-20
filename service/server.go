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

	// ServerHealth is used to check the health status of the server.
	ServerHealth() ServerHealth
}

type ServerInfo struct {
	Name string

	ServiceInfo Info
}

// ServerHealth holds information about a server's health details.
//
// Note that this structure does not have indication for liveness because
// when a server is able to respond to health query,
// the server is assumed to be alive.
type ServerHealth struct {
	// Ready indicates that the server is open to accepting clients
	Ready bool

	// Components holds information about server components' health.
	Components map[string]ServerComponentHealth
}

type ServerComponentHealth struct {
	OK bool
}
