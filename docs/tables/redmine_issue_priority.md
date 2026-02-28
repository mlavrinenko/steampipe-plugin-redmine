# Table: redmine_issue_priority

Issue priority definitions in the Redmine instance. This is a reference table for looking up priority names by ID.

## Examples

### List all priorities

```sql
select
  id,
  name,
  is_default,
  active
from
  redmine_issue_priority;
```

### Get the default priority

```sql
select
  id,
  name
from
  redmine_issue_priority
where
  is_default = true;
```

### List issues with their priority names (via join)

```sql
select
  i.id,
  i.subject,
  p.name as priority
from
  redmine_issue i
  join redmine_issue_priority p on i.priority_id = p.id
where
  i.status_is_closed = false
order by
  p.id desc;
```

### Get a priority by ID

```sql
select
  id,
  name,
  active
from
  redmine_issue_priority
where
  id = 4;
```
