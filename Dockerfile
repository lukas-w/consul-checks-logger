FROM golang:1.17 as builder

ARG CGO_ENABLED=0
WORKDIR /app

COPY  go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags "-s"


FROM scratch
COPY --from=builder /app/consul-checks-logger /consul-checks-logger
CMD ["/consul-checks-logger"]
