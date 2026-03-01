# Table: redmine_issue

Issues in the Redmine instance. By default, both open and closed issues are returned.

**Note:** The `assigned_to_me` column is a filter-only column. It can be used in `WHERE` clauses to filter issues assigned to the API key owner, but it always returns `false` in query results. To see the actual assignee, use `assigned_to_id` and `assigned_to_name`.

## Examples

### Basic info

```sql
select
  id,
  subject,
  project_name,
  status_name,
  priority_name,
  assigned_to_name,
  updated_on
from
  redmine_issue;
```

### List open issues in a project

```sql
select
  id,
  subject,
  status_name,
  assigned_to_name,
  updated_on
from
  redmine_issue
where
  project_id = 1
  and status_is_closed = false
order by
  updated_on desc;
```

### Get a single issue by ID

```sql
select
  id,
  subject,
  project_name,
  status_name,
  priority_name,
  assigned_to_name,
  description,
  created_on,
  updated_on
from
  redmine_issue
where
  id = 12345;
```

### List issues assigned to me (the API key owner)

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

### List overdue issues

```sql
select
  id,
  subject,
  project_name,
  due_date,
  assigned_to_name
from
  redmine_issue
where
  status_is_closed = false
  and due_date < now();
```

### List high-priority issues

```sql
select
  id,
  subject,
  project_name,
  status_name,
  assigned_to_name
from
  redmine_issue
where
  priority_id = 4;
```

### Join issues with journals to find activity by user email

```sql
select
  j.issue_id,
  j.issue_subject,
  j.notes,
  j.created_on,
  u.mail
from
  redmine_issue_journal j
  join redmine_user u on j.user_id = u.id
where
  j.created_on >= '2026-02-01'
  and j.created_on < '2026-03-01'
  and u.mail = 'alice@example.com'
order by
  j.created_on;
```
