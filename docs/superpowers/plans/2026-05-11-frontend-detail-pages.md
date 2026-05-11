# Frontend Detail Pages Enhancement Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Enhance the Hackathon Detail and Project Detail pages with richer metadata and detailed evaluation reasoning, using a tabbed layout and inline criteria breakdown.

**Architecture:** We will add single-entity fetching endpoints to the backend (`/api/hackathons/:id` and `/api/projects/:id`). The frontend will use `swr` for parallel data fetching. Component rendering will be updated to display this new data without changing the existing routing structure.

**Tech Stack:** Go (Backend), React 19, Vite, TypeScript, Tailwind CSS, SWR.

---

### Task 1: Update Frontend Models

**Files:**
- Modify: `frontend/src/types/models.ts`

- [ ] **Step 1: Add new properties to interfaces**

Update the `Hackathon` interface and the `CriteriaScore` interface. We also need to define the `Criterion` type.

```typescript
// frontend/src/types/models.ts
export interface Criterion {
  name: string;
  weight: number;
  description: string;
  max_score: number;
}

export interface Hackathon {
  id: string;
  title: string;
  date: string;
  description: string;
  goal: string;
  status: string;
  criteria?: Criterion[];
  bonus_criteria?: Criterion[];
}

export interface Project {
  id: string;
  name: string;
  title: string;
  url: string;
  github_url: string;
  team_name: string;
  document: string;
  date: string;
  hackathon_id: string;
  score: number;
}

export interface CriteriaScore {
  name: string;
  score: number;
  weight: number;
  reasoning?: string;
}

export interface Evaluation {
  id: string;
  project_id: string;
  judge_id: string;
  status: string; // 'UNKNOWN' | 'RUNNING' | 'SUCCESS' | 'FAILED'
  criteria?: CriteriaScore[];
  total_score: number;
  comment?: string;
  created_at: string;
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/types/models.ts
git commit -m "feat(frontend): update models for detail pages"
```

### Task 2: Add Backend Service Methods

**Files:**
- Modify: `backend/internal/service/hackathon_service.go`

- [ ] **Step 1: Add GetHackathon and GetProject to HackathonService interface**

Open `backend/internal/service/hackathon_service.go` and update the interface.

```go
type HackathonService interface {
	ListHackathons() ([]domain.Hackathon, error)
	GetHackathon(id string) (domain.Hackathon, error)
	ListProjectsByHackathon(id string) ([]domain.Project, error)
	GetProject(id string) (domain.Project, error)
	AddEvaluation(eval domain.Evaluation) error
	ListEvaluationsByProject(projectID string) ([]domain.Evaluation, error)
	TriggerJudging(projectID string) (string, error)
}
```

- [ ] **Step 2: Implement GetHackathon and GetProject methods**

Add these methods to the `hackathonService` struct implementation.

```go
func (s *hackathonService) GetHackathon(id string) (domain.Hackathon, error) {
	return s.repo.GetByID(id)
}

func (s *hackathonService) GetProject(id string) (domain.Project, error) {
	return s.projectRepo.GetProjectByID(id)
}
```

- [ ] **Step 3: Run go build to verify**

```bash
cd backend && go build ./...
```
Expected: No errors.

- [ ] **Step 4: Commit**

```bash
git add backend/internal/service/hackathon_service.go
git commit -m "feat(backend): add GetHackathon and GetProject to service"
```

### Task 3: Add Backend Handler Endpoints

**Files:**
- Modify: `backend/internal/handler/hackathon_handler.go`

- [ ] **Step 1: Register new routes**

Update `RegisterRoutes`.

```go
func (h *HackathonHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/hackathons", h.GetHackathons)
		api.GET("/hackathons/:id", h.GetHackathon)
		api.GET("/hackathons/:id/projects", h.GetProjects)
		api.GET("/projects/:id", h.GetProject)
		api.GET("/projects/:id/evaluations", h.GetEvaluations)
		api.POST("/projects/:id/judge", h.TriggerJudging)
	}
}
```

- [ ] **Step 2: Implement GetHackathon handler**

```go
func (h *HackathonHandler) GetHackathon(c *gin.Context) {
	id := c.Param("id")
	res, err := h.svc.GetHackathon(id)
	if err != nil {
		log.Printf("[ERROR] GetHackathon failed for %s: %v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Hackathon not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}
```

- [ ] **Step 3: Implement GetProject handler**

```go
func (h *HackathonHandler) GetProject(c *gin.Context) {
	id := c.Param("id")
	res, err := h.svc.GetProject(id)
	if err != nil {
		log.Printf("[ERROR] GetProject failed for %s: %v\n", id, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}
```

- [ ] **Step 4: Run go test/build to verify**

