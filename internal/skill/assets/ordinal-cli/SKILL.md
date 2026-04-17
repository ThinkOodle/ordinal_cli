---
name: ordinal-cli
description: |
  Use this skill when an agent needs to operate the Ordinal (tryordinal.com)
  API through the `ordinal` CLI. Covers authentication, command discovery,
  output formats, pagination, JSON-heavy flags, supported resources, and safe
  command patterns for posts, ideas, approvals, comments, engagements,
  subscribers, labels, analytics, webhooks, Slack boosts, and related
  Ordinal workspace operations.
invocable: true
argument-hint: "[resource] [action] [flags]"
---

# Ordinal CLI

Use this skill to work with the installed `ordinal` program.

This is an operator skill, not a repository-maintenance skill. The goal is to help an agent use the CLI correctly to work with the Ordinal API.

## What This Skill Is For

Use it when you need to:

1. Read or modify Ordinal workspace data through the CLI.
2. Translate an Ordinal API task into the right `ordinal <resource> <action>` command.
3. Discover the exact flags or actions supported by the installed CLI.

Do not assume the CLI exposes every API endpoint. Always confirm the available commands from the binary itself.

## Ground Truth

Use this order when deciding how to act:

1. `ordinal --help`
2. `ordinal <resource> --help`
3. `ordinal <resource> <action> --help`
4. Official docs: `https://docs.tryordinal.com/api/introduction`

If the CLI and the API docs differ, the CLI help is the source of truth for what the installed binary can do.

## Current API Facts

- Base URL: `https://app.tryordinal.com/api/v1`
- Auth: Bearer token via `Authorization: Bearer <api_key>` header
- Rate limit: `100 requests per minute` per API key (`429 Too Many Requests` when exceeded)
- API access requires the Pro plan
- API keys are generated from workspace settings; each key is scoped to a single workspace

Useful docs:

- Introduction: `https://docs.tryordinal.com/api/introduction`
- Authentication: `https://docs.tryordinal.com/api/authentication`
- Errors: `https://docs.tryordinal.com/api/errors`
- File uploads: `https://docs.tryordinal.com/api/file-uploads`
- LinkedIn mentions: `https://docs.tryordinal.com/api/linkedin-mentions`
- OpenAPI spec: `https://docs.tryordinal.com/api/openapi.json`

## First Steps

Before running a mutating command, inspect help:

```bash
ordinal --help
ordinal post --help
ordinal post create --help
```

Use help-driven discovery instead of guessing action names or flag names.

## Authentication

The CLI supports three common ways to supply the API key:

1. `--api-key`
2. `ORDINAL_API_KEY`
3. Saved config via `ordinal auth <api-key>`

Examples:

```bash
ordinal auth <api-key>
ordinal post list
ordinal --api-key <api-key> post list
ORDINAL_API_KEY=<api-key> ordinal label list
```

If the CLI reports that an API key is required, either pass `--api-key` explicitly or save it with `ordinal auth`.

## Global Flags

Common global flags:

- `--api-key`, `-k`
- `--output`, `-o`: `json`, `table`, `csv`
- `--no-color`
- `--verbose`, `-v`

Default operating posture:

- Use `--output json` for agent workflows, parsing, or follow-up automation.
- Use `--output table` for quick human inspection.
- Use `--verbose` when debugging request/response failures.

## Core Command Pattern

Most commands follow this pattern:

```bash
ordinal <resource> list
ordinal <resource> get --id <uuid>
ordinal <resource> create ...
ordinal <resource> update --id <uuid> ...
ordinal <resource> delete --id <uuid>
```

Notable exceptions:

- `post` has extra actions: `schedule`, `unschedule`, `archive`, `unarchive`.
- `idea` has extra actions: `archive`, `unarchive`, `add-to-calendar`.
- Some resources (comments, approvals, engagements, subscribers, inline-comments, slack-boost) are scoped by a parent post ID via `--post-id`.
- `profile` has two list actions: `list-scheduling` and `list-engagement`.
- `label` supports `create`, `list`, and `delete` only (no update).
- Mutating success may return a JSON object or a plain confirmation line.

## Supported Top-Level Resources

