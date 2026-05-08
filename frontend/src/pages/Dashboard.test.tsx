// src/pages/Dashboard.test.tsx
import { render, screen, waitFor } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { describe, it, expect, vi } from 'vitest';
import Dashboard from './Dashboard';
import useSWR from 'swr';

vi.mock('swr');

const mockUseSWR = useSWR as unknown as ReturnType<typeof vi.fn>;

describe('Dashboard Component', () => {
  it('renders loading state initially', () => {
    mockUseSWR.mockReturnValue({ data: undefined, error: undefined, isLoading: true });
    render(<BrowserRouter><Dashboard /></BrowserRouter>);
    expect(screen.getByText('Loading hackathons...')).toBeInTheDocument();
  });

  it('renders error state', () => {
    mockUseSWR.mockReturnValue({ data: undefined, error: new Error('Failed to load'), isLoading: false });
    render(<BrowserRouter><Dashboard /></BrowserRouter>);
    expect(screen.getByText('Failed to load hackathons.')).toBeInTheDocument();
  });

  it('renders hackathons data', async () => {
    const mockData = [
      { id: '1', title: 'Spring Hackathon', date: '2026-05-15T00:00:00Z', status: 'upcoming', goal: 'Build cool stuff' },
    ];
    mockUseSWR.mockReturnValue({ data: mockData, error: undefined, isLoading: false });
    
    render(<BrowserRouter><Dashboard /></BrowserRouter>);
    
    await waitFor(() => {
      expect(screen.getByText('Spring Hackathon')).toBeInTheDocument();
      expect(screen.getByText('Status: upcoming')).toBeInTheDocument();
    });
  });
});
