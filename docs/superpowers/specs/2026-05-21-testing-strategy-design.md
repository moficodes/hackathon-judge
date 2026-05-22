# Testing Strategy Design

## Overview
Implement a hybrid testing strategy across the project to address coverage gaps, primarily in the backend repository layer and frontend E2E flows, while fixing a currently failing component test.

## Scope

### Backend
*   **Target:** `internal/repository/bigquery_repo.go`
*   **Action:** Add unit tests for internal mapping logic (`mapBQHackathon`, `mapBQProject`).
*   **Rationale:** The actual `bigquery.Client` is difficult to mock reliably without hitting actual infrastructure. We will isolate the pure data transformation logic for unit testing to ensure data from BigQuery is correctly mapped into our domain models.

### Frontend
*   **Target:** `src/pages/HackathonDetail.test.tsx`
*   **Action:** Fix the failing test `HackathonDetail Component > renders projects data`.
*   **Rationale:** The test is failing due to strict text matching ("Score: 95.50") where the DOM structure splits the label and the value into different elements. The assertion needs to be updated to be more flexible.
*   **Target:** `tests/e2e/hackathon_flow.spec.ts`
*   **Action:** Implement a critical user journey E2E test.
*   **Rationale:** We need to verify that the core application flow works. The test will cover: Home Page -> Dashboard -> Hackathon Detail -> Project Detail.

### Agent
*   No new tests are planned. Existing unit tests cover models and adapter behavior adequately for this phase.

## Exclusions
*   Full integration testing against a real BigQuery instance is out of scope for this task and should be addressed in a future infrastructure testing effort.