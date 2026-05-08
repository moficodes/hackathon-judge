// src/pages/HackathonDetail.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { describe, it, expect, vi } from 'vitest';
import HackathonDetail from './HackathonDetail';
import useSWR from 'swr';

vi.mock('swr');

const mockUseSWR = useSWR as unknown as ReturnType<typeof vi.fn>;

describe('HackathonDetail Component', () => {
  it('renders loading state', () => {
    mockUseSWR.mockReturnValue({ data: undefined, error: undefined, isLoading: true });
    render(
      <MemoryRouter initialEntries={['/hackathons/1']}>
        <Routes>
          <Route path="/hackathons/:id" element={<HackathonDetail />} />
        </Routes>
      </MemoryRouter>
    );
    expect(screen.getByText('Loading projects...')).toBeInTheDocument();
  });

  it('renders projects data', async () => {
    const mockData = [
      { id: 'p1', name: 'AI Builder', team_name: 'Team Alpha', score: 95.5 },
    ];
    mockUseSWR.mockReturnValue({ data: mockData, error: undefined, isLoading: false });
    
    render(
      <MemoryRouter initialEntries={['/hackathons/1']}>
        <Routes>
          <Route path="/hackathons/:id" element={<HackathonDetail />} />
        </Routes>
      </MemoryRouter>
    );
    
    await waitFor(() => {
      expect(screen.getByText('AI Builder')).toBeInTheDocument();
      expect(screen.getByText('Team: Team Alpha')).toBeInTheDocument();
      expect(screen.getByText('Score: 95.5')).toBeInTheDocument();
    });
  });
});
