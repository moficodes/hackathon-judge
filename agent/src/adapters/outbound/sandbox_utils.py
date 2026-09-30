# Copyright 2026 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import os
from k8s_agent_sandbox import SandboxClient
from k8s_agent_sandbox.models import (
    SandboxDirectConnectionConfig,
    SandboxLocalTunnelConnectionConfig,
)

DEFAULT_ROUTER_URL = "http://sandbox-router-svc.agent-sandbox-system.svc.cluster.local:8080"


def get_sandbox_client() -> SandboxClient:
    """Builds a SandboxClient for the configured connection method.

    SANDBOX_CONNECTION_METHOD:
      - "router" (default): in-cluster, via the sandbox-router Service.
        Override the URL with SANDBOX_ROUTER_URL.
      - "tunnel": local development; the SDK port-forwards to the router
        Service in agent-sandbox-system using the local kubeconfig.
    """
    method = os.getenv("SANDBOX_CONNECTION_METHOD", "router").lower()

    if method == "router":
        config = SandboxDirectConnectionConfig(
            api_url=os.getenv("SANDBOX_ROUTER_URL", DEFAULT_ROUTER_URL)
        )
    elif method == "tunnel":
        config = SandboxLocalTunnelConnectionConfig()
    else:
        raise ValueError(f"Unknown SANDBOX_CONNECTION_METHOD: {method}")

    return SandboxClient(connection_config=config)
