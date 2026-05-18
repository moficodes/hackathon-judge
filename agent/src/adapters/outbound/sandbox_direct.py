# agent/src/adapters/outbound/sandbox_direct.py
import json
import logging
import os
import shlex
import traceback
from typing import List

from .sandbox_utils import get_sandbox_client
from pydantic import ValidationError

from src.core.models.message import AgentRequest, AgentResponse, CategoryScore
from src.core.ports.agent_service import AgentService
from pydantic import BaseModel, Field
from typing import List

class EvaluationScore(BaseModel):
    name: str = Field(description="The name of the scoring criteria category.")
    score: float = Field(description="The awarded score.")
    reasoning: str = Field(description="The reasoning behind this score.")

class EvaluationOutput(BaseModel):
    scores: List[EvaluationScore]
    total_score: float = Field(description="The total score summing up all category scores.")
    overall_comments: str = Field(description="Overall comments on the project.")
    confidence_score: float = Field(description="Score between 0.0 and 1.0 indicating confidence in the evaluation.")

logger = logging.getLogger(__name__)

class SandboxDirectAdapter(AgentService):
    async def process_message(self, request: AgentRequest) -> AgentResponse:
        try:
            # Format criteria
            criteria_text = "\n".join([f"- {c.name} (Weight: {c.weight}, Max Score: {c.max_score}): {c.description}" for c in request.scoring_criteria])
            
            judging_criteria = f"""
Project Name: {request.project_name}
GitHub URL: {request.github_url}

Submission Description:
{request.submission_text}

Judging Rubric:
{request.judging_rubric}

Scoring Criteria:
{criteria_text}
"""
            import asyncio
            result_json = await asyncio.to_thread(self._evaluate_repository, request.github_url, judging_criteria)
            
            try:
                eval_output_dict = json.loads(result_json)
                if "error" in eval_output_dict:
                    logger.error(f"Sandbox returned an error: {eval_output_dict['error']}")
                    return AgentResponse(
                        task_id=request.task_id,
                        status="error",
                        error_message=eval_output_dict['error']
                    )
                eval_output = EvaluationOutput.model_validate(eval_output_dict)
            except (json.JSONDecodeError, ValidationError) as e:
                logger.error(f"Failed to parse sandbox output: {result_json}")
                return AgentResponse(
                    task_id=request.task_id,
                    status="error",
                    error_message=f"Failed to parse sandbox output: {e}"
                )

            # Create lookups for max scores and weights
            max_score_lookup = {c.name: c.max_score for c in request.scoring_criteria}
            weight_lookup = {c.name: c.weight for c in request.scoring_criteria}
            
            # Map to response schema
            category_scores = [
                CategoryScore(
                    name=s.name, 
                    score=s.score, 
                    reasoning=s.reasoning,
                    max_score=max_score_lookup.get(s.name, 10.0),
                    weight=weight_lookup.get(s.name, 1.0)
                )
                for s in eval_output.scores
            ]
            
            return AgentResponse(
                task_id=request.task_id,
                status="success",
                scores=category_scores,
                total_score=eval_output.total_score,
                overall_comments=eval_output.overall_comments,
                confidence_score=eval_output.confidence_score
            )
            
        except Exception as e:
            traceback.print_exc()
            return AgentResponse(
                task_id=request.task_id,
                status="error",
                error_message=str(e)
            )

    def _evaluate_repository(self, github_url: str, judging_criteria: str) -> str:
        template = os.getenv("SANDBOX_TEMPLATE")
        if not template:
            raise ValueError("SANDBOX_TEMPLATE environment variable is required")
        namespace = os.getenv("SANDBOX_NAMESPACE")
        if not namespace:
            raise ValueError("SANDBOX_NAMESPACE environment variable is required")
        client = get_sandbox_client()
        sandbox = client.create_sandbox(
            template=template,
            namespace=namespace,
        )
        try:
            # Clone the repo
            safe_url = shlex.quote(github_url)
            sandbox.commands.run(f"git clone {safe_url} repo")
            # Write criteria to markdown file
            sandbox.files.write("criteria.md", judging_criteria)
            # Invoke gemini CLI
            prompt = (
"""
Identity & Objective:
You are an elite, highly critical Hackathon Judge and Technical Auditor. Your core mission is to aggressively cross-examine the project's markdown documentation/claims against the actual source code present in the repository. You are looking for inflation, missing implementations, or outright contradictions between what the team *claims* they built and what is *actually* written in the code.

Evaluation Framework:
1. Load and thoroughly read the evaluation rubric provided in `criteria.md`.
2. Inspect the codebase in the repository root directory. 
3. Execute the project's build, run, or test scripts if available to verify functional claims.

Verification Checklist (Look for these common discrepancies):
- "Ghost Features": Claims of advanced functionality in the README that are actually empty stubs, hardcoded mock data, or TODO comments in the code.
- "Plagiarism/Boilerplate": Distinguish between core hackathon work and pre-existing templates or heavy library scaffolding.
- "Fragility": Code that technically works but is hardcoded to a single user, API response, or test case, violating the spirit of the criteria.

Execution Steps:
1. Static Analysis: Scan all major files, architecture, and logic pathways.
2. Dynamic Validation: Locate any test suites or execution scripts. Attempt to run them. Note if they pass, fail, or if no tests exist.
3. Discrepancy Log: For every score, explicitly tie your reasoning to specific files, functions, or lines of code.

Output Specification:
Write your final evaluation strictly as a valid JSON object directly to `evaluation.json`. Do not include any markdown formatting wrappers (like ```json) inside the file itself. 

The JSON must strictly adhere to this schema:
{
  "scores": [
    {
      "name": "String (The specific criteria name from criteria.md)",
      "score": "Number (Based on the criteria scale, e.g., 1-10 or 1-5)",
      "reasoning": "String (Detailed critique. MUST include explicit file paths, function names, or line numbers verifying or refuting claims. Document test execution results here.)"
    }
  ],
  "total_score": "Number (The sum or calculated average as dictated by criteria.md)",
  "overall_comments": "String (A holistic summary of the project's engineering quality, documentation accuracy, and overall execution.)",
  "confidence_score": "Number (Scale of 0.0 to 1.0, reflecting how confidently you were able to review, build, or verify the codebase based on available files/tests.)"
}
"""                    
            )
            safe_prompt = shlex.quote(prompt)
            
            # sandbox.files.write('evaluation.json', mock_json)
            sandbox.files.write('prompt.md', safe_prompt)
            prompt = 'gemini --yolo -p "Use @criteria.md as a guide and follow @prompt.md for the prompt. Make sure the output is in JSON format and written to evaluation.json"'
            response = sandbox.commands.run(prompt, timeout=600)
            # Read the JSON evaluation
            print(response.stdout)
            print(response.stderr)
            result_json = sandbox.files.read("evaluation.json")
            print(result_json)
            return result_json
        except Exception as e:
            logger.error(f"Sandbox evaluation failed: {str(e)}")
            return json.dumps({"error": f"Sandbox evaluation failed: {str(e)}"})
        finally:
            print("exiting...")
            sandbox.terminate()
