# Table: redmine_issue_journal

Journal entries (comments and field changes) on Redmine issues. Because the Redmine API only returns journals on the single-issue endpoint, this table performs N+1 requests under the hood.

You can get a specific journal entry by `issue_id` and `journal_id`, or list journals with a required `created_on` date range qualifier.

## Examples

### Get a specific journal entry

```sql
select
  journal_id,
  issue_id,
  user_name,
  notes,
  created_on
from
  redmine_issue_journal
where
  issue_id = 12345
  and journal_id = 67890;
```

### Basic info

```sql
select
  journal_id,
  issue_id,
  issue_subject,
  user_name,
  notes,
  created_on
from
  redmine_issue_journal
where
  created_on >= '2026-02-01'
  and created_on < '2026-03-01';
```

### List journal entries for a specific issue

```sql
select
  journal_id,
  user_name,
  notes,
  created_on
from
  redmine_issue_journal
where
  issue_id = 12345
  and created_on >= '2026-01-01';
```

### List journal entries in a specific project

```sql
select
  issue_id,
  issue_subject,
  user_name,
  notes,
  created_on
from
  redmine_issue_journal
where
  project_id = 1
  and created_on >= '2026-02-01'
  and created_on < '2026-03-01'
order by
  created_on desc;
```

### List journal entries with field changes (no comment text)

```sql
select
  journal_id,
  issue_id,
  user_name,
  details,
  created_on
from
  redmine_issue_journal
where
  created_on >= '2026-02-01'
  and created_on < '2026-03-01'
  and notes = '';
```

### Count journal entries per user in the last 7 days

```sql
select
  user_name,
  count(*) as entry_count
from
  redmine_issue_journal
where
  created_on >= current_date - interval '7 days'
group by
  user_name
order by
  entry_count desc;
```
