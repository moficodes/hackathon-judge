import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import Sidebar from './Sidebar';
import { expect, test } from 'vitest';

test('renders brand and navigation links', () => {
  render(
    <BrowserRouter>
      <Sidebar />
    </BrowserRouter>
  );
  
  expect(screen.getByText(/Hackathon Judge/i)).toBeInTheDocument();
  expect(screen.getByText(/Home/i)).toBeInTheDocument();
  expect(screen.getByText(/Dashboard/i)).toBeInTheDocument();
  expect(screen.getByText(/About/i)).toBeInTheDocument();
});
