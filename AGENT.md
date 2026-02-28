# steampipe-plugin-redmine Workspace Instructions

## Project Overview

Steampipe plugin for querying Redmine instances via SQL. Primary use case: fetching
activity (journal entries) from issues to help write time-entry reports.

## Architecture

- **Go module**: `github.com/tank-io/steampipe-plugin-redmine`
- **Redmine client**: `github.com/nixys/nxs-go-redmine/v5` (chosen for type safety, header auth, `UserCurrentGet()`, named include constants)
- **Plugin SDK**: `github.com/turbot/steampipe-plugin-sdk/v5` (requires Go 1.26+)
- **Build system**: Justfile (not Makefile), Nix flake for dev environment

## Project Structure

```
steampipe-plugin-redmine/
  main.go                              # Entry point: plugin.Serve()
  redmine/
    plugin.go                          # Plugin definition, table map (11 tables)
    connection_config.go               # Config struct (endpoint, api_key)
    client.go                          # Redmine client factory with caching
    errors.go                          # Retry/ignore error predicates
    helpers.go                         # Time/date parsing utilities
    table_redmine_issue.go             # Issues table
    table_redmine_issue_journal.go     # Denormalized issue+journal view
    table_redmine_issue_priority.go    # Issue priority reference table
    table_redmine_issue_status.go      # Issue status reference table
    table_redmine_my_account.go        # Current authenticated user (singleton)
    table_redmine_project.go           # Projects table
    table_redmine_time_entry.go        # Time entries table
    table_redmine_time_entry_activity.go # Time entry activity reference table
    table_redmine_tracker.go           # Tracker reference table
    table_redmine_user.go              # Users table
    table_redmine_version.go           # Versions (milestones) table
    *_test.go                          # Unit tests
  docs/
    index.md                           # Plugin documentation
    tables/                            # Per-table documentation with examples
  config/
    redmine.spc                        # Example connection config
  flake.nix                            # Dev shell + plugin package
  flake.lock
  Justfile                             # Build/install/test/lint commands
  go.mod / go.sum
```

## Key References (in .res/)

| Resource | Path | Purpose |
|----------|------|---------|
| Steampipe SDK source | `.res/steampipe-plugin-sdk/` | Plugin SDK (v5), Go 1.26 |
| GitHub plugin example | `.res/steampipe-plugin-github/` | Reference implementation |
| nxs-go-redmine | `.res/nxs-go-redmine/` | Redmine API client (chosen) |
| go-redmine | `.res/go-redmine/` | Alternative client (not used) |
| Steampipe docs | `.res/steampipe-docs/docs/develop/` | Plugin development guides |

## Development Patterns

### Steampipe Plugin Conventions
- Table names: `redmine_{resource}` (snake_case, singular)
- Table files: `table_redmine_{resource}.go`
- Table functions: `tableRedmine{Resource}()` returning `*plugin.Table`
- Hydrate signature: `func(ctx, *plugin.QueryData, *plugin.HydrateData) (interface{}, error)`
- Use `d.StreamListItem(ctx, item)` in List, check `d.RowsRemaining(ctx) == 0`
- Use `transform.FromField("Nested.Path")` for nested structs
- Cache API client in `d.ConnectionManager.Cache`

### Redmine API Key Facts
- Journals (`include=journals`) only work on single-issue endpoint (`/issues/{id}.json`)
- Journals are NOT returned on list endpoint (`/issues.json`) -- N+1 pattern required
- Date range filtering uses `><` operator: `updated_on=%3E%3C2024-01-01|2024-12-31`
- Pagination: `offset` + `limit` (max 100), response includes `total_count`
- No built-in rate limiting -- implement client-side

### Nix / Build
- Environment: NixOS, all tooling via `flake.nix` devShell
- Go 1.26+ required (available as `go_1_26` in nixpkgs unstable)
- Use `nix-shell -p $package --run "$command"` for ad-hoc tooling
- Plugin binary: `steampipe-plugin-redmine.plugin` (`.plugin` extension required)
- Local install path: `~/.steampipe/plugins/local/redmine/`
- Nix plugin output: flat `$out/` with `.plugin` binary + `docs/` + `config/`

## Testing

- Unit test pure logic (filtering, matching, config parsing, error predicates)
- Extract filtering logic into testable functions separate from hydrate functions
- Integration: `just install` then run SQL queries against real Redmine
- Run: `just test` (wraps `go test ./...`)
