# Changelog

## 0.4.1 (2026-04-24)

### Bug Fixes

- Fix `assigned_to_me` filter returning zero rows. The column was declared with a constant `false` transform, so even when the API filter pushed down correctly, Postgres FDW dropped every row when evaluating `assigned_to_me = true` against the constant `false` result. The column is now computed from the actual assignee against the API key owner's user ID, so `WHERE assigned_to_me = true` returns the matching issues and the column value reflects reality in the result set.

## 0.4.0 (2026-03-24)

### Tables

- `redmine_attachment` - Attachment metadata by ID. Get-only (Redmine has no list-all-attachments endpoint).
- `redmine_custom_field` - Custom field definitions (schema discovery). Supports Get by `id`.
- `redmine_document_category` - Document category enumeration (reference table). Supports Get by `id`.
- `redmine_group` - User groups. Supports Get by `id`.
- `redmine_wiki_page` - Wiki pages per project. Supports List by `project_id` and Get by `project_id` + `title` (single-get returns full page text).

### Bug Fixes

- **devShell**: Fix shellhook to make postgres binary patching survive `nix-collect-garbage`.

## 0.3.1 (2026-03-17)

### Bug Fixes

- **redmine_issue**: Add `status_is_closed` filter to query closed issues without knowing specific status IDs.
- **redmine_issue**: Default status filter changed from `*` (all) to `open`, matching Redmine's own default behavior.
- **Makefile**: Fix install target to use correct GHCR path (`ghcr.io/mlavrinenko/steampipe-plugin-redmine@latest/`).

## 0.3.0 (2026-03-16)

### Tables

- `redmine_project_membership` - Project memberships with user/group roles. Supports List by `project_id`.

## 0.2.0 (2026-03-12)

### Tables

- `redmine_search` - Full-text search across Redmine resources via REST API.

### Bug Fixes

- Use full GHCR plugin path for install target and config.

## 0.1.2 (2026-03-12)

### Bug Fixes

- Use GHCR plugin path instead of `local/redmine` in config and docs.

### Documentation

- Add installation command to README.

## 0.1.1 (2026-03-02)

### Bug Fixes

- Switch to `golangci-lint-action@v7` in CI.

## 0.1.0 (2026-02-28)

_Initial release with 11 tables._

### Tables

- `redmine_issue` - Query issues with filters for project, tracker, status, priority, assignee, and date ranges. Supports `assigned_to_me` for filtering by API key owner.
- `redmine_issue_journal` - Journal entries (comments and field changes) on issues. Supports Get by `issue_id` + `journal_id` and List with required `created_on` date range.
- `redmine_issue_priority` - Issue priority definitions (reference table). Supports Get by `id`.
- `redmine_issue_status` - Issue status definitions (reference table). Supports Get by `id`.
- `redmine_my_account` - The currently authenticated user (API key owner). Singleton table.
- `redmine_project` - Projects with trackers, categories, enabled modules, and time entry activities. Supports Get by `id` or `identifier`.
- `redmine_time_entry` - Time entries with filters for project, issue, user, activity, and spent_on date range.
- `redmine_time_entry_activity` - Time entry activity definitions (reference table). Supports Get by `id`.
- `redmine_tracker` - Tracker definitions (reference table). Supports Get by `id`.
- `redmine_user` - Users with groups and memberships (listing requires admin privileges).
- `redmine_version` - Versions (milestones) per project. Supports Get by `id` and List by `project_id`.

### Features

- Environment variable support: `REDMINE_ENDPOINT` and `REDMINE_API_KEY`.
- Connection config validation at plugin load time.
- Retry with exponential backoff for 429 (rate limit) and 503 (unavailable) errors.
- Rate limiter (10 req/s per connection) to prevent API throttling.
- `title` standard column on all tables.
- GitHub Actions CI workflow for tests and builds.
