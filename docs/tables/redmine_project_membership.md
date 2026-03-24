# Table: redmine_project_membership

Project memberships in Redmine. Each row is one user or group assigned to a project with specific roles.

Listing requires either `project_id` or `project_identifier` in the WHERE clause.

## Examples

### List all members of a project by identifier

```sql
select
  id,
  user_name,
  group_name,
  roles
from
  redmine_project_membership
where
  project_identifier = 'my-project';
```

### List all members of a project by ID

```sql
select
  id,
  user_name,
  roles
from
  redmine_project_membership
where
  project_id = 42;
```

### Find members with a specific role

```sql
select
  m.user_id,
  m.user_name,
  r ->> 'name' as role_name
from
  redmine_project_membership m,
  jsonb_array_elements(m.roles) as r
where
  m.project_identifier = 'my-project'
  and (r ->> 'id')::int = 5;
```

### Get a single membership by ID

```sql
select
  id,
  project_name,
  user_name,
  roles
from
  redmine_project_membership
where
  id = 123;
```
