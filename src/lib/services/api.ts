const API_URL = typeof process !== 'undefined' && process.env.PUBLIC_API_URL
  ? process.env.PUBLIC_API_URL
  : 'http://localhost:8080/api';

function getHeaders() {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  };
  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }
  return headers;
}

async function request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  const url = `${API_URL}${endpoint}`;
  const response = await fetch(url, {
    ...options,
    headers: {
      ...getHeaders(),
      ...options.headers,
    },
  });

  const data = await response.json();

  if (!response.ok) {
    throw new Error(data.message || 'Request failed');
  }

  return data;
}

export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface PaginatedData<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export const api = {
  auth: {
    login: (email: string, password: string) =>
      request<ApiResponse<any>>('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password }),
      }),
    register: (data: any) =>
      request<ApiResponse<any>>('/auth/register', {
        method: 'POST',
        body: JSON.stringify(data),
      }),
    getMe: () =>
      request<ApiResponse<any>>('/auth/me'),
    logout: () =>
      request<ApiResponse<any>>('/auth/logout', { method: 'POST' }),
  },

  users: {
    getAll: (params?: { page?: number; limit?: number; search?: string }) => {
      const qs = new URLSearchParams();
      if (params?.page) qs.set('page', String(params.page));
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.search) qs.set('search', params.search);
      return request<ApiResponse<PaginatedData<any>>>(`/users?${qs}`);
    },
    getById: (id: string) => request<ApiResponse<any>>(`/users/${id}`),
    create: (data: any) => request<ApiResponse<any>>('/users', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: any) => request<ApiResponse<any>>(`/users/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) => request<ApiResponse<any>>(`/users/${id}`, { method: 'DELETE' }),
  },

  products: {
    getAll: (params?: { page?: number; limit?: number; search?: string; category_id?: string; brand_id?: string; type?: string }) => {
      const qs = new URLSearchParams();
      if (params?.page) qs.set('page', String(params.page));
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.search) qs.set('search', params.search);
      if (params?.category_id) qs.set('category_id', params.category_id);
      if (params?.brand_id) qs.set('brand_id', params.brand_id);
      if (params?.type) qs.set('type', params.type);
      return request<ApiResponse<PaginatedData<any>>>(`/products?${qs}`);
    },
    getLowStock: () => request<ApiResponse<any[]>>('/products/low-stock'),
    getById: (id: string) => request<ApiResponse<any>>(`/products/${id}`),
    create: (data: any) => request<ApiResponse<any>>('/products', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: any) => request<ApiResponse<any>>(`/products/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) => request<ApiResponse<any>>(`/products/${id}`, { method: 'DELETE' }),
  },

  categories: {
    getAll: () => request<ApiResponse<any[]>>('/categories'),
    getById: (id: string) => request<ApiResponse<any>>(`/categories/${id}`),
    create: (data: any) => request<ApiResponse<any>>('/categories', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: any) => request<ApiResponse<any>>(`/categories/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) => request<ApiResponse<any>>(`/categories/${id}`, { method: 'DELETE' }),
  },

  brands: {
    getAll: () => request<ApiResponse<any[]>>('/brands'),
    getById: (id: string) => request<ApiResponse<any>>(`/brands/${id}`),
    create: (data: any) => request<ApiResponse<any>>('/brands', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: any) => request<ApiResponse<any>>(`/brands/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) => request<ApiResponse<any>>(`/brands/${id}`, { method: 'DELETE' }),
  },

  suppliers: {
    getAll: () => request<ApiResponse<any[]>>('/suppliers'),
    getById: (id: string) => request<ApiResponse<any>>(`/suppliers/${id}`),
    create: (data: any) => request<ApiResponse<any>>('/suppliers', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: any) => request<ApiResponse<any>>(`/suppliers/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) => request<ApiResponse<any>>(`/suppliers/${id}`, { method: 'DELETE' }),
  },

  customers: {
    getAll: (params?: { page?: number; limit?: number; search?: string }) => {
      const qs = new URLSearchParams();
      if (params?.page) qs.set('page', String(params.page));
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.search) qs.set('search', params.search);
      return request<ApiResponse<PaginatedData<any>>>(`/customers?${qs}`);
    },
    getById: (id: string) => request<ApiResponse<any>>(`/customers/${id}`),
    create: (data: any) => request<ApiResponse<any>>('/customers', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: any) => request<ApiResponse<any>>(`/customers/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    delete: (id: string) => request<ApiResponse<any>>(`/customers/${id}`, { method: 'DELETE' }),
  },

  transactions: {
    getAll: (params?: { page?: number; limit?: number; search?: string; status?: string; date_from?: string; date_to?: string }) => {
      const qs = new URLSearchParams();
      if (params?.page) qs.set('page', String(params.page));
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.search) qs.set('search', params.search);
      if (params?.status) qs.set('status', params.status);
      if (params?.date_from) qs.set('date_from', params.date_from);
      if (params?.date_to) qs.set('date_to', params.date_to);
      return request<ApiResponse<PaginatedData<any>>>(`/transactions?${qs}`);
    },
    getStats: () => request<ApiResponse<any>>('/transactions/stats'),
    getById: (id: string) => request<ApiResponse<any>>(`/transactions/${id}`),
    create: (data: any) => request<ApiResponse<any>>('/transactions', { method: 'POST', body: JSON.stringify(data) }),
  },

  services: {
    getAll: (params?: { page?: number; limit?: number; search?: string; status?: string }) => {
      const qs = new URLSearchParams();
      if (params?.page) qs.set('page', String(params.page));
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.search) qs.set('search', params.search);
      if (params?.status) qs.set('status', params.status);
      return request<ApiResponse<PaginatedData<any>>>(`/services?${qs}`);
    },
    getById: (id: string) => request<ApiResponse<any>>(`/services/${id}`),
    create: (data: any) => request<ApiResponse<any>>('/services', { method: 'POST', body: JSON.stringify(data) }),
    update: (id: string, data: any) => request<ApiResponse<any>>(`/services/${id}`, { method: 'PUT', body: JSON.stringify(data) }),
    addNote: (id: string, note: string) => request<ApiResponse<any>>(`/services/${id}/notes`, { method: 'POST', body: JSON.stringify({ note }) }),
    delete: (id: string) => request<ApiResponse<any>>(`/services/${id}`, { method: 'DELETE' }),
  },

  receivables: {
    getAll: (params?: { page?: number; limit?: number; search?: string; status?: string }) => {
      const qs = new URLSearchParams();
      if (params?.page) qs.set('page', String(params.page));
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.search) qs.set('search', params.search);
      if (params?.status) qs.set('status', params.status);
      return request<ApiResponse<PaginatedData<any>>>(`/receivables?${qs}`);
    },
    getStats: () => request<ApiResponse<any>>('/receivables/stats'),
    getById: (id: string) => request<ApiResponse<any>>(`/receivables/${id}`),
    makePayment: (id: string, data: { amount: number; method: string; notes?: string }) =>
      request<ApiResponse<any>>(`/receivables/${id}/payment`, { method: 'POST', body: JSON.stringify(data) }),
  },

  inventory: {
    stockIn: (data: { product_id: string; quantity: number; notes?: string }) =>
      request<ApiResponse<any>>('/inventory/stock-in', { method: 'POST', body: JSON.stringify(data) }),
    stockOut: (data: { product_id: string; quantity: number; notes?: string }) =>
      request<ApiResponse<any>>('/inventory/stock-out', { method: 'POST', body: JSON.stringify(data) }),
    adjustment: (data: { product_id: string; new_stock: number; notes?: string }) =>
      request<ApiResponse<any>>('/inventory/adjustment', { method: 'POST', body: JSON.stringify(data) }),
    getMovements: (params?: { page?: number; limit?: number; product_id?: string; type?: string }) => {
      const qs = new URLSearchParams();
      if (params?.page) qs.set('page', String(params.page));
      if (params?.limit) qs.set('limit', String(params.limit));
      if (params?.product_id) qs.set('product_id', params.product_id);
      if (params?.type) qs.set('type', params.type);
      return request<ApiResponse<PaginatedData<any>>>(`/inventory/movements?${qs}`);
    },
  },

  reports: {
    sales: (params?: { date_from?: string; date_to?: string; group_by?: string }) => {
      const qs = new URLSearchParams();
      if (params?.date_from) qs.set('date_from', params.date_from);
      if (params?.date_to) qs.set('date_to', params.date_to);
      if (params?.group_by) qs.set('group_by', params.group_by);
      return request<ApiResponse<any>>(`/reports/sales?${qs}`);
    },
    products: () => request<ApiResponse<any>>('/reports/products'),
    services: () => request<ApiResponse<any>>('/reports/services'),
    financial: (params?: { date_from?: string; date_to?: string }) => {
      const qs = new URLSearchParams();
      if (params?.date_from) qs.set('date_from', params.date_from);
      if (params?.date_to) qs.set('date_to', params.date_to);
      return request<ApiResponse<any>>(`/reports/financial?${qs}`);
    },
  },

  settings: {
    get: () => request<ApiResponse<any>>('/settings'),
    update: (data: any) => request<ApiResponse<any>>('/settings', { method: 'PUT', body: JSON.stringify(data) }),
  },

  seed: () => request<ApiResponse<any>>('/seed', { method: 'POST' }),
};
