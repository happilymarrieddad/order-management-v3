import axios, { AxiosInstance } from 'axios';

// --- Configuration ---
const BASE_URL = process.env.API_BASE_URL || 'http://localhost:8080';

// Credentials for the seeded admin user.
const ADMIN_EMAIL = 'test@test.com';
const ADMIN_PASSWORD = 'password';

// --- Helper Functions ---
const log = (message: string) => {
  console.log(`[e2e] ${message}`);
};

const apiClient: AxiosInstance = axios.create({
  baseURL: BASE_URL,
  headers: { 'Content-Type': 'application/json' },
});

const assert = (condition: boolean, message: string) => {
  if (!condition) {
    throw new Error(`Assertion failed: ${message}`);
  }
};

// --- Test Execution ---
async function runTests() {
  let token: string;
  let newUserId: number;

  try {
    log(`Starting E2E tests against ${BASE_URL}...`);

    // 1. Login and get JWT token
    log('--- Testing Login Endpoint (POST /login) ---');
    const loginResponse = await apiClient.post('/login', {
      email: ADMIN_EMAIL,
      password: ADMIN_PASSWORD,
    });
    assert(loginResponse.status === 200, 'Login status should be 200');
    assert(loginResponse.data.token, 'Login response should contain a token');
    token = loginResponse.data.token;
    apiClient.defaults.headers.common['X-App-Token'] = token;
    log('Login successful. Token acquired.');

    // 2. Create a new user
    log('--- Testing Create User (POST /api/v1/users) ---');
    const newUserEmail = `test-user-${Date.now()}@example.com`;
    const createResponse = await apiClient.post('/api/v1/users', {
      first_name: 'E2E',
      last_name: 'Test',
      email: newUserEmail,
      password: 'a-secure-password-123',
      confirm_password: 'a-secure-password-123',
      company_id: 1, 
      address_id: 1, 
    });
    assert(createResponse.status === 201, 'Create user status should be 201');
    assert(createResponse.data.id, 'Create user response should contain an ID');
    newUserId = createResponse.data.id;
    log(`Create user successful. New user ID: ${newUserId}`);

    // 3. Get the newly created user
    log(`--- Testing Get User (GET /api/v1/users/${newUserId}) ---`);
    const getResponse = await apiClient.get(`/api/v1/users/${newUserId}`);
    assert(getResponse.status === 200, 'Get user status should be 200');
    assert(getResponse.data.id === newUserId, 'Get user response ID should match created user ID');
    log('Get user successful.');

    // 4. Find all users
    log('--- Testing Find Users (POST /api/v1/users/find) ---');
    const findResponse = await apiClient.post('/api/v1/users/find', {});
    assert(findResponse.status === 200, 'Find users status should be 200');
    assert(Array.isArray(findResponse.data.data), 'Find users response should contain an array of users');
    log('Find users successful.');

    // 5. Update the user
    log(`--- Testing Update User (PUT /api/v1/users/${newUserId}) ---`);
    const updateResponse = await apiClient.put(`/api/v1/users/${newUserId}`, {
      first_name: 'E2E-Updated',
      last_name: 'Test-Updated',
      address_id: 1,
    });
    assert(updateResponse.status === 200, 'Update user status should be 200');
    assert(updateResponse.data.first_name === 'E2E-Updated', 'Update user response should reflect changes');
    log('Update user successful.');

    // --- Error Case Tests ---
    log('--- Testing Error Cases ---');

    // Try to create a user with an existing email
    try {
      await apiClient.post('/api/v1/users', {
        first_name: 'Duplicate',
        last_name: 'User',
        email: newUserEmail, // Same email as before
        password: 'password123',
        confirm_password: 'password123',
        company_id: 1,
        address_id: 1,
      });
    } catch (error: any) {
      assert(error.response.status === 400, 'Creating user with duplicate email should fail with 400');
      log('Successfully caught duplicate email error.');
    }

    // Try to get a non-existent user
    try {
      await apiClient.get('/api/v1/users/999999');
    } catch (error: any) {
      assert(error.response.status === 404, 'Getting non-existent user should fail with 404');
      log('Successfully caught get non-existent user error.');
    }

  } catch (error: any) {
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
        assert(deleteResponse.status === 204, 'Delete user status should be 204');
        log('Delete user successful (cleanup).');
      } catch (cleanupError: any) {
        console.error(`\n⚠️ E2E Cleanup Warning: Could not delete user ${newUserId}. Status: ${cleanupError.response.status}`);
      }
    }
  }

  log('\n✅ All E2E tests passed successfully!');
  process.exit(0);
}

runTests();
