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
    expect(screen.getByText('Loading details...')).toBeInTheDocument();
  });

  it('renders projects data', async () => {
    mockUseSWR.mockImplementation((key) => {
      if (typeof key === 'string' && key.includes('/projects')) {
        if (key.endsWith('/evaluations')) {
           return { data: [], error: undefined, isLoading: false };
        }
        return {
          data: [
            { id: 'p1', title: 'AI Builder', team_name: 'Team Alpha', score: 95.5 },
          ],
          error: undefined,
          isLoading: false,
        };
      }
      if (typeof key === 'string' && key.includes('/hackathons/')) {
        return {
           data: { id: '1', title: 'Test Hackathon', date: '2026-05-11', status: 'ACTIVE', goal: 'Test Goal' },
           error: undefined,
           isLoading: false,
        };
      }
      return { data: undefined, error: undefined, isLoading: true };
    });
    
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
      expect(screen.getByText('Score: 95.50')).toBeInTheDocument();
    });
  });
});
