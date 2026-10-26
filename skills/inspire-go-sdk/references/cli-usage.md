# CLI usage

## Commands and parameters

```shell
inspire --help --cli-output json
inspire train CreateJobConsole --help --cli-output json
inspire user GetUserDetail
```

Commands use `inspire <product> <Action>`. API flags retain SDK field names, such as `--WorkspaceId` and `--JobId`. Nested fields use dots and one-based array indexes: `--Filter.WorkspaceId ws-id`, `--FilterBy.UserId.1 user-id`. Complex parameters also accept JSON.

For complete requests use `--body-file request.json`, or `--body-file=-` with JSON on stdin. The equals form also works through the Windows PowerShell launcher. Body keys may use native JSON names or SDK field names. Do not mix a body with API parameter flags. Unknown fields, duplicate members, overlapping paths and invalid types are rejected locally.

## Credentials and profiles

The CLI uses existing management credentials from `INSPIRE_COOKIE` or `INSPIRE_ACCESS_TOKEN`; choose one. A cookie is the full `inspire-session=...` value. An inference key is not a management credential. Reuse configured credentials without printing them or writing them into source files.

```shell
inspire configure set --profile research --cookie-file session.txt --socks5-proxy socks5://127.0.0.1:1080
inspire configure switch --profile research
inspire configure list
```

The profile file defaults to `~/.inspire/config.json`; `--config-path FILE` or `INSPIRE_CONFIG_FILE` selects another file. Profiles store a cookie file path, not the cookie. Profile views redact proxy authentication. `--profile` or `INSPIRE_PROFILE` selects a named profile. A configured cookie file selects that credential source over credential environment variables.

Connection flags override environment values, which override profile connection settings. `INSPIRE_ENDPOINT`, `INSPIRE_PROXY` and `INSPIRE_COOKIE_FILE` provide environment overrides. When a proxy is required, keep using it; the CLI does not switch to direct access after failure. SDK proxy selection and TLS validation remain active.

`--read-timeout`, `--connect-timeout` and `--timeout` use seconds; the CLI converts SDK timeouts to milliseconds. `--timeout` covers the complete command (default 600 seconds). `--retry-count` adds bounded attempts to read-only requests; writes are never retried. Authentication rejection stops retries.

## Output and errors

Success writes one JSON value to stdout, preserving `ResponseMetadata` and `Result`. Treat actual returned fields as authoritative. Errors leave stdout empty, write a JSON diagnostic to stderr and return a nonzero exit status: 2 for invalid local input, 1 for execution/API failures.

Diagnostics use `schemaVersion`, `message`, `error_code`, optional `status_code`/`request_id`, and `recovery.action`/`recovery.hint`. `recovery.action=reconcile` means a write may have been accepted. Query before resubmitting; a matching name alone does not prove identity. Do not print raw request bodies or credential values in summaries.

The launchers only fetch and run released executables. A `BootstrapError` happens before the CLI starts; inspect release availability, GitHub access, cache integrity and platform support. `LauncherError` with `recovery.action=reconcile` means execution or output forwarding failed after startup; check platform state before repeating a write. Do not replace the workflow with ad-hoc Go code.
