import os
import pytest
from unittest.mock import patch, MagicMock
from src.adapters.outbound.sandbox_utils import get_sandbox_client
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
        with patch("src.adapters.outbound.sandbox_utils.SandboxClient") as mock_client:
            get_sandbox_client()
            mock_client.assert_called_once()
            config = mock_client.call_args[1]["connection_config"]
            assert isinstance(config, SandboxGatewayConnectionConfig)
            assert config.gateway_name == "my-gateway"
            assert config.gateway_namespace == "test-ns"

def test_get_sandbox_client_router_dns():
    with patch.dict(os.environ, {
        "SANDBOX_CONNECTION_METHOD": "router_dns",
        "SANDBOX_NAMESPACE": "my-cool-namespace"
    }):
        with patch("src.adapters.outbound.sandbox_utils.SandboxClient") as mock_client:
            get_sandbox_client()
            config = mock_client.call_args[1]["connection_config"]
            assert isinstance(config, SandboxDirectConnectionConfig)
            assert config.api_url == "http://sandbox-router-svc.my-cool-namespace.svc.cluster.local:8080"

def test_get_sandbox_client_router_dns_missing_namespace():
    with patch.dict(os.environ, {
        "SANDBOX_CONNECTION_METHOD": "router_dns"
    }, clear=True):
        with pytest.raises(ValueError, match="SANDBOX_NAMESPACE is required for router_dns connection method"):
            get_sandbox_client()

def test_get_sandbox_client_in_cluster_dns():
    with patch.dict(os.environ, {
        "SANDBOX_CONNECTION_METHOD": "in_cluster_dns"
    }):
        with patch("src.adapters.outbound.sandbox_utils.SandboxClient") as mock_client:
            get_sandbox_client()
            config = mock_client.call_args[1]["connection_config"]
            assert isinstance(config, SandboxInClusterConnectionConfig)
            assert config.use_pod_ip is False

def test_get_sandbox_client_in_cluster_ip():
    with patch.dict(os.environ, {
        "SANDBOX_CONNECTION_METHOD": "in_cluster_ip"
    }):
        with patch("src.adapters.outbound.sandbox_utils.SandboxClient") as mock_client:
            get_sandbox_client()
            config = mock_client.call_args[1]["connection_config"]
            assert isinstance(config, SandboxInClusterConnectionConfig)
            assert config.use_pod_ip is True