```bash
cd backend && go build ./...
```
Expected: No errors.

- [ ] **Step 5: Commit**

```bash
git add backend/internal/handler/hackathon_handler.go
git commit -m "feat(backend): add GET /hackathons/:id and GET /projects/:id handlers"
```

### Task 4: Update HackathonDetail Component

**Files:**
- Modify: `frontend/src/pages/HackathonDetail.tsx`

- [ ] **Step 1: Replace HackathonDetail component with tabbed layout**

Update the file to fetch the hackathon data, manage tab state, and display criteria. Replace the *entire* `HackathonDetail` default export function (keep `ProjectCard` as is, but update its header to use `project.title`).

```tsx
// frontend/src/pages/HackathonDetail.tsx
import { useState, useEffect } from 'react';
import useSWR from 'swr';
import { useParams, Link } from 'react-router-dom';
import { fetcher } from '../utils/fetcher';
import type { Project, Evaluation, Hackathon } from '../types/models';

function ProjectCard({ project }: { project: Project }) {
  const { data: evaluations, mutate } = useSWR<Evaluation[]>(`/api/projects/${project.id}/evaluations`, fetcher);
  const [isTriggering, setIsTriggering] = useState(false);
  const [judgeMessage, setJudgeMessage] = useState<{ text: string, type: 'success' | 'error' } | null>(null);

  const isRunning = evaluations?.some(e => e.status === 'RUNNING');
  const hasEvaluations = evaluations && evaluations.length > 0;

  useEffect(() => {
    let interval: ReturnType<typeof setInterval>;
    if (isRunning) {
      interval = setInterval(() => {
        mutate();
      }, 5000);
    }
    return () => clearInterval(interval);
  }, [isRunning, mutate]);

  const handleJudge = async () => {
    if (hasEvaluations && !window.confirm('Evaluations already exist for this project. Are you sure you want to run another evaluation?')) {
      return;
    }

    setIsTriggering(true);
    setJudgeMessage(null);
    try {
      const response = await fetch(`/api/projects/${project.id}/judge`, { method: 'POST' });
      if (!response.ok) throw new Error('Failed to start judging');
      setJudgeMessage({ text: 'Judging task started!', type: 'success' });
      setTimeout(() => mutate(), 1000);
    } catch {
      setJudgeMessage({ text: 'Error starting judging', type: 'error' });
    } finally {
      setIsTriggering(false);
      setTimeout(() => setJudgeMessage(null), 3000);
    }
  };

  return (
    <div className="border border-slate-200 rounded-lg p-4 bg-white hover:border-blue-600 transition-colors">
      <h3 className="text-xl font-semibold mb-2">{project.title}</h3>
      <p className="text-gray-600 mb-1">Team: {project.team_name}</p>
      <p className="text-gray-600 mb-4 font-medium text-lg">Score: {project.score}</p>
      
      <div className="flex flex-col gap-2">
        <div className="flex gap-2">
          <Link 
            to={`/projects/${project.id}`}
            className="inline-block border border-blue-600 text-blue-600 px-4 py-2 rounded hover:bg-blue-50 transition-colors"
          >
            View Evaluations
          </Link>
          <button
            onClick={handleJudge}
            disabled={isTriggering || isRunning}
            className={`inline-block border px-4 py-2 rounded transition-colors disabled:opacity-50 disabled:cursor-not-allowed ${
              isRunning ? 'border-yellow-600 bg-yellow-600 text-white' : 
              'border-green-600 bg-green-600 text-white hover:bg-green-700'
            }`}
          >
            {isTriggering ? 'Starting...' : isRunning ? 'Judging in progress...' : hasEvaluations ? 'Rerun Judge' : 'Judge Project'}
          </button>
        </div>
        {judgeMessage && (
          <p className={`text-sm mt-1 ${judgeMessage.type === 'error' ? 'text-red-600' : 'text-green-600'}`}>
            {judgeMessage.text}
          </p>
        )}
      </div>
    </div>
  );
}

export default function HackathonDetail() {
  const { id } = useParams<{ id: string }>();
  const [activeTab, setActiveTab] = useState<'projects' | 'criteria'>('projects');
  
  const { data: hackathon, error: hackathonError, isLoading: isHackathonLoading } = useSWR<Hackathon>(id ? `/api/hackathons/${id}` : null, fetcher);
  const { data: projects, error: projectsError, isLoading: isProjectsLoading } = useSWR<Project[]>(id ? `/api/hackathons/${id}/projects` : null, fetcher);

  if (isHackathonLoading || isProjectsLoading) return <div className="p-4">Loading details...</div>;
  if (hackathonError) return <div className="p-4 text-red-500">Failed to load hackathon details.</div>;
  if (projectsError) return <div className="p-4 text-red-500">Failed to load projects.</div>;
  
  return (
    <div className="p-4">
      <div className="mb-6">
        <Link to="/dashboard" className="text-blue-600 hover:underline mb-4 inline-block">&larr; Back to Dashboard</Link>
        {hackathon && (
          <div className="bg-white p-6 rounded-lg border border-slate-200 mb-6 shadow-sm">
            <h2 className="text-3xl font-bold mb-2">{hackathon.title}</h2>
            <p className="text-gray-500 mb-4">{new Date(hackathon.date).toLocaleDateString()} &middot; Status: <span className="font-medium text-slate-700">{hackathon.status}</span></p>
            <div className="prose max-w-none text-gray-700">
              <p className="font-semibold text-lg mb-2">Goal:</p>
              <p>{hackathon.goal}</p>
            </div>
          </div>
        )}
      </div>

      <div className="mb-6 border-b border-gray-200">
        <nav className="-mb-px flex space-x-8">
          <button
            onClick={() => setActiveTab('projects')}
            className={`${
              activeTab === 'projects'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm`}
          >
            Projects
          </button>
          <button
            onClick={() => setActiveTab('criteria')}
            className={`${
              activeTab === 'criteria'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            } whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm`}
          >
            Judging Criteria
          </button>
        </nav>
      </div>

      {activeTab === 'projects' && (
        <>
          {!projects || projects.length === 0 ? (
            <div>No projects found for this hackathon.</div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {projects.map((p) => (
                <ProjectCard key={p.id} project={p} />
              ))}
            </div>
          )}
        </>
      )}

      {activeTab === 'criteria' && hackathon && (
        <div className="space-y-8">
          <section>
            <h3 className="text-xl font-bold mb-4">Standard Criteria</h3>
            {(!hackathon.criteria || hackathon.criteria.length === 0) ? (
              <p className="text-gray-500">No standard criteria defined.</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {hackathon.criteria.map((c, idx) => (
                  <div key={idx} className="bg-white p-4 border rounded-lg shadow-sm">
                    <div className="flex justify-between items-start mb-2">
                      <h4 className="font-bold text-lg text-slate-800">{c.name}</h4>
                      <span className="bg-blue-100 text-blue-800 text-xs px-2 py-1 rounded-full font-medium">Weight: {c.weight} &middot; Max: {c.max_score}</span>
                    </div>
                    <p className="text-gray-600 text-sm">{c.description}</p>
                  </div>
                ))}
              </div>
            )}
          </section>

          <section>
            <h3 className="text-xl font-bold mb-4 text-purple-700">Bonus Criteria</h3>
            {(!hackathon.bonus_criteria || hackathon.bonus_criteria.length === 0) ? (
              <p className="text-gray-500">No bonus criteria defined.</p>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {hackathon.bonus_criteria.map((c, idx) => (
                  <div key={idx} className="bg-purple-50 p-4 border border-purple-100 rounded-lg shadow-sm">
                    <div className="flex justify-between items-start mb-2">
                      <h4 className="font-bold text-lg text-purple-900">{c.name}</h4>
                      <span className="bg-purple-200 text-purple-900 text-xs px-2 py-1 rounded-full font-medium">Weight: {c.weight} &middot; Max: {c.max_score}</span>
                    </div>
                    <p className="text-purple-700 text-sm">{c.description}</p>
                  </div>
                ))}
              </div>
            )}
          </section>
        </div>
      )}
    </div>
  );
}
```

- [ ] **Step 2: Run frontend linter**

```bash
cd frontend && npm run lint
```
Expected: No errors related to `HackathonDetail.tsx`.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/pages/HackathonDetail.tsx
git commit -m "feat(frontend): enhance HackathonDetail with tabs and criteria"
```

### Task 5: Update ProjectDetail Component

**Files:**
- Modify: `frontend/src/pages/ProjectDetail.tsx`

- [ ] **Step 1: Replace ProjectDetail component**

Update the component to fetch project data and render the inline list layout for criteria reasoning.

```tsx
// src/pages/ProjectDetail.tsx
import useSWR from 'swr';
import { useParams, useNavigate } from 'react-router-dom';
import { fetcher } from '../utils/fetcher';
import type { Evaluation, Project } from '../types/models';

export default function ProjectDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  
  const { data: project, error: projectError, isLoading: isProjectLoading } = useSWR<Project>(
    id ? `/api/projects/${id}` : null,
    fetcher
  );

  const { data: evaluations, error: evalError, isLoading: isEvalLoading } = useSWR<Evaluation[]>(
    id ? `/api/projects/${id}/evaluations` : null, 
    fetcher
  );

  if (isProjectLoading || isEvalLoading) return <div className="p-4">Loading details...</div>;
  if (projectError) return <div className="p-4 text-red-500">Failed to load project details.</div>;
  if (evalError) return <div className="p-4 text-red-500">Failed to load evaluations.</div>;

  return (
    <div className="p-4">
      <div className="mb-6">
        <button 
          onClick={() => navigate(-1)} 
          className="text-blue-600 hover:underline mb-4 inline-block bg-transparent border-none cursor-pointer"
        >
          &larr; Back
        </button>
        
        {project && (
          <div className="bg-white p-6 rounded-lg border border-slate-200 mb-6 shadow-sm">
            <h2 className="text-3xl font-bold mb-2">{project.title}</h2>
            <div className="text-gray-600 space-y-1">
              <p><span className="font-semibold text-gray-700">Team:</span> {project.team_name}</p>
              {project.url && (
                <p><span className="font-semibold text-gray-700">Project URL:</span> <a href={project.url} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">{project.url}</a></p>
              )}
              {project.github_url && (
                <p><span className="font-semibold text-gray-700">GitHub:</span> <a href={project.github_url} target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:underline">{project.github_url}</a></p>
              )}
            </div>
          </div>
        )}
        
        <h2 className="text-2xl font-bold">Project Evaluations</h2>
      </div>

      {!evaluations || evaluations.length === 0 ? (
        <div>No evaluations found for this project.</div>
      ) : (
        <div className="space-y-6">
          {evaluations.map((e) => (
            <div key={e.id} className="border border-slate-200 rounded-lg p-6 bg-white shadow-sm">
              <div className="flex justify-between items-start mb-6 pb-4 border-b">
                <div>
                  <span className="text-xs text-gray-500 uppercase tracking-wider font-bold">Judge ID</span>
                  <p className="font-mono text-sm mt-1 bg-gray-100 px-2 py-1 rounded inline-block">{e.judge_id}</p>
                </div>
                <div>
                  <span className="text-xs text-gray-500 uppercase tracking-wider font-bold block mb-1">Status</span>
                  <p className="font-semibold">
                    {e.status === 'RUNNING' && <span className="bg-yellow-100 text-yellow-800 text-xs px-3 py-1.5 rounded-full">RUNNING</span>}
                    {e.status === 'SUCCESS' && <span className="bg-green-100 text-green-800 text-xs px-3 py-1.5 rounded-full">SUCCESS</span>}
                    {e.status === 'FAILED' && <span className="bg-red-100 text-red-800 text-xs px-3 py-1.5 rounded-full">FAILED</span>}
                    {!['RUNNING', 'SUCCESS', 'FAILED'].includes(e.status) && <span className="bg-gray-100 text-gray-800 text-xs px-3 py-1.5 rounded-full">{e.status || 'UNKNOWN'}</span>}
                  </p>
                </div>
                <div className="text-right">
                  <span className="text-xs text-gray-500 uppercase tracking-wider font-bold">Total Score</span>
                  <p className="text-3xl font-black text-blue-600 mt-1">{e.total_score}</p>
                </div>
              </div>
              
              {e.criteria && e.criteria.length > 0 && (
                <div className="mb-6">
                  <h4 className="text-lg font-bold mb-4 text-slate-800">Criteria Breakdown</h4>
                  <div className="space-y-4">
                    {e.criteria.map((c, idx) => (
                      <div key={idx} className="bg-gray-50 p-4 rounded-lg border border-gray-100">
                        <div className="flex justify-between items-center mb-2 pb-2 border-b border-gray-200">
                          <span className="font-bold text-slate-700 text-lg">{c.name}</span>
                          <div className="text-right">
                            <span className="font-black text-blue-600 text-lg">{c.score}</span>
                            <span className="text-gray-400 text-xs ml-2">(Weight: {c.weight})</span>
                          </div>
                        </div>
                        {c.reasoning ? (
                          <p className="text-sm text-gray-600 italic leading-relaxed">"{c.reasoning}"</p>
                        ) : (
                          <p className="text-sm text-gray-400 italic">No detailed reasoning provided.</p>
                        )}
                      </div>
                    ))}
                  </div>
                </div>
              )}

              <div className="bg-blue-50 p-4 rounded-lg border border-blue-100">
                <h4 className="text-sm font-bold text-blue-900 mb-2 uppercase tracking-wide">Overall Comment</h4>
                <p className="text-blue-800 leading-relaxed">
                  {e.comment || "No overall comment provided."}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
```

- [ ] **Step 2: Run frontend linter**

```bash
cd frontend && npm run lint
```
Expected: No errors related to `ProjectDetail.tsx`.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/pages/ProjectDetail.tsx
git commit -m "feat(frontend): enhance ProjectDetail with metadata and inline reasoning"
```
