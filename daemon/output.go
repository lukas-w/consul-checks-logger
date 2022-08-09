package daemon

import (
	"bufio"
	"encoding/json"
	"github.com/hashicorp/consul/api"
	"os"
	"sync"
	"time"
)

type HealthCheckLogger interface {
	log(checks *api.HealthCheck) error
}

type JsonHealthCheckLogger struct {
	mu      sync.Mutex
	encoder *json.Encoder
	writer  *bufio.Writer
}

func NewJsonHealthCheckLogger() *JsonHealthCheckLogger {
	w := bufio.NewWriter(os.Stdout)
	j := &JsonHealthCheckLogger{
		mu:      sync.Mutex{},
		encoder: json.NewEncoder(w),
		writer:  w,
	}
	j.encoder.SetEscapeHTML(false)
	return j
}

func (l *JsonHealthCheckLogger) log(check *api.HealthCheck) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	err := l.encoder.Encode(map[string]interface{}{
		"Node":        check.Node,
		"CheckID":     check.CheckID,
		"Name":        check.Name,
		"Status":      check.Status,
		"Output":      check.Output,
		"ServiceID":   check.ServiceID,
		"ServiceName": check.ServiceName,
		"ServiceTags": check.ServiceTags,
		"Type":        check.Type,
		"Time":        time.Now().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return l.writer.Flush()
}
