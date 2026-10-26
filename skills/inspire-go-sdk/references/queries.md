# Queries

## Identity

```shell
inspire user GetUserDetail
```

Use the returned `Result.id` or `Result.user_id` for owner filters; do not substitute a display name. An empty identity does not establish authentication. Report only fields actually returned.

## Resource commands

```shell
inspire project ListProjects --Filter.WorkspaceId ws-id --Page 1 --PageSize 100
inspire workspace ListLogicComputeGroups --Filter.WorkspaceId ws-id --PageNum 1 --PageSize -1
inspire train ListJobs --WorkspaceId ws-id --CreatedBy user-id --PageNum 1 --PageSize 50
inspire train GetJob --JobId job-id
inspire notebook ListNotebooks --WorkspaceId ws-id --FilterBy.UserId.1 user-id --Page 1 --PageSize 100
inspire notebook GetNotebook --NotebookId notebook-id
inspire openapi GetDiscovery
```

Replace identifiers with confirmed resources. Page numbers are one-based. Training PageSize is at most 999. Other operations retain their own SDK limits. Compute groups require an explicit PageSize; -1 selects all. Notebook queries must explicitly select users.

Projects may be queried without a workspace filter. Inspect returned workspace associations instead of inventing a ListWorkspaces API. Compute-group names do not establish free GPU capacity or quota entitlement. The current CLI does not have full quota/image/dataset discovery commands.

## Training pagination

```shell
inspire train ListJobs --WorkspaceId ws-id --CreatedBy user-id --PageSize 50 --pager --max-pages 100
```

`--pager` starts at page 1, reuses the SDK paginator and merges jobs into the native Result. It checks repeated pages/IDs, changing totals and truncation. If the page bound or another check fails, the command returns an error with no partial success JSON. Do not call an errored traversal a complete list.

Other list commands use explicit pages. Preserve the server total and disclose the visited page range; returned row count is not the total. Discovery results do not prove that unlisted commands are implemented.
