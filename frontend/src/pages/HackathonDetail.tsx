import { useState, useEffect } from 'react';
import useSWR from 'swr';
import { useParams, Link } from 'react-router-dom';
import { fetcher } from '../utils/fetcher';
import type { Project, Evaluation } from '../types/models';

function ProjectCard({ project }: { project: Project }) {
  const { data: evaluations, mutate } = useSWR<Evaluation[]>(`/api/projects/${project.id}/evaluations`, fetcher);
  const [isTriggering, setIsTriggering] = useState(false);
  const [judgeMessage, setJudgeMessage] = useState<{ text: string, type: 'success' | 'error' } | null>(null);

  const isRunning = evaluations?.some(e => e.status === 'RUNNING');
  const hasEvaluations = evaluations && evaluations.length > 0;

  // Polling logic if running
  useEffect(() => {
    let interval: ReturnType<typeof setInterval>;
    if (isRunning) {
      interval = setInterval(() => {
        mutate(); // Re-fetch evaluations every 5 seconds
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
      setTimeout(() => mutate(), 1000); // refresh evaluations to show running status
    } catch {
      setJudgeMessage({ text: 'Error starting judging', type: 'error' });
    } finally {
      setIsTriggering(false);
      setTimeout(() => setJudgeMessage(null), 3000);
    }
  };

  return (
    <div className="border border-slate-200 rounded-lg p-4 bg-white hover:border-blue-600 transition-colors">
      <h3 className="text-xl font-semibold mb-2">{project.name}</h3>
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
            <ProjectCard key={p.id} project={p} />
          ))}
        </div>
      )}
    </div>
  );
}
