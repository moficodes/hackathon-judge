---
id: agent-substrate-code-evaluation
summary: Deploy an Automatic Code Evaluator on GKE with Agent Substrate and Antigravity SDK.
authors: Mofi Rahman
categories: GKE, Antigravity SDK, Agent Substrate, BigQuery
tags: GKE, Antigravity SDK, Agent Substrate
keywords: docType:Codelab,skill:Intermediate,language:Python,language:GoLang,category:Cloud,category:AiAndMachineLearning,product:GoogleKubernetesEngine,product:BigQuery,product:AgentPlatform,event:GoogleCloudNext2026
feedback_link: https://github.com/googlecodelabs/feedback/issues/new?title=[agent-substrate-code-evaluation]:&labels[]=content-platform&labels[]=cloud
analytics_account: UA-66226300-1
layout: paginated

---

# Automatic Code Evaluation with Agent Substrate and Antigravity SDK on GKE

## Introduction
Duration: 02:00

In this codelab, you will deploy the Hackathon Judge application on Google
Kubernetes Engine (GKE) and use **Agent Substrate** (`v1.0.0`) and the **Google
Antigravity SDK** (`google-antigravity`) to run agentic code evaluation
workloads safely and at sub-second startup speeds.

The platform automates the process of reviewing, testing, and grading hackathon
projects using AI agents. Since judging requires cloning, compiling, and testing
untrusted participant code submissions, each evaluation runs inside an
ephemeral, zero-credential **Judge Actor** restored from a memory-fidelity
**Golden Snapshot** inside a gVisor sandbox.

### What you'll do

- Provision the target Google Cloud services and enable required APIs.
- Initialize a GKE Standard cluster with nested virtualization (`--enable-nested-virtualization`), and install **Agent Substrate** (`v1.0.0`) and `kubectl-ate`.
- Configure the Agent Substrate **WorkerPool** (`judge-workerpool`), **AccessPolicy**, and **ActorTemplate** (`judge-sandbox`) to create a pre-initialized **Golden Snapshot**.
- Deploy the REST Backend API, the Python **Judge Orchestrator** (`agent`) with its internal **Model Proxy**, and the React Frontend UI.
- Wire external load-balanced routing and watch ephemeral **Judge Actors** and per-actor **EgressPolicy** rules run autonomous **Antigravity SDK** evaluations through `atenet-router`.

### What you'll need

