# Table: redmine_time_entry_activity

Time entry activity definitions in the Redmine instance (e.g., Development, Design, Support). This is a reference table useful for joining with time entries.

## Examples

### List all activities

```sql
select
  id,
  name,
  is_default,
  active
from
  redmine_time_entry_activity;
```

### Get the default activity

```sql
select
  id,
  name
from
  redmine_time_entry_activity
where
  is_default = true;
```

### Get an activity by ID

```sql
select
  id,
  name,
  active
from
  redmine_time_entry_activity
where
  id = 9;
```

### Total hours per activity in a project

```sql
select
  a.name as activity,
  sum(te.hours) as total_hours
from
  redmine_time_entry te
  join redmine_time_entry_activity a on te.activity_id = a.id
where
  te.project_id = 1
group by
  a.name
order by
  total_hours desc;
```
