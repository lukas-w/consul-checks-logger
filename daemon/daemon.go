package daemon

import (
	"context"
	"github.com/hashicorp/consul/api"
)

type Daemon struct {
	context  context.Context
	cancel   context.CancelFunc
	config   Config
	client   *api.Client
	watchers map[string]*ChecksWatcher
	logger   HealthCheckLogger
}

func NewDaemon(config Config) (*Daemon, error) {
	client, err := api.NewClient(config.ConsulConfig)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &Daemon{
		context:  ctx,
		cancel:   cancel,
		config:   config,
		client:   client,
		watchers: make(map[string]*ChecksWatcher),
		logger:   NewJsonHealthCheckLogger(),
	}, nil
}

func (d *Daemon) Run() {
	newService, delService := make(chan string), make(chan string)
	go WatchServiceChanges(d.client, d.context, newService, delService, d.config)

out:
	for {
		select {
		case name := <-newService:
			d.watchService(name)
		case name := <-delService:
			d.unwatchService(name)
		case <-d.context.Done():
			break out
		}
	}
}

func (d *Daemon) Stop() {
	d.cancel()
}

func (d *Daemon) watchService(id string) {
	w := NewChecksWatcher(id, d.client, d.context, d.config)
	go w.Run(d.logger)
	d.watchers[id] = w
}

func (d *Daemon) unwatchService(id string) {
	w := d.watchers[id]
	delete(d.watchers, id)
	w.Stop()
}
