//

package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/alloyzeus/go-azfl/v2/errors"
	"github.com/kadisoka/kad-fl/app"
	"github.com/kadisoka/kad-fl/service"
)

// App abstracts the application itself. There should be only one instance
// for a running instance of an app.
type App interface {
	app.App

	InstanceID() string

	AddService(service.Service)
	IsAllServicesReady() bool

	Run(context.Context)

	// Status reports the app's overall status including all services.
	// This method is designed to be used in health checks.
	Status() service.Status
}

// AppBase is the base layer for an app.
type AppBase struct {
	mutex sync.RWMutex
	log   *slog.Logger

	appInfo    app.Info
	instanceID string
	isStarted  bool

	services []service.Service
}

var _ App = &AppBase{}

func (appBase *AppBase) AppInfo() app.Info { return appBase.appInfo }

func (appBase *AppBase) InstanceID() string { return appBase.instanceID }

// AddService adds a service to be run simultaneously. Do NOT call this
// method after the app has been started.
func (appBase *AppBase) AddService(srv service.Service) {
	appBase.mutex.Lock()
	defer appBase.mutex.Unlock()

	if !appBase.isStarted {
		appBase.services = append(appBase.services, srv)
	}
}

// Run runs all the services. Do NOT add any new service after this method
// was called.
func (appBase *AppBase) Run(ctx context.Context) {
	services := appBase.Services()
	if len(services) == 0 {
		return
	}

	log := appBase.log

	if !func() bool {
		appBase.mutex.Lock()
		defer appBase.mutex.Unlock()

		if appBase.isStarted {
			log.Info("App is already running")
			return false
		}
		appBase.isStarted = true
		return true
	}() {
		return
	}

	var shutdownSignal <-chan os.Signal

	// Start a go-routine to detect that all services are ready
	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			if appBase.IsAllServicesReady() {
				log.Info("All services are ready")
				break
			}
		}
	}()

	// Used to determine if all services have stopped.
	var servicesStopWaiter sync.WaitGroup

	log.Info("Starting all services...")

	// Start the services
	for _, srv := range services {
		servicesStopWaiter.Add(1)
		go func(innerSrv service.Service) {
			srvName := innerSrv.ServiceInfo().Name
			log.Info(fmt.Sprintf("Starting %s...", srvName))
			err := innerSrv.Serve(ctx)
			if err != nil {
				log.Error(fmt.Sprintf("%s: %v", srvName, err))
				os.Exit(-1)
			} else {
				log.Info(fmt.Sprintf("%s stopped", srvName))
			}
			servicesStopWaiter.Done()
		}(srv)
	}

	// We set up the signal handler (interrupt and terminate).
	// We are using the signal to gracefully and forcefully stop the service.
	if shutdownSignal == nil {
		sigChan := make(chan os.Signal, 1)
		// Listen to interrupt and terminal signals
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		shutdownSignal = sigChan
	}

	// Wait for the shutdown signal
	select {
	case <-shutdownSignal:
		log.Info("Got shutdown signal")
	case <-ctx.Done():
		log.Info("Context is done")
	}

	log.Info("Shutting down all services...")

	// Start another go-routine to catch another signal so the shutdown
	// could be forced. If we get another signal, we'll exit immediately.
	go func() {
		<-shutdownSignal

		log.Info("Forced shutdown.")
		os.Exit(0)
	}()

	// Gracefully shutdown the services
	shutdownCtx, shutdownCtxCancel := context.WithTimeout(
		context.Background(), 15*time.Second)
	defer shutdownCtxCancel()

	for _, srv := range services {
		go func(innerSrv service.Service) {
			srvName := innerSrv.ServiceInfo().Name
			log.Info(fmt.Sprintf("Shutting down %s...", srvName))
			err := innerSrv.ShutdownService(shutdownCtx)
			if err != nil {
				log.Error(fmt.Sprintf("Service %s shutdown with error: %v", srvName, err))
			}
		}(srv)
	}

	// Wait for all services to stop.
	servicesStopWaiter.Wait()

	log.Info("Services are gracefully stopped.")
}

func (appBase *AppBase) Status() service.Status {
	services := appBase.Services()

	components := map[string]service.ComponentStatus{}
	allReady := true

	for _, svc := range services {
		status := svc.ServiceStatus()
		if !status.Ready {
			allReady = false
		}
		comp := service.ComponentStatus{
			IsInterprocess: false,
			Ready:          &status.Ready,
			Components:     status.Components,
		}
		components[svc.ServiceInfo().Name] = comp
	}

	return service.Status{
		Ready:      allReady,
		Components: components,
	}
}

// IsAllServicesReady checks if every service is ready to accept clients.
func (appBase *AppBase) IsAllServicesReady() bool {
	services := appBase.Services()
	for _, srv := range services {
		if status := srv.ServiceStatus(); !status.Ready {
			return false
		}
	}
	return true
}

// Services returns an array of services added to this app.
func (appBase *AppBase) Services() []service.Service {
	out := make([]service.Service, len(appBase.services))
	appBase.mutex.RLock()
	copy(out, appBase.services)
	appBase.mutex.RUnlock()
	return out
}

var (
	defApp     App
	defAppOnce sync.Once
)

func Instance() App {
	if defApp == nil {
		panic("App has not been initialized. Call app.Init to initialize App.")
	}
	return defApp
}

type InitOpts struct {
	Logger *slog.Logger `env:"-"`
}

// Instantiate global instance of App with the default implementation.
func Init(appInfo app.Info, opts InitOpts) (App, error) {
	err := errors.Msg("app instance already initialized")
	defAppOnce.Do(func() {
		err = nil

		if appInfo.BuildInfo.RevisionID == "" {
			err = errors.ArgMsg("appInfo.BuildInfo.RevisionID", "empty")
			return
		}
		if appInfo.BuildInfo.Timestamp == "" {
			err = errors.ArgMsg("appInfo.BuildInfo.RevisionID", "empty")
			return
		}

		var unameStr string
		unameStr, err = os.Hostname()
		if err != nil {
			return
		}

		instanceID := fmt.Sprintf("%v@%s", os.Getpid(), unameStr)

		logger := opts.Logger
		if logger == nil {
			logger = slog.Default()
		}
		logger = logger.With("svc", "_app")

		defApp = &AppBase{
			appInfo:    appInfo,
			instanceID: instanceID,
			log:        logger,
		}
	})

	if err != nil {
		return nil, err
	}

	return defApp, nil
}

// InitCustom
func InitCustom(customApp App) error {
	err := errors.Msg("app instance already initialized")
	defAppOnce.Do(func() {
		err = nil
		defApp = customApp
	})
	return err
}
