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
