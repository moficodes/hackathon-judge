# src/adapters/outbound/adk_agent.py
from src.core.ports.agent_service import AgentService
from src.core.models.message import AgentRequest, AgentResponse, CategoryScore
import asyncio
import os

from google.adk.agents import Agent
from google.adk.runners import Runner
from google.adk.sessions import InMemorySessionService
from google.genai import types as genai_types
from pydantic import BaseModel, Field
from typing import List
from k8s_agent_sandbox import AsyncSandboxClient
from k8s_agent_sandbox.models import SandboxDirectConnectionConfig

async def evaluate_repository(github_url: str, judging_criteria: str) -> str:
    """
    Clones a repository into a secure sandbox and uses an autonomous agent to evaluate it 
    against the provided judging criteria.
    
    Args:
        github_url: The URL of the GitHub repository to clone.
        judging_criteria: A markdown formatted string containing the judging criteria and rubric.
        
    Returns:
        A JSON string containing the evaluation results.
    """
    template = os.getenv("SANDBOX_TEMPLATE")
    if not template:
        raise ValueError("SANDBOX_TEMPLATE environment variable is required")
    namespace = os.getenv("SANDBOX_NAMESPACE")
    if not namespace:
        raise ValueError("SANDBOX_NAMESPACE environment variable is required")
    config = SandboxDirectConnectionConfig(
        api_url=f"http://sandbox-router-svc.{namespace}.cluster.svc.local:8080"
    )
    async with AsyncSandboxClient(connection_config=config) as client:
        sandbox = await client.create_sandbox(
            template=template,
            namespace=namespace,
        )
        try:
            # Clone the repo
            await sandbox.commands.run(f"git clone {github_url} repo")

            # Write criteria to markdown file
            await sandbox.files.write("criteria.md", judging_criteria)

            # Invoke gemini CLI
            prompt = (
                "Evaluate the codebase in the repo directory against the criteria in criteria.md. "
                "Analyze the code, run it if necessary, and write your final findings as a JSON object "
                "to evaluation.json. The JSON must have 'scores' (list of {name, score, reasoning}), "
                "'total_score', 'overall_comments', and 'confidence_score'."
            )
            await sandbox.commands.run(f"gemini --yolo '{prompt}'")

            # Read the JSON evaluation
            result_json = await sandbox.files.read("evaluation.json")
            return result_json
        except Exception as e:
            return f"{{ 'error': 'Sandbox evaluation failed: {str(e)}' }}"
        finally:
            await client.delete_sandbox(claim_name=sandbox.claim_name, namespace=namespace)

class EvaluationScore(BaseModel):
    name: str = Field(description="The name of the scoring criteria category.")
    score: float = Field(description="The awarded score.")
    reasoning: str = Field(description="The reasoning behind this score.")

class EvaluationOutput(BaseModel):
    scores: List[EvaluationScore]
    total_score: float = Field(description="The total score summing up all category scores.")
    overall_comments: str = Field(description="Overall comments on the project.")
    confidence_score: float = Field(description="Score between 0.0 and 1.0 indicating confidence in the evaluation.")

class ADKAgentAdapter(AgentService):
    def __init__(self):
        self.agent = Agent(
            name="hackathon_judge",
            model="gemini-3-flash-preview",
            instruction="""
You are an expert hackathon judge evaluating a project.
Analyze the provided submission and evaluate it against the given criteria and rubric.
Be highly objective and provide detailed reasoning for each score.

You must use the `evaluate_repository` tool to delegate the deep codebase analysis to a sandboxed agent.
Pass the GitHub URL and the scoring criteria formatted as markdown to the tool.
The tool will return a JSON string with the evaluation results. Use this data to formulate your final response.
            """,
            tools=[evaluate_repository],
            output_schema=EvaluationOutput,
            output_key="evaluation_result",
        )
        self.session_service = InMemorySessionService()
        self.runner = Runner(
            agent=self.agent,
            app_name="hackathon_judge_app",
            session_service=self.session_service
        )

    async def process_message(self, request: AgentRequest) -> AgentResponse:
        try:
            # Format criteria
            criteria_text = "\n".join([f"- {c.name} (Weight: {c.weight}, Max Score: {c.max_score}): {c.description}" for c in request.scoring_criteria])
            
            prompt = f"""
Project Name: {request.project_name}
GitHub URL: {request.github_url}

Submission Description:
{request.submission_text}

Judging Rubric:
{request.judging_rubric}

Scoring Criteria:
{criteria_text}

Please evaluate this project and provide scores for each category.
"""
            session_id = f"session_{request.task_id}"
            user_id = "system"
            
            await self.session_service.create_session(
                app_name="hackathon_judge_app", 
                user_id=user_id, 
                session_id=session_id
            )

            # Run the agent
            async for _ in self.runner.run_async(
                user_id=user_id,
                session_id=session_id,
                new_message=genai_types.Content(role="user", parts=[genai_types.Part.from_text(text=prompt)])
            ):
                pass
                
            session = await self.session_service.get_session(app_name="hackathon_judge_app", user_id=user_id, session_id=session_id)
            eval_output_dict = session.state.get("evaluation_result")
            
            if not eval_output_dict:
                return AgentResponse(
                    task_id=request.task_id,
                    status="error",
                    error_message="Agent failed to produce evaluation_result in state."
                )
                
            # Parse output
            if isinstance(eval_output_dict, EvaluationOutput):
                eval_output = eval_output_dict
            else:
                eval_output = EvaluationOutput.model_validate(eval_output_dict)
            
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
            import traceback
            traceback.print_exc()
            return AgentResponse(
                task_id=request.task_id,
                status="error",
                error_message=str(e)
            )
