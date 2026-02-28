# Redmine Plugin for Steampipe

Use SQL to query issues, journals, and other resources from a [Redmine](https://www.redmine.org) instance.

## Configuration

Connection configuration is defined in a `.spc` file:

```hcl
connection "redmine" {
  plugin = "local/redmine"

  # Redmine instance URL (required).
  endpoint = "https://www.redmine.org"

  # API key from My Account -> API access key (required).
  api_key = "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2"
}
```

- `endpoint` (required) - The base URL of your Redmine instance.
- `api_key` (required) - Your Redmine API key, found under My Account -> API access key.

## Get Started

Install the plugin locally:

```bash
just install
```

Copy and edit the configuration file:

```bash
cp config/redmine.spc ~/.steampipe/config/redmine.spc
vi ~/.steampipe/config/redmine.spc
```

Run a query:

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
