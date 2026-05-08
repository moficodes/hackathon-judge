import { render, screen } from '@testing-library/react';
import { BrowserRouter } from 'react-router-dom';
import { describe, it, expect } from 'vitest';
import Layout from './Layout';

describe('Layout Component', () => {
  it('renders navigation links', () => {
    render(
      <BrowserRouter>
        <Layout />
      </BrowserRouter>
    );
    expect(screen.getByRole('link', { name: 'Home' })).toHaveAttribute('href', '/');
    expect(screen.getByRole('link', { name: 'Dashboard' })).toHaveAttribute('href', '/dashboard');
    expect(screen.getByRole('link', { name: 'About' })).toHaveAttribute('href', '/about');
  });
});
