import React, { useState, useEffect } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import AuthAPI from '../apis/auth';
import { City } from '../apis/models/auth';

const RegisterPage: React.FC = () => {
  const [formData, setFormData] = useState({
    username: '',
    email: '',
    password: '',
    confirmPassword: '',
    age: '',
    cityID: '',
    street: '',
  });
  const [cities, setCities] = useState<City[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loadingCities, setLoadingCities] = useState(false);
  const { state, register } = useAuth();
  const navigate = useNavigate();

  // Redirect if already authenticated
  useEffect(() => {
    if (state.isAuthenticated) {
      navigate('/');
    }
  }, [state.isAuthenticated, navigate]);

  // Load cities on component mount
  useEffect(() => {
    const loadCities = async () => {
      setLoadingCities(true);
      try {
        const citiesData = await AuthAPI.getCities();
        setCities(citiesData);
      } catch (error) {
        console.error('Failed to load cities:', error);
      } finally {
        setLoadingCities(false);
      }
    };

    loadCities();
  }, []);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    // Clear error when user starts typing
    if (error) setError(null);
  };

  const validateForm = (): string | null => {
    // Required fields validation
    // if (!formData.username.trim()) return '請輸入使用者名稱';
    if (!formData.email.trim()) return '請輸入電子郵件';
    if (!formData.password) return '請輸入密碼';
    if (!formData.confirmPassword) return '請確認密碼';

    // Username validation
    // if (formData.username && formData.username.length < 2) return '使用者名稱至少需要2個字符';
    if (formData.username.length > 50) return '使用者名稱不能超過50個字符';

    // Email format validation
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(formData.email)) return '請輸入有效的電子郵件地址';

    // Password validation
    if (formData.password.length < 6) return '密碼至少需要6個字符';
    if (formData.password !== formData.confirmPassword) return '密碼確認不匹配';

    // Age validation (if provided)
    if (formData.age && (isNaN(Number(formData.age)) || Number(formData.age) < 1 || Number(formData.age) > 120)) {
      return '請輸入有效的年齡 (1-120)';
    }

    // Address validation (if cityID is provided, street is required)
    if (formData.cityID && !formData.street.trim()) {
      return '選擇城市後請輸入街道地址';
    }

    return null;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    const validationError = validateForm();
    if (validationError) {
      setError(validationError);
      return;
    }

    try {
      const registerRequest = {
        email: formData.email.trim(),
        password: formData.password,
        ...(formData.username && { username: formData.username.trim() }),
        ...(formData.age && { age: Number(formData.age) }),
        ...(formData.cityID && formData.street && {
          address: {
            cityID: formData.cityID,
            street: formData.street.trim()
          }
        })
      };

      await register(registerRequest);
      navigate('/');
    } catch (error) {
      setError(error instanceof Error ? error.message : '註冊失敗，請稍後再試');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        <div className="text-center">
          <h2 className="mt-6 text-3xl font-extrabold text-gray-900 dark:text-white">
            註冊新帳戶
          </h2>
          <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
            或{' '}
            <Link
              to="/login"
              className="font-medium text-indigo-600 hover:text-indigo-500 dark:text-indigo-400 dark:hover:text-indigo-300"
            >
              登入現有帳戶
            </Link>
          </p>
        </div>

        <form className="mt-8 space-y-6" onSubmit={handleSubmit}>
          <div className="space-y-4">


            {/* Email */}
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                電子郵件 <span className="text-red-500">*</span>
              </label>
              <input
                id="email"
                name="email"
                type="email"
                autoComplete="email"
                required
                value={formData.email}
                onChange={handleInputChange}
                className="relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-800 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="請輸入您的電子郵件"
                disabled={state.isLoading}
              />
            </div>

            {/* Password */}
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                密碼 <span className="text-red-500">*</span>
              </label>
              <input
                id="password"
                name="password"
                type="password"
                autoComplete="new-password"
                required
                value={formData.password}
                onChange={handleInputChange}
                className="relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-800 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="請輸入密碼（至少6個字符）"
                disabled={state.isLoading}
              />
            </div>

            {/* Confirm Password */}
            <div>
              <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                確認密碼 <span className="text-red-500">*</span>
              </label>
              <input
                id="confirmPassword"
                name="confirmPassword"
                type="password"
                autoComplete="new-password"
                required
                value={formData.confirmPassword}
                onChange={handleInputChange}
                className="relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-800 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="請再次輸入密碼"
                disabled={state.isLoading}
              />
            </div>

            {/* Username */}
            <div>
              <label htmlFor="username" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                使用者名稱 <span className="text-gray-500 text-xs">(選填)</span>
              </label>
              <input
                id="username"
                name="username"
                type="text"
                value={formData.username}
                onChange={handleInputChange}
                className="relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-800 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="請輸入使用者名稱"
                disabled={state.isLoading}
              />
            </div>

            {/* Age - Optional */}
            <div>
              <label htmlFor="age" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                年齡 <span className="text-gray-500 text-xs">(選填)</span>
              </label>
              <input
                id="age"
                name="age"
                type="number"
                min="1"
                max="120"
                value={formData.age}
                onChange={handleInputChange}
                className="relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-800 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                placeholder="請輸入您的年齡"
                disabled={state.isLoading}
              />
            </div>

            {/* City - Optional */}
            <div>
              <label htmlFor="cityID" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                所在城市 <span className="text-gray-500 text-xs">(選填)</span>
              </label>
              <select
                id="cityID"
                name="cityID"
                value={formData.cityID}
                onChange={handleInputChange}
                className="relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 text-gray-900 dark:text-white dark:bg-gray-800 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                disabled={state.isLoading || loadingCities}
              >
                <option value="">請選擇城市</option>
                {cities.map(city => (
                  <option key={city.id} value={city.id}>
                    {city.name}
                  </option>
                ))}
              </select>
              {loadingCities && (
                <p className="mt-1 text-xs text-gray-500">載入城市列表中...</p>
              )}
            </div>

            {/* Street - Required if city is selected */}
            {formData.cityID && (
              <div>
                <label htmlFor="street" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                  街道地址 <span className="text-red-500">*</span>
                </label>
                <input
                  id="street"
                  name="street"
                  type="text"
                  value={formData.street}
                  onChange={handleInputChange}
                  className="relative block w-full px-3 py-2 border border-gray-300 dark:border-gray-600 placeholder-gray-500 dark:placeholder-gray-400 text-gray-900 dark:text-white dark:bg-gray-800 rounded-md focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                  placeholder="請輸入詳細地址"
                  disabled={state.isLoading}
                />
              </div>
            )}
          </div>

          {error && (
            <div className="rounded-md bg-red-50 dark:bg-red-900/20 p-4">
              <div className="text-sm text-red-800 dark:text-red-200">
                {error}
              </div>
            </div>
          )}

          <div>
            <button
              type="submit"
              disabled={state.isLoading}
              className="group relative w-full flex justify-center py-2 px-4 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed dark:focus:ring-offset-gray-900"
            >
              {state.isLoading ? (
                <div className="flex items-center">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></div>
                  註冊中...
                </div>
              ) : (
                '註冊'
              )}
            </button>
          </div>

          {/* OAuth2 placeholder for future implementation */}
          <div className="mt-6">
            <div className="relative">
              <div className="absolute inset-0 flex items-center">
                <div className="w-full border-t border-gray-300 dark:border-gray-600" />
              </div>
              <div className="relative flex justify-center text-sm">
                <span className="px-2 bg-gray-50 dark:bg-gray-900 text-gray-500 dark:text-gray-400">或使用</span>
              </div>
            </div>

            <div className="mt-6 grid grid-cols-1 gap-3">
              {/* Future OAuth2 buttons will be added here */}
              <div className="text-center text-xs text-gray-500 dark:text-gray-400">
                更多註冊方式即將推出
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  );
};

export default RegisterPage;