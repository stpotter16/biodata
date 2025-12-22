import { test, expect } from '@playwright/test';

test.describe('Example test', () => {
  test('homepage loads', async ({ page }) => {
    // This test assumes your server is running on localhost:8080
    await page.goto('/');

    // Simple check - adjust based on what your homepage shows
    await expect(page).toHaveTitle(/biotrak/i);
  });
});
