# Ordinal CLI

A command-line interface for the [Ordinal API](https://docs.tryordinal.com/api/introduction). Implements the full in-scope Core REST API surface: posts, ideas, approvals, comments, engagements, subscribers, labels, analytics, webhooks, Slack boosts, and related workspace resources.

## Installation

```bash
go build -o ordinal .
```

## Configuration

The CLI resolves configuration in this priority order:

1. CLI flags
2. Environment variables (`ORDINAL_API_KEY`, `ORDINAL_OUTPUT_FORMAT`)
3. Config file (`~/.config/ordinal/config.yaml`)

Save your API key:

```bash
./ordinal auth <your-api-key>
```

Or pass it directly:

```bash
./ordinal --api-key <key> post list
```

## Usage

Commands follow the pattern: `ordinal <resource> <action> [flags]`

```bash
# Authenticate once, then use the saved key
ordinal auth <your-api-key>

# Posts
ordinal post list --limit 25
ordinal post get --id <uuid>
ordinal post create --body-file ./post.json
ordinal post schedule --id <uuid> --publish-at 2026-05-01T10:00:00Z
ordinal post archive --id <uuid>

# Ideas
ordinal idea list --limit 50
ordinal idea create --title "Launch teaser" --body-file ./idea.json
ordinal idea add-to-calendar --id <uuid> --publish-date 2026-05-02

# Post-scoped resources
ordinal comment list --post-id <uuid>
ordinal comment create --post-id <uuid> --message "Looks good"
ordinal approval list --post-id <uuid>
ordinal subscriber create --post-id <uuid> --user-ids <uuid1>,<uuid2>

# Labels, webhooks, workspace
ordinal label create --name "Thought Leadership" --color purple
ordinal webhook create --name "hook" --url https://example.com/hook --topics post.created,post.published
ordinal workspace get

# Inspect exact flags for complex commands
ordinal post create --help
ordinal webhook create --help
ordinal engagement create --help
```

### Utility Commands

These ship with the CLI but do not map to API resource groups:

- `auth` — Save an API key to `~/.config/ordinal/config.yaml`
- `skill` — Install the bundled `SKILL.md` for agent tools
- `completion` — Generate shell completion scripts

### Agent Skill Installation

The CLI includes a bundled AI-agent skill for tools that support `SKILL.md`-based workflows.

```bash
ordinal skill install
ordinal skill install --target claude
ordinal skill install --target codex --force
```

By default, the installer writes `SKILL.md` into:

- `~/.agents/skills/ordinal-cli/`
- `~/.claude/skills/ordinal-cli/`
- `~/.codex/skills/ordinal-cli/`

### Global Flags

| Flag | Short | Description |
|------|-------|-------------|
| `--api-key` | `-k` | API key (also via `ORDINAL_API_KEY` env var) |
| `--output` | `-o` | Output format: `json`, `table`, `csv` |
| `--no-color` | | Disable colored output |
| `--verbose` | `-v` | Show request/response details |

### Pagination Flags

The `post list` and `idea list` commands use Ordinal's cursor pagination:

- `--limit` — Number of items per page (1–100)
- `--cursor` — Cursor for the next page (returned as `nextCursor`)
- `--all` — Auto-paginate and return all results

### Body-Driven Flags

Posts and ideas have deeply-nested channel-specific content (LinkedIn, X, Instagram). Those commands accept inline or file-based JSON bodies:

- `--body-json '<json>'` — Inline JSON
- `--body-file <path>` — Path to a JSON file, or `-` to read from stdin

Individual top-level flags (`--title`, `--publish-at`, `--status`, etc.) merge into the body when provided.

## Supported Resources

| Resource | CLI Subcommand | API Path Prefix |
|----------|----------------|-----------------|
| Analytics | `analytics` | `/analytics/*` |
| Approvals | `approval` | `/approvals`, `/posts/{id}/approvals` |
| Comments | `comment` | `/comments`, `/posts/{id}/comments` |
| Engagements | `engagement` | `/engagements`, `/posts/{id}/engagements` |
| File Uploads | `upload` | `/uploads` |
| Ideas | `idea` | `/ideas` |
| Inline Comments | `inline-comment` | `/posts/{id}/inline-comments` |
| Instagram | `instagram` | `/instagram/*` |
| Invites | `invite` | `/invites` |
| Labels | `label` | `/labels` |
| LinkedIn | `linkedin` | `/linkedin/*` |
| LinkedIn Leads | `linkedin-leads` | `/linkedin/leads/*` |
| Posts | `post` | `/posts` |
| Profiles | `profile` | `/profiles/*` |
| Slack Boosts | `slack-boost` | `/slack-boosts`, `/posts/{id}/slack-boosts` |
| Slack Webhooks | `slack-webhook` | `/slack-webhooks` |
| Subscribers | `subscriber` | `/subscribers`, `/posts/{id}/subscribers` |
| Users | `user` | `/users` |
| Webhooks | `webhook` | `/webhooks` |
| Workspace | `workspace` | `/workspace` |

### Intentionally Excluded API Groups

These Ordinal API areas are currently out of scope for this CLI:

- `agency-api/*` (multi-workspace/company-scoped endpoints — requires company API key)
- MCP-related endpoints

## Development

```bash
# Build
go build -o ordinal .

# Run tests
go test ./...

# Run tests (verbose)
go test -v ./...
```

## Tech Stack

- **Go** 1.24+
- [Cobra](https://github.com/spf13/cobra) for CLI commands
- [Viper](https://github.com/spf13/viper) for configuration
- Standard library `net/http` for API calls

## API Reference

- **Base URL:** `https://app.tryordinal.com/api/v1`
- **Auth:** Bearer token via `Authorization: Bearer <api_key>` header
- **Rate Limit:** 100 requests per minute per API key (returns `429` on overage)
- **Docs:** https://docs.tryordinal.com/api/introduction
- **OpenAPI:** https://docs.tryordinal.com/api/openapi.json

## License

MIT
