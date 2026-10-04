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

import asyncio
import json

import pytest

from app import judge
from app.judge_schema import Evaluation
from google.antigravity import BuiltinTools
from google.antigravity.types import ToolCall

VALID = {
    "scores": [{"name": "Innovation", "score": 7, "reasoning": "repo/main.go:10"}],
    "total_score": 7,
    "overall_comments": "ok",
    "confidence_score": 0.8,
}


class FakeResponse:
    def __init__(self, output):
        self._output = output

    @property
    def chunks(self):
        async def gen():
            yield ToolCall(name=BuiltinTools.RUN_COMMAND, args={"command": "ls"})
        return gen()

    async def structured_output(self):
        return self._output


@pytest.fixture
def fake_agent(monkeypatch):
    """Replaces app.judge.Agent; returns a dict to configure behaviour and inspect calls."""
    state = {"output": VALID, "delay": 0, "error": None, "configs": [], "messages": []}

    class FakeAgent:
        def __init__(self, config):
            state["configs"].append(config)

        async def __aenter__(self):
            if state["error"]:
                raise state["error"]
            return self

        async def __aexit__(self, *exc):
            return False

        async def chat(self, message):
            state["messages"].append(message)
            if state["delay"]:
                await asyncio.sleep(state["delay"])
            return FakeResponse(state["output"])

    monkeypatch.setattr(judge, "Agent", FakeAgent)
    return state


@pytest.fixture
def workdir(tmp_path):
    (tmp_path / "prompt.md").write_text("You are a judge.", encoding="utf-8")
    (tmp_path / "criteria.md").write_text("# Innovation 1-10", encoding="utf-8")
    return tmp_path


ENV = {"GOOGLE_CLOUD_PROJECT": "proj"}


def test_build_config_defaults():
    cfg = judge.build_config("sys", "/w", {"GOOGLE_CLOUD_PROJECT": "proj"})
    assert cfg.vertex is True
    assert cfg.project == "proj"
    assert cfg.location == "global"
    assert cfg.model is None
    assert cfg.system_instructions == "sys"
    assert cfg.workspaces == ["/w"]
    # The SDK normalizes a pydantic model class to its JSON schema string.
    assert json.loads(cfg.response_schema) == Evaluation.model_json_schema()
    assert len(cfg.policies) == 1
    assert cfg.capabilities.disabled_tools == [
        BuiltinTools.SEARCH_WEB,
        BuiltinTools.GENERATE_IMAGE,
        BuiltinTools.START_SUBAGENT,
    ]


def test_build_config_env_overrides():
    cfg = judge.build_config(
        "sys", "/w",
        {"GOOGLE_CLOUD_PROJECT": "proj", "GOOGLE_CLOUD_LOCATION": "us-central1", "JUDGE_MODEL": "gemini-x"},
    )
    assert cfg.location == "us-central1"
    assert cfg.model == "gemini-x"


def test_empty_judge_model_uses_sdk_default():
    cfg = judge.build_config("sys", "/w", {"GOOGLE_CLOUD_PROJECT": "proj", "JUDGE_MODEL": ""})
    assert cfg.model is None


def test_success_writes_evaluation(fake_agent, workdir):
    rc = asyncio.run(judge.run(workdir, ENV))
    assert rc == 0
    data = json.loads((workdir / "evaluation.json").read_text(encoding="utf-8"))
    assert data["total_score"] == 7
    assert data["scores"][0]["name"] == "Innovation"
    assert fake_agent["configs"][0].system_instructions == "You are a judge."
    assert "criteria.md" in fake_agent["messages"][0]
    assert "repo/" in fake_agent["messages"][0]


def test_success_accepts_pydantic_instance(fake_agent, workdir):
    fake_agent["output"] = Evaluation.model_validate(VALID)
    assert asyncio.run(judge.run(workdir, ENV)) == 0


def test_success_accepts_json_string(fake_agent, workdir):
    fake_agent["output"] = json.dumps(VALID)
    assert asyncio.run(judge.run(workdir, ENV)) == 0


@pytest.mark.parametrize("missing", ["prompt.md", "criteria.md"])
def test_missing_input_exits_2(fake_agent, workdir, missing, capsys):
    (workdir / missing).unlink()
    assert asyncio.run(judge.run(workdir, ENV)) == judge.EXIT_MISSING_INPUT == 2
    assert missing in capsys.readouterr().err
    assert fake_agent["configs"] == []


def test_timeout_exits_3(fake_agent, workdir, capsys):
    fake_agent["delay"] = 5
    rc = asyncio.run(judge.run(workdir, {**ENV, "JUDGE_TIMEOUT_SECONDS": "0.05"}))
    assert rc == judge.EXIT_TIMEOUT == 3
    assert "timed out" in capsys.readouterr().err
    assert not (workdir / "evaluation.json").exists()


@pytest.mark.parametrize("output", [None, {"scores": []}, "not json"])
def test_bad_output_exits_4(fake_agent, workdir, output, capsys):
    fake_agent["output"] = output
    assert asyncio.run(judge.run(workdir, ENV)) == judge.EXIT_BAD_OUTPUT == 4
    assert "structured output" in capsys.readouterr().err
    assert not (workdir / "evaluation.json").exists()


def test_sdk_error_exits_1(fake_agent, workdir, capsys):
    fake_agent["error"] = RuntimeError("boom")
    assert asyncio.run(judge.run(workdir, ENV)) == judge.EXIT_SDK_ERROR == 1
    assert "boom" in capsys.readouterr().err


def test_logs_tool_calls(fake_agent, workdir, capsys):
    asyncio.run(judge.run(workdir, ENV))
    assert "tool: run_command" in capsys.readouterr().out
