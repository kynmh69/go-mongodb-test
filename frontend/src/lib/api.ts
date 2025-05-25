import { User, CreateUserRequest, UpdateUserRequest, UsersListResponse } from './types';

const API_BASE_URL = '/api/v1';

class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message);
    this.name = 'ApiError';
  }
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  const config: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
    ...options,
  };

  try {
    const response = await fetch(url, config);
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({}));
      throw new ApiError(
        response.status,
        errorData.error || `HTTP ${response.status}: ${response.statusText}`
      );
    }

    return await response.json();
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }
    throw new ApiError(500, 'Network error occurred');
  }
}

export const userApi = {
  // Get all users
  async getUsers(): Promise<UsersListResponse> {
    return apiRequest<UsersListResponse>('/users');
  },

  // Get user by MongoDB ID
  async getUserById(id: string): Promise<User> {
    return apiRequest<User>(`/users/${id}`);
  },

  // Search user by user_id
  async getUserByUserId(userId: string): Promise<User> {
    return apiRequest<User>(`/users/search?user_id=${encodeURIComponent(userId)}`);
  },

  // Search user by email
  async getUserByEmail(email: string): Promise<User> {
    return apiRequest<User>(`/users/search/email?email=${encodeURIComponent(email)}`);
  },

  // Create new user
  async createUser(userData: CreateUserRequest): Promise<User> {
    return apiRequest<User>('/users', {
      method: 'POST',
      body: JSON.stringify(userData),
    });
  },

  // Update user
  async updateUser(id: string, userData: UpdateUserRequest): Promise<User> {
    return apiRequest<User>(`/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(userData),
    });
  },

  // Delete user
  async deleteUser(id: string): Promise<{ message: string }> {
    return apiRequest<{ message: string }>(`/users/${id}`, {
      method: 'DELETE',
    });
  },
};

export const healthApi = {
  async checkHealth(): Promise<{ status: string; message: string }> {
    return apiRequest<{ status: string; message: string }>('/health');
  },
};

export { ApiError };