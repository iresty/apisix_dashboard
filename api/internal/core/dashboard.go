package core

import (
	"github.com/apache/apisix-dashboard/api/internal/config"
	server "github.com/apache/apisix-dashboard/api/internal/core/server2"
)

type Dashboard struct {
	initialized bool
	config      config.Config
	server      *server.Server
}

func NewDashboard() *Dashboard {
	return &Dashboard{}
}

func (d *Dashboard) SetConfig(cfg config.Config) *Dashboard {
	d.config = cfg
	return d
}

func (d *Dashboard) Start() {
	if !d.initialized {
		d.init()
	}
}

func (d *Dashboard) init() {
	d.initServer()
}

func (d *Dashboard) initServer() {
	d.server = server.NewServer(d.config)
}
