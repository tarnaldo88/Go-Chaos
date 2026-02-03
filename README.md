# Go Chaos

<img src="images/GoChaosLogo.png" alt="Go Chaos Logo" width="150">

Go Chaos is a reverse proxy that injects controlled failures and latency into HTTP traffic to validate API resilience. It supports live configuration updates, per-route rules, transport-level failures, and a built-in admin UI.

## Features
- HTTP reverse proxy with configurable target
- Inject latency, HTTP errors, and connection drops
- Transport-level chaos: DNS failures and upstream timeouts
- Per-route rules with include and exclude path prefixes
- Live config updates via admin API
- Built-in admin UI for live tuning

## Requirements
- Go 1.22 or later

## Quick Start

### 1) Create config
Create `config/config.yaml`:

```yaml
listen_addr: ":8080"
target_url: "http://localhost:9000"
chaos:
  error_rate: 0.05
  disconnect_rate: 0.02
  latency_min_ms: 100
  latency_max_ms: 300
  upstream_timeout_rate: 0.10
  dns_failure_rate: 0.05
  include_paths:
    - "/api/"
  exclude_paths:
    - "/healthz"
    - "/admin/"
    - "/images/"
```

### 2) Build and run
```bash
go build -o go-chaos ./cmd/go-chaos
./go-chaos
```

Or run directly:
```bash
go run ./cmd/go-chaos
```

### 3) Admin UI
Open:
```
http://localhost:8080/admin
```

### 4) Health check
```
http://localhost:8080/healthz
```

## Admin UI Screenshot

<img src="images/AdminPage.png" alt="Go Chaos Admin Page" width="450">

## Configuration

### Root
- `listen_addr` string, required
- `target_url` string, required

### Chaos
- `error_rate` float, 0.0 to 1.0
- `disconnect_rate` float, 0.0 to 1.0
- `latency_min_ms` int, minimum latency in ms
- `latency_max_ms` int, maximum latency in ms
- `latency_ms` int, fixed latency in ms, used when min and max are zero
- `upstream_timeout_rate` float, 0.0 to 1.0
- `dns_failure_rate` float, 0.0 to 1.0
- `include_paths` list of path prefixes to apply chaos
- `exclude_paths` list of path prefixes to bypass chaos

Rules:
- If `include_paths` is empty, chaos applies to all paths except those in `exclude_paths`.
- If `include_paths` is set, chaos applies only to matching prefixes.

## Admin API

### Get current config
```
GET /admin/config
```

### Update config
```
POST /admin/config
Content-Type: application/x-yaml
```

Example:
```bash
curl -X POST http://localhost:8080/admin/config \
  -H "Content-Type: application/x-yaml" \
  --data-binary @config/config.yaml
```

## Admin UI

The UI allows live editing of all chaos settings without restarting the server.

- Load current config from the server
- Edit key fields in a form
- Sync raw YAML and form values
- Validate inputs before saving

## Static Assets

If you serve a favicon or logo locally, add routes for them and exclude from chaos:

- Serve: `/images/*` or `/favicon.ico`
- Exclude: `/images/` and `/favicon.ico`

## Tests
```bash
go test ./...
```

## Project Layout
```
cmd/go-chaos/            main entrypoint
internal/config/         config schema and store
internal/server/         HTTP server and admin endpoints
internal/chaos/          chaos injection logic
internal/proxy/          reverse proxy and transport
internal/observability/  logging utilities
```

## License
MIT
