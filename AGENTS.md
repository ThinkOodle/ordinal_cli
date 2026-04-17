# AGENTS.md - Ordinal CLI

## Project Overview

A command-line interface (CLI) tool for the [Ordinal API](https://docs.tryordinal.com/api/introduction). Built in Go, the CLI provides full coverage of the in-scope Core REST API surface, allowing users to manage social posts, ideas, approvals, comments, engagements, subscribers, labels, analytics, webhooks, Slack boosts, and related workspace entities directly from the terminal.

## API Reference

- **Base URL:** `https://app.tryordinal.com/api/v1`
- **Auth:** Bearer token via `Authorization: Bearer <api_key>` header
- **Rate Limit:** 100 requests per minute per API key (`429 Too Many Requests` on overage)
- **Docs:** https://docs.tryordinal.com/api/introduction
- **OpenAPI:** https://docs.tryordinal.com/api/openapi.json

## Architecture

```
ordinal_cli/
  cmd/                     # Cobra command definitions
    root.go                # Root command, global flags, JSON body helpers
    auth.go                # Save API key to config
    skill.go               # Install the bundled SKILL.md for agent tools
    posts.go               # Post commands
    ideas.go               # Idea commands
    ...                    # One file per resource group
  internal/
    api/                   # Typed API methods (one file per resource group)
      posts.go
      ideas.go
      ...
    client/                # HTTP client, retries, auth, pagination helpers
      client.go
      pagination.go
    config/                # Config loading and persistence
      config.go
    models/                # Request/response structs
      posts.go
      ideas.go
      ...
    output/                # JSON, table, and CSV formatting
      formatter.go
    skill/                 # Bundled agent skill asset and installer helpers
      assets/ordinal-cli/SKILL.md
      install.go
  main.go                  # Entry point
  README.md                # User-facing CLI documentation
  COVERAGE.md              # Current API endpoint coverage snapshot
  go.mod
  go.sum
```

## Tech Stack

- **Language:** Go (1.24+)
- **CLI Framework:** [Cobra](https://github.com/spf13/cobra)
- **Configuration:** [Viper](https://github.com/spf13/viper) (config files + env vars)
- **HTTP:** Standard library `net/http`
- **Output:** JSON (default), table, and CSV formats
- **Testing:** Standard library `testing` + `net/http/httptest` for API mocking

## API Resource Coverage

The CLI implements every in-scope Core REST resource group. The `agency-api/*` (multi-workspace / company-scoped) and MCP-related endpoints are intentionally out of scope.

Implemented resource groups:

| Resource | API Path Prefix | CLI Subcommand |
|---|---|---|
| Analytics | `/analytics/*` | `analytics` |
| Approvals | `/approvals`, `/posts/{id}/approvals` | `approval` |
| Comments | `/comments`, `/posts/{id}/comments` | `comment` |
| Engagements | `/engagements`, `/posts/{id}/engagements` | `engagement` |
| File Uploads | `/uploads` | `upload` |
| Ideas | `/ideas` | `idea` |
| Inline Comments | `/posts/{id}/inline-comments` | `inline-comment` |
| Instagram | `/instagram/*` | `instagram` |
| Invites | `/invites` | `invite` |
| Labels | `/labels` | `label` |
| LinkedIn | `/linkedin/*` | `linkedin` |
| LinkedIn Leads | `/linkedin/leads/*` | `linkedin-leads` |
| Posts | `/posts` | `post` |
| Profiles | `/profiles/*` | `profile` |
| Slack Boosts | `/slack-boosts`, `/posts/{id}/slack-boosts` | `slack-boost` |
| Slack Webhooks | `/slack-webhooks` | `slack-webhook` |
| Subscribers | `/subscribers`, `/posts/{id}/subscribers` | `subscriber` |
| Users | `/users` | `user` |
| Webhooks | `/webhooks` | `webhook` |
| Workspace | `/workspace` | `workspace` |

Out-of-scope resource groups:

| Resource | API Path Prefix | Reason |
|---|---|---|
| Agency API | `/agency/*` | Requires company-scoped API key; out of scope for this CLI |
| MCP | — | Handled by Ordinal's MCP integration, not this CLI |

## CLI Design Conventions

### Command Structure

Commands follow the pattern: `ordinal <resource> <action> [flags]`

```bash
ordinal auth <api-key>

ordinal post list --limit 50
ordinal post get --id <uuid>
ordinal post create --body-file ./post.json
ordinal post schedule --id <uuid> --publish-at 2026-05-01T10:00:00Z

ordinal comment create --post-id <uuid> --message "Looks good"
ordinal label create --name "Product" --color green
ordinal webhook create --name "hook" --url https://example.com --topics post.created
```

For complex operations (post/idea creation, webhook update, engagement create), direct users and agents to command help:

```bash
ordinal post create --help
ordinal idea create --help
ordinal webhook update --help
```

### Standard Actions

For resources that support CRUD, use consistent action names:

- `list` — List resources
- `get` — Get a single resource
- `create` — Create a new resource
- `update` — Update an existing resource
- `delete` — Delete a resource

Non-CRUD actions keep their API semantics, for example `archive`, `unarchive`, `schedule`, `unschedule`, `add-to-calendar`, `cpm-get`, `cpm-update`, `get-profile`, `get-mention`, `search-locations`, `list-scheduling`, `list-engagement`.

Post-scoped resources (comments, approvals, engagements, subscribers, inline-comments, slack-boosts) are scoped via `--post-id` on `list` and `create`; `delete` and `get` use the resource's own ID.

### Global Flags

- `--api-key` / `-k` — API key (also via `ORDINAL_API_KEY` or saved config)
- `--output` / `-o` — Output format: `json`, `table`, `csv`
- `--no-color` — Disable colored output
- `--verbose` / `-v` — Verbose output (shows request/response details)

### Pagination Flags

The `post list` and `idea list` commands use cursor-based pagination:

- `--limit` (1–100)
- `--cursor`
- `--all` (auto-paginate to the end)

### JSON Body Flags

Post and idea creation accept nested channel-specific content (LinkedIn, X, Instagram). Those commands accept:

- `--body-json '<inline JSON>'`
- `--body-file <path>` (or `-` to read stdin)

Individual top-level flags (`--title`, `--publish-at`, `--status`, etc.) merge into the body when provided.

Some other commands (webhook update, engagement create/update, approval create, analytics cpm-update, slack-boost update) also expose `--body-json` / `--body-file` for arbitrary partial bodies.

### Configuration

Priority order (highest to lowest):

1. CLI flags
2. Environment variables (`ORDINAL_API_KEY`, `ORDINAL_OUTPUT_FORMAT`)
3. Config file (`~/.config/ordinal/config.yaml`)

### Utility Commands

These top-level commands are part of the CLI but do not map to API resource groups:

- `auth`
- `skill`
- `completion`

## Coding Standards

- All exported types, functions, and methods must have Go doc comments.
- Error messages must be lowercase and not end with punctuation, per Go conventions.
- Use `fmt.Errorf` with `%w` for error wrapping.
- HTTP responses must always have their body closed (handled in the shared `client.do`).
- All API calls must include proper error handling for non-2xx status codes, returning the API error `{error:{message,code}}` when available.
- Rate limit errors (429) should be handled with automatic retry and backoff.
- Tests must use `httptest.NewServer` (or a roundtrip func against a custom `http.Client`) for API mocking; never call the real API in tests.
- No third-party dependencies beyond Cobra, Viper, and `go-yaml/yaml` unless discussed first.

## Build & Run

```bash
# Build
go build -o ordinal .

# Run
./ordinal --api-key <key> post list

# Test
go test ./...

# Test verbose
go test -v ./...
```
