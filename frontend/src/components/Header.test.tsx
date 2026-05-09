import { render, screen } from '@testing-library/react';
import Header from './Header';
import { expect, test } from 'vitest';

test('renders search, default status, and default user profile', () => {
  render(<Header />);
  
  expect(screen.getByPlaceholderText(/Search projects/i)).toBeInTheDocument();
  expect(screen.getByLabelText(/Search/i)).toBeInTheDocument();
  expect(screen.getByText(/ACTIVE JUDGE/i)).toBeInTheDocument();
  expect(screen.getByText(/Judge User/i)).toBeInTheDocument();
  expect(screen.getByText(/JD/i)).toBeInTheDocument();
});

test('renders with custom user and status', () => {
  const customUser = {
    name: "Jane Doe",
    initials: "JD",
    role: "Admin"
  };
  const customStatus = "Idle";

  render(<Header user={customUser} status={customStatus} />);

  expect(screen.getByText(/IDLE/i)).toBeInTheDocument();
  expect(screen.getByText(/Jane Doe/i)).toBeInTheDocument();
  expect(screen.getByText(/Admin/i)).toBeInTheDocument();
});
