# Hackathon Judge - Frontend Detail Pages Enhancement Design

## 1. Overview
The goal of this project is to enhance the frontend detail pages (`HackathonDetail` and `ProjectDetail`) to provide richer context to users. Currently, these pages lack important metadata. This design outlines the necessary backend API additions and frontend component updates to display full hackathon criteria, project metadata, and detailed evaluation reasoning.

## 2. Architecture & Data Flow

The frontend will use parallel data fetching (Option 1) to retrieve the necessary data for the `ProjectDetail` page.

### 2.1 Backend API Additions
To support the detailed views, two new `GET` endpoints are required:
- `GET /api/hackathons/:id`
  - Fetches a single `Hackathon` by ID.
  - Returns the full `domain.Hackathon` object, crucially including `Criteria` and `BonusCriteria`.
- `GET /api/projects/:id`
  - Fetches a single `Project` by ID.
  - Returns the full `domain.Project` object, including `URL`, `GitHubURL`, and `TeamName`.

### 2.2 Frontend Model Updates (`src/types/models.ts`)
The TypeScript interfaces must be updated to align with the backend `domain` models:
- **`Hackathon` interface:** Add `criteria` and `bonus_criteria` properties. These will be arrays of objects containing `name`, `weight`, `description`, and `max_score`.
- **`CriteriaScore` interface:** Add the `reasoning: string` property to capture the AI's justification for the score.

## 3. UI Component Enhancements

### 3.1 HackathonDetail.tsx
- **Layout:** Implement a Tabbed Layout (Option A).
- **Header:** Display the Hackathon `title`, `date`, and `goal`.
- **Tabs State:** Switch between "Projects" and "Judging Criteria".
- **Projects Tab:** Displays the existing grid of `ProjectCard` components.
- **Criteria Tab:** Renders the list of `criteria` and `bonus_criteria`, showing their descriptions, weights, and max scores.
- **ProjectCard Component:** Update the card to display `project.title` instead of `project.name` in the header.

### 3.2 ProjectDetail.tsx
- **Data Fetching:** Implement a new `useSWR` hook to fetch `/api/projects/:id` alongside the existing evaluations fetch.
- **Project Metadata Header:** Add a section at the top displaying the project's `team_name`, a clickable `url` (if present), and a clickable `github_url`.
- **Evaluation Breakdown Layout:** Update the rendering of the `CriteriaScore` array inside each evaluation to use the Inline List format (Option X).
- **Evaluation Details:** For each criterion, display the `name`, `score`, `weight`, and the new `reasoning` text in an italicized block beneath the score.

## 4. Error Handling and Edge Cases
- **Loading States:** Ensure both new API calls have appropriate loading indicators in the UI before rendering the data.
- **Missing Data:** Handle cases where optional fields (like `URL` or `reasoning`) might be empty or null gracefully in the UI.
- **API Failures:** Handle 404s or 500s from the new endpoints by displaying standard error messages.

## 5. Scope and Boundaries
- This design strictly focuses on data display. No changes will be made to the judging trigger logic or evaluation creation flow.
- Routing remains unchanged; we are only enhancing the components rendered at existing paths.
