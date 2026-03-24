# Table Coverage

Tracks which Redmine REST API resources have steampipe tables.

## Implemented (18 tables)

| Table | API Resource | List | Get | Notes |
|---|---|---|---|---|
| `redmine_attachment` | `/attachments/:id` | - | by id | No list endpoint in Redmine API |
| `redmine_custom_field` | `/custom_fields` | all | by id | Admin-only endpoint; cached |
| `redmine_document_category` | `/enumerations/document_categories` | all | by id | Cached enumeration |
| `redmine_group` | `/groups` | all | by id | Cached |
| `redmine_issue` | `/issues` | filtered | by id | Server-side filtering |
| `redmine_issue_journal` | `/issues/:id` (include=journals) | by issue_id | - | Denormalized; N+1 pattern |
| `redmine_issue_priority` | `/enumerations/issue_priorities` | all | by id | Cached enumeration |
| `redmine_issue_status` | `/issue_statuses` | all | by id | Cached enumeration |
| `redmine_my_account` | `/users/current` | singleton | - | Current API key owner |
| `redmine_project` | `/projects` | all | by id | Server-side filtering |
| `redmine_project_membership` | `/projects/:id/memberships` | by project_id | by id | |
| `redmine_search` | `/search` | by query | - | Full-text search |
| `redmine_time_entry` | `/time_entries` | filtered | by id | Date range filtering |
| `redmine_time_entry_activity` | `/enumerations/time_entry_activities` | all | by id | Cached enumeration |
| `redmine_tracker` | `/trackers` | all | by id | Cached |
| `redmine_user` | `/users` | filtered | by id | Server-side filtering |
| `redmine_version` | `/projects/:id/versions` | by project_id | by id | Custom impl (not in library) |
| `redmine_wiki_page` | `/projects/:id/wiki` | by project_id | by project_id+title | Single-get returns full text |

## Not Yet Implemented

These Redmine API resources are not covered. Some may lack upstream library support.

| Resource | API Endpoint | Library Support | Priority | Notes |
|---|---|---|---|---|
| Issue categories | `/projects/:id/issue_categories` | No | Medium | Per-project issue categorization |
| Issue relations | `/issues/:id/relations` | No | Medium | Tracks issue dependencies |
| News | `/news`, `/projects/:id/news` | No | Low | Project news/announcements |
| Roles | `/roles` | No | Medium | Role definitions for RBAC |
| Queries (saved) | `/queries` | No | Low | Saved issue filters |
| Files | `/projects/:id/files` | No | Low | Project file listings |
| Documents | `/projects/:id/documents` | No | Low | Requires documents module |

Resources marked "No" for library support would need custom HTTP implementations
(like `redmine_version` and `redmine_search`).
