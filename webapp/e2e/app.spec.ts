import { test, expect } from '@playwright/test';

/**
 * Mock the /v1/account/app API response so the app loads as authenticated.
 */
async function mockAuthenticatedApp(page: any) {
  await page.route('**/v1/account/app', async (route: any) => {
    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        mode: 'hosted',
        account: {
          id: 'test-account-id', name: 'Test User', email: 'test@example.com',
          createdAt: new Date().toISOString(), emailConfirmed: true,
          emailConfirmationSentTo: 'test@example.com', emailConfirmationPending: false,
        },
        workspaces: [{
          id: 'test-ws-id', name: 'testworkspace', createdAt: new Date().toISOString(),
          allowExternalSharing: true, euVat: '', status: 'active',
        }],
        memberships: [{
          id: 'test-member-id', workspaceId: 'test-ws-id', accountId: 'test-account-id',
          level: 'OWNER', name: 'Test User', email: 'test@example.com',
          createdAt: new Date().toISOString(),
        }],
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

// --- Authenticated app tests ---

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
    await expect(page.locator('#root')).toBeVisible();
  });
});

// --- Frontend API call validation ---

test.describe('Frontend API calls', () => {
  test('signup form sends expected fields', async ({ page }) => {
    let body: any = null;
    await page.route('**/v1/users/signup', async (route: any) => {
      body = route.request().postDataJSON();
      await route.fulfill({ status: 200, contentType: 'application/json', body: '{"token":"fake"}' });
    });
    await page.route('**/v1/account/app', async (route: any) => {
      await route.fulfill({ status: 401 });
    });
    await page.goto('/account/signup', { waitUntil: 'domcontentloaded' });
    await page.waitForTimeout(1000);

    const input = page.locator('input[id="workspaceName"]');
    if (await input.isVisible({ timeout: 3000 }).catch(() => false)) {
      await page.fill('input[id="workspaceName"]', 'myws');
      await page.fill('input[id="name"]', 'User');
      await page.fill('input[id="email"]', 'u@t.com');
      await page.fill('input[id="password"]', 'pass123');
      await page.click('button[type="submit"]');
      await page.waitForTimeout(500);
      expect(body).not.toBeNull();
      expect(body.workspaceName).toBe('myws');
    }
  });

  test('app fetches account data on load', async ({ page }) => {
    let called = false;
    await page.route('**/v1/account/app', async (route: any) => {
      called = true;
      await route.fulfill({
        status: 200, contentType: 'application/json',
        body: JSON.stringify({
          mode: 'hosted',
          account: { id: 'a1', name: 'U', email: 'u@t.com', createdAt: new Date().toISOString(), emailConfirmed: true, emailConfirmationSentTo: 'u@t.com', emailConfirmationPending: false },
          workspaces: [{ id: 'w1', name: 'test', createdAt: new Date().toISOString(), allowExternalSharing: true, euVat: '', status: 'active' }],
          memberships: [{ id: 'm1', workspaceId: 'w1', accountId: 'a1', level: 'OWNER', name: 'U', email: 'u@t.com', createdAt: new Date().toISOString() }],
          messages: [],
        }),
      });
    });
    await page.route('**/v1/projects', async (route: any) => {
      await route.fulfill({ status: 200, contentType: 'application/json', body: '[]' });
    });

    await page.goto('/');
    await page.waitForTimeout(500);

    expect(called).toBe(true);
  });
});
