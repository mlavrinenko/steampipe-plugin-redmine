---
organization: "tank-io"
category: ["software development"]
brand_color: "#B32024"
display_name: "Redmine"
name: "redmine"
description: "Steampipe plugin for querying projects, issues, users, time entries, and journals from Redmine."
og_description: "Query Redmine with SQL! Open source CLI. No DB required."
---

# Redmine Plugin for Steampipe

[Redmine](https://www.redmine.org) is a flexible project management web application written using the Ruby on Rails framework.

Use SQL to query projects, issues, users, time entries, and journals from a Redmine instance.

## Tables

- [redmine_issue](tables/redmine_issue.md)
- [redmine_issue_journal](tables/redmine_issue_journal.md)
- [redmine_issue_status](tables/redmine_issue_status.md)
- [redmine_project](tables/redmine_project.md)
- [redmine_time_entry](tables/redmine_time_entry.md)
- [redmine_tracker](tables/redmine_tracker.md)
- [redmine_user](tables/redmine_user.md)

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

### Arguments

| Argument   | Description |
|------------|-------------|
| `endpoint` | Required. The base URL of your Redmine instance. Can also be set with the `REDMINE_ENDPOINT` environment variable. |
| `api_key`  | Required. Your Redmine API key, found under My Account -> API access key. Can also be set with the `REDMINE_API_KEY` environment variable. |

### Credentials

This plugin requires a Redmine API key. To obtain one:

1. Log in to your Redmine instance.
2. Go to **My Account** (top-right menu).
3. In the **API access key** section, click **Show** or **Reset** to get your key.
4. Ensure that **Enable REST web service** is checked under **Administration -> Settings -> API**.

### Environment Variables

| Environment Variable | Setting   |
|---------------------|-----------|
| `REDMINE_ENDPOINT`  | `endpoint` |
| `REDMINE_API_KEY`   | `api_key`  |

## Multiple Connections

You may create multiple connections to different Redmine instances:

```hcl
connection "redmine_production" {
  plugin   = "local/redmine"
  endpoint = "https://redmine.example.com"
  api_key  = "abc123..."
}

connection "redmine_staging" {
  plugin   = "local/redmine"
  endpoint = "https://staging-redmine.example.com"
  api_key  = "def456..."
}
```

You can also create an [aggregator connection](https://steampipe.io/docs/managing/connections#using-aggregators) to query across all Redmine instances:

```hcl
connection "redmine_all" {
  plugin      = "local/redmine"
  type        = "aggregator"
  connections = ["redmine_production", "redmine_staging"]
}
```
