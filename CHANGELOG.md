# Changelog

## 0.1.0 (2026-02-28)

_Initial release with 7 tables._

### Tables

- `redmine_issue` - Query issues with filters for project, tracker, status, priority, assignee, and date ranges. Supports `assigned_to_me` for filtering by API key owner.
- `redmine_issue_journal` - Journal entries (comments and field changes) on issues. Supports Get by `issue_id` + `journal_id` and List with required `created_on` date range.
- `redmine_issue_status` - Issue status definitions (reference table). Supports Get by `id`.
- `redmine_project` - Projects with trackers, categories, enabled modules, and time entry activities. Supports Get by `id` or `identifier`.
- `redmine_time_entry` - Time entries with filters for project, issue, user, activity, and spent_on date range.
- `redmine_tracker` - Tracker definitions (reference table). Supports Get by `id`.
- `redmine_user` - Users with groups and memberships (listing requires admin privileges).

### Features

- Environment variable support: `REDMINE_ENDPOINT` and `REDMINE_API_KEY`.
- Connection config validation at plugin load time.
- Retry with exponential backoff for 429 (rate limit) and 503 (unavailable) errors.
- Rate limiter (10 req/s per connection) to prevent API throttling.
- `title` standard column on issue, project, and user tables.
