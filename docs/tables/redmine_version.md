# Table: redmine_version

Versions (milestones) in the Redmine instance. Listing versions requires a `project_id` qualifier.

## Examples

### List all versions for a project

```sql
select
  id,
  name,
  status,
  due_date,
  description
from
  redmine_version
where
  project_id = 1;
```

### Get a single version by ID

```sql
select
  id,
  name,
  project_name,
  status,
  due_date,
  sharing
from
  redmine_version
where
  id = 10;
```

### List open versions

```sql
select
  id,
  name,
  due_date
from
  redmine_version
where
  project_id = 1
  and status = 'open'
order by
  due_date;
```

### List issues targeting a specific version

```sql
select
  i.id,
  i.subject,
  i.status_name,
  v.name as version_name
from
  redmine_issue i
  join redmine_version v on i.fixed_version_id = v.id
where
  v.project_id = 1
  and v.name = 'v1.0';
```
