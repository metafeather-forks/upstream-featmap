import { test, expect } from '@playwright/test';

/**
 * Mock the /v1/account/app API response so the app loads as authenticated
 * with a test workspace and account.
 */
async function mockAuthenticatedApp(page: any) {
  await page.route('**/v1/account/app', async (route: any) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        mode: 'hosted',
        account: {
          id: 'test-account-id',
          name: 'Test User',
          email: 'test@example.com',
          createdAt: new Date().toISOString(),
          emailConfirmed: true,
          emailConfirmationSentTo: 'test@example.com',
          emailConfirmationPending: false,
        },
        workspaces: [
          {
            id: 'test-ws-id',
            name: 'testworkspace',
            createdAt: new Date().toISOString(),
            allowExternalSharing: true,
            euVat: '',
            status: 'active',
          },
        ],
        memberships: [
          {
            id: 'test-member-id',
            workspaceId: 'test-ws-id',
            accountId: 'test-account-id',
            level: 'OWNER',
            name: 'Test User',
            email: 'test@example.com',
            createdAt: new Date().toISOString(),
          },
        ],
        messages: [],
      }),
    });
  });
}

async function mockProjectData(page: any) {
  await page.route('**/v1/projects', async (route: any) => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: '[]' });
  });
  await page.route('**/v1/members', async (route: any) => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: '[]' });
  });
  await page.route('**/v1/invites', async (route: any) => {
    await route.fulfill({ status: 200, contentType: 'application/json', body: '[]' });
  });
}

// --- Authenticated app tests (mocked API) ---

test.describe('Authenticated app renders', () => {
  test('root redirects to workspace when authenticated', async ({ page }) => {
    await mockAuthenticatedApp(page);
    await mockProjectData(page);
    await page.goto('/');
    await expect(page).toHaveURL(/\/testworkspace/, { timeout: 10000 });
  });

  test('workspace page renders without crash', async ({ page }) => {
    await mockAuthenticatedApp(page);
    await mockProjectData(page);
    await page.goto('/testworkspace');
    await expect(page.locator('#root')).not.toBeEmpty({ timeout: 10000 });
  });

  test('settings page loads successfully', async ({ page }) => {
    await mockAuthenticatedApp(page);
    await mockProjectData(page);
    await page.goto('/testworkspace/settings');
    await expect(page.locator('#root')).not.toBeEmpty({ timeout: 10000 });
  });

  test('document has correct title', async ({ page }) => {
    await page.goto('/');
    await expect(page).toHaveTitle(/Featmap/);
  });
});

// --- Error resilience tests ---

test.describe('Error resilience', () => {
  test('app does not crash on 401 API response', async ({ page }) => {
    await page.route('**/v1/account/app', async (route: any) => {
      await route.fulfill({ status: 401 });
    });
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('#root')).not.toBeEmpty({ timeout: 5000 });
  });

  test('app does not crash on network failure', async ({ page }) => {
    await page.route('**/v1/account/app', async (route: any) => {
      await route.abort('connectionrefused');
    });
    await page.goto('/');
    await page.waitForLoadState('networkidle');
    // App should not white-screen even with total API failure
    await expect(page.locator('#root')).toBeVisible();
  });
});

// --- Backend API tests (requires Go server on port 5000) ---

test.describe('Backend API', () => {
  test.skip('unauthenticated routes return 401/400', async ({ request }) => {
    const API = 'http://localhost:5000/v1';
    const resp = await request.get(`${API}/account/app`);
    expect(resp.status()).toBeGreaterThanOrEqual(400);
  });

  test.skip('POST /users/login validates credentials', async ({ request }) => {
    const API = 'http://localhost:5000/v1';
    const resp = await request.post(`${API}/users/login`, {
      data: { email: 'noone@example.com', password: 'wrong' },
    });
    expect(resp.status()).toBeGreaterThanOrEqual(400);
  });

  test.skip('POST /users/signup validates input', async ({ request }) => {
    const API = 'http://localhost:5000/v1';
    const resp = await request.post(`${API}/users/signup`, {
      data: { workspaceName: '', name: '', email: '', password: '' },
    });
    expect(resp.status()).toBeGreaterThanOrEqual(400);
  });
});
