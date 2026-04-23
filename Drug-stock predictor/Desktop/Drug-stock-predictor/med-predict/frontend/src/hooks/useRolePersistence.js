// ============================================================
// src/hooks/useRolePersistence.js
//
// Manages role selection and auth state across sessions.
// localStorage stores:
//   med_predict_role   — the selected role ('data_entrant'|'admin'|'dho')
//   med_predict_token  — JWT token
//   med_predict_user   — serialized user object
// ============================================================
import { useState, useEffect, useCallback, createContext, useContext } from 'react';

const AUTH_KEYS = {
  ROLE:  'med_predict_role',
  TOKEN: 'med_predict_token',
  USER:  'med_predict_user',
};

// ─── Auth Context ─────────────────────────────────────────────
export const AuthContext = createContext(null);

export function AuthProvider({ children }) {
  const auth = useAuthLogic();
  return <AuthContext.Provider value={auth}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be used within AuthProvider');
  return ctx;
}

// ─── Core logic ───────────────────────────────────────────────
function useAuthLogic() {
  // Load persisted state synchronously to avoid flash
  const [selectedRole, setSelectedRole] = useState(
    () => localStorage.getItem(AUTH_KEYS.ROLE) || null
  );
  const [token, setToken] = useState(
    () => localStorage.getItem(AUTH_KEYS.TOKEN) || null
  );
  const [user, setUser] = useState(() => {
    try {
      const raw = localStorage.getItem(AUTH_KEYS.USER);
      return raw ? JSON.parse(raw) : null;
    } catch { return null; }
  });
  const [loading, setLoading] = useState(false);

  // ── Derived: is the user authenticated? ───────────────────
  const isAuthenticated = Boolean(token && user);

  // ── Role selection (shows the correct login portal) ────────
  // Once set, the user only sees the login form for that role
  const selectRole = useCallback((role) => {
    if (!['data_entrant', 'admin', 'dho'].includes(role)) return;
    localStorage.setItem(AUTH_KEYS.ROLE, role);
    setSelectedRole(role);
  }, []);

  const clearRoleSelection = useCallback(() => {
    localStorage.removeItem(AUTH_KEYS.ROLE);
    setSelectedRole(null);
  }, []);

  // ── Login ───────────────────────────────────────────────────
  const login = useCallback(async (email, password) => {
    setLoading(true);
    try {
      const res = await fetch(`${import.meta.env.VITE_API_URL}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email, password })
      });

      const data = await res.json();
      if (!res.ok) throw new Error(data.error || 'Login failed');

      // Validate role matches the selected portal
      if (data.user.role !== selectedRole) {
        throw new Error(
          `This login portal is for ${roleLabel(selectedRole)} only. ` +
          `Your account has role: ${roleLabel(data.user.role)}.`
        );
      }

      // Persist session
      localStorage.setItem(AUTH_KEYS.TOKEN, data.token);
      localStorage.setItem(AUTH_KEYS.USER, JSON.stringify(data.user));
      setToken(data.token);
      setUser(data.user);

      return data.user;
    } finally {
      setLoading(false);
    }
  }, [selectedRole]);

  // ── Logout ──────────────────────────────────────────────────
  const logout = useCallback(() => {
    localStorage.removeItem(AUTH_KEYS.TOKEN);
    localStorage.removeItem(AUTH_KEYS.USER);
    // Keep role selection — user returns to same login portal
    setToken(null);
    setUser(null);
  }, []);

  // ── Full reset (switch portal) ───────────────────────────────
  const fullReset = useCallback(() => {
    Object.values(AUTH_KEYS).forEach(k => localStorage.removeItem(k));
    setSelectedRole(null);
    setToken(null);
    setUser(null);
  }, []);

  // ── Auto-logout on token expiry ──────────────────────────────
  useEffect(() => {
    if (!token) return;
    try {
      // Decode JWT payload (no verification — that's the server's job)
      const payload = JSON.parse(atob(token.split('.')[1]));
      const expiresIn = payload.exp * 1000 - Date.now();
      if (expiresIn <= 0) { logout(); return; }
      const timer = setTimeout(logout, expiresIn);
      return () => clearTimeout(timer);
    } catch { logout(); }
  }, [token, logout]);

  // ── DE session wipe after batch upload ───────────────────────
  // Call this ONLY when the DE's data has been confirmed submitted
  const wipeDESession = useCallback(() => {
    sessionStorage.removeItem('med_predict_de_session');
    sessionStorage.removeItem('med_predict_de_records');
    sessionStorage.removeItem('med_predict_de_batch_draft');
  }, []);

  return {
    selectedRole,
    user,
    token,
    loading,
    isAuthenticated,
    selectRole,
    clearRoleSelection,
    login,
    logout,
    fullReset,
    wipeDESession,
  };
}

// ─── Helpers ─────────────────────────────────────────────────
export function roleLabel(role) {
  const labels = { data_entrant: 'Data Entrant', admin: 'Admin/Manager', dho: 'DHO (Top Admin)' };
  return labels[role] || role;
}

export function roleMeta(role) {
  return {
    data_entrant: {
      label: 'Data Entrant',
      shortLabel: 'DE',
      description: 'Enter daily patient and stock data',
      icon: '⚕️',
      color: 'teal',
      accentClass: 'from-teal-600 to-cyan-700'
    },
    admin: {
      label: 'Admin / Manager',
      shortLabel: 'Admin',
      description: 'Approve data, view AI analytics',
      icon: '📊',
      color: 'indigo',
      accentClass: 'from-indigo-600 to-violet-700'
    },
    dho: {
      label: 'District Health Officer',
      shortLabel: 'DHO',
      description: 'Regional oversight and AI predictions',
      icon: '🗺️',
      color: 'amber',
      accentClass: 'from-amber-600 to-orange-700'
    }
  }[role] || {};
}
