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

"""In-sandbox judge runner.

Run from the sandbox working directory (which contains repo/, criteria.md and
prompt.md) as:

    uv run python -m app.judge

Uses the Google Antigravity SDK against Vertex AI (Workload Identity via ADC),
and writes the validated result to evaluation.json.

Exit codes: 0 success, 1 SDK/model error, 2 missing input, 3 timeout,
4 missing or invalid structured output.
"""

import asyncio
import json
import os
import sys
from collections.abc import Mapping
from pathlib import Path

import pydantic
from google.antigravity import Agent, BuiltinTools, CapabilitiesConfig, LocalAgentConfig
from google.antigravity.types import ToolCall
from google.antigravity.hooks import policy

from app.judge_schema import Evaluation

DEFAULT_LOCATION = "global"  # gemini-3.8-flash is only served from the global endpoint.
DEFAULT_TIMEOUT_SECONDS = 540  # Below the agent's 600s command timeout.

DISABLED_TOOLS = [
    BuiltinTools.SEARCH_WEB,
    BuiltinTools.GENERATE_IMAGE,
    BuiltinTools.START_SUBAGENT,
]

USER_MESSAGE = (
    "Judge the project in the `repo/` directory against the rubric in `criteria.md`, "
    "following your instructions. When you are done, finish with your evaluation."
)

EXIT_OK = 0
EXIT_SDK_ERROR = 1
EXIT_MISSING_INPUT = 2
EXIT_TIMEOUT = 3
EXIT_BAD_OUTPUT = 4


def _err(message: str) -> None:
    print(f"judge: {message}", file=sys.stderr, flush=True)


def build_config(system_instructions: str, workspace: str, env: Mapping[str, str]) -> LocalAgentConfig:
    kwargs = dict(
        vertex=True,
        project=env.get("GOOGLE_CLOUD_PROJECT"),
        location=env.get("GOOGLE_CLOUD_LOCATION") or DEFAULT_LOCATION,
        system_instructions=system_instructions,
        workspaces=[workspace],
        response_schema=Evaluation,
        # The gVisor sandbox pod is the isolation boundary, so every tool call is approved.
        policies=[policy.allow_all()],
        capabilities=CapabilitiesConfig(disabled_tools=DISABLED_TOOLS),
    )
    model = env.get("JUDGE_MODEL")
    if model:
        kwargs["model"] = model
    return LocalAgentConfig(**kwargs)


async def _evaluate(config: LocalAgentConfig):
    async with Agent(config) as agent:
        response = await agent.chat(USER_MESSAGE)
        async for chunk in response.chunks:
            if isinstance(chunk, ToolCall):
                name = getattr(chunk.name, "value", chunk.name)
                print(f"judge: tool: {name}", flush=True)
        return await response.structured_output()


def _to_evaluation(raw) -> Evaluation:
    if raw is None:
        raise ValueError("agent returned no structured output")
    if isinstance(raw, pydantic.BaseModel):
        raw = raw.model_dump()
    if isinstance(raw, str):
        raw = json.loads(raw)
    return Evaluation.model_validate(raw)


async def run(workdir: Path, env: Mapping[str, str]) -> int:
    prompt_path = workdir / "prompt.md"
    criteria_path = workdir / "criteria.md"
    for path in (prompt_path, criteria_path):
        if not path.is_file():
            _err(f"missing input file {path.name}")
            return EXIT_MISSING_INPUT

    timeout = float(env.get("JUDGE_TIMEOUT_SECONDS") or DEFAULT_TIMEOUT_SECONDS)
    config = build_config(prompt_path.read_text(encoding="utf-8"), str(workdir), env)

    try:
        raw = await asyncio.wait_for(_evaluate(config), timeout=timeout)
    except TimeoutError:
        _err(f"timed out after {timeout:g}s")
        return EXIT_TIMEOUT
    except Exception as e:  # SDK, auth or model errors
        _err(f"agent error: {type(e).__name__}: {e}")
        return EXIT_SDK_ERROR

    try:
        evaluation = _to_evaluation(raw)
    except (ValueError, pydantic.ValidationError) as e:
        _err(f"invalid structured output: {e}")
        return EXIT_BAD_OUTPUT

    (workdir / "evaluation.json").write_text(evaluation.model_dump_json(indent=2), encoding="utf-8")
    print("judge: wrote evaluation.json", flush=True)
    return EXIT_OK


def main() -> None:
    sys.exit(asyncio.run(run(Path.cwd(), os.environ)))


if __name__ == "__main__":
    main()