The installed CLI currently exposes these top-level commands:

- `analytics`
- `approval`
- `comment`
- `engagement`
- `idea`
- `inline-comment`
- `instagram`
- `invite`
- `label`
- `linkedin`
- `linkedin-leads`
- `post`
- `profile`
- `slack-boost`
- `slack-webhook`
- `subscriber`
- `upload`
- `user`
- `webhook`
- `workspace`

If a needed resource is missing from `ordinal --help`, the binary does not currently support it.

## Pagination

Ordinal uses cursor-based pagination on list endpoints that support it (`post list`, `idea list`). The cursor is an item ID returned as `nextCursor` on the previous page.

List-style commands commonly support:

- `--limit` (1-100)
- `--cursor`
- `--all` (auto-paginate)

Guidance:

- Start with a bounded `--limit` unless you truly need the full dataset.
- Use `--all` only when you need every page.
- Preserve cursors if you are building a multi-step workflow.

Examples:

```bash
ordinal post list --limit 25
ordinal post list --cursor <post-id>
ordinal idea list --all --output json
```

## JSON-Heavy Flags

Post and idea creation support multi-channel content (LinkedIn, X, Instagram). These use nested objects that are hard to express purely as flags, so those commands accept `--body-json` (inline JSON) or `--body-file` (path to a JSON file, `-` for stdin) for the channel configs.

Common JSON-driven patterns:

- `post create --body-file ./post.json`
- `post update --id <uuid> --body-json '<json>'`
- `idea create --body-file ./idea.json`
- `idea update --id <uuid> --body-json '<json>'`
- `engagement create --post-id <uuid> --channel linkedin --body-file ./engagements.json`
- `approval create --post-id <uuid> --body-json '<approvals-array-json>'`

Do not invent body shapes. Use the official docs and command help to confirm the expected structure. In particular:

- Post / idea channel configs: see `POST /posts` and `POST /ideas` in the OpenAPI spec.
- Engagement inputs: see `EngagementInput` schema.

Use single quotes around JSON on Unix-like shells:

```bash
ordinal label create --name "Thought Leadership" --color purple
ordinal subscriber create --post-id <uuid> --user-ids <uuid1>,<uuid2>
ordinal post create --body-file ./post.json
```

## Safe Discovery Workflow

When you know the business task but not the exact command:

1. Find the resource from `ordinal --help`
2. Inspect the resource with `ordinal <resource> --help`
3. Inspect the action with `ordinal <resource> <action> --help`
4. Only then build the command

## Common Workflows

### Posts

```bash
ordinal post list --limit 25 --output json
ordinal post get --id <post-uuid>
ordinal post create --body-file ./post.json
ordinal post update --id <post-uuid> --body-json '{"title":"New title"}'
ordinal post schedule --id <post-uuid> --publish-at 2026-05-01T10:00:00Z
ordinal post unschedule --id <post-uuid>
ordinal post archive --id <post-uuid>
ordinal post unarchive --id <post-uuid>
```

Notes:

- Posts require at least one channel (`linkedIn`, `x`, or `instagram`) in the body.
- `publishAt` is UTC ISO 8601.

### Ideas

```bash
ordinal idea list --limit 50 --output json
ordinal idea get --id <idea-uuid>
ordinal idea create --title "Launch teaser" --body-file ./idea.json
ordinal idea update --id <idea-uuid> --body-json '{"title":"Renamed"}'
ordinal idea archive --id <idea-uuid>
ordinal idea unarchive --id <idea-uuid>
ordinal idea add-to-calendar --id <idea-uuid> --publish-at 2026-05-02T09:00:00Z
```

### Approvals, Comments, Engagements, Subscribers, Inline Comments, Slack Boosts

All are scoped by a parent post for listing/creating:

