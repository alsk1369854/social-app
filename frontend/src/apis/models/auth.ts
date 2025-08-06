export interface UserLoginRequest {
  email: string;
  password: string;
}

export interface UserLoginResponse {
  accessToken: string;
  email: string;
  id: string;
  username: string;
}

export interface UserRegisterRequestAddress {
  cityID: string;
  street: string;
}

export interface UserRegisterRequest {
  email: string;
  password: string;
  username: string;
  address?: UserRegisterRequestAddress;
  age?: number;
}

export interface UserRegisterResponseAddress {
  cityID: string;
  street: string;
}

export interface UserRegisterResponse {
  address?: UserRegisterResponseAddress;
  age?: number;
  email: string;
  id: string;
  username: string;
}

export interface ErrorResponse {
  error: string;
}

export interface City {
  id: string;
  name: string;
}