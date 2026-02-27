# steampipe-plugin-redmine

A [Steampipe](https://steampipe.io) plugin to query a [Redmine](https://www.redmine.org) instance.

# Query Example

```sql
select issue_id, issue_subject, project_name, user_name, notes, created_on
from redmine_issue_journal
where created_on >= '2026-02-01' and created_on < '2026-03-01'
order by created_on;
```

or in CLI:

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
