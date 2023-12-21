package service

import (
	"context"
)

// Service abstracts all services
type Service interface {
	// ServiceInfo returns basic information about the service.
	ServiceInfo() Info

	// Serve starts the service. This method is blocking and won't return
	// until the service is stopped, e.g., through ShutdownService.
	Serve(context.Context) error

	// ShutdownService gracefully stops the service.
	ShutdownService(context.Context) error

	// ServiceStatus is used to check the status of the service.
	ServiceStatus() Status
}

// Status holds information about a service's status.
//
// Note that this structure does not have indication for liveness because
// when a service is able to respond to status query,
// the service is assumed to be alive.
type Status struct {
	// Ready indicates readiness, i.e., the service is ready to accept clients.
	Ready bool `json:"ready"`

	// Components holds information about service components' status.
	Components map[string]ComponentStatus `json:"components"`
}

// TODO: if component is interprocess, we might want to provide the ping value.
type ComponentStatus struct {
	// IsInterprocess indicates whether the component is on different process.
	IsInterprocess bool `json:"is_interprocess"`

	// Ready indicates the component's readiness.
	//
	// There are three possible values: nil, true, and false.
	//
	// - If the value is nil, it means that component's readiness can not be or
	//   was not be determined.
	// - If the value is false, the component was reporting that it's not ready.
	// - If the value is true, the component was reporting that it's ready.
	Ready *bool `json:"ready,omitempty"`

	Components map[string]ComponentStatus `json:"components"`
}
