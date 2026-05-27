# SQL Seeds & Idempotent Deployment Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Transition BigQuery data initialization to use SQL seeds primarily and ensure the deployment process is idempotent.

**Architecture:** 
- Use BigQuery SQL (`seeds.sql`) for all initial data (Hackathons, Projects, Evaluations).
- Implement `INSERT ... SELECT ... WHERE NOT EXISTS` patterns in `seeds.sql` to prevent duplicates on re-deployment.
- Remove redundant CSV loading logic from `deploy.sh` and `deploy-all.sh`.

**Tech Stack:** Bash, Google Cloud SDK (bq), BigQuery SQL.

---

### Task 1: Make seeds.sql Idempotent

**Files:**
- Modify: `backend/internal/repository/seeds.sql`

- [ ] **Step 1: Refactor Hackathons INSERT**
  Change the `INSERT` statement for `hackathons` to use a subquery check.

```sql
INSERT INTO `hackathons` (id, title, date, description, goal, status, criteria, bonus_criteria)
SELECT 'CFB272DC-BDA6-45EF-B899-343B0EAB85E1', 'Productivity Hackathon', '2026-05-11T09:00:00Z', 'A hackathon focused on building tools to enhance productivity and efficiency.', 'Build tools that drastically improve productivity and efficiency.', 'active', 
  [STRUCT('crit_innov' as id, 'Innovation & Originality' as name, '...', 0.2, 5.0, 5.0), ...], 
  [...]
WHERE NOT EXISTS (SELECT 1 FROM `hackathons` WHERE id = 'CFB272DC-BDA6-45EF-B899-343B0EAB85E1');
```

- [ ] **Step 2: Refactor Projects INSERTs**
  Convert all 9 project `INSERT` statements to use the same `WHERE NOT EXISTS` pattern.

- [ ] **Step 3: Add Evaluations Seed (Optional/Sample)**
  Add a placeholder or sample evaluation if needed to verify the table works.

- [ ] **Step 4: Commit**
```bash
git add backend/internal/repository/seeds.sql
git commit -m "refactor: make seeds.sql idempotent using WHERE NOT EXISTS"
```

### Task 2: Update deploy.sh and deploy-all.sh

**Files:**
- Modify: `deploy.sh`
- Modify: `deploy-all.sh`

- [ ] **Step 1: Remove CSV loading logic from deploy.sh**
  Remove the `CSV_TABLE_MAP` array and the `for` loop that calls `bq load`.

- [ ] **Step 2: Remove CSV loading logic from deploy-all.sh**
  Apply the same removal to `deploy-all.sh`.

- [ ] **Step 3: Update README downloader comments**
  Clarify that `projects.csv` is used as a manifest for README downloads but BigQuery data comes from `seeds.sql`.

- [ ] **Step 4: Commit**
```bash
git add deploy.sh deploy-all.sh
git commit -m "chore: remove redundant CSV loading logic from deployment scripts"
```

### Task 3: Verification

- [ ] **Step 1: Dry run seeds.sql**
  Run the `seeds.sql` against a test dataset or manually verify the syntax of one `INSERT` statement.

- [ ] **Step 2: Verify idempotency logic**
  Simulate a second run of the seed query to ensure it doesn't fail or add duplicates.

- [ ] **Step 3: Final check of deploy scripts**
  Ensure no other part of the script depends on the CSV loading logic being present.
