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
import pytest
from unittest.mock import patch
from src.adapters.outbound.sandbox_utils import get_sandbox_client, DEFAULT_ROUTER_URL
from k8s_agent_sandbox.models import (
    SandboxDirectConnectionConfig,
    SandboxLocalTunnelConnectionConfig,
)


def _config_for(env):
    with patch.dict(os.environ, env, clear=True):
        with patch("src.adapters.outbound.sandbox_utils.SandboxClient") as mock_client:
            get_sandbox_client()
            mock_client.assert_called_once()
            return mock_client.call_args[1]["connection_config"]


def test_default_router_url_points_at_agent_sandbox_system():
    assert DEFAULT_ROUTER_URL == "http://sandbox-router-svc.agent-sandbox-system.svc.cluster.local:8080"


def test_router_is_default_method():
    config = _config_for({})
    assert isinstance(config, SandboxDirectConnectionConfig)
    assert config.api_url == DEFAULT_ROUTER_URL


def test_router_method_explicit():
    config = _config_for({"SANDBOX_CONNECTION_METHOD": "router"})
    assert isinstance(config, SandboxDirectConnectionConfig)
    assert config.api_url == DEFAULT_ROUTER_URL


def test_router_method_url_override():
    config = _config_for({
        "SANDBOX_CONNECTION_METHOD": "router",
        "SANDBOX_ROUTER_URL": "http://my-router.example:9000",
    })
    assert isinstance(config, SandboxDirectConnectionConfig)
    assert config.api_url == "http://my-router.example:9000"


def test_method_is_case_insensitive():
    config = _config_for({"SANDBOX_CONNECTION_METHOD": "ROUTER"})
    assert isinstance(config, SandboxDirectConnectionConfig)


def test_tunnel_method():
    config = _config_for({"SANDBOX_CONNECTION_METHOD": "tunnel"})
    assert isinstance(config, SandboxLocalTunnelConnectionConfig)


@pytest.mark.parametrize("method", ["gateway", "router_dns", "in_cluster_dns", "in_cluster_ip", "bogus"])
def test_unknown_or_removed_methods_raise(method):
    with patch.dict(os.environ, {"SANDBOX_CONNECTION_METHOD": method}, clear=True):
        with pytest.raises(ValueError, match=f"Unknown SANDBOX_CONNECTION_METHOD: {method}"):
            get_sandbox_client()
