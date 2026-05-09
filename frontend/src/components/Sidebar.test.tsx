import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import Sidebar from './Sidebar';
import { expect, test } from 'vitest';

test('renders brand and navigation links with correct hrefs', () => {
  render(
    <BrowserRouter>
      <Sidebar />
    </BrowserRouter>
  );
  
  expect(screen.getByText(/Hackathon Judge/i)).toBeInTheDocument();
  
  const homeLink = screen.getByRole('link', { name: /home/i });
  const dashboardLink = screen.getByRole('link', { name: /dashboard/i });
  const aboutLink = screen.getByRole('link', { name: /about/i });

  expect(homeLink).toHaveAttribute('href', '/');
  expect(dashboardLink).toHaveAttribute('href', '/dashboard');
  expect(aboutLink).toHaveAttribute('href', '/about');
});
