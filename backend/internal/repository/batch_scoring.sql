-- This query is intended to be run on-demand for a specific project.
-- The calling application (e.g. Go backend) should provide the @project_id parameter.

INSERT INTO `hackathon_judge.evaluations` (id, project_id, judge_id, status, criteria_json, total_score, comment, created_at)
WITH ProjectInfo AS (
  SELECT p.id, p.readme_ref, p.hackathon_id
  FROM `hackathon_judge.projects` p
  WHERE p.id = @project_id
),
CriteriaToScore AS (
  SELECT c.id, c.name, c.description, c.max_score, c.weight
  FROM `hackathon_judge.hackathons` h, UNNEST(h.criteria) c
  JOIN ProjectInfo pi ON h.id = pi.hackathon_id
),
Scored AS (
  SELECT
    id, name, description, weight, max_score,
    AI.SCORE(
      prompt => ('Evaluate this project against the following rubric:', description, OBJ.GET_ACCESS_URL((SELECT readme_ref FROM ProjectInfo), 'r'))
    ) as score
  FROM CriteriaToScore
),
NewEval AS (
  SELECT
    GENERATE_UUID() as id,
    (SELECT id FROM ProjectInfo) as project_id,
    'BQ_JUDGE' as judge_id,
    'COMPLETED' as status,
    ARRAY(SELECT AS STRUCT name, description, weight, score, max_score FROM Scored) as criteria_json,
    (SELECT SUM(score) FROM Scored) as total_score,
    'Automated evaluation' as comment,
    CURRENT_TIMESTAMP() as created_at
)
SELECT * FROM NewEval;

-- Update the project score (optional, but often desired after scoring)
UPDATE `hackathon_judge.projects`
SET score = (SELECT total_score FROM (
    SELECT total_score FROM `hackathon_judge.evaluations` 
    WHERE project_id = @project_id 
    ORDER BY created_at DESC LIMIT 1
))
WHERE id = @project_id;
