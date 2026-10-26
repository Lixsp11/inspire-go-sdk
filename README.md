<div align="center">

# Inspire Go SDK

**An agent-ready CLI and Go SDK for the Inspire computing platform.**

<sub>Prebuilt CLI · resource queries · training jobs · bundled skill</sub>

<br>

![Agent ready](https://img.shields.io/badge/Agent_Ready-Skill_Included-7C3AED?style=for-the-badge)
![Go](https://img.shields.io/badge/Go_1.23%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![HTTP](https://img.shields.io/badge/HTTP-JSON-2563EB?style=for-the-badge)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=githubactions&logoColor=white)

</div>

---

## Install

With a Codex CLI that supports plugins:

```shell
codex plugin marketplace add Lixsp11/inspire-go-sdk --ref v0.2.0
codex plugin add inspire-go-sdk@inspire-go-sdk
```

Start a new Codex session and ask it to use **Inspire Go SDK**.

The skill uses the matching prebuilt CLI, downloaded and checksum-verified on first use. No Go installation is required. Set `INSPIRE_COOKIE` to your existing management session and `INSPIRE_PROXY` to your proxy URL when needed.

## Tech stack

| | |
| --- | --- |
| **Agent integration** | Plugin loads the [`skill`](skills/inspire-go-sdk/SKILL.md); the skill calls the compiled CLI |
| **CLI** | `inspire <product> <Action>`, SDK parameter names, JSON output, profiles, pagination and job waiters |
| **SDK** | **Go 1.23+**, typed requests and responses, `context` |
| **Transport** | **net/http**, **encoding/json**, HTTP and SOCKS5 proxies |
| **Testing** | **testing**, **httptest**, race detector |
| **Shipping** | **GitHub Actions**, Windows / Linux / macOS, amd64 / arm64 |

[Skill](skills/inspire-go-sdk/SKILL.md) | [Releases](https://github.com/Lixsp11/inspire-go-sdk/releases)
