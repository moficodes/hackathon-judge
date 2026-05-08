import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import { MemoryRouter } from 'react-router-dom';
import App from './App';

describe('App Component', () => {
  it('renders the layout and home page', () => {
    render(
      <MemoryRouter initialEntries={['/']}>
        <App />
      </MemoryRouter>
    );
    expect(screen.getByText(/Welcome to Hackathon Judge/i)).toBeInTheDocument();
  });

  it('renders the about page', () => {
    render(
      <MemoryRouter initialEntries={['/about']}>
        <App />
      </MemoryRouter>
    );
    expect(screen.getByText(/This is a tool for judging hackathons/i)).toBeInTheDocument();
  });

  it('renders the dashboard page', () => {
    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <App />
      </MemoryRouter>
    );
    expect(screen.getByText(/Loading hackathons\.\.\./i)).toBeInTheDocument();
  });

  it('renders the hackathon detail page', () => {
    render(
      <MemoryRouter initialEntries={['/hackathons/1']}>
        <App />
      </MemoryRouter>
    );
    expect(screen.getByText(/Loading projects\.\.\./i)).toBeInTheDocument();
  });

  it('renders the project detail page', () => {
    render(
      <MemoryRouter initialEntries={['/projects/1']}>
        <App />
      </MemoryRouter>
    );
    expect(screen.getByText(/Loading evaluations\.\.\./i)).toBeInTheDocument();
  });
});
