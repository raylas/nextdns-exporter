# nextdns-exporter

[![Main](https://github.com/raylas/nextdns-exporter/actions/workflows/main.yaml/badge.svg)](https://github.com/raylas/nextdns-exporter/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/raylas/nextdns-exporter)](https://goreportcard.com/report/github.com/raylas/nextdns-exporter)

A [Prometheus exporter](https://prometheus.io/docs/instrumenting/exporters/) for [NextDNS data](https://nextdns.github.io/api/#analytics).

## Configuration

- `LOG_LEVEL` (default: `INFO`)
- `METRICS_PORT` (default: `9948`)
- `METRICS_PATH` (default: `/metrics`)
- `NEXTDNS_RESULT_WINDOW` (default: `-5m`)
- `NEXTDNS_RESULT_LIMIT` (default: `50`)
- `NEXTDNS_PROFILE` (required: [docs](https://nextdns.github.io/api/#profile))
- `NEXTDNS_API_KEY` (required: [docs](https://nextdns.github.io/api/#authentication))

For most accurate data, the scrape interval _should_ match the value set via `NEXTDNS_RESULT_WINDOW`.

Recommended scrape timeout is `10` seconds.

## Usage

### Binary

Either [download a recent release](https://github.com/raylas/nextdns-exporter/releases) or compile the binary yourself (`go build -o nextdns-exporter`) and run:
```sh
export NEXTDNS_PROFILE=<profile_id>
export NEXTDNS_API_KEY=<api_key>
./nextdns-exporter
2022-11-19T12:39:34.479-0800 [INFO]  starting exporter: port=:9948 path=/metrics
```

### Docker

```sh
docker run -d \
  -p 9948:9948 \
  -e NEXTDNS_PROFILE=<profile_id> \
  -e NEXTDNS_API_KEY=<api_key> \
  ghcr.io/raylas/nextdns-exporter
```

### Compose

The following will create a local stack of the exporter, Prometheus, and Grafana:
```sh
NEXTDNS_PROFILE=<profile_id> \
NEXTDNS_API_KEY=<api_key> \
docker-compose up -d
```

Access Grafana by navigating to [http://localhost:3000](http://localhost:3000).

**Note:** Data will take a few minutes to trickle in.

## Dashboard

A basic Grafana dashboard is found [here](/grafana/dashboards/nextdns.json).

![Grafana dashboard](/grafana/dashboards/nextdns.png)
