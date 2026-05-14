-- This query is intended to be run on-demand for a specific project.
-- The calling application (e.g. Go backend) should provide the @project_id parameter.

MERGE `hackathon_judge.projects` t
USING (
  WITH ProjectInfo AS (
    SELECT p.id, p.readme_ref, p.hackathon_id, p.evaluations
    FROM `hackathon_judge.projects` p
    WHERE p.id = @project_id
  ),
  CriteriaToScore AS (
    SELECT c.id, c.name, c.description, c.max_score
    FROM `hackathon_judge.hackathons` h, UNNEST(h.criteria) c
    JOIN ProjectInfo pi ON h.id = pi.hackathon_id
  ),
  Scored AS (
    SELECT
      id, name, description, max_score,
      AI.SCORE(
        prompt => (description, OBJ.GET_ACCESS_URL((SELECT readme_ref FROM ProjectInfo), 'r'))
      ) as score
    FROM CriteriaToScore
  ),
  NewEval AS (
    SELECT
      'BQ_JUDGE' as judge_id,
      ARRAY(SELECT AS STRUCT id, name, description, score, max_score FROM Scored) as criteria,
      (SELECT SUM(score) FROM Scored) as total_score,
      'Automated evaluation' as comment,
      CURRENT_TIMESTAMP() as created_at
  )
  SELECT pi.id, ARRAY_CONCAT(pi.evaluations, [ne]) as new_evaluations
  FROM ProjectInfo pi
  CROSS JOIN NewEval ne
) s
ON t.id = s.id
WHEN MATCHED THEN
  UPDATE SET evaluations = s.new_evaluations;

-- sample Hardcoded scoring for one project
UPDATE `hackathon_judge.projects` p
SET evaluations = ARRAY_CONCAT(evaluations, [
  STRUCT(
    'BQ_JUDGE' as judge_id,
    [STRUCT('crit_clean' as id, 'Clean Code' as name, 'Repository is well-documented with a clear README.' as description, AI.SCORE(prompt => ('Repository is well-documented with a clear README.', OBJ.GET_ACCESS_URL(p.readme_ref, 'r'))) as score, 2.0 as max_score)],
    2.0 as total_score,
    'Hardcoded test evaluation' as comment,
    CURRENT_TIMESTAMP() as created_at
  )
])
WHERE p.id = '14710f6b-dbf7-4055-8106-9dbed109dae2';

