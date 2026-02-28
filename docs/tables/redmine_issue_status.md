# Table: redmine_issue_status

Issue statuses defined in the Redmine instance. This is a reference table useful for joining with issues.

## Examples

### List all statuses

```sql
select
  id,
  name,
  is_closed
from
  redmine_issue_status;
```

### List only closed statuses

```sql
select
  id,
  name
from
  redmine_issue_status
where
  is_closed = true;
```

### Count issues per status

```sql
select
  s.name as status,
  s.is_closed,
  count(i.id) as issue_count
from
  redmine_issue_status s
  left join redmine_issue i on s.id = i.status_id
group by
  s.name, s.is_closed
order by
  issue_count desc;
```
