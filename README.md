# consul-checks-logger

A simple daemon that watches Consul health checks and writes them as JSON to standard output when they change.
Available as Docker image at `ghcr.io/lukas-w/consul-checks-logger`.
Configure Consul connection using [standard environment variables](https://www.consul.io/commands#environment-variables) `CONSUL_HTTP_ADDR`, `CONSUL_HTTP_TOKEN`, etc.
