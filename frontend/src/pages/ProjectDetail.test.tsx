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
    expect(screen.getByText('Loading evaluations...')).toBeInTheDocument();
  });

  it('renders evaluations data', async () => {
    const mockData = [
      { id: 'e1', judge_id: 'j1', total_score: 85, comment: 'Good work', criteria: [] },
    ];
    mockUseSWR.mockReturnValue({ data: mockData, error: undefined, isLoading: false });
    
    render(
      <MemoryRouter initialEntries={['/projects/p1']}>
        <Routes>
          <Route path="/projects/:id" element={<ProjectDetail />} />
        </Routes>
      </MemoryRouter>
    );
    
    await waitFor(() => {
      expect(screen.getByText(/Total Score/)).toBeInTheDocument();
      expect(screen.getByText('Good work')).toBeInTheDocument();
    });
  });
});
