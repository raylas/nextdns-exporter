# nextdns-exporter

[![Main](https://github.com/raylas/nextdns-exporter/actions/workflows/main.yml/badge.svg)](https://github.com/raylas/nextdns-exporter/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/raylas/nextdns-exporter)](https://goreportcard.com/report/github.com/raylas/nextdns-exporter)

A [Prometheus exporter](https://prometheus.io/docs/instrumenting/exporters/) for [NextDNS data](https://nextdns.github.io/api/#analytics).

## Configuration

- `LOG_LEVEL` (default: `INFO`)
- `METRICS_PORT` (default: `9948`)
- `METRICS_PATH` (default: `/metrics`)
- `NEXTDNS_PROFILE` (required: [docs](https://nextdns.github.io/api/#profile))
- `NEXTDNS_API_KEY` (required: [docs](https://nextdns.github.io/api/#authentication))
- `NEXTDNS_FILTER_FROM` (default: `-1d`)
- `NEXTDNS_RESULT_LIMIT` (default: `50`)

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
  -e NEXDTDNS_PROFILE=<profile_id> \
  -e NEXTDNS_API_KEY=<api_key> \
  ghcr.io/raylas/nextdns-exporter
```
