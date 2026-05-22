import { test, expect } from '@playwright/test';

test.describe('Hackathon Judging Flow', () => {
  test('should navigate from home to project details', async ({ page }) => {
    // 1. Start at home
    await page.goto('/');
    await expect(page).toHaveTitle(/Hackathon Judge/);
    
    // 2. Go to Dashboard
    await page.click('text=Go to Dashboard');
    await expect(page.url()).toContain('/dashboard');
    
    // Wait for mock data to load (assuming there's a hackathon title visible)
    // The exact text will depend on the mocked data or backend state if running against real API
    // We assume the page loads a hackathon
    await page.waitForSelector('text=Hackathons');
    
    // 3. Click first View button in the dashboard list
    // This relies on the "View" link existing for a hackathon.
    const viewButtons = page.locator('text=View').first();
    await viewButtons.click();
    
    // 4. We should be on the Hackathon Detail page
    await expect(page.url()).toContain('/hackathons/');
    
    // Wait for projects tab
    await page.waitForSelector('text=Projects');

    // 5. Click the first View button for a project
    // Since there are multiple view buttons (some for hackathons, some for projects if we are on the detail page)
    // we need to be careful. In the HackathonDetail page, project links say "View".
    const projectViewButtons = page.locator('text=View').first();
    await projectViewButtons.click();

    // 6. We should be on Project Detail page
    await expect(page.url()).toContain('/projects/');
    await expect(page.locator('text=Back to Hackathon')).toBeVisible();
  });
});
