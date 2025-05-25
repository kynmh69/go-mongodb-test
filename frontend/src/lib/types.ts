export interface User {
  id: string;
  user_id: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface CreateUserRequest {
  user_id: string;
  email: string;
  password: string;
}

export interface UpdateUserRequest {
  user_id?: string;
  email?: string;
  password?: string;
}

export interface ApiResponse<T> {
  data?: T;
  error?: string;
}

export interface UsersListResponse {
  users: User[];
  count: number;
}