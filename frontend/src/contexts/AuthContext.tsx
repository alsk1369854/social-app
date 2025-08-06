import React, { createContext, useContext, useReducer, useEffect, ReactNode } from 'react';
import { User } from '../models/User';
import { UserLoginRequest, UserRegisterRequest } from '../apis/models/auth';
import AuthAPI from '../apis/auth';

interface AuthState {
  user: User | null;
  accessToken: string | null;
  isLoading: boolean;
  isAuthenticated: boolean;
}

type AuthAction = 
  | { type: 'LOGIN_START' }
  | { type: 'LOGIN_SUCCESS'; payload: { user: User; accessToken: string } }
  | { type: 'LOGIN_ERROR' }
  | { type: 'REGISTER_START' }
  | { type: 'REGISTER_SUCCESS'; payload: { user: User; accessToken: string } }
  | { type: 'REGISTER_ERROR' }
  | { type: 'LOGOUT' }
  | { type: 'RESTORE_AUTH'; payload: { user: User; accessToken: string } };

interface AuthContextType {
  state: AuthState;
  login: (request: UserLoginRequest) => Promise<void>;
  register: (request: UserRegisterRequest) => Promise<void>;
  logout: () => void;
}

const initialState: AuthState = {
  user: null,
  accessToken: null,
  isLoading: false,
  isAuthenticated: false,
};

const authReducer = (state: AuthState, action: AuthAction): AuthState => {
  switch (action.type) {
    case 'LOGIN_START':
    case 'REGISTER_START':
      return {
        ...state,
        isLoading: true,
      };
    case 'LOGIN_SUCCESS':
    case 'REGISTER_SUCCESS':
      return {
        ...state,
        user: action.payload.user,
        accessToken: action.payload.accessToken,
        isAuthenticated: true,
        isLoading: false,
      };
    case 'LOGIN_ERROR':
    case 'REGISTER_ERROR':
      return {
        ...state,
        isLoading: false,
      };
    case 'LOGOUT':
      return {
        ...initialState,
      };
    case 'RESTORE_AUTH':
      return {
        ...state,
        user: action.payload.user,
        accessToken: action.payload.accessToken,
        isAuthenticated: true,
      };
    default:
      return state;
  }
};

const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [state, dispatch] = useReducer(authReducer, initialState);

  // Restore auth from localStorage on app start
  useEffect(() => {
    const storedToken = localStorage.getItem('accessToken');
    const storedUser = localStorage.getItem('user');
    
    if (storedToken && storedUser) {
      try {
        const user = JSON.parse(storedUser);
        dispatch({
          type: 'RESTORE_AUTH',
          payload: { user, accessToken: storedToken }
        });
      } catch (error) {
        // Clear invalid stored data
        localStorage.removeItem('accessToken');
        localStorage.removeItem('user');
      }
    }
  }, []);

  const login = async (request: UserLoginRequest) => {
    dispatch({ type: 'LOGIN_START' });
    
    try {
      const response = await AuthAPI.login(request);
      const user: User = {
        id: response.id,
        username: response.username,
        email: response.email,
      };

      // Store in localStorage
      localStorage.setItem('accessToken', response.accessToken);
      localStorage.setItem('user', JSON.stringify(user));

      dispatch({
        type: 'LOGIN_SUCCESS',
        payload: {
          user,
          accessToken: response.accessToken
        }
      });
    } catch (error) {
      dispatch({ type: 'LOGIN_ERROR' });
      throw error;
    }
  };

  const register = async (request: UserRegisterRequest) => {
    dispatch({ type: 'REGISTER_START' });
    
    try {
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      const response = await AuthAPI.register(request);
      
      // After registration, we need to login to get the access token
      // Since registration doesn't return accessToken in the current API
      await login({ email: request.email, password: request.password });
      
    } catch (error) {
      dispatch({ type: 'REGISTER_ERROR' });
      throw error;
    }
  };

  const logout = () => {
    // Clear localStorage
    localStorage.removeItem('accessToken');
    localStorage.removeItem('user');
    
    dispatch({ type: 'LOGOUT' });
  };

  // OAuth2 extensibility - placeholder for future implementation
  // const loginWithOAuth = async (provider: 'google' | 'facebook' | 'github') => {
  //   // Future implementation
  // };

  return (
    <AuthContext.Provider 
      value={{
        state,
        login,
        register,
        logout,
        // loginWithOAuth, // Future OAuth2 support
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};