from k8s_agent_sandbox import AsyncSandboxClient
from k8s_agent_sandbox.models import SandboxDirectConnectionConfig
import asyncio
import os

async def main():
    template = os.getenv("SANDBOX_TEMPLATE", "sandbox-hackathon-judge-template")
    namespace = os.getenv("SANDBOX_NAMESPACE", "hackathon-judge")
    config = SandboxDirectConnectionConfig(
        api_url=f"http://localhost:9000"
    )
    async with AsyncSandboxClient(connection_config=config) as client:
        sandbox = await client.create_sandbox(
            template=template,
            namespace=namespace,
        )
        result = await sandbox.commands.run("gemini --yolo -p 'hello'")
        print(result.stdout)
        print(result.stderr)
        print(result)


if __name__ == "__main__":
    asyncio.run(main())
