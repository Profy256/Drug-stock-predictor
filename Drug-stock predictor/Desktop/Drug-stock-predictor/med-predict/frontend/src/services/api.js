// ============================================================
// src/services/api.js — Centralized API client
// All requests attach Bearer token automatically
// ============================================================

const BASE = import.meta.env.VITE_API_URL || 'http://localhost:4000';

function getToken() {
  return localStorage.getItem('med_predict_token');
}

async function request(path, options = {}) {
  const token = getToken();
  const res = await fetch(`${BASE}/api/v1${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(options.headers || {}),
    },
    body: options.body ? JSON.stringify(options.body) : undefined,
  });

  if (res.status === 401) {
    // Token expired — clear session
    ['med_predict_token', 'med_predict_user'].forEach(k => localStorage.removeItem(k));
    window.location.href = '/'; // Force re-login
    return;
  }

  const data = await res.json().catch(() => ({}));
  if (!res.ok) throw new Error(data.error || `Request failed (${res.status})`);
  return data;
}

// ─── Auth ─────────────────────────────────────────────────────
export const api = {
  auth: {
    login: (email, password) =>
      request('/auth/login', { method: 'POST', body: { email, password } }),
    me: () => request('/auth/me'),
    register: (data) =>
      request('/auth/register', { method: 'POST', body: data }),
  },

  // ─── Stock ──────────────────────────────────────────────────
  stock: {
    list: (params = {}) =>
      request('/stock?' + new URLSearchParams(params)),
    search: (q) =>
      request(`/stock/search?q=${encodeURIComponent(q)}`),
    expiring: (days) =>
      request(`/stock/expiring${days ? `?days=${days}` : ''}`),
    add: (data) =>
      request('/stock', { method: 'POST', body: data }),
    adjust: (id, adjustment, reason) =>
      request(`/stock/${id}`, { method: 'PUT', body: { adjustment, reason } }),
  },

  // ─── Batches / Approval queue ────────────────────────────────
  batches: {
    submit: (records) =>
      request('/batches', { method: 'POST', body: { records } }),
    list: (params = {}) =>
      request('/batches?' + new URLSearchParams(params)),
    getById: (id) =>
      request(`/batches/${id}`),
    approve: (id) =>
      request(`/batches/${id}/approve`, { method: 'POST' }),
    reject: (id, reason) =>
      request(`/batches/${id}/reject`, { method: 'POST', body: { reason } }),
  },

  // ─── Patient search ──────────────────────────────────────────
  patients: {
    search: (credential) =>
      request(`/patients/search?credential=${encodeURIComponent(credential)}`),
  },

  // ─── Analytics ───────────────────────────────────────────────
  analytics: {
    trends: (from, to) =>
      request(`/analytics/trends?from=${from}&to=${to}`),
    aiSummary: (from, to) =>
      request(`/analytics/ai-summary?from=${from}&to=${to}`),
    stockoutRisk: (pharmacyId) =>
      request(`/analytics/stockout-risk${pharmacyId ? `?pharmacyId=${pharmacyId}` : ''}`),
    regional: () =>
      request('/analytics/regional'),
  },

  // ─── Admin ───────────────────────────────────────────────────
  admin: {
    getFormFields: () => request('/admin/form-fields'),
    addFormField: (data) =>
      request('/admin/form-fields', { method: 'POST', body: data }),
    updateFormField: (id, data) =>
      request(`/admin/form-fields/${id}`, { method: 'PUT', body: data }),
    deleteFormField: (id) =>
      request(`/admin/form-fields/${id}`, { method: 'DELETE' }),
    getUsers: () => request('/admin/users'),
    updateUser: (id, data) =>
      request(`/admin/users/${id}`, { method: 'PUT', body: data }),
    getAuditLogs: (page) =>
      request(`/admin/audit-logs?page=${page || 1}`),
  },

  // ─── DHO ────────────────────────────────────────────────────
  dho: {
    getPharmacies: () => request('/dho/pharmacies'),
    registerPharmacy: (data) =>
      request('/dho/pharmacies', { method: 'POST', body: data }),
    getRegionalMap: () => request('/dho/regional-map'),
  },
};
