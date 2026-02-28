# steampipe-plugin-redmine

A [Steampipe](https://steampipe.io) plugin to query a [Redmine](https://www.redmine.org) instance.

# Query Examples

List all journal entries for a date range:

```sql
select issue_id, issue_subject, project_name, user_name, notes, created_on
from redmine_issue_journal
where created_on >= '2026-02-01' and created_on < '2026-03-01'
order by created_on;
```

Filter journal entries from issues where you participated, by email:

```sql
select issue_id, issue_subject, project_name, user_name, notes, created_on
from (
  select j.*,
    bool_or(u.mail = 'me@example.com') over (partition by j.issue_id) as i_participated
  from redmine_issue_journal j
  join redmine_user u on j.user_id = u.id
  where j.created_on >= '2026-02-01' and j.created_on < '2026-03-01'
) sub
where i_participated
order by created_on;
```

Quick CLI check:

```bash
steampipe query "select issue_id, issue_subject, created_on from redmine_issue_journal where created_on >= current_date - interval '7 days' limit 5"
```

## In-Development Usage

```bash
nix develop
just install
nano ~/.steampipe/config/redmine.spc # configure Redmine endpoint to use
```

## Development

```bash
nix build       # builds the plugin as a Nix package
nix develop     # drops you into a dev shell with all tools
just test       # run unit tests
```

## TODO

- Add `Get` config to `redmine_issue_journal` for efficient single-row lookup by `issue_id` + `journal_id`.
- Add standard columns (`title`, `akas`, `tags`) — debatable for this denormalized journal view, but recommended by Steampipe standards for resource tables.
- Add `MaxConcurrency` control for the N+1 `IssueSingleGet` calls (requires refactoring to a proper hydrate function or using SDK rate limiters).
