import { render, screen, waitFor } from '@testing-library/react';
import { describe, it, expect, vi } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import App from './App';
import useSWR from 'swr';

vi.mock('swr');
const mockUseSWR = useSWR as unknown as ReturnType<typeof vi.fn>;

describe('App Component', () => {
  it('renders the layout and home page', () => {
    mockUseSWR.mockReturnValue({ data: undefined, error: undefined, isLoading: true });
    render(
      <MemoryRouter initialEntries={['/']}>
        <App />
      </MemoryRouter>
    );
    expect(screen.getByText(/Welcome to Hackathon Judge/i)).toBeInTheDocument();
  });

  it('renders the about page', () => {
    mockUseSWR.mockReturnValue({ data: undefined, error: undefined, isLoading: true });
    render(
      <MemoryRouter initialEntries={['/about']}>
        <App />
      </MemoryRouter>
    );
    expect(screen.getByText(/This is a tool for judging hackathons/i)).toBeInTheDocument();
  });

  it('renders the dashboard page', () => {
    mockUseSWR.mockReturnValue({ data: undefined, error: undefined, isLoading: true });
    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <App />
      </MemoryRouter>
    );
    expect(screen.getByText(/Loading hackathons\.\.\./i)).toBeInTheDocument();
  });

  it('renders the hackathon detail page', async () => {
    mockUseSWR.mockReturnValue({ data: undefined, error: undefined, isLoading: true });
    render(
      <MemoryRouter initialEntries={['/hackathons/1']}>
        <App />
      </MemoryRouter>
    );
    await waitFor(() => {
        expect(screen.getByText(/Loading details\.\.\./i)).toBeInTheDocument();
    });
  });

  it('renders the project detail page', async () => {
    mockUseSWR.mockReturnValue({ data: undefined, error: undefined, isLoading: true });
    render(
      <MemoryRouter initialEntries={['/projects/1']}>
        <App />
      </MemoryRouter>
    );
    await waitFor(() => {
        expect(screen.getByText(/Loading details\.\.\./i)).toBeInTheDocument();
    });
  });
});
