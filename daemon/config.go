package daemon

import (
	"github.com/hashicorp/consul/api"
	"time"
)

type Config struct {
	ConsulConfig *api.Config
	ErrorTimeout time.Duration
}

func DefaultConfig() Config {
	c := Config{
		ErrorTimeout: 5 * time.Second,
	}
	c.ConsulConfig = api.DefaultConfig()
	return c
}
