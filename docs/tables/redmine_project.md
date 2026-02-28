# Table: redmine_project

Projects in the Redmine instance.

## Examples

### List all active projects

```sql
select
  id,
  identifier,
  name,
  status
from
  redmine_project
where
  status = 1;
```

### Get a project by identifier

```sql
select
  id,
  name,
  description,
  is_public,
  created_on
from
  redmine_project
where
  identifier = 'my-project';
```

### List projects with their trackers

```sql
select
  name,
  identifier,
  trackers
from
  redmine_project;
```

### List sub-projects of a parent

Note: `parent_id` is filtered client-side since the Redmine API does not support it as a query parameter.

```sql
select
  id,
  name,
  identifier,
  parent_id
from
  redmine_project
where
  parent_id = 1;
```
