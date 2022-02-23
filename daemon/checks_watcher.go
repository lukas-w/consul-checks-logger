package daemon

import (
	"context"
	"github.com/hashicorp/consul/api"
	"log"
	"time"
)

type checksMap map[string]*api.HealthCheck

type ChecksWatcher struct {
	service string
	client  *api.Client
	ctx     context.Context
	cancel  context.CancelFunc
	config  Config
	done    chan int
}

func NewChecksWatcher(service string, client *api.Client, ctx context.Context, config Config) *ChecksWatcher {
	ctx, cancel := context.WithCancel(ctx)

	return &ChecksWatcher{
		service: service,
		client:  client,
		ctx:     ctx,
		cancel:  cancel,
		config:  config,
		done:    make(chan int),
	}
}

func (w *ChecksWatcher) Run(logger HealthCheckLogger) {
	health := w.client.Health()
	q := (&api.QueryOptions{}).WithContext(w.ctx)

	checksState := make(checksMap)

out:
	for {
		checks, qm, err := health.Checks(w.service, q)
		select {
		case <-w.ctx.Done():
			break out
		default:
		}

		logChecksChange(checksState, checks, logger)
		checksState = make(checksMap)
		for _, check := range checks {
			checksState[check.CheckID] = &api.HealthCheck{
				CheckID: check.CheckID,
				Output:  check.Output,
				Status:  check.Status,
			}
		}

		if err != nil {
			log.Printf("failed getting consul checks for service %s, retrying in %s: %s", w.service, w.config.ErrorTimeout, err)
			time.Sleep(w.config.ErrorTimeout * time.Second)
			continue
		}
		q.WaitIndex = qm.LastIndex
	}

	w.done <- 1
}

func (w *ChecksWatcher) Stop() {
	w.cancel()
	<-w.done
}

func logChecksChange(a checksMap, b api.HealthChecks, logger HealthCheckLogger) {
	for _, check := range b {
		id := check.CheckID
		lastCheck := a[id]
		if lastCheck != nil && check.Status == lastCheck.Status && check.Output == lastCheck.Output {
			continue
		}
		if check.Status == api.HealthPassing && (lastCheck == nil || lastCheck.Status == api.HealthPassing) {
			continue
		}
		err := logger.log(check)
		if err != nil {
			log.Printf("failed logging check %s: %s", id, err)
		}
	}
}
