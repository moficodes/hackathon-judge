import { render, screen } from '@testing-library/react';
import Header from './Header';
import { expect, test } from 'vitest';

test('renders search, status, and user profile', () => {
  render(<Header />);
  
  expect(screen.getByPlaceholderText(/Search projects/i)).toBeInTheDocument();
  expect(screen.getByText(/ACTIVE JUDGE/i)).toBeInTheDocument();
  expect(screen.getByText(/Judge User/i)).toBeInTheDocument();
  expect(screen.getByText(/JD/i)).toBeInTheDocument();
});
