package daemon

import (
	"context"
	"github.com/hashicorp/consul/api"
	"log"
	"time"
)

type serviceMap map[string][]string

func WatchServices(client *api.Client, ctx context.Context, servicesChan chan<- serviceMap, errorTimeout time.Duration) {
	catalog := client.Catalog()
	q := (&api.QueryOptions{}).WithContext(ctx)

	for {
		services, qm, err := catalog.Services(q)
		select {
		case <-ctx.Done():
			break
		default:
		}

		if err != nil {
			log.Printf("failed getting consul services, retrying in %s: %s", errorTimeout, err)
			servicesChan <- serviceMap{}
			time.Sleep(errorTimeout)
		} else {
			q.WaitIndex = qm.LastIndex
			servicesChan <- services
		}
	}
}

func WatchServiceChanges(client *api.Client, ctx context.Context, newService chan<- string, delService chan<- string, config Config) {
	servicesChan := make(chan serviceMap)
	go WatchServices(client, ctx, servicesChan, config.ErrorTimeout)

	lastServices := make(serviceMap)

	for {
		select {
		case services := <-servicesChan:
			for service, _ := range services {
				if _, found := lastServices[service]; !found {
					newService <- service
				}
			}
			for service := range lastServices {
				if _, found := services[service]; !found {
					delService <- service
				}
			}
			lastServices = services
		case <-ctx.Done():
			break
		}
	}
}
