# Changelog

## 0.1.0 (Unreleased)

### Tables

- `redmine_issue` - Query issues with filters for project, tracker, status, assignee, and date ranges. Supports `assigned_to_me` for filtering by API key owner.
- `redmine_issue_journal` - Journal entries (comments and field changes) on issues. Supports Get by `issue_id` + `journal_id` and List with required `created_on` date range.
- `redmine_issue_status` - Issue status definitions (reference table).
- `redmine_project` - Projects with trackers, categories, enabled modules, and time entry activities.
- `redmine_time_entry` - Time entries with filters for project, user, activity, and spent_on date range.
- `redmine_tracker` - Tracker definitions (reference table).
- `redmine_user` - Users with groups and memberships (listing requires admin privileges).
