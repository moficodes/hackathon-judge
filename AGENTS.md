# AGENTS.md

Context for AI coding agents working in the `moficodes/hackathon-judge` repository.

## Project overview

`hackathon-judge` is a multi-component platform on Google Kubernetes Engine (GKE) that automates hackathon project code review, scoring, and structured feedback generation using AI agents:

- **`frontend/`**: React + TypeScript + Vite web UI for administrators and judges.
- **`backend/`**: Go REST API managing hackathons, projects, and evaluations via Google Cloud BigQuery and Google Cloud Pub/Sub.
- **`agent/`**: Python Google ADK orchestrator subscribing to Pub/Sub judging tasks, coordinating sandboxed evaluation runs, and publishing results.
- **`agent-sandbox/`**: Sandboxed execution runner that clones untrusted project repositories and runs the Antigravity SDK judge inside an isolated sandbox.

## Agent skills

### Issue tracker

Issues and specs for this repo live in GitHub Issues (`moficodes/hackathon-judge`) managed via the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

Uses the five canonical triage labels (`needs-triage`, `needs-info`, `ready-for-agent`, `ready-for-human`, `wontfix`). See `docs/agents/triage-labels.md`.

### Domain docs

Single-context layout (`GLOSSARY.md` at repo root and ADRs in `docs/adr/`). See `docs/agents/domain.md`.
