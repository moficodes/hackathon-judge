# Configurable Sandbox Connection Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Make the sandbox connection method configurable via environment variables in the agent.

**Architecture:** Introduce a utility function `get_sandbox_client` that reads environment variables to determine the connection strategy (Gateway, Router DNS, In-cluster IP, or In-cluster DNS) and returns a `SandboxClient` with the appropriate configuration.

**Tech Stack:** Python, k8s-agent-sandbox

---

### Task 1: Create Sandbox Utility for Configurable Connection

**Files:**
- Create: `agent/src/adapters/outbound/sandbox_utils.py`
- Test: `agent/tests/adapters/test_sandbox_utils.py`

- [ ] **Step 1: Write the failing test**

```python
import os
from unittest.mock import patch, MagicMock
from agent.src.adapters.outbound.sandbox_utils import get_sandbox_client
from k8s_agent_sandbox.models import (
    SandboxGatewayConnectionConfig,
    SandboxInClusterConnectionConfig,
    SandboxDirectConnectionConfig
)

def test_get_sandbox_client_gateway():
    with patch.dict(os.environ, {
        "SANDBOX_CONNECTION_METHOD": "gateway",
        "SANDBOX_GATEWAY_NAME": "my-gateway",
        "SANDBOX_NAMESPACE": "test-ns"
    }):
        with patch("agent.src.adapters.outbound.sandbox_utils.SandboxClient") as mock_client:
            get_sandbox_client()
            mock_client.assert_called_once()
            config = mock_client.call_args[1]["connection_config"]
            assert isinstance(config, SandboxGatewayConnectionConfig)
            assert config.gateway_name == "my-gateway"
            assert config.gateway_namespace == "test-ns"

def test_get_sandbox_client_router_dns():
    with patch.dict(os.environ, {
        "SANDBOX_CONNECTION_METHOD": "router_dns",
        "SANDBOX_ROUTER_URL": "http://router.svc"
    }):
        with patch("agent.src.adapters.outbound.sandbox_utils.SandboxClient") as mock_client:
            get_sandbox_client()
            config = mock_client.call_args[1]["connection_config"]
            assert isinstance(config, SandboxDirectConnectionConfig)
            assert config.api_url == "http://router.svc"

def test_get_sandbox_client_in_cluster_dns():
    with patch.dict(os.environ, {
        "SANDBOX_CONNECTION_METHOD": "in_cluster_dns"
    }):
        with patch("agent.src.adapters.outbound.sandbox_utils.SandboxClient") as mock_client:
            get_sandbox_client()
            config = mock_client.call_args[1]["connection_config"]
            assert isinstance(config, SandboxInClusterConnectionConfig)
            assert config.use_pod_ip is False

def test_get_sandbox_client_in_cluster_ip():
    with patch.dict(os.environ, {
        "SANDBOX_CONNECTION_METHOD": "in_cluster_ip"
    }):
        with patch("agent.src.adapters.outbound.sandbox_utils.SandboxClient") as mock_client:
            get_sandbox_client()
            config = mock_client.call_args[1]["connection_config"]
            assert isinstance(config, SandboxInClusterConnectionConfig)
            assert config.use_pod_ip is True
```

- [ ] **Step 2: Run test to verify it fails**

Run: `uv run pytest agent/tests/adapters/test_sandbox_utils.py`
Expected: FAIL (ModuleNotFoundError)

- [ ] **Step 3: Implement `get_sandbox_client`**

```python
import os
from k8s_agent_sandbox import SandboxClient
from k8s_agent_sandbox.models import (
    SandboxGatewayConnectionConfig,
    SandboxInClusterConnectionConfig,
    SandboxDirectConnectionConfig
)

def get_sandbox_client():
    method = os.getenv("SANDBOX_CONNECTION_METHOD", "gateway").lower()
    
    if method == "gateway":
        gateway_name = os.getenv("SANDBOX_GATEWAY_NAME", "sandbox-router-gateway")
        gateway_namespace = os.getenv("SANDBOX_GATEWAY_NAMESPACE", os.getenv("SANDBOX_NAMESPACE"))
        config = SandboxGatewayConnectionConfig(
            gateway_name=gateway_name,
            gateway_namespace=gateway_namespace
        )
    elif method == "router_dns":
        router_url = os.getenv("SANDBOX_ROUTER_URL")
        if not router_url:
            raise ValueError("SANDBOX_ROUTER_URL is required for router_dns connection method")
        config = SandboxDirectConnectionConfig(api_url=router_url)
    elif method == "in_cluster_dns":
        config = SandboxInClusterConnectionConfig(use_pod_ip=False)
    elif method == "in_cluster_ip":
        config = SandboxInClusterConnectionConfig(use_pod_ip=True)
    else:
        raise ValueError(f"Unknown SANDBOX_CONNECTION_METHOD: {method}")
        
    return SandboxClient(connection_config=config)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `uv run pytest agent/tests/adapters/test_sandbox_utils.py`
Expected: PASS

### Task 2: Update SandboxDirectAdapter to use the utility

**Files:**
- Modify: `agent/src/adapters/outbound/sandbox_direct.py`

- [ ] **Step 1: Update imports and use `get_sandbox_client`**

```python
# Replace imports
from k8s_agent_sandbox import SandboxClient
from k8s_agent_sandbox.models import SandboxGatewayConnectionConfig
# with
from .sandbox_utils import get_sandbox_client

# In _evaluate_repository, replace:
client = SandboxClient(
    connection_config=SandboxGatewayConnectionConfig(
        gateway_name="sandbox-router-gateway",
        gateway_namespace=namespace
    )
)
# with
client = get_sandbox_client()
```

- [ ] **Step 2: Verify existing tests still pass**

Run: `uv run pytest agent/tests/adapters/test_sandbox_direct.py`
Expected: PASS (since it mocks `_evaluate_repository` which hasn't changed its external interface)

### Task 3: Update ADKAgentAdapter to use the utility

**Files:**
- Modify: `agent/src/adapters/outbound/adk_agent.py`

- [ ] **Step 1: Update imports and use `get_sandbox_client`**

```python
# Replace imports
from k8s_agent_sandbox import SandboxClient
from k8s_agent_sandbox.models import SandboxDirectConnectionConfig, SandboxGatewayConnectionConfig
# with
from .sandbox_utils import get_sandbox_client

# In evaluate_repository, replace:
client = SandboxClient(
    connection_config=SandboxGatewayConnectionConfig(
        gateway_name="sandbox-router-gateway",
        gateway_namespace=namespace
    )
)
# with
client = get_sandbox_client()
```

- [ ] **Step 2: Verify existing tests still pass**

Run: `uv run pytest agent/tests/adapters/test_outbound.py` (assuming tests are there)
Expected: PASS
