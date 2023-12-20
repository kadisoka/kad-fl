//

package app

import (
	"context"
	"sync"

	"github.com/alloyzeus/go-azfl/errors"
	"github.com/kadisoka/kad-fl/app"
	"github.com/kadisoka/kad-fl/service"
)

// App abstracts the application itself. There should be only one instance
// for a running instance of an app.
type App interface {
	app.App

	InstanceID() string

	AddServiceServer(service.Server)
	IsAllServiceServersReady() bool

	Run(context.Context)
}

// AppBase is the base layer for an app.
type AppBase struct {
	appInfo    app.Info
	instanceID string

	servers   []service.Server
	serversMu sync.RWMutex
}

var _ App = &AppBase{}

func (appBase *AppBase) AppInfo() app.Info { return appBase.appInfo }

func (appBase *AppBase) InstanceID() string { return appBase.instanceID }

// AddServiceServer adds a server to be run simultaneously. Do NOT call this
// method after the app has been started.
func (appBase *AppBase) AddServiceServer(srv service.Server) {
	appBase.serversMu.Lock()
	appBase.servers = append(appBase.servers, srv)
	appBase.serversMu.Unlock()
}

// Run runs all the servers. Do NOT add any new server after this method
// was called.
func (appBase *AppBase) Run(ctx context.Context) {
	service.RunServers(ctx, appBase.ServiceServers(), nil)
}

// IsAllServiceServersReady checks if every server is ready to accept clients.
func (appBase *AppBase) IsAllServiceServersReady() bool {
	servers := appBase.ServiceServers()
	for _, srv := range servers {
		if health := srv.ServerHealth(); !health.Ready {
			return false
		}
	}
	return true
}

// ServiceServers returns an array of servers added to this app.
func (appBase *AppBase) ServiceServers() []service.Server {
	out := make([]service.Server, len(appBase.servers))
	appBase.serversMu.RLock()
	copy(out, appBase.servers)
	appBase.serversMu.RUnlock()
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

// Instantiate global instance of App with the default implementation.
func Init(appInfo app.Info) (App, error) {
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
		unameStr, err = unameString()
		if err != nil {
			return
		}

		var taskID string
		taskID, _, err = getECSTaskID()
		if err != nil {
			return
		}

		var instanceID string
		if taskID != "" {
			if unameStr != "" {
				instanceID = taskID + " (" + unameStr + ")"
			} else {
				instanceID = taskID
			}
		} else {
			instanceID = unameStr
		}

		defApp = &AppBase{
			appInfo:    appInfo,
			instanceID: instanceID,
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
