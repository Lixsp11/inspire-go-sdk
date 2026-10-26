---
name: inspire-go-sdk
description: Use the prebuilt Inspire CLI to query accounts, projects, compute groups, notebooks and training jobs on qz.sii.edu.cn; create, stop and monitor training jobs through the native APIs.
---

# Inspire CLI

Use the compiled `inspire` CLI. The CLI calls the Go SDK; these workflows do not require writing or compiling Go code. This skill targets CLI 0.2.x.

## Prepare the CLI

Run `inspire version` when the executable is already available. Otherwise use the bundled launcher, passing the same CLI arguments:

- Windows: `& <skill-directory>/scripts/inspire.ps1 version`
- Linux/macOS: `sh <skill-directory>/scripts/inspire.sh version`

The launcher fetches the pinned release, verifies its SHA256 checksum and caches the executable. It never retries a CLI invocation. `INSPIRE_CLI_CACHE_DIR` can select a writable cache. Resolve the actual skill directory from the installed package; do not assume a global install path. No Go compiler or Python runtime is required.

## Choose a workflow

| Task | Reference |
| --- | --- |
| Commands, profiles, credentials, proxy and errors | [CLI usage](references/cli-usage.md) |
| Account and resource queries, pagination | [Queries](references/queries.md) |
| Create, stop and monitor training jobs | [Jobs](references/jobs.md) |

Inspect `inspire <product> <Action> --help --cli-output json` before using unfamiliar parameters. JSON help includes native request fields, CLI names and nested object/array shapes. Read-only help and `--cli-dry-run-json` perform no platform calls.

Reuse the authorized session and proxy. Query existing state to resolve IDs. Perform writes only within the user's requested work. Distinguish a rejected request, an uncertain write, an accepted job and a running job. Never replay `MutationUncertain`; query platform state first.

Only use operations listed by CLI help. Discovery describes server capabilities, not additional CLI commands. This package does not provide SSH, SCP, tunnels, container execution or inference APIs.
