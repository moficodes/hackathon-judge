import { render, screen, fireEvent } from '@testing-library/react';
import Header from './Header';
import { expect, test } from 'vitest';
import { BrowserRouter } from 'react-router-dom';

test('renders search, default status, and default user profile', () => {
  render(
    <BrowserRouter>
      <Header />
    </BrowserRouter>
  );
  
  expect(screen.getByPlaceholderText(/Search projects/i)).toBeInTheDocument();
  expect(screen.getByLabelText(/Search/i)).toBeInTheDocument();
  expect(screen.getByText(/ACTIVE JUDGE/i)).toBeInTheDocument();
  expect(screen.getByText(/Judge User/i)).toBeInTheDocument();
  expect(screen.getByText(/JD/i)).toBeInTheDocument();
});

test('opens notification dropdown on click', () => {
  render(
    <BrowserRouter>
      <Header />
    </BrowserRouter>
  );

  const notificationBtn = screen.getByLabelText(/Notifications/i);
  fireEvent.click(notificationBtn);

  expect(screen.getByText(/No new notifications/i)).toBeInTheDocument();
  expect(screen.getByText(/View all/i)).toBeInTheDocument();
});

test('renders with custom user and status', () => {
  const customUser = {
    name: "Jane Doe",
    initials: "JD",
    role: "Admin"
  };
  const customStatus = "Idle";

  render(
    <BrowserRouter>
      <Header user={customUser} status={customStatus} />
    </BrowserRouter>
  );

  expect(screen.getByText(/IDLE/i)).toBeInTheDocument();
  expect(screen.getByText(/Jane Doe/i)).toBeInTheDocument();
  expect(screen.getByText(/Admin/i)).toBeInTheDocument();
});
