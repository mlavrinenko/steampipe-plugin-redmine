# steampipe-plugin-redmine

A [Steampipe](https://steampipe.io) plugin to query a [Redmine](https://www.redmine.org) instance via SQL.

## Tables

| Table | Description |
|-------|-------------|
| [redmine_issue](docs/tables/redmine_issue.md) | Issues with filtering by project, tracker, status, priority, assignee, and date ranges. |
| [redmine_issue_journal](docs/tables/redmine_issue_journal.md) | Journal entries (comments and field changes) on issues. |
| [redmine_issue_priority](docs/tables/redmine_issue_priority.md) | Issue priority definitions (Low, Normal, High, Urgent, Immediate). |
| [redmine_issue_status](docs/tables/redmine_issue_status.md) | Issue status definitions (New, In Progress, Closed, etc.). |
| [redmine_my_account](docs/tables/redmine_my_account.md) | The currently authenticated user (API key owner). |
| [redmine_project](docs/tables/redmine_project.md) | Projects with trackers, categories, and modules. |
| [redmine_time_entry](docs/tables/redmine_time_entry.md) | Time entries with filtering by project, issue, user, activity, and date. |
| [redmine_tracker](docs/tables/redmine_tracker.md) | Tracker definitions (Bug, Feature, Support, etc.). |
| [redmine_user](docs/tables/redmine_user.md) | Users (listing requires admin privileges). |
| [redmine_version](docs/tables/redmine_version.md) | Versions (milestones) per project. |

## Query Examples

List all journal entries for a date range:

```sql
select
  issue_id,
  issue_subject,
  project_name,
  user_name,
  notes,
  created_on
from
  redmine_issue_journal
where
  created_on >= '2026-02-01'
  and created_on < '2026-03-01'
order by
  created_on;
```

List issues assigned to me:

```sql
select
  id,
  subject,
  project_name,
  status_name,
  priority_name
from
  redmine_issue
where
  assigned_to_me = true;
```

Total hours per user in the last 7 days:

```sql
select
  user_name,
  sum(hours) as total_hours
from
  redmine_time_entry
where
  spent_on >= current_date - interval '7 days'
group by
  user_name
order by
  total_hours desc;
```

Quick CLI check:

```bash
steampipe query "select issue_id, issue_subject, created_on from redmine_issue_journal where created_on >= current_date - interval '7 days' limit 5"
```

## Configuration

Connection configuration is defined in a `.spc` file:

```hcl
connection "redmine" {
  plugin = "local/redmine"

  # Redmine instance URL (required).
  # Can also be set with the REDMINE_ENDPOINT environment variable.
  endpoint = "https://www.redmine.org"

  # API key from My Account -> API access key (required).
  # Can also be set with the REDMINE_API_KEY environment variable.
  api_key = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
}
```

## Getting Started

```bash
nix develop
just install
cp config/redmine.spc ~/.steampipe/config/redmine.spc
vi ~/.steampipe/config/redmine.spc  # configure endpoint and api_key
```

## Development

```bash
nix build       # builds the plugin as a Nix package
nix develop     # drops you into a dev shell with all tools
just test       # run unit tests
just lint       # run golangci-lint
just build      # build the plugin binary
just install    # build and install to ~/.steampipe/plugins/
```
