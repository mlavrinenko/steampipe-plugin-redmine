# Table: redmine_my_account

The currently authenticated user (API key owner). This is a singleton table that always returns exactly one row.

## Examples

### Get current user info

```sql
select
  id,
  login,
  first_name,
  last_name,
  mail,
  admin
from
  redmine_my_account;
```

### Check if current user is admin

```sql
select
  login,
  admin
from
  redmine_my_account;
```

### List issues assigned to the current user (by ID)

```sql
select
  i.id,
  i.subject,
  i.project_name,
  i.status_name
from
  redmine_issue i
  join redmine_my_account me on i.assigned_to_id = me.id;
```

### List current user's project memberships

```sql
select
  login,
  memberships
from
  redmine_my_account;
```
