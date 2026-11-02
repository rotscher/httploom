# httploom

A Go workspace with two modules built around a small OpenAPI-described HTTP
API:

- `loomserver` — the API server, using the generated stdlib `net/http` handler.
- `loomclient` — a CLI using the generated typed client.
- `api/openapi.yaml` — the single source of truth for the API.

The API stores texts and "weaves" them into fixed-width lines.

## Layout

```
api/
  openapi.yaml                  OpenAPI 3.0 specification
  oapi-codegen.server.yaml      server codegen config
  oapi-codegen.client.yaml      client codegen config
loomserver/
  main.go                       server entrypoint
  generate.go                   //go:generate directive
  internal/api/api.gen.go       generated models + server interface
  internal/handler/             ServerInterface implementation
  internal/weave/               text wrapping + in-memory store
loomclient/
  main.go                       CLI entrypoint
  generate.go                   //go:generate directive
  internal/api/api.gen.go       generated models + typed client
```

## Requirements

- Go 1.22 or newer (the workspace declares Go 1.27.3)
- Network access the first time the generator runs (it is pinned to
  `oapi-codegen` v2.8.2)
- `golangci-lint` v2 for linting (`go install
  github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest`)

## Generate

The generated files are checked in. After editing `api/openapi.yaml`,
regenerate both sides:

```bash
make generate
```

This runs `go generate ./...` in each module, which invokes
`oapi-codegen` with the matching config file.

## Run

Start the server:

```bash
make run-server
# or: cd loomserver && go run . -addr :8080
```

Use the client from another terminal:

```bash
cd loomclient

go run . health
go run . weave -width 20 "The quick brown fox jumps over the lazy dog"
go run . list
go run . get <id>

# text can also come from stdin
echo "some text to weave" | go run . weave -width 30
```

The client targets `http://localhost:8080` by default. Override it with
`-server` or the `LOOM_SERVER` environment variable:

```bash
go run . -server http://example.com:8080 health
LOOM_SERVER=http://example.com:8080 go run . list
```

Flags must come before the command.

## API

| Method | Path          | Operation      | Description                       |
| ------ | ------------- | -------------- | --------------------------------- |
| GET    | `/health`     | `getHealth`    | Health and version                |
| GET    | `/weaves`     | `listWeaves`   | List stored weaves                |
| POST   | `/weaves`     | `createWeave`  | Weave text to a width (1–200)     |
| GET    | `/weaves/{id}`| `getWeave`     | Fetch a single weave              |

## Tests

```bash
make test
```

## Lint

The repository uses `golangci-lint` v2 with its default linter set
(`errcheck`, `govet`, `ineffassign`, `staticcheck`, `unused`), configured in
`.golangci.yml`. Both modules are linted:

```bash
make lint
```
