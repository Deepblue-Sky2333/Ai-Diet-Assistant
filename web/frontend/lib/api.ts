// API 基础 URL 配置
// 生产模式：使用相对路径（前后端集成在同一端口）
// 开发模式：使用完整 URL（前后端分离）
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 
  (typeof window !== 'undefined' && window.location.hostname !== 'localhost'
    ? '/api/v1'  // 生产环境使用相对路径
    : 'http://localhost:9090/api/v1'); // 开发环境使用完整 URL

const DEMO_MODE = typeof window !== 'undefined' 
  ? process.env.NEXT_PUBLIC_DEMO_MODE === 'true'
  : false;

console.log('[v0] API Client initialized - Demo Mode:', DEMO_MODE, 'Base URL:', API_BASE_URL);

interface ApiResponse<T> {
  code: number;
  message: string;
  data?: T;
  error?: string;
  timestamp?: number;
  pagination?: {
    page: number;
    page_size: number;
    total: number;
    total_pages: number;
  };
}

class ApiClient {
  private baseUrl: string;
  private demoMode: boolean;
  private refreshPromise: Promise<boolean> | null = null;
  private readonly defaultTimeout = 30000; // 30秒默认超时

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
    this.demoMode = typeof window !== 'undefined' 
      ? process.env.NEXT_PUBLIC_DEMO_MODE === 'true'
      : false;
    console.log('[v0] ApiClient constructor - Demo mode set to:', this.demoMode);
  }

  /**
   * 创建带超时的fetch请求
   * @param url 请求URL
   * @param options fetch选项
   * @param timeout 超时时间（毫秒），默认30秒
   * @returns Promise<Response>
   */
  private async fetchWithTimeout(
    url: string,
    options: RequestInit = {},
    timeout: number = this.defaultTimeout
  ): Promise<Response> {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);

    try {
      const response = await fetch(url, {
        ...options,
        signal: controller.signal,
      });
      clearTimeout(timeoutId);
      return response;
    } catch (error) {
      clearTimeout(timeoutId);
      
      // 检查是否是超时错误
      if (error instanceof Error && error.name === 'AbortError') {
        throw new Error(`Request timeout after ${timeout}ms`);
      }
      
      throw error;
    }
  }

  private getAuthHeaders(): HeadersInit {
    if (typeof window === 'undefined') {
      return { 'Content-Type': 'application/json' };
    }
    
    const token = localStorage.getItem('access_token');
    console.log('[v0] Getting auth headers, token exists:', !!token);
    return {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    };
  }

  private async handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
    console.log('[v0] Response status:', response.status);
    const contentType = response.headers.get('content-type');
    if (contentType && !contentType.includes('application/json')) {
      console.error('[v0] Response is not JSON, content-type:', contentType);
      throw new Error('Backend returned non-JSON response. Please check your API configuration.');
    }
    
    const result: ApiResponse<T> = await response.json();
    console.log('[v0] API response:', { code: result.code, message: result.message });
    
    if (result.code === 40101) {
      // Token expired, try refresh
      console.log('[v0] Token expired, attempting refresh...');
      const refreshed = await this.refreshToken();
      if (!refreshed) {
        // Refresh failed, redirect to login
        console.log('[v0] Token refresh failed, clearing storage and redirecting');
        if (typeof window !== 'undefined') {
          localStorage.clear();
          window.location.href = '/login';
        }
        throw new Error('Authentication required');
      }
      throw new Error('TOKEN_REFRESH_RETRY');
    }
    
    return result;
  }

  private async handleDemoRequest<T>(endpoint: string, options: RequestInit = {}): Promise<ApiResponse<T>> {
    await new Promise(resolve => setTimeout(resolve, 500));

    if (endpoint === '/auth/login') {
      if (typeof window !== 'undefined') {
        localStorage.setItem('access_token', 'demo_token');
        localStorage.setItem('refresh_token', 'demo_refresh_token');
        localStorage.setItem('demo_mode', 'true');
      }
      return {
        code: 0,
        message: 'Success',
        data: {
          access_token: 'demo_token',
          refresh_token: 'demo_refresh_token',
          expires_in: 3600,
        } as T,
      };
    }

    if (endpoint === '/dashboard') {
      return {
        code: 0,
        message: 'Success',
        data: {
          today_nutrition: {
            date: new Date().toISOString().split('T')[0],
            total_calories: 1850,
            total_protein: 78,
            total_fat: 62,
            total_carbohydrates: 210,
            total_fiber: 28,
          },
          goals: {
            daily_calories: 2000,
            daily_protein: 150,
            daily_fat: 65,
            daily_carbohydrates: 250,
          },
          recent_meals: [
            { id: 1, meal_type: 'breakfast', time: '08:00', total_calories: 450 },
            { id: 2, meal_type: 'lunch', time: '12:30', total_calories: 650 },
            { id: 3, meal_type: 'snack', time: '15:00', total_calories: 200 },
          ],
        } as T,
      };
    }

    if (endpoint.startsWith('/foods')) {
      if (options.method === 'POST' || options.method === 'PUT' || options.method === 'DELETE') {
        return { code: 0, message: 'Success', data: {} as T };
      }
      return {
        code: 0,
        message: 'Success',
        data: [] as T,
        pagination: { page: 1, page_size: 20, total: 0, total_pages: 0 },
      };
    }

    if (endpoint.startsWith('/meals')) {
      return {
        code: 0,
        message: 'Success',
        data: [] as T,
      };
    }

    if (endpoint.startsWith('/plans')) {
      return {
        code: 0,
        message: 'Success',
        data: [] as T,
      };
    }

    if (endpoint.startsWith('/nutrition')) {
      return {
        code: 0,
        message: 'Success',
        data: {} as T,
      };
    }

    if (endpoint.startsWith('/ai')) {
      return {
        code: 0,
        message: 'Success',
        data: { message: 'This is a demo response. Please configure your AI provider in settings.', message_id: Date.now() } as T,
      };
    }

    if (endpoint.startsWith('/settings')) {
      return {
        code: 0,
        message: 'Success',
        data: {} as T,
      };
    }

    return {
      code: 0,
      message: 'Success',
      data: {} as T,
    };
  }

  async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    console.log('[v0] Request called - Demo mode:', this.demoMode, 'Endpoint:', endpoint);
    
    if (this.demoMode) {
      console.log('[v0] Demo mode active, returning mock response');
      return this.handleDemoRequest<T>(endpoint, options);
    }

    console.log('[v0] Making real API call (demo mode is OFF)');
    try {
      const url = `${this.baseUrl}${endpoint}`;
      console.log('[v0] API request:', { url, method: options.method || 'GET' });
      
      const response = await this.fetchWithTimeout(url, {
        ...options,
        headers: {
          ...this.getAuthHeaders(),
          ...options.headers,
        },
        mode: 'cors',
        credentials: 'omit',
      });

      return await this.handleResponse<T>(response);
    } catch (error) {
      console.error('[v0] API request error:', error);
      
      // 处理超时错误
      if (error instanceof Error && error.message.startsWith('Request timeout')) {
        throw new Error('Request timeout. Please check your network connection and try again.');
      }
      
      if (error instanceof TypeError && error.message === 'Failed to fetch') {
        throw new Error('Cannot connect to backend API. Please ensure the backend server is running and NEXT_PUBLIC_API_URL is configured correctly.');
      }
      
      if ((error as Error).message === 'TOKEN_REFRESH_RETRY') {
        const response = await this.fetchWithTimeout(`${this.baseUrl}${endpoint}`, {
          ...options,
          headers: {
            ...this.getAuthHeaders(),
            ...options.headers,
          },
          mode: 'cors',
          credentials: 'omit',
        });
        return await this.handleResponse<T>(response);
      }
      throw error;
    }
  }

  async login(username: string, password: string) {
    console.log('[v0] Login method called, demo mode:', this.demoMode);
    const result = await this.request<{
      access_token: string;
      refresh_token: string;
      expires_in: number;
    }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    });

    if (result.code === 0 && result.data) {
      console.log('[v0] Login successful, storing tokens');
      if (typeof window !== 'undefined') {
        localStorage.setItem('access_token', result.data.access_token);
        localStorage.setItem('refresh_token', result.data.refresh_token);
      }
    }

    return result;
  }

  async loginWithTestAccount() {
    console.log('[v0] loginWithTestAccount called, demo mode:', this.demoMode);
    const result = await this.request<{
      access_token: string;
      refresh_token: string;
      expires_in: number;
    }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username: 'test', password: '114514' }),
    });

    if (result.code === 0 && result.data) {
      console.log('[v0] Test login successful, storing tokens');
      if (typeof window !== 'undefined') {
        localStorage.setItem('access_token', result.data.access_token);
        localStorage.setItem('refresh_token', result.data.refresh_token);
      }
    }

    return result;
  }

  async refreshToken(): Promise<boolean> {
    // 如果已经有刷新请求在进行中，返回同一个Promise
    if (this.refreshPromise) {
      console.log('[v0] Token refresh already in progress, reusing existing promise');
      return this.refreshPromise;
    }

    // 创建新的刷新Promise并缓存
    this.refreshPromise = (async () => {
      try {
        if (typeof window === 'undefined') return false;
        
        const refreshToken = localStorage.getItem('refresh_token');
        if (!refreshToken) {
          console.log('[v0] No refresh token found');
          return false;
        }

        console.log('[v0] Attempting token refresh');
        const response = await this.fetchWithTimeout(`${this.baseUrl}/auth/refresh`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ refresh_token: refreshToken }),
        });

        const result = await response.json();

        if (result.code === 0 && result.data) {
          console.log('[v0] Token refresh successful');
          localStorage.setItem('access_token', result.data.access_token);
          return true;
        }
        console.log('[v0] Token refresh failed');
        return false;
      } catch (error) {
        console.error('[v0] Token refresh error:', error);
        return false;
      } finally {
        // 刷新完成后清除缓存，允许下次刷新
        console.log('[v0] Clearing refresh promise cache');
        this.refreshPromise = null;
      }
    })();

    return this.refreshPromise;
  }

  async logout() {
    console.log('[v0] Logging out');
    await this.request('/auth/logout', { method: 'POST' });
    if (typeof window !== 'undefined') {
      localStorage.clear();
    }
  }

  // Auth methods
  async changePassword(oldPassword: string, newPassword: string) {
    return this.request('/auth/password', {
      method: 'PUT',
      body: JSON.stringify({ old_password: oldPassword, new_password: newPassword }),
    });
  }

  // Food methods
  async getFoods(page = 1, pageSize = 20, category = '') {
    const params = new URLSearchParams({ page: String(page), page_size: String(pageSize) });
    if (category) params.append('category', category);
    return this.request(`/foods?${params.toString()}`);
  }

  async getFood(id: number) {
    return this.request(`/foods/${id}`);
  }

  async createFood(food: any) {
    return this.request('/foods', {
      method: 'POST',
      body: JSON.stringify(food),
    });
  }

  async updateFood(id: number, food: any) {
    return this.request(`/foods/${id}`, {
      method: 'PUT',
      body: JSON.stringify(food),
    });
  }

  async deleteFood(id: number) {
    return this.request(`/foods/${id}`, { method: 'DELETE' });
  }

  async batchImportFoods(foods: any[]) {
    return this.request('/foods/batch', {
      method: 'POST',
      body: JSON.stringify({ foods }),
    });
  }

  // Meal methods
  async getMeals(date?: string) {
    const params = date ? `?date=${date}` : '';
    return this.request(`/meals${params}`);
  }

  async createMeal(meal: any) {
    return this.request('/meals', {
      method: 'POST',
      body: JSON.stringify(meal),
    });
  }

  async updateMeal(id: number, meal: any) {
    return this.request(`/meals/${id}`, {
      method: 'PUT',
      body: JSON.stringify(meal),
    });
  }

  async deleteMeal(id: number) {
    return this.request(`/meals/${id}`, { method: 'DELETE' });
  }

  // Plan methods
  async generatePlan(days: number, preferences: string) {
    return this.request('/plans/generate', {
      method: 'POST',
      body: JSON.stringify({ days, preferences }),
    });
  }

  async getPlans() {
    return this.request('/plans');
  }

  async getPlan(id: number) {
    return this.request(`/plans/${id}`);
  }

  async updatePlan(id: number, plan: any) {
    return this.request(`/plans/${id}`, {
      method: 'PUT',
      body: JSON.stringify(plan),
    });
  }

  async deletePlan(id: number) {
    return this.request(`/plans/${id}`, { method: 'DELETE' });
  }

  async completePlan(id: number) {
    return this.request(`/plans/${id}/complete`, { method: 'POST' });
  }

  // Nutrition methods
  async getDailyNutrition(date: string) {
    return this.request(`/nutrition/daily/${date}`);
  }

  async getMonthlyNutrition(year: number, month: number) {
    return this.request(`/nutrition/monthly?year=${year}&month=${month}`);
  }

  async compareNutrition(startDate: string, endDate: string) {
    return this.request(`/nutrition/compare?start_date=${startDate}&end_date=${endDate}`);
  }

  // AI methods
  async chat(message: string) {
    return this.request('/ai/chat', {
      method: 'POST',
      body: JSON.stringify({ message }),
    });
  }

  async getChatHistory(page = 1, pageSize = 20) {
    return this.request(`/ai/history?page=${page}&page_size=${pageSize}`);
  }

  // Dashboard
  async getDashboard() {
    return this.request('/dashboard');
  }

  // Settings
  async getSettings() {
    return this.request('/settings');
  }

  async updateAISettings(provider: string, apiKey?: string, apiEndpoint?: string) {
    // 只在提供新密钥时才发送
    const body: any = { provider };
    if (apiKey) {
      body.api_key = apiKey;
    }
    if (apiEndpoint) {
      body.api_endpoint = apiEndpoint;
    }
    
    return this.request('/settings/ai', {
      method: 'PUT',
      body: JSON.stringify(body),
    });
  }

  async testAIConnection() {
    return this.request('/settings/ai/test');
  }

  async getUserProfile() {
    return this.request('/user/profile');
  }

  async updateUserPreferences(preferences: any) {
    return this.request('/user/preferences', {
      method: 'PUT',
      body: JSON.stringify(preferences),
    });
  }
}

export const apiClient = new ApiClient(API_BASE_URL);
export type { ApiResponse };
