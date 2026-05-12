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
