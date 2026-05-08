// src/pages/HackathonDetail.tsx
import useSWR from 'swr';
import { useParams, Link } from 'react-router-dom';
import { fetcher } from '../utils/fetcher';
import type { Project } from '../types/models';

export default function HackathonDetail() {
  const { id } = useParams<{ id: string }>();
  const { data: projects, error, isLoading } = useSWR<Project[]>(id ? `/api/hackathons/${id}/projects` : null, fetcher);

  if (isLoading) return <div className="p-4">Loading projects...</div>;
  if (error) return <div className="p-4 text-red-500">Failed to load projects.</div>;
  
  return (
    <div className="p-4">
      <div className="mb-6">
        <Link to="/dashboard" className="text-blue-600 hover:underline mb-4 inline-block">&larr; Back to Dashboard</Link>
        <h2 className="text-2xl font-bold">Hackathon Projects</h2>
      </div>

      {!projects || projects.length === 0 ? (
        <div>No projects found for this hackathon.</div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {projects.map((p) => (
            <div key={p.id} className="border border-slate-200 rounded-lg p-4 bg-white hover:border-blue-600 transition-colors">
              <h3 className="text-xl font-semibold mb-2">{p.name}</h3>
              <p className="text-gray-600 mb-1">Team: {p.team_name}</p>
              <p className="text-gray-600 mb-4 font-medium text-lg">Score: {p.score}</p>
              <Link 
                to={`/projects/${p.id}`}
                className="inline-block border border-blue-600 text-blue-600 px-4 py-2 rounded hover:bg-blue-50 transition-colors"
              >
                View Evaluations
              </Link>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
