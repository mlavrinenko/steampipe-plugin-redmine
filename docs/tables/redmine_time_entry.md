# Table: redmine_time_entry

Time entries in the Redmine instance. Can be filtered by project, issue, user, activity, and date range.

## Examples

### Basic info

```sql
select
  id,
  project_name,
  issue_id,
  user_name,
  hours,
  comments,
  spent_on
from
  redmine_time_entry;
```

### Get a time entry by ID

```sql
select
  id,
  project_name,
  issue_id,
  user_name,
  hours,
  comments,
  spent_on
from
  redmine_time_entry
where
  id = 123;
```

### List time entries for a specific issue

```sql
select
  id,
  user_name,
  activity_name,
  hours,
  comments,
  spent_on
from
  redmine_time_entry
where
  issue_id = 12345
order by
  spent_on;
```

### List time entries for a project in a date range

```sql
select
  id,
  issue_id,
  user_name,
  activity_name,
  hours,
  comments,
  spent_on
from
  redmine_time_entry
where
  project_id = 1
  and spent_on >= '2026-02-01'
  and spent_on <= '2026-02-28'
order by
  spent_on;
```

### Total hours per user in the last 7 days

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

### Compare logged time with journal activity

```sql
select
  t.user_name,
  t.spent_on,
  sum(t.hours) as logged_hours,
  count(distinct j.journal_id) as journal_entries
from
  redmine_time_entry t
  left join redmine_issue_journal j
    on t.issue_id = j.issue_id
    and t.user_id = j.user_id
    and j.created_on >= '2026-02-01'
    and j.created_on < '2026-03-01'
where
  t.spent_on >= '2026-02-01'
  and t.spent_on <= '2026-02-28'
group by
  t.user_name, t.spent_on
order by
  t.spent_on;
```
