# Table: redmine_group_member

Group memberships in Redmine. Each row is one user in a group.

Listing requires `group_id` in the WHERE clause. Backed by Redmine's
`GET /groups/:id.json?include=users`, so a single query hits the API once
per group regardless of member count.

## Examples

### List all users in a group

```sql
select
  user_id,
  user_name
from
  redmine_group_member
where
  group_id = 3;
```

### Join with redmine_user to enrich members

```sql
select
  u.id,
  u.login,
  u.mail,
  u.first_name,
  u.last_name
from
  redmine_group_member gm
  join redmine_user u on u.id = gm.user_id
where
  gm.group_id = 3
  and u.status = 1
order by
  u.last_name, u.first_name;
```

### Count members per group

```sql
select
  g.id,
  g.name,
  count(gm.user_id) as member_count
from
  redmine_group g
  left join redmine_group_member gm on gm.group_id = g.id
group by
  g.id, g.name
order by
  member_count desc;
```
