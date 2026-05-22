# Testing Strategy Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the hybrid testing strategy: unit testing BigQuery repository mappers, fixing the broken frontend test, and adding a basic E2E Playwright test.

**Architecture:** We will add a new test file for the backend mappers, modify the existing failing frontend test to use more robust queries, and create a new Playwright spec for the user journey.

**Tech Stack:** Go (testing), React Testing Library (frontend unit), Playwright (frontend E2E).

---

### Task 1: Backend - Unit Test BigQuery Mappers

**Files:**
- Create: `backend/internal/repository/bigquery_repo_test.go`
- Test: `backend/internal/repository/bigquery_repo_test.go`

- [ ] **Step 1: Write the failing test for mapBQHackathon**

```go
package repository

import (
	"testing"
	"time"

	"github.com/moficodes/hackathon-judge/backend/internal/domain"
)

func TestMapBQHackathon(t *testing.T) {
	repo := &BigQueryRepo{} // We only need the method, not a real client

	date := time.Now()
	bqh := bqHackathon{
		ID:          "h1",
		Title:       "Hackathon 1",
		Date:        date,
		Description: "Desc",
		Goal:        "Goal",
		Status:      "ACTIVE",
		Criteria: []bqCriterion{
			{Name: "C1", Description: "D1", Weight: 0.5, MaxScore: 10},
		},
		BonusCriteria: []bqCriterion{},
	}

	h, err := repo.mapBQHackathon(bqh)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h.ID != "h1" || h.Title != "Hackathon 1" || h.Description != "Desc" || h.Goal != "Goal" || h.Status != "ACTIVE" || !h.Date.Equal(date) {
		t.Errorf("basic fields mapped incorrectly. got: %+v", h)
	}

	if len(h.Criteria) != 1 || h.Criteria[0].Name != "C1" || h.Criteria[0].Weight != 0.5 || h.Criteria[0].MaxScore != 10 {
		t.Errorf("criteria mapped incorrectly. got: %+v", h.Criteria)
	}
    
    // Deliberate failure to ensure test runs
    t.Errorf("failing test")
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `cd backend && go test ./internal/repository -run TestMapBQHackathon -v`
Expected: FAIL with "failing test"

- [ ] **Step 3: Remove deliberate failure**

```go
package repository

import (
	"testing"
	"time"
)

func TestMapBQHackathon(t *testing.T) {
	repo := &BigQueryRepo{} // We only need the method, not a real client

	date := time.Now()
	bqh := bqHackathon{
		ID:          "h1",
		Title:       "Hackathon 1",
		Date:        date,
		Description: "Desc",
		Goal:        "Goal",
		Status:      "ACTIVE",
		Criteria: []bqCriterion{
			{Name: "C1", Description: "D1", Weight: 0.5, MaxScore: 10},
		},
		BonusCriteria: []bqCriterion{},
	}

	h, err := repo.mapBQHackathon(bqh)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if h.ID != "h1" || h.Title != "Hackathon 1" || h.Description != "Desc" || h.Goal != "Goal" || h.Status != "ACTIVE" || !h.Date.Equal(date) {
		t.Errorf("basic fields mapped incorrectly. got: %+v", h)
	}

	if len(h.Criteria) != 1 || h.Criteria[0].Name != "C1" || h.Criteria[0].Weight != 0.5 || h.Criteria[0].MaxScore != 10 {
		t.Errorf("criteria mapped incorrectly. got: %+v", h.Criteria)
	}
}
```

- [ ] **Step 4: Write test for mapBQProject**

```go
func TestMapBQProject(t *testing.T) {
	repo := &BigQueryRepo{}

    // Note: bigquery.NullString is handled internally in mapBQProject, we need to pass a mocked version if needed or just use the struct
	bqp := bqProject{
		ID:          "p1",
		Name:        "Project 1",
		Title:       "Title 1",
		URL:         "url",
		GitHubURL:   "gh",
		TeamName:    "Team",
		HackathonID: "h1",
		Score:       95.5,
	}

	p := repo.mapBQProject(bqp)

	if p.ID != "p1" || p.Name != "Project 1" || p.Title != "Title 1" || p.URL != "url" || p.GitHubURL != "gh" || p.TeamName != "Team" || p.HackathonID != "h1" || p.Score != 95.5 {
		t.Errorf("fields mapped incorrectly. got: %+v", p)
	}
}
```

- [ ] **Step 5: Run tests to verify they pass**

Run: `cd backend && go test ./internal/repository -v`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add backend/internal/repository/bigquery_repo_test.go
git commit -m "test(backend): add mapping unit tests for bigquery repo"
```

