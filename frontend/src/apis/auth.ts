import {
  UserLoginRequest,
  UserLoginResponse,
  UserRegisterRequest,
  UserRegisterResponse,
  ErrorResponse,
  City
} from './models/auth';

const API_BASE_URL = process.env.REACT_APP_API_BASE_URL || "";

class AuthAPI {
  static async login(request: UserLoginRequest): Promise<UserLoginResponse> {
    const response = await fetch(`${API_BASE_URL}/user/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const errorData: ErrorResponse = await response.json();
      throw new Error(errorData.error || 'Login failed');
    }

    return response.json();
  }

  static async register(request: UserRegisterRequest): Promise<UserRegisterResponse> {
    const response = await fetch(`${API_BASE_URL}/user/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const errorData: ErrorResponse = await response.json();
      throw new Error(errorData.error || 'Registration failed');
    }

    return response.json();
  }

  static async getCities(): Promise<City[]> {
    const response = await fetch(`${API_BASE_URL}/city/all`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
      },
    });

    if (!response.ok) {
      const errorData: ErrorResponse = await response.json();
      throw new Error(errorData.error || 'Failed to fetch cities');
    }

    return response.json();
  }
}

export default AuthAPI;