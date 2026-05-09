import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import Layout from './Layout';
import { expect, test, vi } from 'vitest';

// Mock components to focus on Layout's structural role
vi.mock('./Sidebar', () => ({
  default: () => <div data-testid="sidebar">Sidebar</div>
}));
vi.mock('./Header', () => ({
  default: () => <div data-testid="header">Header</div>
}));

test('renders Sidebar, Header, and main content area', () => {
  render(
    <BrowserRouter>
      <Layout />
    </BrowserRouter>
  );
  
  expect(screen.getByTestId('sidebar')).toBeInTheDocument();
  expect(screen.getByTestId('header')).toBeInTheDocument();
  expect(screen.getByRole('main')).toBeInTheDocument();
});
