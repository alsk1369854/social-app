export interface User {
  id: string;
  username: string;
  email?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  address?: {
    city?: string;
    district?: string;
    detail?: string;
  };
}

export interface AuthResponse {
  user: User;
  token: string;
}