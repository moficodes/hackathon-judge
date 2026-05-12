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
                            <span className="font-black text-blue-600 text-lg">{c.score}{c.max_score ? ` / ${c.max_score}` : ""}</span>
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
