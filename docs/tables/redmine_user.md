# Table: redmine_user

Users in the Redmine instance. Listing users requires admin privileges.

Groups and memberships are only populated when querying a single user by ID (`where id = X`).

## Examples

### Basic info

```sql
select
  id,
  login,
  first_name,
  last_name,
  mail,
  status
from
  redmine_user;
```

### List all active users

```sql
select
  id,
  login,
  first_name,
  last_name,
  mail
from
  redmine_user
where
  status = 1;
```

### Get a user by ID (includes groups and memberships)

```sql
select
  id,
  login,
  mail,
  groups,
  memberships
from
  redmine_user
where
  id = 42;
```

### Find users by email domain

```sql
select
  login,
  first_name,
  last_name,
  mail
from
  redmine_user
where
  status = 1
  and mail like '%@example.com';
```

### List admin users

```sql
select
  login,
  first_name,
  last_name,
  mail
from
  redmine_user
where
  admin = true;
```