```bash
ordinal approval list --post-id <post-uuid>
ordinal approval create --post-id <post-uuid> --body-file ./approvals.json
ordinal approval delete --id <approval-uuid>

ordinal comment list --post-id <post-uuid>
ordinal comment create --post-id <post-uuid> --message "Looks good"
ordinal comment delete --id <comment-uuid>

ordinal engagement list --post-id <post-uuid>
ordinal engagement create --post-id <post-uuid> --channel linkedin --body-file ./engagements.json
ordinal engagement update --id <engagement-uuid> --body-json '{"copy":"Nice"}'
ordinal engagement delete --id <engagement-uuid>

ordinal subscriber list --post-id <post-uuid>
ordinal subscriber create --post-id <post-uuid> --user-ids <uuid1>,<uuid2>
ordinal subscriber delete --id <subscriber-uuid>

ordinal inline-comment list --post-id <post-uuid>

ordinal slack-boost list --post-id <post-uuid>
ordinal slack-boost create --post-id <post-uuid> --slack-webhook-id <uuid>
ordinal slack-boost get --id <boost-uuid>
ordinal slack-boost update --id <boost-uuid> --copy "Please boost"
ordinal slack-boost delete --id <boost-uuid>
ordinal slack-webhook list
```

### Labels, Uploads, Profiles

```bash
ordinal label list
ordinal label create --name "Product" --color green
ordinal label delete --id <label-uuid>

ordinal upload create --url https://example.com/image.png
ordinal upload get --id <upload-uuid>

ordinal profile list-scheduling
ordinal profile list-engagement
```

### Users, Invites, Workspace

```bash
ordinal user list
ordinal invite list
ordinal invite create --email someone@example.com
ordinal invite delete --id <invite-uuid>
ordinal workspace get
```

### Webhooks

```bash
ordinal webhook list
ordinal webhook get --id <webhook-uuid>
ordinal webhook create --name "My hook" --url https://example.com/hook --topics post.created,post.published
ordinal webhook update --id <webhook-uuid> --body-json '{"topics":["post.archived"]}'
ordinal webhook delete --id <webhook-uuid>
```

### Analytics

```bash
ordinal analytics cpm-get
ordinal analytics cpm-update --body-json '<cpm-json>'
ordinal analytics linkedin-followers --profile-id <uuid>
ordinal analytics linkedin-posts --profile-id <uuid> --start 2026-01-01 --end 2026-02-01
ordinal analytics x-followers --profile-id <uuid>
ordinal analytics x-posts --profile-id <uuid> --start 2026-01-01 --end 2026-02-01
```

### LinkedIn / Instagram

```bash
ordinal linkedin get-profile --urn <urn>
ordinal linkedin get-mention --username <username>
ordinal linkedin-leads list-posts --profile-id <uuid>
ordinal linkedin-leads get-leads --profile-id <uuid> --post-id <linkedin-post-id>
ordinal instagram search-locations --query "Brooklyn"
```

## Output Expectations

Expect three broad response styles:

1. JSON objects or arrays, especially with `--output json`
2. Human-readable tables with `--output table`
3. Plain success lines for certain mutate/toggle/delete commands

If a workflow needs machine-readability, force JSON output.

## Failure Handling

If a command fails:

1. Re-run with `--verbose`
2. Confirm auth and workspace scope
3. Confirm the resource/action exists in CLI help
4. Confirm flag names and JSON shape against the OpenAPI spec
5. Respect the 100 req/min rate limit; back off on `429`

Typical failure sources:

- Missing or wrong API key
- API key is for a different workspace than the resource
- Invalid resource IDs (UUIDs must match existing records)
- Bad JSON in `--body-json` / `--body-file`
- Using a feature the installed CLI does not expose
- Rate limiting (100 requests/minute per key)

## Agent Operating Rules

1. Prefer `--output json` unless the user clearly wants tables.
2. Inspect help before using unfamiliar actions.
3. Many list endpoints (comments, approvals, subscribers, engagements, inline-comments, slack-boosts) require `--post-id`.
4. Do not invent JSON shapes for complex flags; confirm against the OpenAPI spec.
5. Use `--all` carefully; list operations count toward the 100/min rate limit.
6. When a command is absent from help, say the installed CLI does not currently support it.

## Minimal Workflow Template

For most tasks, use this sequence:

```bash
ordinal --help
ordinal <resource> --help
ordinal <resource> <action> --help
ordinal <resource> <action> ... --output json
```

That is the default safe path for agents using the `ordinal` CLI.
