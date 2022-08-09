# consul-checks-logger

A simple daemon that watches Consul health checks and writes them as JSON to standard output when they change.
Available as Docker image at `ghcr.io/lukas-w/consul-checks-logger`.
Configure Consul connection using [standard environment variables](https://www.consul.io/commands#environment-variables) `CONSUL_HTTP_ADDR`, `CONSUL_HTTP_TOKEN`, etc.

Sample output:
```json
{
  "CheckID": "oxygen",
  "Name": "Oxygen status",
  "Node": "apollo",
  "Output": "we've had a problem here",
  "ServiceID": "life-support",
  "ServiceName": "Life Support",
  "ServiceTags": [],
  "Status": "warning",
  "Time": "1970-04-22:08:20-05:00",
  "Type": "ttl"
}
```
