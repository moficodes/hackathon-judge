SET @@dataset_project_id = '<<YOUR_PROJECT_ID>>';
SET @@dataset_id = 'hackathon-judge';

DECLARE current_run_time TIMESTAMP;
SET current_run_time = CURRENT_TIMESTAMP();

INSERT INTO `hackathon_judge.evaluations` (id, project_id, judge_id, criteria, created_at)
WITH LastRun AS (
  SELECT COALESCE(MAX(processing_date), TIMESTAMP('1970-01-01')) AS last_date
  FROM `hackathon_judge.projects`
),
NewObjects AS (
  SELECT uri, ref, updated, REGEXP_EXTRACT(uri, r'([^/]+)$') AS filename
  FROM `hackathon_judge.submissions_objects`
  WHERE updated > (SELECT last_date FROM LastRun)
    AND updated <= current_run_time
    AND uri LIKE '%README%' -- TODO: Confirm this is what we want
),
NewProjects AS (
  SELECT p.id AS project_id, p.name, o.uri, o.updated
  FROM `hackathon_judge.projects` p
  JOIN NewObjects o ON p.name = o.filename
),
CriteriaToUse AS (
  SELECT DISTINCT c.name, c.prompt, c.weight
  FROM `hackathon_judge.evaluations`, UNNEST(criteria) c
)
SELECT
  GENERATE_UUID() AS id,
  np.project_id,
  'AI_JUDGE' AS judge_id,
  ARRAY(
    SELECT AS STRUCT
      ctu.name,
      ctu.prompt,
      AI.SCORE(
        prompt => (ctu.prompt, OBJ.GET_ACCESS_URL(np.ref, 'r'))
      ) AS score,
      ctu.weight
    FROM CriteriaToUse ctu
  ) AS criteria,
  current_run_time AS created_at
FROM NewProjects np;

UPDATE `hackathon_judge.projects`
SET processing_date = current_run_time
WHERE id IN (
  SELECT project_id
  FROM `hackathon_judge.evaluations`
  WHERE judge_id = 'AI_JUDGE' AND created_at = current_run_time
);

