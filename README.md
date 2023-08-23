# ArgoCD Exporter

The purpose of this project is to provide a Prometheus metric for each ArgoCD application, which can be used to monitor versions of the application.

## Versions codemap

| Status | Codename |
| ------- | -------- |
| `up-to-date` | 0 |
| `pending-upgrade`  | 1 |
| `pending-rollback` | 2 |
| `unknown` | 3 |

## Flags

| Flag | Description | Default |
| ---- | ----------- | ------- |
| `address` | Address to listen on for web interface and telemetry. | `:` |
| `port` | Port to listen on for web interface and telemetry. | `8080` |
| `metrics-path` | Path under which to expose metrics. | `/metrics` |
| `verbosity` | Log level (debug, info, warn, error, fatal, panic). | `info` |
| `update-interval` | Interval between updates of the metrics. | `1` |
