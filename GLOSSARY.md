# Hackathon Judge

Automated code-review, testing, and scoring platform that evaluates hackathon project submissions inside isolated Agent Substrate actors on Google Kubernetes Engine.

## Language

**Judging Task**:
An asynchronous request to evaluate a submitted hackathon project repository against a hackathon's rubric and scoring criteria.
_Avoid_: Job, ticket, evaluation request

**Judge Orchestrator**:
The always-on worker service (`agent`) that consumes Judging Tasks, manages Judge Actor lifecycles, proxies authenticated model requests, and publishes Evaluation Results.
_Avoid_: Supervisor, controller, master agent

**Judge Actor**:
An ephemeral, isolated Agent Substrate actor booted from a Golden Snapshot to clone, build, test, and audit a single untrusted project repository.
_Avoid_: Sandbox pod, sandbox claim, container instance

**Golden Snapshot**:
A pre-initialized memory and filesystem checkpoint of the Judge Actor template from which every new Judge Actor restores prior to receiving a repository.
_Avoid_: Warm pool pod, base image, standby replica

**Worker Pool**:
The set of gVisor-enabled Kubernetes pods (`judge-workerpool`) that provide compute seats for running Judge Actors.
_Avoid_: Sandbox warm pool, node pool

**Egress Policy**:
The per-actor network allowlist governing which internal cluster services and external internet hosts a Judge Actor may contact.
_Avoid_: Sandbox network policy, firewall rule

**Model Proxy**:
The boundary endpoint on the Judge Orchestrator that mediates LLM calls on behalf of a credential-free Judge Actor.
_Avoid_: Metadata server, direct Vertex client, token broker

**Evaluation Result**:
The structured scorecard produced by a Judge Actor containing per-criterion scores, code-cited reasoning, total score, overall comments, and confidence score.
_Avoid_: Review, grade, audit log