---

### Task 2: Frontend - Fix Component Test

**Files:**
- Modify: `frontend/src/pages/HackathonDetail.test.tsx`
- Test: `frontend/src/pages/HackathonDetail.test.tsx`

- [ ] **Step 1: Run the test to confirm it fails**

Run: `cd frontend && npm run test -- src/pages/HackathonDetail.test.tsx`
Expected: FAIL due to "Unable to find an element with the text: Score: 95.50"

- [ ] **Step 2: Modify the assertion in HackathonDetail.test.tsx**

Replace the strict text match with a regular expression or a custom matcher. Look for the `waitFor` block.

```tsx
    await waitFor(() => {
      expect(screen.getByText('AI Builder')).toBeInTheDocument();
      expect(screen.getByText('Team: Team Alpha')).toBeInTheDocument();
      // Use regex to match the text across elements
      expect(screen.getByText(/Score:/)).toBeInTheDocument();
      expect(screen.getByText('95.50')).toBeInTheDocument();
    });
```

- [ ] **Step 3: Run the test to verify it passes**

Run: `cd frontend && npm run test -- src/pages/HackathonDetail.test.tsx`
Expected: PASS

- [ ] **Step 4: Commit**

```bash
git add frontend/src/pages/HackathonDetail.test.tsx
git commit -m "test(frontend): fix strict text matching in hackathon detail test"
```

---

### Task 3: Frontend - E2E User Flow Test

**Files:**
- Create: `frontend/tests/e2e/hackathon_flow.spec.ts`
- Test: `frontend/tests/e2e/hackathon_flow.spec.ts`

- [ ] **Step 1: Write the basic E2E test**

```typescript
import { test, expect } from '@playwright/test';

test.describe('Hackathon Judging Flow', () => {
  test('should navigate from home to project details', async ({ page }) => {
    // 1. Start at home
    await page.goto('/');
    await expect(page).toHaveTitle(/Hackathon Judge/);
    
    // 2. Go to Dashboard
    await page.click('text=Go to Dashboard');
    await expect(page.url()).toContain('/dashboard');
    
    // Wait for mock data to load (assuming there's a hackathon title visible)
    // The exact text will depend on the mocked data or backend state if running against real API
    // We assume the page loads a hackathon
    await page.waitForSelector('text=Hackathons');
    
    // 3. Click first View button in the dashboard list
    // This relies on the "View" link existing for a hackathon.
    const viewButtons = page.locator('text=View').first();
    await viewButtons.click();
    
    // 4. We should be on the Hackathon Detail page
    await expect(page.url()).toContain('/hackathons/');
    
    // Wait for projects tab
    await page.waitForSelector('text=Projects');

    // 5. Click the first View button for a project
    // Since there are multiple view buttons (some for hackathons, some for projects if we are on the detail page)
    // we need to be careful. In the HackathonDetail page, project links say "View".
    const projectViewButtons = page.locator('text=View').first();
    await projectViewButtons.click();

    // 6. We should be on Project Detail page
    await expect(page.url()).toContain('/projects/');
    await expect(page.locator('text=Back to Hackathon')).toBeVisible();
  });
});
```

- [ ] **Step 2: Run the test to verify it fails (or passes if the app is running and mock data matches)**

Note: Playwright tests need the dev server running. The plan assumes the executor will start the frontend before running this test or use Playwright's `webServer` config. We will run it to see if it catches the flow.

Run: `cd frontend && npm run test:e2e`
Expected: It might fail if the dev server isn't running or if the mock data doesn't render fast enough.

- [ ] **Step 3: Commit**

```bash
git add frontend/tests/e2e/hackathon_flow.spec.ts
git commit -m "test(frontend): add e2e user flow test"
```