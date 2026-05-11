// src/pages/ProjectDetail.tsx
import useSWR from 'swr';
import { useParams, useNavigate } from 'react-router-dom';
import { fetcher } from '../utils/fetcher';
import type { Evaluation } from '../types/models';

export default function ProjectDetail() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { data: evaluations, error, isLoading } = useSWR<Evaluation[]>(
    id ? `/api/projects/${id}/evaluations` : null, 
    fetcher
  );

  if (isLoading) return <div className="p-4">Loading evaluations...</div>;
  if (error) return <div className="p-4 text-red-500">Failed to load evaluations.</div>;

  return (
    <div className="p-4">
      <div className="mb-6">
        <button 
          onClick={() => navigate(-1)} 
          className="text-blue-600 hover:underline mb-4 inline-block bg-transparent border-none cursor-pointer"
        >
          &larr; Back
        </button>
        <h2 className="text-2xl font-bold">Project Evaluations</h2>
      </div>

      {!evaluations || evaluations.length === 0 ? (
        <div>No evaluations found for this project.</div>
      ) : (
        <div className="space-y-4">
          {evaluations.map((e) => (
            <div key={e.id} className="border border-slate-200 rounded-lg p-5 bg-white">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <span className="text-sm text-gray-500 uppercase tracking-wide">Judge ID</span>
                  <p className="font-mono">{e.judge_id}</p>
                </div>
                <div>
                  <span className="text-sm text-gray-500 uppercase tracking-wide">Status</span>
                  <p className="font-semibold mt-1">
                    {e.status === 'RUNNING' && <span className="bg-yellow-100 text-yellow-800 text-xs px-2 py-1 rounded">RUNNING</span>}
                    {e.status === 'SUCCESS' && <span className="bg-green-100 text-green-800 text-xs px-2 py-1 rounded">SUCCESS</span>}
                    {e.status === 'FAILED' && <span className="bg-red-100 text-red-800 text-xs px-2 py-1 rounded">FAILED</span>}
                    {!['RUNNING', 'SUCCESS', 'FAILED'].includes(e.status) && <span className="bg-gray-100 text-gray-800 text-xs px-2 py-1 rounded">{e.status || 'UNKNOWN'}</span>}
                  </p>
                </div>
                <div className="text-right">
                  <span className="text-sm text-gray-500 uppercase tracking-wide">Total Score</span>
                  <p className="text-2xl font-bold text-blue-600">{e.total_score}</p>
                </div>
              </div>
              
              {e.criteria && e.criteria.length > 0 && (
                <div className="mb-4 bg-gray-50 p-3 rounded">
                  <h4 className="text-sm font-semibold mb-2">Criteria Breakdown</h4>
                  <ul className="space-y-1">
                    {e.criteria.map((c, idx) => (
                      <li key={idx} className="flex justify-between text-sm">
                        <span className="text-gray-700">{c.name}</span>
                        <span className="font-medium">{c.score} <span className="text-gray-400 text-xs">(w:{c.weight})</span></span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}

              <div>
                <h4 className="text-sm font-semibold text-gray-700 mb-1">Comment</h4>
                <p className="text-gray-800 italic bg-gray-50 p-3 rounded border-l-4 border-gray-300">
                  {e.comment || "No comment provided."}
                </p>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
