// src/pages/ProjectDetail.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, it, expect, vi } from 'vitest';
import ProjectDetail from './ProjectDetail';
import useSWR from 'swr';

vi.mock('swr');

const mockUseSWR = useSWR as unknown as ReturnType<typeof vi.fn>;

describe('ProjectDetail Component', () => {
  it('renders loading state', () => {
    mockUseSWR.mockReturnValue({ data: undefined, error: undefined, isLoading: true });
    render(
      <MemoryRouter initialEntries={['/projects/p1']}>
        <Routes>
          <Route path="/projects/:id" element={<ProjectDetail />} />
        </Routes>
      </MemoryRouter>
    );
    expect(screen.getByText('Loading details...')).toBeInTheDocument();
  });

  it('renders evaluations data and project details', async () => {
    const mockEvaluations = [
      { id: 'e1', judge_id: 'j1', total_score: 85, comment: 'Good work', criteria: [] },
    ];
    const mockProject = {
      id: 'p1', title: 'Awesome App', team_name: 'Team Alpha', url: 'https://example.com'
    };

    mockUseSWR.mockImplementation((url) => {
      if (url?.includes('/evaluations')) {
        return { data: mockEvaluations, error: undefined, isLoading: false };
      }
      if (url?.includes('/projects/')) {
        return { data: mockProject, error: undefined, isLoading: false };
      }
      return { data: undefined, error: undefined, isLoading: true };
    });
    
    render(
      <MemoryRouter initialEntries={['/projects/p1']}>
        <Routes>
          <Route path="/projects/:id" element={<ProjectDetail />} />
        </Routes>
      </MemoryRouter>
    );
    
    await waitFor(() => {
      expect(screen.getByText('Awesome App')).toBeInTheDocument();
      expect(screen.getByText('Team Alpha')).toBeInTheDocument();
      expect(screen.getByText(/Total Score/)).toBeInTheDocument();
      expect(screen.getByText('Good work')).toBeInTheDocument();
    });
  });
});
