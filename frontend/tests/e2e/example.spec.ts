import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('/');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/frontend/);
});

test('navigation to about page', async ({ page }) => {
  await page.goto('/');

  // Click the about link.
  await page.getByRole('link', { name: 'About', exact: true }).click();

  // Expect the about page text.
  await expect(page.getByText('This is a tool for judging hackathons.')).toBeVisible();
});