- A web browser such as [Chrome](https://www.google.com/chrome/).
- A Google Cloud project with billing enabled.

The resources created in this codelab should cost less than $5 in total runtime charges.

## The Problem: Securely Evaluating Untrusted Code
Duration: 02:00

Hackathons are fast-paced innovation events where participants build and submit projects—often including source code—for evaluation. Manually judging these submissions is time-consuming and resource-intensive. Using AI agents to automate grading is a promising solution, but it introduces a significant security challenge: how do you safely run participant-provided code that could be buggy, malicious, or resource-intensive?

Running untrusted code directly on your infrastructure exposes you to risks
like:

- **Code Injection**: Malicious scripts could try to access or modify sensitive data.
- **Privilege Escalation**: Code might attempt to gain unauthorized access to other systems or network resources.
- **Resource Abuse**: Poorly written code or denial-of-service attacks could consume excessive CPU, memory, or network bandwidth, impacting other operations.
- **Credential Exfiltration**: Untrusted build scripts or prompt-injected agent tool calls could attempt to read ambient cloud credentials or query the GKE metadata server to abuse cloud APIs.

To automate hackathon judging with AI, you need a way to execute submitted code
in an environment that holds zero cloud credentials and is completely isolated
from the rest of your system and other submissions.

## The Solution: Agent Substrate and Antigravity SDK on GKE
Duration: 02:00

**Agent Substrate** on GKE pairs kernel-level sandbox isolation and memory-fidelity snapshotting with the **Google Antigravity SDK** (`google-antigravity`) to solve this exact challenge. It manages isolated, stateful, single-submission **Judge Actors** on GKE, purpose-built for AI agent runtimes that execute untrusted code securely and efficiently.

Key benefits of this architecture include:

* **Kernel-level gVisor isolation via `WorkerPool`**: Runs Judge Actors inside `ateworker-gvisor` worker pods backed by [gVisor](https://gvisor.dev/), intercepting system calls in user space so untrusted code cannot access the host kernel or sibling workloads.
* **Sub-second Golden Snapshot restore via `ActorTemplate`**: Pre-boots the sandbox runtime once, verifies its `/health` wakeup probe, and captures a memory-fidelity **Golden Snapshot** (`SNAPSHOT_FIDELITY_MEMORY`) in Cloud Storage. Every evaluation restores a fresh ephemeral Judge Actor from this snapshot in ~1–2 seconds and deletes it immediately after grading.
* **Zero-credential Judge Actors with `EgressPolicy` & Model Proxy**: Judge Actors run with zero Google Cloud credentials and no access to the GKE metadata server. A per-actor `EgressPolicy` permits only public HTTPS (`:443`) to clone repositories and install packages, plus internal HTTP to the Judge Orchestrator's `/vertex-proxy` endpoint, which enforces a strict API path allowlist and injects its own Workload Identity token on the wire.
* **Autonomous code evaluation with the Antigravity SDK**: Inside each Judge Actor, `app/judge.py` uses `google-antigravity` to inspect the repository, install dependencies, run test suites, cross-check README claims against actual source code, and emit a structured Pydantic `Evaluation` scorecard to `evaluation.json`.

By combining Agent Substrate with the Antigravity SDK, you create on-demand, zero-credential environments for each hackathon submission that start in seconds and leave zero residual state behind.

## Before You Begin
Duration: 05:00

### Select or Create a Project

Ensure you are working within the intended Google Cloud project:

1. **Select or Create a Project**: Choose an existing project or [create a new project](https://console.cloud.google.com/projectcreate) in the Google Cloud Console.
2. **Set the Active Project**: Note your **Project ID** to configure Cloud Shell to use it:
   ```bash
   gcloud config set project [YOUR_PROJECT_ID]
   ```

### Verify Billing

Make sure that billing is enabled for your Google Cloud project. You can learn more on how to do this by [following this guide](https://docs.cloud.google.com/billing/docs/how-to/verify-billing-enabled).

### Start Cloud Shell

Click the button below to start Google Cloud Shell, which comes pre-configured with the required developer and cloud command-line utilities.

<button>[Start Cloud Shell](https://console.cloud.google.com/?cloudshell=true)</button>

### Enable APIs

Run the following command in Cloud Shell to enable all target Google Cloud APIs required to run the platform:

```bash
gcloud services enable \
  container.googleapis.com \
  generativelanguage.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com \
  pubsub.googleapis.com \
  aiplatform.googleapis.com \
  cloudresourcemanager.googleapis.com \
  iam.googleapis.com \
  bigquery.googleapis.com \
  bigqueryconnection.googleapis.com \
  storage.googleapis.com
```

> **Why we enable these APIs:**
> Google Cloud services are disabled by default to prevent unauthorized access and charges. We enable these specific APIs to activate container orchestration (GKE), secure container storage (Artifact Registry), serverless build packaging (Cloud Build), reliable messaging queues (Pub/Sub), AI model services (**Gemini platform**), project configuration (Cloud Resource Manager & IAM), serverless data analytics (BigQuery), database-level AI bindings (BigQuery Connection), and Golden Snapshot object storage (Cloud Storage).

## Set up Infrastructure
Duration: 15:00

In this step, you will clone the code repository and execute the automated setup script to deploy target cloud architecture, Agent Substrate control plane components, and baseline cluster policies.

### Clone the Repository

Clone the repository containing all the application services, setup scripts, and Kubernetes manifest declarations:

```bash
git clone --depth 1 --filter=blob:none --sparse https://github.com/GoogleCloudPlatform/devrel-demos.git
cd devrel-demos
git sparse-checkout set codelabs/agent-substrate-code-evaluation/hackathon-judge
cd codelabs/agent-substrate-code-evaluation/hackathon-judge
```

### Run Deployment Script

The foundational setup of your cloud resources, database models, GKE cluster, and Agent Substrate (`v1.0.0`) control plane is automated by the `deploy.sh` script.

Run the script:

```bash
./deploy.sh
```

Follow the interactive shell prompts to supply configurations like your active Project ID and Target Region (`us-central1`). The script automatically generates a local `.env` configuration, provisions cloud resources in parallel, compiles container images, and installs Agent Substrate (`v1.0.0`).

Here are the target operations performed by the script:

#### 1. Environment Configuration Setup
The script creates a `.env` configuration file to store GKE, Pub/Sub, BigQuery, and project variable parameters. Sourcing this file dynamically populates all subsequent Kubernetes and Agent Substrate manifest definitions.
> **Why we configure this environment file:**
> The `.env` file centralizes configuration parameters, ensuring that the GKE and Agent Substrate manifests you apply in subsequent steps use identical regional settings, project names, and image digests, strictly decoupling environmental configuration from source code.

#### 2. Google Cloud CLI & Target Project Configuration
The script verifies that `gcloud`, `bq`, `kubectl`, and `envsubst` utilities are installed, checks authentication state, and locks active configuration targets onto your active Google Cloud project.
> **Why we target the active project:**
> Setting the active target project prevents CLI commands from affecting other projects in your account and performs pre-flight authentication checks, ensuring that setup commands do not fail mid-deployment due to invalid credentials.

#### 3. Enabling Target Google Cloud APIs
The script performs an idempotent check to verify and enable the target Google Cloud service APIs (GKE, Artifact Registry, Cloud Build, Pub/Sub, **Gemini platform**, BigQuery, Cloud Storage, and IAM).
> **Why we enable Google Cloud APIs:**
> Managed cloud services must be activated before their endpoints are reachable or resources can be created. Enabling them at the start prepares the regional Google Cloud API gateway to handle subsequent resource provisioning commands.

#### 4. Provisioning the Artifact Registry Docker Repository
The script provisions a Docker container repository named `hackathon-judge-repo` in the selected target location.
> **Why we create an Artifact Registry repository:**
> GKE clusters require secure access to private container registries in the same regional network to pull application images quickly. Artifact Registry provides a secure, private host to catalog, scan, and store Docker container images.

#### 5. Provisioning the GKE Standard Cluster with Nested Virtualization
The script provisions a Google Kubernetes Engine (GKE) Standard cluster named `hackathon-judge-cluster` on the `rapid` release channel (`1.37`) using `c3-standard-4` nodes configured with `--enable-nested-virtualization`, `--gateway-api=standard`, and Workload Identity (`[PROJECT_ID].svc.id.goog`).
> **Why we deploy a GKE Standard cluster with nested virtualization:**
> Agent Substrate hosts isolated gVisor actor sandboxes inside long-running `ateworker-gvisor` worker pods managed by a `WorkerPool` and node-level `atelet` DaemonSet. Provisioning `c3-standard-4` nodes with `--enable-nested-virtualization` lets Agent Substrate's `ateworker-gvisor` and `atelet` runtime deliver user-space syscall isolation and sub-second memory snapshot restoration.

#### 6. Configuring Pub/Sub Topics and Subscriptions
The script provisions the message topics (`judging-tasks` and `judging-results`) along with their corresponding worker (`judging-tasks-agent-sub`) and API (`judging-results-backend-sub`) subscriptions.
> **Why we deploy Pub/Sub topics and subscriptions:**
> Evaluating code submissions takes several minutes and is resource-intensive. Using a message-queue architecture decouples the synchronous user-facing API from worker nodes. The API backend pushes jobs to the `judging-tasks` topic, and the Judge Orchestrator pulls tasks as capacity is available, preventing API blocking and providing automatic retry capabilities.

#### 7. Configuring BigQuery Datasets and Tables
The script creates the `hackathon_judge` dataset, applies structural SQL database schemas (`backend/internal/repository/schema.sql`), and loads idempotent seed records (`backend/internal/repository/seeds.sql`).
> **Why we provision BigQuery datasets and seed tables:**
> BigQuery stores the hackathon rubrics, project submissions, and structured evaluation scorecards, pre-populating the platform with sample hackathon submissions ready for immediate grading.

#### 8. Triggering Container Builds using Cloud Build
The script triggers the `cloudbuild.yaml` definition to compile the React UI (`frontend`), Go REST server (`backend`), Python **Judge Orchestrator** (`agent`, bundled with `kubectl-ate` `v1.0.0`), and FastAPI + **Antigravity SDK** sandbox runtime (`agent-sandbox`), packaging them into container images tagged with the active repository Git commit SHA and pushing them to Artifact Registry.

#### 9. Provisioning Cloud Storage Snapshot Buckets & Installing Agent Substrate (`v1.0.0`)
The script creates the Cloud Storage snapshot buckets (`[PROJECT_ID]-ate-snapshots` and `[PROJECT_ID]-ate-snapshots-assets`) and the `ate-snapshot-transfer` IAM service account, binds storage and token-creator permissions for `ate-snapshot-controller` and `ate-api-server`, builds `kubectl-ate` (`v1.0.0`) into `~/.local/bin`, resolves the digest-pinned `ateworker-gvisor` worker image into `.env`, and runs `ate-setup deploy ate-system`. It configures `--sandbox-classes="gvisor"`, `--credential-provider='{"enabled":false}'`, and `--ingress-allowed-spiffe-ids` so `hackathon-judge-sa` can route mTLS traffic through `atenet-router`.
> **Why we install the Agent Substrate control plane and snapshot storage:**
> Standard Kubernetes clusters cannot checkpoint or restore running container memory states in sub-second times. Agent Substrate extends GKE with `ate-api-server`, `atenet-router`, `atenet-egress`, `ate-snapshot-controller`, and the `atelet` DaemonSet to orchestrate memory-snapshotted gVisor actors backed by Cloud Storage.

#### 10. Configuring Namespaces, Service Accounts, and Zero-Credential Workload Identity
The script applies `k8s/namespace.yaml` (labeling `hackathon-judge` with `ate.dev/api-access: "true"`), registers the `hackathon-judge-sa` Kubernetes Service Account (`k8s/service-account.yaml`), applies RBAC (`k8s/rbac.yaml`) and `k8s/configmap.yaml`, and binds Google Cloud IAM roles (`roles/bigquery.dataEditor`, `roles/bigquery.jobUser`, `roles/pubsub.publisher`, `roles/pubsub.subscriber`, and `roles/aiplatform.user`) strictly to `hackathon-judge-sa`.
> **Why Judge Actors receive zero Google Cloud credentials:**
> Only the trusted Judge Orchestrator (`hackathon-judge-sa`) receives Workload Identity IAM bindings. Judge Actors run without a Kubernetes Service Account or metadata server access, ensuring untrusted participant code can never exfiltrate cloud tokens.

---

## Configure Agent Substrate WorkerPool, AccessPolicy, and ActorTemplate
Duration: 05:00

In this step, you will configure the Agent Substrate **WorkerPool**, **AccessPolicy**, and **ActorTemplate** (**Golden Snapshot**) resources that power sub-second, zero-credential code evaluations.

### Source Environment Variables

Before applying templates that require environment variables, source the `setup-env.sh` script to add `~/.local/bin` (`kubectl-ate`) to your `PATH` and export all necessary variables (including the digest-pinned `SUBSTRATE_WORKER_IMAGE` and `AGENT_SANDBOX_IMAGE`) to your shell:

```bash
source ./setup-env.sh
```

### Apply WorkerPool

Use `envsubst` to populate `k8s/substrate.yaml` and apply the `WorkerPool` (`judge-workerpool`), which maintains 3 `ateworker-gvisor` worker pods in the `hackathon-judge` namespace:

```bash
envsubst < k8s/substrate.yaml | kubectl apply -f -
```

> **Why we deploy the WorkerPool with 4608Mi memory:**
> The `WorkerPool` manages warm `ateworker-gvisor` host pods on GKE nodes. Because `ateworker-gvisor` reserves `64Mi` of memory overhead per worker pod from `limits.memory`, configuring `4608Mi` (`4.5 GiB`) on each worker pod provides `512Mi` (at least `64Mi`) of headroom above the `4Gi` (`4096Mi`) Judge Actor memory limit so actors schedule cleanly.

### Create Atespace and Apply AccessPolicy

Ensure the `hackathon-judge` atespace exists and apply `k8s/access-policy.yaml` to grant the Judge Orchestrator service account (`hackathon-judge-sa`) OpenFGA `editor` permissions within the atespace:

```bash
kubectl-ate create atespace hackathon-judge 2>/dev/null || true
kubectl-ate create access-policy -a hackathon-judge -f k8s/access-policy.yaml
```

> **Why we configure the Atespace AccessPolicy:**
> Agent Substrate enforces fine-grained authorization on every atespace. Binding `user:system:serviceaccount:hackathon-judge:hackathon-judge-sa` to the `editor` role allows the Judge Orchestrator pod to create and delete ephemeral Judge Actors and per-actor `EgressPolicy` resources on demand.

### Apply ActorTemplate

Use `envsubst` to populate `k8s/actortemplate.yaml` with your active environment variables and register the `judge-sandbox` `ActorTemplate`:

```bash
source ./setup-env.sh
envsubst < k8s/actortemplate.yaml | kubectl-ate create actor-template -f -
```

Verify that the `WorkerPool` standby pods are running and inspect the `ActorTemplate` Golden Snapshot status:

```bash
kubectl get pods -n hackathon-judge -l workload=hackathon-judge
kubectl-ate get actor-template judge-sandbox -a hackathon-judge -o json
```

> aside positive
> **Golden Snapshot pre-warming latency:**
> When `judge-sandbox` is first created, Agent Substrate boots a builder actor, waits for `GET /health` on port `80`, and checkpoints its memory state (`SNAPSHOT_FIDELITY_MEMORY`) to Cloud Storage. This initial snapshot takes roughly 30–60 seconds; once `status.goldenSnapshotStatus.snapshotTime` is populated, every new Judge Actor restores in ~1–2 seconds.

> **Why we deploy an ActorTemplate and Golden Snapshot:**
> The `ActorTemplate` (`k8s/actortemplate.yaml`) defines the blueprint for every Judge Actor: it runs the FastAPI sandbox server on port `80` with `NET_BIND_SERVICE` (matching `atenet-router` ingress), mounts the `egress-mitm.ate.dev` and `system-roots.ate.dev` CA certificates at `/run/ate/trust-bundle.pem` for outbound TLS inspection, sets `MODEL_PROXY_BASE_URL` to `http://agent.hackathon-judge.svc.cluster.local:8000/vertex-proxy`, and enables `SNAPSHOT_FIDELITY_MEMORY` so actors restore from pre-initialized memory rather than cold-booting containers.

---

## Deploy Application Components
Duration: 05:00

With the Agent Substrate `WorkerPool` and `ActorTemplate` Golden Snapshot active, you will now deploy the central backend API, Judge Orchestrator (`agent`), React web interface, and external ingress Gateway mapping.

### Deploy Backend

Deploy the orchestrator REST API backend:

```bash
source ./setup-env.sh
envsubst < k8s/backend.yaml | kubectl apply -f -
```

### Deploy Agent

Deploy the Python **Judge Orchestrator** (`agent`) and its internal cluster `Service`:

```bash
source ./setup-env.sh
envsubst < k8s/agent.yaml | kubectl apply -f -
```

> **Why the Agent deployment mounts router-client credentials and exposes a Model Proxy service:**
> `k8s/agent.yaml` mounts a projected volume at `/run/router-client` containing an `api.ate-system.svc` ServiceAccount token (`ate-api-token`), the `servicedns.podcert.ate.dev/identity` trust bundle (`trust-bundle.pem`), and an ECDSA P-256 `podidentity.podcert.ate.dev/identity` client certificate (`credential-bundle.pem`) so the orchestrator can authenticate to `atenet-router` over mTLS. It also creates `Service/agent` on ports `8000` and `8082` to expose the internal `/vertex-proxy` endpoint to zero-credential Judge Actors.

### Deploy Frontend

Deploy the interactive web user interface:

```bash
source ./setup-env.sh
envsubst < k8s/frontend.yaml | kubectl apply -f -
```

### Configure External Gateway and Routing

Deploy the main Gateway and ingress HTTP routes mapping external client traffic:

```bash
kubectl apply -f k8s/gateway.yaml
```

> **Why we deploy the External Ingress Gateway:**
> The external Gateway exposes our services using the Kubernetes Gateway API. It provisions a load-balanced public IP address and maps routes based on path rules—directing API requests under `/api/*` to the Go Backend, and mapping all other client web traffic (`/`) to the React Frontend, securing public cluster access.

### Verify Rollouts

Wait until all three core service deployments have reached a healthy, ready rollout status:

```bash
kubectl rollout status deployment/backend -n hackathon-judge --timeout=300s
kubectl rollout status deployment/agent -n hackathon-judge --timeout=300s
kubectl rollout status deployment/frontend -n hackathon-judge --timeout=300s
```

---

## Verify and Use the Application
Duration: 05:00

### Access the UI

Fetch the external public IP address of the newly provisioned main load balancer gateway:

> aside negative
> **Wait for Load Balancer Provisioning:**
> Google Cloud's global external load balancer takes **2 to 5 minutes** to allocate IP resources, set up global forwarding rules, and configure the global routing path. If the `ADDRESS` column is empty or shows `[pending]`, wait a few minutes and query again.

To watch the provisioning status in real time, run the command with the watch flag (`-w`) and wait until a public IP address is populated in the `ADDRESS` field:

```bash
kubectl get gateway -n hackathon-judge hackathon-judge-gateway -w
```

When successfully provisioned, you should see output similar to the following:

```console
NAME                      CLASS    ADDRESS          PROGRAMMED   AGE
hackathon-judge-gateway   gke-l7   34.120.120.120   True         3m
```

Once you see a valid public IP address in the `ADDRESS` column and the `PROGRAMMED` status is `True`, press `Ctrl+C` to stop the watch.

> **Why we get the Gateway status:**
> The Gateway API handles public ingress. Checking the Gateway status returns the public, load-balanced external IP address allocated by Google Cloud's external global load balancer to our cluster, which represents the public address of our platform.

Open the allocated public IP address in your browser to load the Hackathon Judge dashboard.

### Submit Tasks

- Use the frontend UI to navigate to the dashboard and select the hackathon.

![Dashboard](images/hackathon-dashboard.png)

- On any of the projects, click `Run Agent` to start the agent judging the project against the rubric.

![Projects](images/hackathon-judge.png)

### Watch Judge Actor Kickoff

Monitor the active Judge Actors inside the `hackathon-judge` atespace to see an ephemeral Judge Actor restored from the `judge-sandbox` Golden Snapshot for judging execution:

```bash
kubectl-ate get actors -a hackathon-judge
```

While an evaluation is in progress, you should see output similar to the following:

```console
NAME             TEMPLATE                        STATE
judge-b1d9d23b   hackathon-judge/judge-sandbox   ACTOR_STATE_RUNNING
```

> aside positive
> **Ephemeral Golden Snapshot Restore & Automatic Teardown by Design:**
> When the Judge Orchestrator pulls a task from Pub/Sub, it calls `kubectl-ate create actor` and `kubectl-ate create egress-policy` to restore a fresh `judge-[8-hex]` actor from the `judge-sandbox` Golden Snapshot in ~1–2 seconds. The moment the evaluation completes (or if an error occurs), the orchestrator's `finally` block deletes both the `EgressPolicy` and the Judge Actor so no files, processes, or network rules ever persist between submissions.

Check the logs of the Judge Orchestrator pod to witness Pub/Sub task pickup, asynchronous command polling over `atenet-router`, zero-credential Model Proxy requests, and Judge Actor teardown:

```bash
kubectl logs -l app=agent -n hackathon-judge
```

You should see log output similar to the following:

```console
INFO:src.adapters.inbound.pubsub_subscriber:Received message: 14829104582910432
INFO:src.adapters.inbound.pubsub_subscriber:Processing message 14829104582910432
INFO:httpx:HTTP Request: GET https://atenet-router.ate-system.svc:443/health "HTTP/1.1 200 OK"
INFO:httpx:HTTP Request: POST https://atenet-router.ate-system.svc:443/execute/start "HTTP/1.1 200 OK"
INFO:httpx:HTTP Request: GET https://atenet-router.ate-system.svc:443/execute/status/9f4c2a18e7b04d12a6c3e5f701234567 "HTTP/1.1 200 OK"
INFO:     10.24.2.8:51342 - "POST /vertex-proxy/v1beta1/projects/my-project/locations/global/publishers/google/models/gemini-3.8-flash:streamGenerateContent?alt=sse HTTP/1.1" 200 OK
INFO:httpx:HTTP Request: GET https://atenet-router.ate-system.svc:443/download/evaluation.json "HTTP/1.1 200 OK"
INFO:src.adapters.outbound.shared_sandbox:Terminating Judge Actor...
INFO:src.adapters.inbound.pubsub_subscriber:Successfully processed and acked message 14829104582910432
```

> **Why we inspect Agent logs:**
> Inspecting the Judge Orchestrator logs displays both sides of the evaluation lifecycle in real time: outbound `httpx` calls through `atenet-router` to verify Judge Actor readiness (`/health`), start and poll non-blocking repository clone and `uv run python -m app.judge` jobs (`/execute/start`, `/execute/status/{job_id}`), and download `evaluation.json`, alongside inbound uvicorn requests served by the orchestrator's internal Model Proxy (`POST /vertex-proxy/...`) forwarding streaming **Gemini platform** calls before terminating the Judge Actor.

---

## (Optional) How this works
Duration: 03:00

### Agent Substrate & Antigravity SDK Architecture

Judging an engineering project requires cloning untrusted repositories, installing third-party libraries, building source code, and running real test suites.

Running raw user code poses severe security risks, including host compromise, container breakouts, lateral movement, and cloud credential exfiltration. **Agent Substrate** (`v1.0.0`) and the **Google Antigravity SDK** mitigate these risks by executing every evaluation inside an ephemeral, zero-credential gVisor **Judge Actor** restored from a memory-fidelity **Golden Snapshot**.

#### System Interaction Flow

The diagram below maps how the components of our event-driven system communicate during a secure, sandboxed judging execution:

![Hackathon Judge Agent Substrate and Antigravity SDK architecture diagram](images/agent-substrate-hackathon-judge-flow.png)

#### How the Engaged Tools and Components Work Together

* **React Frontend UI**: Exposes an interactive interface where users configure rubric criteria, register teams, submit project repository URLs, and review finalized grading scorecards, including file-level discrepancies and engineering critiques.
* **Go REST Backend API**: Manages global API endpoints, stores project configurations and scorecards in BigQuery, and publishes judging tasks to Pub/Sub to decouple heavy evaluation pipelines.
* **Google Pub/Sub**: The message broker (`judging-tasks` and `judging-results`) that buffers task and result messages asynchronously between the Go Backend API and the Judge Orchestrator.
* **Python Judge Orchestrator & Model Proxy (`agent`)**: A FastAPI background worker that pulls tasks from Pub/Sub and runs `SandboxDirectAdapter` by default (`AGENT_TYPE=sandbox_direct`, with an optional Google ADK coordinator adapter available via `AGENT_TYPE=adk`). It invokes `evaluate_repository` (`SubstrateActorClient`) to manage ephemeral Judge Actors via `kubectl-ate` and hosts the internal `/vertex-proxy` FastAPI route (`model_proxy.py`).
* **Agent Substrate `WorkerPool` (`judge-workerpool`)**: Maintains warm `ateworker-gvisor` pods on GKE Standard nodes with nested virtualization and gVisor (`runsc`) isolation, intercepting container system calls in user space to protect the host kernel.
* **`ActorTemplate` (`judge-sandbox`) & Golden Snapshot (`SNAPSHOT_FIDELITY_MEMORY`)**: Pre-initializes the `agent-sandbox` container image once, waits for `/health` on port `80`, and checkpoints the running memory and filesystem state to Cloud Storage so ephemeral `judge-[8-hex]` actors restore in ~1–2 seconds.
* **`atenet-router` & Asynchronous Command Polling**: Routes mTLS requests from the Judge Orchestrator to the target Judge Actor using `X-Ate-Target-Atespace`, `X-Ate-Target-Actor`, and `X-Ate-Target-Port: 80` headers. Because repository builds and multi-turn LLM grading can run for several minutes, the orchestrator starts commands asynchronously via `POST /execute/start` and polls `GET /execute/status/{job_id}` every 2 seconds so no single HTTP request ever hits an idle proxy timeout.
* **Per-Actor `EgressPolicy` & `atenet-egress`**: Each Judge Actor is paired with a dedicated `EgressPolicy` that allows only outbound HTTPS (`:443`, inspected via the `egress-mitm.ate.dev` CA bundle) to clone public GitHub repositories and fetch packages, plus internal HTTP (`:8000` and `:8082`) to `agent.hackathon-judge.svc.cluster.local`. All other cluster and metadata traffic is blocked.
* **FastAPI Sandbox Runtime (`app/main.py`)**: A lightweight Python server listening on port `80` (with `NET_BIND_SERVICE`) inside the Judge Actor. It exposes `/health`, `/execute/start`, `/execute/status/{job_id}`, `/upload`, and `/download/{path}` so the orchestrator can upload `criteria.md` and `prompt.md`, trigger evaluations, and retrieve `evaluation.json`.
* **Antigravity SDK Judge (`app/judge.py`)**: Invoked inside the Judge Actor via `uv run python -m app.judge`. It configures `google-antigravity` (`LocalAgentConfig`) with `policy.allow_all()` (since the gVisor Judge Actor is the isolation boundary), a `RunCommandConfig(timeout_seconds=180)` cap, a Pydantic `Evaluation` `response_schema`, and a `ModelTarget` whose `VertexEndpoint` points at `MODEL_PROXY_BASE_URL` to:
  - Inspect the repository layout and source files in `repo/`.
  - Install dependencies (`npm`, `pip`, `uv`, `go`) and run project build and test scripts.
  - Cross-examine README claims against actual source code to catch ghost features or hardcoded stubs, and write the validated scorecard to `evaluation.json`.
* **Standard Developer Environments**: Bundles `node`, `npm`, `yarn`, `pnpm`, `python`, `pip`, `uv`, `go`, `git`, `tree`, `ripgrep`, and `playwright` in the `agent-sandbox` container image, giving the Antigravity SDK agent a complete toolchain.

> aside positive
> **Zero-Credential Model Proxy security in action:**
> Unlike traditional sandboxes that mount a cloud service account inside the untrusted pod, Judge Actors hold **zero Google Cloud credentials**. When the Antigravity SDK inside the Judge Actor calls `gemini-3.8-flash`, it sends an unauthenticated request to `http://agent.hackathon-judge.svc.cluster.local:8000/vertex-proxy`. The Judge Orchestrator verifies that the upstream path matches `^v1(beta1)?/projects/.../locations/.../publishers/google/models/...:(generateContent|streamGenerateContent)$` (rejecting any other path with `403 Forbidden`), attaches its own Workload Identity access token (`hackathon-judge-sa`), and streams the **Gemini platform** response back into the actor. Even if untrusted code completely compromises the sandbox user space, there are no cloud credentials to steal!

---

## Clean Up
Duration: 03:00

To avoid ongoing charges to your Google Cloud account, delete the resources created during this codelab.

```bash
./destroy.sh
```

> **Why we clean up resources:**
> Google Cloud bills on a resource utilization model. Active resources like GKE Standard clusters, network load balancers, persistent volumes, and Cloud Storage buckets incur ongoing charges even when they are idle. Running `./destroy.sh` deletes the Pub/Sub topics and subscriptions, releases `ate-system` PersistentVolumeClaims and deletes the GKE Standard cluster, removes the Artifact Registry repository and BigQuery dataset, and deletes the Agent Substrate Cloud Storage snapshot buckets (`[PROJECT_ID]-ate-snapshots` and `[PROJECT_ID]-ate-snapshots-assets`) and `ate-snapshot-transfer` service account.

---

## Congratulations
Duration: 01:00

Congratulations! You have successfully deployed the Hackathon Judge application with **Agent Substrate** and the **Antigravity SDK** on GKE!

You have implemented a secure, event-driven AI evaluation platform capable of cloning, building, testing, and grading untrusted codebase submissions inside zero-credential gVisor sandboxes restored from memory-fidelity Golden Snapshots.

### What you've learned

- **GKE & Agent Substrate Infrastructure**: How to provision a GKE Standard cluster with nested virtualization (`--enable-nested-virtualization`), and deploy Agent Substrate (`v1.0.0`) with Cloud Storage snapshot buckets.
- **WorkerPools & Golden Snapshots**: How to configure a `WorkerPool` (`judge-workerpool`), `AccessPolicy`, and `ActorTemplate` (`judge-sandbox` with `SNAPSHOT_FIDELITY_MEMORY`) for sub-second Judge Actor restoration and automatic teardown.
- **Zero-Credential Sandboxing & Model Proxy**: How to isolate untrusted code with zero Google Cloud credentials, restrict outbound traffic with per-actor `EgressPolicy` rules, and proxy strictly validated model requests through the Judge Orchestrator to **Gemini platform**.
- **Antigravity SDK Code Evaluation**: How to run an autonomous `google-antigravity` judging agent over `atenet-router` mTLS using asynchronous `/execute/start` and `/execute/status/{job_id}` polling to produce structured Pydantic scorecards.

### Next steps

- Check out the [Gemini platform documentation](https://docs.cloud.google.com/gemini-enterprise-agent-platform).
- Explore [Agent Substrate on GitHub](https://github.com/agent-substrate/substrate).
- Learn more about [GKE Sandbox and gVisor overview](https://cloud.google.com/kubernetes-engine/docs/concepts/sandbox-pods).

### Reference docs

- [Kubernetes Gateway API Specification](https://gateway-api.sigs.k8s.io/)
- [gVisor Sandboxing Runtime](https://gvisor.dev/)
