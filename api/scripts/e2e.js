// e2e.js
//
// This script runs a series of end-to-end tests against the running API.
// It simulates a user logging in and performing CRUD operations on the /users endpoint.
//
// Dependencies:
//   - axios: For making HTTP requests.
//
// Usage:
//   npm install
//   npm run test:e2e

const axios = require('axios');

// --- Configuration ---
const BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';

// Credentials for the seeded admin user.
const ADMIN_EMAIL = 'test@test.com';
const ADMIN_PASSWORD = 'password';

// --- Helper Functions ---
const log = (message) => {
  console.log(`[e2e] ${message}`);
};

const apiClient = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' },
});

// --- Test Execution ---
async function runTests() {
  let token;
  let newUserId;

  try {
    log(`Starting E2E tests against ${BASE_URL}...`);

    // 1. Login and get JWT token
    log('--- Testing Login Endpoint (POST /login) ---');
    const loginResponse = await apiClient.post('/login', {
      email: ADMIN_EMAIL,
      password: ADMIN_PASSWORD,
    });

    if (loginResponse.status !== 200 || !loginResponse.data.token) {
      throw new Error(`Login failed! Status: ${loginResponse.status}`);
    }
    token = loginResponse.data.token;
    // Set the token on the X-App-Token header for all subsequent requests
    apiClient.defaults.headers.common['X-App-Token'] = token;
    log('Login successful. Token acquired.');

    // 2. Create a new user
    log('--- Testing Create User (POST /api/v1/users) ---');
    const newUserEmail = `test-user-${Date.now()}@example.com`;
    const createResponse = await apiClient.post('/api/v1/users', {
      // The Go handler expects snake_case as defined by the `json` tags.
      first_name: 'E2E',
      last_name: 'Test',
      email: newUserEmail,
      password: 'a-secure-password-123',
      confirm_password: 'a-secure-password-123',
      company_id: 1, // Assumes a company with ID 1 exists from seeded data.
      address_id: 1, // Assumes an address with ID 1 exists from seeded data.
    });

    if (createResponse.status !== 201 || !createResponse.data.id) {
      throw new Error(`Create user failed! Status: ${createResponse.status}`);
    }
    newUserId = createResponse.data.id;
    log(`Create user successful. New user ID: ${newUserId}`);

    // 3. Get the newly created user
    log(`--- Testing Get User (GET /api/v1/users/${newUserId}) ---`);
    const getResponse = await apiClient.get(`/api/v1/users/${newUserId}`);
    if (getResponse.status !== 200) {
      throw new Error(`Get user failed! Status: ${getResponse.status}`);
    }
    log('Get user successful.');

    // 4. Find all users
    log('--- Testing Find Users (GET /api/v1/users) ---');
    const findResponse = await apiClient.get('/api/v1/users');
    if (findResponse.status !== 200) {
      throw new Error(`Find users failed! Status: ${findResponse.status}`);
    }
    log('Find users successful.');

    // 5. Update the user
    log(`--- Testing Update User (PUT /api/v1/users/:id) ---`);
    const updateResponse = await apiClient.put(`/api/v1/users/${newUserId}`, {
      // The update endpoint likely also expects snake_case for consistency.
      first_name: 'E2E-Updated',
    });
    if (updateResponse.status !== 200) {
      throw new Error(`Update user failed! Status: ${updateResponse.status}`);
    }
    log('Update user successful.');

  } catch (error) {
    console.error('\n❌ E2E Test Failed!');
    if (error.response) {
      console.error(`Status: ${error.response.status}`);
      console.error('Data:', JSON.stringify(error.response.data, null, 2));
    } else if (error.request) {
      console.error('Error: No response received from the server. Is it running?');
    } else {
      console.error('Error:', error.message);
    }
    process.exit(1);
  } finally {
    // 6. Delete the user (cleanup)
    if (token && newUserId) {
      log(`--- Testing Delete User (DELETE /api/v1/users/${newUserId}) ---`);
      try {
        const deleteResponse = await apiClient.delete(`/api/v1/users/${newUserId}`);
        if (deleteResponse.status === 204) {
          log('Delete user successful (cleanup).');
        } else {
          // Log a warning if cleanup fails, but don't fail the whole script
          console.error(`\n⚠️ E2E Cleanup Warning: Could not delete user ${newUserId}. Status: ${deleteResponse.status}`);
        }
      } catch (cleanupError) {
        console.error(`\n⚠️ E2E Cleanup Warning: Error during user deletion:`, cleanupError.message);
      }
    }
  }

  log('\n✅ All E2E tests passed successfully!');
  process.exit(0);
}

runTests();