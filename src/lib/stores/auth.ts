import { writable, derived } from 'svelte/store';
import type { User, Role } from '$lib/types';
import { api } from '$lib/services/api';

const THEME_KEY = 'theme';

export const currentUser = writable<User | null>(null);
export const isAuthenticated = writable(false);
export const sidebarCollapsed = writable(false);
export const currentTheme = writable('light');

export const userRole = derived(currentUser, ($user) => $user?.role || null);

export function hasAccess(role: Role | Role[], userRole: Role | null): boolean {
  if (!userRole) return false;
  if (userRole === 'super_admin') return true;
  if (Array.isArray(role)) return role.includes(userRole);
  return role === userRole;
}

export function applyTheme(theme: string) {
  currentTheme.set(theme);
  if (typeof window !== 'undefined') {
    try {
      localStorage.setItem(THEME_KEY, theme);
      document.documentElement.setAttribute('data-theme', theme);
    } catch {}
  }
}

if (typeof window !== 'undefined') {
  try {
    const saved = localStorage.getItem(THEME_KEY);
    if (saved) {
      currentTheme.set(saved);
      document.documentElement.setAttribute('data-theme', saved);
    } else {
      document.documentElement.setAttribute('data-theme', 'light');
    }
  } catch {}
}

export async function login(email: string, password: string): Promise<{ success: boolean; error?: string }> {
  try {
    const response = await api.auth.login(email, password);
    if (response.success && response.data.token) {
      localStorage.setItem('token', response.data.token);
      currentUser.set(response.data.user);
      isAuthenticated.set(true);
      return { success: true };
    }
    return { success: false, error: 'Login failed' };
  } catch (err: any) {
    return { success: false, error: err.message || 'Login failed' };
  }
}

export async function loadUser() {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;
  if (!token) return;

  try {
    const response = await api.auth.getMe();
    if (response.success) {
      currentUser.set(response.data);
      isAuthenticated.set(true);
    }
  } catch {
    localStorage.removeItem('token');
    currentUser.set(null);
    isAuthenticated.set(false);
  }
}

export function logout() {
  localStorage.removeItem('token');
  currentUser.set(null);
  isAuthenticated.set(false);
}
