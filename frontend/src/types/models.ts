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
  max_score?: number;
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
