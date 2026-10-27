<div align="center">

# Inspire Go SDK

**Go SDK → CLI → Skill / Plugin**

<sub>The Inspire computing platform, from typed API calls to agent workflows.</sub>

<br>

![Agent ready](https://img.shields.io/badge/Agent_Ready-Skill_%2B_Plugin-7C3AED?style=for-the-badge)
![Go](https://img.shields.io/badge/Go_1.23%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Prebuilt CLI](https://img.shields.io/badge/CLI-Prebuilt-2563EB?style=for-the-badge)
![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=githubactions&logoColor=white)

[Releases](https://github.com/Lixsp11/inspire-go-sdk/releases) · [Skill](skills/inspire-go-sdk/SKILL.md) · [Plugin](plugin.json)

</div>

---

```mermaid
flowchart LR
    SDK["Go SDK<br/><sub>Typed API calls</sub>"]
    CLI["inspire CLI<br/><sub>Prebuilt binaries · JSON I/O</sub>"]
    AGENT["Skill / Plugin<br/><sub>Agent workflows · plugin installation</sub>"]
    SDK -->|compiled into| CLI
    CLI -->|called by| AGENT
    classDef sdk fill:#0c2d36,stroke:#22d3ee,color:#ecfeff
    classDef cli fill:#10264b,stroke:#60a5fa,color:#eff6ff
    classDef agent fill:#2e1a47,stroke:#c084fc,color:#faf5ff
    class SDK sdk
    class CLI cli
    class AGENT agent
```

## Install

With a Codex CLI that supports plugins:

```shell
codex plugin marketplace add Lixsp11/inspire-go-sdk --ref v0.2.0
codex plugin add inspire-go-sdk@inspire-go-sdk
```

Start a new Codex session and ask it to use **Inspire Go SDK**.

The plugin bundles the skill; the skill runs the prebuilt CLI; the CLI calls the Go SDK. The matching binary is downloaded and SHA256-verified on first use. **No Go installation is required.**

Set `INSPIRE_COOKIE` to your existing management session and `INSPIRE_PROXY` to your proxy URL when needed.

## Tech stack

| | |
| --- | --- |
| **SDK** | **Go 1.23+**, typed requests and responses, `context` |
| **CLI** | Prebuilt **`inspire`**, JSON output, profiles, pagination and job waiters |
| **Skill / Plugin** | CLI workflows for agents, bundled launchers and plugin installation |
| **Transport** | **net/http**, **encoding/json**, HTTP and SOCKS5 proxies |
| **Testing** | **testing**, **httptest**, race detector |
| **Shipping** | **GitHub Actions**, Windows / Linux / macOS, amd64 / arm64 |
