# Table: redmine_tracker

Trackers defined in the Redmine instance (e.g., Bug, Feature, Support). This is a reference table useful for joining with issues.

## Examples

### List all trackers

```sql
select
  id,
  name,
  default_status_name,
  description
from
  redmine_tracker;
```

### List trackers with their enabled fields

```sql
select
  name,
  enabled_standard_fields
from
  redmine_tracker;
```

### Count issues per tracker

```sql
select
  t.name as tracker,
  count(i.id) as issue_count
from
  redmine_tracker t
  left join redmine_issue i on t.id = i.tracker_id
group by
  t.name
order by
  issue_count desc;
```
