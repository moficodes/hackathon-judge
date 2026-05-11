// src/types/models.ts
export interface Hackathon {
  id: string;
  title: string;
  date: string;
  description: string;
  goal: string;
  status: string;
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
}

export interface Evaluation {
  id: string;
  project_id: string;
  judge_id: string;
  criteria: CriteriaScore[];
  total_score: number;
  comment: string;
  created_at: string;
}
