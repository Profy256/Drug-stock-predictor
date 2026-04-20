// ============================================================
// MED PREDICT SYSTEM — Complete React Frontend
// Single-file app with all three role portals
// Aesthetic: clinical precision — dark medical UI
// ============================================================
import { useState, useEffect, useCallback, useRef, createContext, useContext } from "react";

// ─── Constants ────────────────────────────────────────────────
const API_BASE = "http://localhost:4000/api/v1";
const AUTH_KEYS = {
  ROLE: "med_predict_role",
  TOKEN: "med_predict_token",
  USER: "med_predict_user",
};
const DE_SESSION_KEY = "med_predict_de_records";

// ─── Auth Context ─────────────────────────────────────────────
const AuthCtx = createContext(null);
function useAuth() { return useContext(AuthCtx); }

function AuthProvider({ children }) {
  const [role, setRole] = useState(() => localStorage.getItem(AUTH_KEYS.ROLE));
  const [token, setToken] = useState(() => localStorage.getItem(AUTH_KEYS.TOKEN));
  const [user, setUser] = useState(() => {
    try { return JSON.parse(localStorage.getItem(AUTH_KEYS.USER)); } catch { return null; }
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const selectRole = (r) => { localStorage.setItem(AUTH_KEYS.ROLE, r); setRole(r); };

  const login = async (email, password) => {
    setLoading(true); setError("");
    try {
      const res = await fetch(`${API_BASE}/auth/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data.error || "Login failed");
      if (data.user.role !== role)
        throw new Error(`Wrong portal. Your role is: ${ROLE_META[data.user.role]?.label}`);
      localStorage.setItem(AUTH_KEYS.TOKEN, data.token);
      localStorage.setItem(AUTH_KEYS.USER, JSON.stringify(data.user));
      setToken(data.token); setUser(data.user);
    } catch (e) { setError(e.message); }
    finally { setLoading(false); }
  };

  const logout = () => {
    [AUTH_KEYS.TOKEN, AUTH_KEYS.USER].forEach(k => localStorage.removeItem(k));
    setToken(null); setUser(null);
  };

  const fullReset = () => {
    Object.values(AUTH_KEYS).forEach(k => localStorage.removeItem(k));
    sessionStorage.removeItem(DE_SESSION_KEY);
    setRole(null); setToken(null); setUser(null);
  };

  // Auto-logout on token expiry
  useEffect(() => {
    if (!token) return;
    try {
      const { exp } = JSON.parse(atob(token.split(".")[1]));
      const ms = exp * 1000 - Date.now();
      if (ms <= 0) { logout(); return; }
      const t = setTimeout(logout, ms);
      return () => clearTimeout(t);
    } catch { logout(); }
  }, [token]);

  const apiFetch = useCallback(async (path, opts = {}) => {
    const res = await fetch(`${API_BASE}${path}`, {
      ...opts,
      headers: {
        "Content-Type": "application/json",
        ...(token ? { Authorization: `Bearer ${token}` } : {}),
        ...(opts.headers || {}),
      },
      body: opts.body ? JSON.stringify(opts.body) : undefined,
    });
    if (res.status === 401) { logout(); return null; }
    const d = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(d.error || `Error ${res.status}`);
    return d;
  }, [token]);

  return (
    <AuthCtx.Provider value={{ role, user, token, loading, error, isAuth: !!(token && user),
      selectRole, login, logout, fullReset, apiFetch }}>
      {children}
    </AuthCtx.Provider>
  );
}

// ─── Role Metadata ────────────────────────────────────────────
const ROLE_META = {
  data_entrant: { label: "Data Entrant", short: "DE",    color: "#00BFA6", bg: "#003D35",  icon: "⚕" },
  admin:        { label: "Admin / Manager", short: "MGR", color: "#7C6FFF", bg: "#1A1535", icon: "⊞" },
  dho:          { label: "District Health Officer", short: "DHO", color: "#FFB547", bg: "#2A1F00", icon: "◈" },
};

// ─── Root App ─────────────────────────────────────────────────
export default function App() {
  return (
    <AuthProvider>
      <style>{GLOBAL_CSS}</style>
      <Router />
    </AuthProvider>
  );
}

function Router() {
  const { role, isAuth } = useAuth();
  if (!role)       return <RoleSelect />;
  if (!isAuth)     return <LoginPortal />;
  if (role === "data_entrant") return <DEPortal />;
  if (role === "admin")        return <AdminPortal />;
  if (role === "dho")          return <DHOPortal />;
  return <div>Unknown role</div>;
}

// ═══════════════════════════════════════════════════════════════
// ROLE SELECTION SCREEN
// ═══════════════════════════════════════════════════════════════
function RoleSelect() {
  const { selectRole } = useAuth();
  const [hovered, setHovered] = useState(null);

  return (
    <div className="screen-center" style={{ background: "var(--bg-deep)", minHeight: "100vh" }}>
      <div style={{ maxWidth: 640, width: "100%", padding: "0 24px" }}>
        <div className="logo-block" style={{ marginBottom: 48, textAlign: "center" }}>
          <div style={{ fontSize: 40, marginBottom: 8 }}>✚</div>
          <h1 style={{ fontSize: 28, fontWeight: 700, letterSpacing: 2, color: "var(--text-primary)",
            fontFamily: "var(--font-display)", margin: 0 }}>MED PREDICT</h1>
          <p style={{ color: "var(--text-muted)", fontSize: 13, marginTop: 6, letterSpacing: 1 }}>
            HOSPITAL STOCK MANAGEMENT SYSTEM
          </p>
        </div>

        <p style={{ color: "var(--text-muted)", textAlign: "center", marginBottom: 32, fontSize: 14 }}>
          Select your access portal. This device will remember your selection.
        </p>

        <div style={{ display: "flex", flexDirection: "column", gap: 16 }}>
          {Object.entries(ROLE_META).map(([roleKey, meta]) => (
            <button key={roleKey}
              className="role-card"
              style={{
                background: hovered === roleKey ? meta.bg : "var(--bg-card)",
                borderColor: hovered === roleKey ? meta.color : "var(--border)",
                transform: hovered === roleKey ? "translateX(6px)" : "none",
              }}
              onMouseEnter={() => setHovered(roleKey)}
              onMouseLeave={() => setHovered(null)}
              onClick={() => selectRole(roleKey)}>
              <span style={{ fontSize: 24, marginRight: 16, opacity: 0.9 }}>{meta.icon}</span>
              <div style={{ flex: 1, textAlign: "left" }}>
                <div style={{ fontWeight: 600, color: hovered === roleKey ? meta.color : "var(--text-primary)",
                  fontSize: 15, fontFamily: "var(--font-display)" }}>{meta.label}</div>
                <div style={{ color: "var(--text-muted)", fontSize: 12, marginTop: 2 }}>
                  {roleKey === "data_entrant" && "Enter daily patient data and medicine stock"}
                  {roleKey === "admin" && "Approve data batches, view AI analytics dashboard"}
                  {roleKey === "dho" && "Regional oversight, multi-pharmacy AI predictions"}
                </div>
              </div>
              <span style={{ color: meta.color, fontSize: 18 }}>›</span>
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}

// ═══════════════════════════════════════════════════════════════
// LOGIN PORTAL
// ═══════════════════════════════════════════════════════════════
function LoginPortal() {
  const { role, login, loading, error, fullReset } = useAuth();
  const meta = ROLE_META[role];
  const [email, setEmail] = useState("");
  const [pass, setPass]   = useState("");

  const handleSubmit = (e) => { e.preventDefault(); login(email, pass); };

  return (
    <div className="screen-center" style={{ background: "var(--bg-deep)", minHeight: "100vh" }}>
      <div style={{ width: "100%", maxWidth: 420, padding: "0 24px" }}>
        {/* Portal badge */}
        <div style={{ textAlign: "center", marginBottom: 40 }}>
          <div style={{ display: "inline-flex", alignItems: "center", gap: 10,
            background: meta.bg, border: `1px solid ${meta.color}`,
            borderRadius: 100, padding: "8px 20px", marginBottom: 24 }}>
            <span style={{ fontSize: 18 }}>{meta.icon}</span>
            <span style={{ color: meta.color, fontWeight: 600, fontSize: 13,
              letterSpacing: 1, fontFamily: "var(--font-display)" }}>
              {meta.label.toUpperCase()}
            </span>
          </div>
          <h2 style={{ color: "var(--text-primary)", fontSize: 22, margin: 0,
            fontFamily: "var(--font-display)", fontWeight: 700 }}>Sign in to continue</h2>
        </div>

        <form onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="field-label">Email address</label>
            <input className="field-input" type="email" value={email}
              onChange={e => setEmail(e.target.value)} placeholder="you@hospital.org"
              autoFocus required />
          </div>
          <div className="form-group" style={{ marginTop: 16 }}>
            <label className="field-label">Password</label>
            <input className="field-input" type="password" value={pass}
              onChange={e => setPass(e.target.value)} placeholder="••••••••" required />
          </div>

          {error && (
            <div className="alert-danger" style={{ marginTop: 16 }}>{error}</div>
          )}

          <button type="submit" className="btn-primary" disabled={loading}
            style={{ background: meta.color, marginTop: 24, width: "100%" }}>
            {loading ? "Signing in…" : "Sign in"}
          </button>
        </form>

        <button onClick={fullReset}
          style={{ marginTop: 20, width: "100%", background: "none", border: "none",
            color: "var(--text-muted)", cursor: "pointer", fontSize: 12 }}>
          ← Switch portal
        </button>
      </div>
    </div>
  );
}

// ═══════════════════════════════════════════════════════════════
// DATA ENTRANT PORTAL
// ═══════════════════════════════════════════════════════════════
function DEPortal() {
  const { user, logout, apiFetch, wipeDESession } = useAuth();
  const [tab, setTab] = useState("stock");
  const [formFields, setFormFields] = useState([]);
  // Session records — stored in sessionStorage, wiped after submission
  const [sessionRecords, setSessionRecords] = useState(() => {
    try { return JSON.parse(sessionStorage.getItem(DE_SESSION_KEY)) || []; }
    catch { return []; }
  });
  const [medicines, setMedicines] = useState([]);
  const [submitStatus, setSubmitStatus] = useState(null);

  useEffect(() => {
    apiFetch("/admin/form-fields").then(d => d && setFormFields(d.fields || []));
    apiFetch("/stock?").then(d => d && setMedicines(d.medicines || []));
  }, []);

  // Persist session records to sessionStorage
  useEffect(() => {
    sessionStorage.setItem(DE_SESSION_KEY, JSON.stringify(sessionRecords));
  }, [sessionRecords]);

  const addRecord = (record) => {
    setSessionRecords(prev => [...prev, { ...record, id: Date.now() }]);
  };

  const removeRecord = (id) => {
    setSessionRecords(prev => prev.filter(r => r.id !== id));
  };

  const submitBatch = async () => {
    if (!sessionRecords.length) return;
    setSubmitStatus("loading");
    try {
      await apiFetch("/batches", {
        method: "POST",
        body: { records: sessionRecords.map(({ id, ...r }) => r) }
      });
      // ── PII WIPE: clear all local session data ─────────────
      sessionStorage.removeItem(DE_SESSION_KEY);
      setSessionRecords([]);
      setSubmitStatus("success");
      setTimeout(() => setSubmitStatus(null), 4000);
    } catch (e) {
      setSubmitStatus("error:" + e.message);
    }
  };

  return (
    <Shell role="data_entrant" user={user} onLogout={logout}>
      {/* Session badge */}
      <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between",
        padding: "12px 24px", borderBottom: "1px solid var(--border)", background: "var(--bg-card)" }}>
        <div style={{ color: "var(--text-muted)", fontSize: 13 }}>
          Session records: <strong style={{ color: "#00BFA6" }}>{sessionRecords.length}</strong>
          <span style={{ marginLeft: 8, fontSize: 11, opacity: 0.6 }}>(wiped on submission)</span>
        </div>
        {sessionRecords.length > 0 && (
          <button className="btn-primary" onClick={submitBatch}
            disabled={submitStatus === "loading"}
            style={{ background: "#00BFA6", fontSize: 13, padding: "6px 18px" }}>
            {submitStatus === "loading" ? "Submitting…" : `Submit ${sessionRecords.length} record(s) for approval`}
          </button>
        )}
      </div>

      {submitStatus === "success" && (
        <div className="alert-success" style={{ margin: "16px 24px" }}>
          ✓ Batch submitted successfully. All local data has been wiped pending manager approval.
        </div>
      )}
      {submitStatus?.startsWith("error") && (
        <div className="alert-danger" style={{ margin: "16px 24px" }}>
          {submitStatus.replace("error:", "")}
        </div>
      )}

      {/* Tabs */}
      <div style={{ display: "flex", borderBottom: "1px solid var(--border)" }}>
        {[["stock", "Add Stock"], ["patient", "Patient Visit"], ["search", "Patient Search"]].map(
          ([key, label]) => (
            <button key={key} onClick={() => setTab(key)}
              style={{ padding: "14px 24px", background: "none", border: "none",
                borderBottom: tab === key ? "2px solid #00BFA6" : "2px solid transparent",
                color: tab === key ? "#00BFA6" : "var(--text-muted)",
                fontFamily: "var(--font-display)", fontSize: 13, cursor: "pointer",
                fontWeight: tab === key ? 600 : 400, letterSpacing: 0.5 }}>
              {label}
            </button>
          )
        )}
      </div>

      <div style={{ padding: 24 }}>
        {tab === "stock"   && <StockEntryForm onAdd={r => apiFetch("/stock", { method: "POST", body: r })} />}
        {tab === "patient" && <PatientVisitForm fields={formFields} medicines={medicines} onSave={addRecord} />}
        {tab === "search"  && <PatientSearchForm />}
      </div>

      {/* Session record list (current session only) */}
      {sessionRecords.length > 0 && tab === "patient" && (
        <div style={{ padding: "0 24px 24px" }}>
          <h4 style={{ color: "var(--text-muted)", fontSize: 12, letterSpacing: 1, marginBottom: 12 }}>
            CURRENT SESSION — {sessionRecords.length} RECORD(S)
          </h4>
          <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
            {sessionRecords.map((r) => (
              <div key={r.id} className="record-row">
                <div>
                  <span style={{ color: "var(--text-primary)", fontSize: 13, fontWeight: 500 }}>
                    {r.patientData?.credential || "—"}
                  </span>
                  <span style={{ color: "var(--text-muted)", fontSize: 12, marginLeft: 12 }}>
                    {r.diagnosis || "No diagnosis"} •{" "}
                    {r.medicinesGiven?.map(m => m.name).join(", ")}
                  </span>
                </div>
                <button onClick={() => removeRecord(r.id)}
                  style={{ background: "none", border: "none", color: "#E53935",
                    cursor: "pointer", fontSize: 16 }}>✕</button>
              </div>
            ))}
          </div>
        </div>
      )}
    </Shell>
  );
}

function StockEntryForm({ onAdd }) {
  const [form, setForm] = useState({ name: "", quantityTotal: "", expiryDate: "",
    batchNumber: "", category: "", notificationDays: 14 });
  const [status, setStatus] = useState(null);

  const set = (k, v) => setForm(f => ({ ...f, [k]: v }));

  const submit = async (e) => {
    e.preventDefault(); setStatus("loading");
    try {
      await onAdd(form);
      setForm({ name: "", quantityTotal: "", expiryDate: "", batchNumber: "", category: "", notificationDays: 14 });
      setStatus("success");
      setTimeout(() => setStatus(null), 3000);
    } catch (e) { setStatus("error:" + e.message); }
  };

  return (
    <div style={{ maxWidth: 520 }}>
      <h3 className="section-title">Add Incoming Stock</h3>
      <form onSubmit={submit}>
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16 }}>
          <div className="form-group" style={{ gridColumn: "1/-1" }}>
            <label className="field-label">Medicine Name *</label>
            <input className="field-input" value={form.name} onChange={e => set("name", e.target.value)}
              placeholder="e.g. Amoxicillin 500mg" required />
          </div>
          <div className="form-group">
            <label className="field-label">Quantity (boxes) *</label>
            <input className="field-input" type="number" min="1" value={form.quantityTotal}
              onChange={e => set("quantityTotal", e.target.value)} required />
          </div>
          <div className="form-group">
            <label className="field-label">Expiry Date *</label>
            <input className="field-input" type="date" value={form.expiryDate}
              onChange={e => set("expiryDate", e.target.value)} required />
          </div>
          <div className="form-group">
            <label className="field-label">Batch Number</label>
            <input className="field-input" value={form.batchNumber}
              onChange={e => set("batchNumber", e.target.value)} placeholder="Optional" />
          </div>
          <div className="form-group">
            <label className="field-label">Alert me (days before expiry)</label>
            <input className="field-input" type="number" min="1" max="90" value={form.notificationDays}
              onChange={e => set("notificationDays", e.target.value)} />
          </div>
        </div>
        {status === "success" && <div className="alert-success" style={{ marginTop: 12 }}>Stock added ✓</div>}
        {status?.startsWith("error") && <div className="alert-danger" style={{ marginTop: 12 }}>{status.replace("error:", "")}</div>}
        <button type="submit" className="btn-primary" disabled={status === "loading"}
          style={{ marginTop: 20, background: "#00BFA6" }}>
          {status === "loading" ? "Saving…" : "Add Stock"}
        </button>
      </form>
    </div>
  );
}

function PatientVisitForm({ fields, medicines, onSave }) {
  const [patientData, setPatientData] = useState({});
  const [selectedMeds, setSelectedMeds] = useState([]);
  const [diagnosis, setDiagnosis] = useState("");
  const [medSearch, setMedSearch] = useState("");
  const [saved, setSaved] = useState(false);

  const activeFields = fields.filter(f => f.is_active);
  const filteredMeds = medicines.filter(m =>
    m.name.toLowerCase().includes(medSearch.toLowerCase()) && m.quantity_remaining > 0
  );

  const toggleMed = (med) => {
    setSelectedMeds(prev => {
      const exists = prev.find(m => m.medicine_id === med.id);
      if (exists) return prev.filter(m => m.medicine_id !== med.id);
      return [...prev, { medicine_id: med.id, name: med.name, qty: 1 }];
    });
  };

  const setMedQty = (medId, qty) => {
    setSelectedMeds(prev => prev.map(m => m.medicine_id === medId ? { ...m, qty: parseInt(qty) } : m));
  };

  const handleSave = () => {
    if (!selectedMeds.length) return alert("Select at least one medicine");
    onSave({ patientData, medicinesGiven: selectedMeds, diagnosis, visitDate: new Date().toISOString().split("T")[0] });
    // Reset form (but keep field structure)
    setPatientData({}); setSelectedMeds([]); setDiagnosis(""); setMedSearch("");
    setSaved(true); setTimeout(() => setSaved(false), 2500);
  };

  return (
    <div style={{ maxWidth: 580 }}>
      <h3 className="section-title">Record Patient Visit</h3>
      {saved && <div className="alert-success" style={{ marginBottom: 16 }}>Record added to session ✓</div>}

      {/* Dynamic form fields */}
      <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 14, marginBottom: 20 }}>
        {activeFields.map(f => (
          <div key={f.field_key} className="form-group"
            style={f.field_key === "name" || f.field_key === "credential" ? { gridColumn: "1/-1" } : {}}>
            <label className="field-label">
              {f.label} {f.is_required && <span style={{ color: "#F44336" }}>*</span>}
            </label>
            {f.field_type === "select" ? (
              <select className="field-input" value={patientData[f.field_key] || ""}
                onChange={e => setPatientData(p => ({ ...p, [f.field_key]: e.target.value }))}>
                <option value="">Select…</option>
                {(f.options || []).map(o => <option key={o} value={o}>{o}</option>)}
              </select>
            ) : (
              <input className="field-input" type={f.field_type}
                value={patientData[f.field_key] || ""}
                onChange={e => setPatientData(p => ({ ...p, [f.field_key]: e.target.value }))}
                required={f.is_required} />
            )}
          </div>
        ))}
      </div>

      {/* Diagnosis */}
      <div className="form-group" style={{ marginBottom: 20 }}>
        <label className="field-label">Diagnosis / Condition</label>
        <input className="field-input" value={diagnosis}
          onChange={e => setDiagnosis(e.target.value)} placeholder="e.g. Malaria, Hypertension" />
      </div>

      {/* Medicine selection */}
      <div className="form-group">
        <label className="field-label">Medicines Dispensed *</label>
        <input className="field-input" value={medSearch} onChange={e => setMedSearch(e.target.value)}
          placeholder="Search medicines…" style={{ marginBottom: 10 }} />
        <div style={{ maxHeight: 200, overflowY: "auto", border: "1px solid var(--border)",
          borderRadius: 8, background: "var(--bg-deep)" }}>
          {filteredMeds.slice(0, 12).map(m => {
            const sel = selectedMeds.find(s => s.medicine_id === m.id);
            return (
              <div key={m.id} style={{ display: "flex", alignItems: "center",
                padding: "10px 14px", borderBottom: "1px solid var(--border)",
                background: sel ? "rgba(0,191,166,0.08)" : "transparent",
                cursor: "pointer" }} onClick={() => toggleMed(m)}>
                <div style={{ width: 18, height: 18, borderRadius: 4,
                  border: `2px solid ${sel ? "#00BFA6" : "var(--border)"}`,
                  background: sel ? "#00BFA6" : "transparent", marginRight: 12,
                  display: "flex", alignItems: "center", justifyContent: "center",
                  fontSize: 11, color: "white", flexShrink: 0 }}>
                  {sel ? "✓" : ""}
                </div>
                <div style={{ flex: 1 }}>
                  <div style={{ fontSize: 13, color: "var(--text-primary)", fontWeight: 500 }}>{m.name}</div>
                  <div style={{ fontSize: 11, color: "var(--text-muted)" }}>
                    {m.quantity_remaining} available • expires {m.expiry_date}
                  </div>
                </div>
                {sel && (
                  <input type="number" min="1" max={m.quantity_remaining} value={sel.qty}
                    onClick={e => e.stopPropagation()}
                    onChange={e => setMedQty(m.id, e.target.value)}
                    style={{ width: 50, background: "var(--bg-card)", border: "1px solid var(--border)",
                      borderRadius: 6, color: "var(--text-primary)", padding: "4px 8px", fontSize: 12 }} />
                )}
              </div>
            );
          })}
          {filteredMeds.length === 0 && (
            <div style={{ padding: 16, color: "var(--text-muted)", fontSize: 13, textAlign: "center" }}>
              No medicines found
            </div>
          )}
        </div>
      </div>

      <button className="btn-primary" onClick={handleSave}
        style={{ marginTop: 20, background: "#00BFA6" }}>
        + Add to Session
      </button>
    </div>
  );
}

function PatientSearchForm() {
  const { apiFetch } = useAuth();
  const [cred, setCred] = useState("");
  const [results, setResults] = useState(null);
  const [loading, setLoading] = useState(false);

  const search = async (e) => {
    e.preventDefault(); if (!cred.trim()) return;
    setLoading(true);
    try {
      const data = await apiFetch(`/patients/search?credential=${encodeURIComponent(cred)}`);
      setResults(data?.visits || []);
    } catch { setResults([]); }
    finally { setLoading(false); }
  };

  return (
    <div style={{ maxWidth: 560 }}>
      <h3 className="section-title">Patient Dose History</h3>
      <p style={{ color: "var(--text-muted)", fontSize: 13, marginBottom: 20 }}>
        Search by patient credential (ID / NIN) to view prescription history for dose planning.
      </p>
      <form onSubmit={search} style={{ display: "flex", gap: 12, marginBottom: 24 }}>
        <input className="field-input" style={{ flex: 1 }} value={cred}
          onChange={e => setCred(e.target.value)} placeholder="Patient ID / NIN" />
        <button type="submit" className="btn-primary" disabled={loading}
          style={{ background: "#00BFA6", whiteSpace: "nowrap" }}>
          {loading ? "…" : "Search"}
        </button>
      </form>
      {results !== null && (
        results.length === 0
          ? <div className="alert-danger">No records found for this credential.</div>
          : results.map((v, i) => (
            <div key={i} className="record-row" style={{ flexDirection: "column", alignItems: "flex-start", gap: 6 }}>
              <div style={{ display: "flex", justifyContent: "space-between", width: "100%" }}>
                <span style={{ color: "var(--text-primary)", fontWeight: 600, fontSize: 13 }}>
                  {v.visit_date}
                </span>
                <span style={{ color: "var(--text-muted)", fontSize: 12 }}>{v.diagnosis}</span>
              </div>
              <div style={{ color: "#00BFA6", fontSize: 12 }}>
                Medicines: {Array.isArray(v.medicines_given)
                  ? v.medicines_given.map(m => `${m.name} × ${m.qty}`).join(", ")
                  : "—"}
              </div>
            </div>
          ))
      )}
    </div>
  );
}

// ═══════════════════════════════════════════════════════════════
// ADMIN PORTAL
// ═══════════════════════════════════════════════════════════════
function AdminPortal() {
  const { user, logout, apiFetch } = useAuth();
  const [tab, setTab] = useState("queue");

  return (
    <Shell role="admin" user={user} onLogout={logout}>
      <div style={{ display: "flex", borderBottom: "1px solid var(--border)", overflowX: "auto" }}>
        {[["queue","Approval Queue"],["analytics","AI Analytics"],["stock","Stock"],
          ["fields","Form Fields"],["users","Users"]].map(([key, label]) => (
          <button key={key} onClick={() => setTab(key)}
            style={{ padding: "14px 22px", background: "none", border: "none",
              borderBottom: tab === key ? "2px solid #7C6FFF" : "2px solid transparent",
              color: tab === key ? "#7C6FFF" : "var(--text-muted)", fontFamily: "var(--font-display)",
              fontSize: 13, cursor: "pointer", fontWeight: tab === key ? 600 : 400,
              whiteSpace: "nowrap", letterSpacing: 0.5 }}>
            {label}
          </button>
        ))}
      </div>
      <div style={{ padding: 24 }}>
        {tab === "queue"     && <ApprovalQueue apiFetch={apiFetch} />}
        {tab === "analytics" && <AnalyticsDashboard apiFetch={apiFetch} />}
        {tab === "stock"     && <StockPanel apiFetch={apiFetch} />}
        {tab === "fields"    && <FormFieldsManager apiFetch={apiFetch} />}
        {tab === "users"     && <UsersManager apiFetch={apiFetch} />}
      </div>
    </Shell>
  );
}

function ApprovalQueue({ apiFetch }) {
  const [batches, setBatches] = useState([]);
  const [loading, setLoading] = useState(true);
  const [selected, setSelected] = useState(null);
  const [rejectReason, setRejectReason] = useState("");
  const [actionStatus, setActionStatus] = useState({});

  const load = async () => {
    setLoading(true);
    const d = await apiFetch("/batches?status=pending");
    if (d) setBatches(d.batches || []);
    setLoading(false);
  };
  useEffect(() => { load(); }, []);

  const loadDetail = async (batchId) => {
    const d = await apiFetch(`/batches/${batchId}`);
    if (d) setSelected(d);
  };

  const handleApprove = async (batchId) => {
    setActionStatus(s => ({ ...s, [batchId]: "loading" }));
    try {
      await apiFetch(`/batches/${batchId}/approve`, { method: "POST" });
      setActionStatus(s => ({ ...s, [batchId]: "approved" }));
      setBatches(prev => prev.filter(b => b.id !== batchId));
      if (selected?.batch.id === batchId) setSelected(null);
    } catch (e) { setActionStatus(s => ({ ...s, [batchId]: "error:" + e.message })); }
  };

  const handleReject = async (batchId) => {
    if (!rejectReason.trim()) return alert("Provide a rejection reason");
    setActionStatus(s => ({ ...s, [batchId]: "loading" }));
    try {
      await apiFetch(`/batches/${batchId}/reject`, { method: "POST", body: { reason: rejectReason } });
      setActionStatus(s => ({ ...s, [batchId]: "rejected" }));
      setBatches(prev => prev.filter(b => b.id !== batchId));
      setSelected(null); setRejectReason("");
    } catch (e) { setActionStatus(s => ({ ...s, [batchId]: "error:" + e.message })); }
  };

  if (loading) return <Spinner />;

  return (
    <div>
      <h3 className="section-title">Pending Approval Queue</h3>
      {batches.length === 0 ? (
        <div style={{ textAlign: "center", padding: 40, color: "var(--text-muted)", fontSize: 14 }}>
          ✓ No pending batches. All data is up to date.
        </div>
      ) : (
        <div style={{ display: "flex", gap: 20, alignItems: "flex-start" }}>
          {/* Batch list */}
          <div style={{ flex: 1, maxWidth: 380, display: "flex", flexDirection: "column", gap: 10 }}>
            {batches.map(b => (
              <div key={b.id} className="record-row"
                style={{ flexDirection: "column", alignItems: "flex-start",
                  background: selected?.batch.id === b.id ? "rgba(124,111,255,0.1)" : "var(--bg-card)",
                  borderColor: selected?.batch.id === b.id ? "#7C6FFF" : "var(--border)",
                  cursor: "pointer" }}
                onClick={() => loadDetail(b.id)}>
                <div style={{ display: "flex", justifyContent: "space-between", width: "100%" }}>
                  <span style={{ fontWeight: 600, fontSize: 13, color: "var(--text-primary)" }}>
                    {b.submitted_by_name}
                  </span>
                  <span style={{ fontSize: 11, color: "var(--text-muted)" }}>
                    {new Date(b.submitted_at).toLocaleString()}
                  </span>
                </div>
                <div style={{ fontSize: 12, color: "var(--text-muted)" }}>
                  {b.record_count} record(s)
                </div>
                <div style={{ display: "flex", gap: 8, marginTop: 8 }}>
                  <button className="btn-sm-success"
                    disabled={actionStatus[b.id] === "loading"}
                    onClick={e => { e.stopPropagation(); handleApprove(b.id); }}>
                    ✓ Approve
                  </button>
                  <button className="btn-sm-danger"
                    onClick={e => { e.stopPropagation(); loadDetail(b.id); }}>
                    ✕ Reject…
                  </button>
                </div>
              </div>
            ))}
          </div>

          {/* Batch detail */}
          {selected && (
            <div style={{ flex: 1, background: "var(--bg-card)", border: "1px solid var(--border)",
              borderRadius: 12, padding: 20 }}>
              <h4 style={{ color: "var(--text-primary)", marginBottom: 16, fontSize: 14,
                fontFamily: "var(--font-display)" }}>
                Batch Detail — {selected.batch.record_count} record(s)
              </h4>
              <div style={{ display: "flex", flexDirection: "column", gap: 10, maxHeight: 300,
                overflowY: "auto", marginBottom: 16 }}>
                {(selected.records || []).map((r, i) => (
                  <div key={r.id} style={{ background: "var(--bg-deep)", borderRadius: 8, padding: "10px 14px" }}>
                    <div style={{ fontSize: 12, color: "#7C6FFF", marginBottom: 4 }}>Record #{i + 1}</div>
                    <div style={{ fontSize: 13, color: "var(--text-primary)" }}>
                      {r.patient_data?.name || "—"} / {r.patient_data?.credential || "—"}
                    </div>
                    <div style={{ fontSize: 12, color: "var(--text-muted)", marginTop: 2 }}>
                      {r.diagnosis} • {(r.medicines_given || []).map(m => m.name).join(", ")}
                    </div>
                  </div>
                ))}
              </div>
              <textarea className="field-input" rows={2} value={rejectReason}
                onChange={e => setRejectReason(e.target.value)}
                placeholder="Rejection reason (required to reject)…"
                style={{ resize: "none", marginBottom: 12, fontSize: 13 }} />
              <div style={{ display: "flex", gap: 10 }}>
                <button className="btn-sm-success" style={{ flex: 1 }}
                  onClick={() => handleApprove(selected.batch.id)}>
                  ✓ Approve & Commit (PII will be wiped)
                </button>
                <button className="btn-sm-danger" style={{ flex: 1 }}
                  onClick={() => handleReject(selected.batch.id)}>
                  ✕ Reject
                </button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function AnalyticsDashboard({ apiFetch }) {
  const [trends, setTrends] = useState(null);
  const [summary, setSummary] = useState("");
  const [summaryLoading, setSummaryLoading] = useState(false);
  const today = new Date().toISOString().split("T")[0];
  const [from, setFrom] = useState(new Date(Date.now() - 30 * 864e5).toISOString().split("T")[0]);
  const [to, setTo] = useState(today);

  const load = async () => {
    const d = await apiFetch(`/analytics/trends?from=${from}&to=${to}`);
    if (d) setTrends(d.trends);
  };
  useEffect(() => { load(); }, [from, to]);

  const getAISummary = async () => {
    setSummaryLoading(true);
    try {
      const d = await apiFetch(`/analytics/ai-summary?from=${from}&to=${to}`);
      if (d) setSummary(d.summary);
    } finally { setSummaryLoading(false); }
  };

  return (
    <div>
      <div style={{ display: "flex", alignItems: "center", gap: 16, marginBottom: 24 }}>
        <h3 className="section-title" style={{ margin: 0 }}>AI Analytics</h3>
        <div style={{ display: "flex", gap: 10, alignItems: "center" }}>
          <input type="date" value={from} onChange={e => setFrom(e.target.value)}
            className="field-input" style={{ width: 140, fontSize: 12 }} />
          <span style={{ color: "var(--text-muted)", fontSize: 12 }}>to</span>
          <input type="date" value={to} onChange={e => setTo(e.target.value)}
            className="field-input" style={{ width: 140, fontSize: 12 }} />
        </div>
      </div>

      {trends && (
        <>
          {/* Stats cards */}
          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fit, minmax(160px, 1fr))", gap: 16, marginBottom: 24 }}>
            {[
              { label: "Total Visits", value: trends.summary?.currentVisits, color: "#7C6FFF" },
              { label: "vs Previous Period", value: `${trends.summary?.changePercent > 0 ? "+" : ""}${trends.summary?.changePercent || 0}%`,
                color: (trends.summary?.changePercent || 0) >= 0 ? "#00BFA6" : "#F44336" },
              { label: "Expiring Soon", value: trends.stockSnapshot?.expiring_soon || 0, color: "#FFB547" },
              { label: "Low Stock", value: trends.stockSnapshot?.low_stock || 0, color: "#F44336" },
            ].map(c => (
              <div key={c.label} style={{ background: "var(--bg-card)", border: "1px solid var(--border)",
                borderRadius: 12, padding: "16px 20px" }}>
                <div style={{ color: c.color, fontSize: 26, fontWeight: 700, fontFamily: "var(--font-display)" }}>
                  {c.value}
                </div>
                <div style={{ color: "var(--text-muted)", fontSize: 12, marginTop: 4 }}>{c.label}</div>
              </div>
            ))}
          </div>

          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 20, marginBottom: 24 }}>
            {/* Top medicines */}
            <div style={{ background: "var(--bg-card)", border: "1px solid var(--border)", borderRadius: 12, padding: 20 }}>
              <h4 style={{ color: "#7C6FFF", fontSize: 12, letterSpacing: 1, marginBottom: 16 }}>TOP MEDICINES</h4>
              {(trends.topMedicines || []).slice(0, 5).map((m, i) => (
                <div key={i} style={{ display: "flex", alignItems: "center", gap: 12, marginBottom: 12 }}>
                  <span style={{ color: "var(--text-muted)", fontSize: 12, width: 16 }}>{i + 1}</span>
                  <div style={{ flex: 1 }}>
                    <div style={{ fontSize: 13, color: "var(--text-primary)", fontWeight: 500, marginBottom: 4 }}>
                      {m.medicine_name}
                    </div>
                    <div style={{ height: 4, borderRadius: 2, background: "var(--border)",
                      overflow: "hidden" }}>
                      <div style={{ height: "100%", background: "#7C6FFF",
                        width: `${(m.total_dispensed / (trends.topMedicines[0]?.total_dispensed || 1)) * 100}%`,
                        transition: "width 0.8s" }} />
                    </div>
                  </div>
                  <span style={{ color: "var(--text-muted)", fontSize: 12, width: 32, textAlign: "right" }}>
                    {m.total_dispensed}
                  </span>
                </div>
              ))}
            </div>

            {/* Top diseases */}
            <div style={{ background: "var(--bg-card)", border: "1px solid var(--border)", borderRadius: 12, padding: 20 }}>
              <h4 style={{ color: "#00BFA6", fontSize: 12, letterSpacing: 1, marginBottom: 16 }}>TOP DISEASES</h4>
              {(trends.topDiseases || []).slice(0, 5).map((d, i) => (
                <div key={i} style={{ display: "flex", alignItems: "center", gap: 12, marginBottom: 12 }}>
                  <span style={{ color: "var(--text-muted)", fontSize: 12, width: 16 }}>{i + 1}</span>
                  <div style={{ flex: 1 }}>
                    <div style={{ fontSize: 13, color: "var(--text-primary)", fontWeight: 500,
                      textTransform: "capitalize", marginBottom: 4 }}>{d.disease}</div>
                    <div style={{ height: 4, borderRadius: 2, background: "var(--border)", overflow: "hidden" }}>
                      <div style={{ height: "100%", background: "#00BFA6",
                        width: `${d.percentage}%`, transition: "width 0.8s" }} />
                    </div>
                  </div>
                  <span style={{ color: "var(--text-muted)", fontSize: 12, width: 44, textAlign: "right" }}>
                    {d.case_count} ({d.percentage}%)
                  </span>
                </div>
              ))}
            </div>
          </div>

          {/* AI Summary */}
          <div style={{ background: "var(--bg-card)", border: "1px solid var(--border)", borderRadius: 12, padding: 20 }}>
            <div style={{ display: "flex", alignItems: "center", justifyContent: "space-between", marginBottom: 16 }}>
              <h4 style={{ color: "#FFB547", fontSize: 12, letterSpacing: 1, margin: 0 }}>AI MANAGEMENT BRIEFING</h4>
              <button className="btn-primary" onClick={getAISummary} disabled={summaryLoading}
                style={{ background: "#FFB547", color: "#1A1535", fontSize: 12, padding: "6px 16px" }}>
                {summaryLoading ? "Generating…" : "✦ Generate AI Summary"}
              </button>
            </div>
            {summaryLoading && (
              <div style={{ color: "var(--text-muted)", fontSize: 13, fontStyle: "italic" }}>
                AI is analyzing your pharmacy data…
              </div>
            )}
            {summary && (
              <p style={{ color: "var(--text-primary)", fontSize: 13, lineHeight: 1.7,
                whiteSpace: "pre-wrap", margin: 0 }}>{summary}</p>
            )}
            {!summary && !summaryLoading && (
              <div style={{ color: "var(--text-muted)", fontSize: 13, textAlign: "center", padding: 20 }}>
                Click "Generate AI Summary" for an intelligent briefing
              </div>
            )}
          </div>
        </>
      )}
    </div>
  );
}

function StockPanel({ apiFetch }) {
  const [medicines, setMedicines] = useState([]);
  const [loading, setLoading] = useState(true);

  const load = async () => {
    const d = await apiFetch("/stock");
    if (d) setMedicines(d.medicines || []);
    setLoading(false);
  };
  useEffect(() => { load(); }, []);

  const statusColors = { ok: "#00BFA6", expiring_soon: "#FFB547", expired: "#F44336", low_stock: "#FF7043" };

  if (loading) return <Spinner />;

  return (
    <div>
      <h3 className="section-title">Stock Overview</h3>
      <div style={{ overflowX: "auto" }}>
        <table style={{ width: "100%", borderCollapse: "collapse", fontSize: 13 }}>
          <thead>
            <tr>
              {["Medicine", "Category", "Remaining", "Total", "Expires", "Status"].map(h => (
                <th key={h} style={{ textAlign: "left", padding: "10px 14px",
                  borderBottom: "1px solid var(--border)", color: "var(--text-muted)",
                  fontSize: 11, letterSpacing: 0.5, fontWeight: 600 }}>
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {medicines.map(m => (
              <tr key={m.id} style={{ borderBottom: "1px solid var(--border)" }}>
                <td style={{ padding: "12px 14px", color: "var(--text-primary)", fontWeight: 500 }}>{m.name}</td>
                <td style={{ padding: "12px 14px", color: "var(--text-muted)" }}>{m.category || "—"}</td>
                <td style={{ padding: "12px 14px", color: "var(--text-primary)" }}>{m.quantity_remaining}</td>
                <td style={{ padding: "12px 14px", color: "var(--text-muted)" }}>{m.quantity_total}</td>
                <td style={{ padding: "12px 14px", color: "var(--text-muted)" }}>{m.expiry_date}</td>
                <td style={{ padding: "12px 14px" }}>
                  <span style={{ background: `${statusColors[m.status]}22`,
                    color: statusColors[m.status], padding: "3px 10px",
                    borderRadius: 100, fontSize: 11, fontWeight: 600, letterSpacing: 0.5 }}>
                    {m.status.replace("_", " ").toUpperCase()}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

function FormFieldsManager({ apiFetch }) {
  const [fields, setFields] = useState([]);
  const [newField, setNewField] = useState({ fieldKey: "", label: "", fieldType: "text", isRequired: false });

  const load = async () => {
    const d = await apiFetch("/admin/form-fields");
    if (d) setFields(d.fields || []);
  };
  useEffect(() => { load(); }, []);

  const addField = async (e) => {
    e.preventDefault();
    await apiFetch("/admin/form-fields", { method: "POST", body: newField });
    setNewField({ fieldKey: "", label: "", fieldType: "text", isRequired: false });
    load();
  };

  const toggleActive = async (f) => {
    await apiFetch(`/admin/form-fields/${f.id}`, { method: "PUT", body: { isActive: !f.is_active } });
    load();
  };

  const deleteField = async (id) => {
    if (!confirm("Delete this field?")) return;
    await apiFetch(`/admin/form-fields/${id}`, { method: "DELETE" });
    load();
  };

  return (
    <div style={{ maxWidth: 600 }}>
      <h3 className="section-title">Configure Patient Form Fields</h3>
      <p style={{ color: "var(--text-muted)", fontSize: 13, marginBottom: 20 }}>
        These fields appear on the Data Entrant's patient visit form.
      </p>

      {/* Existing fields */}
      <div style={{ display: "flex", flexDirection: "column", gap: 8, marginBottom: 24 }}>
        {fields.map(f => (
          <div key={f.id} style={{ display: "flex", alignItems: "center", gap: 12,
            background: "var(--bg-card)", border: "1px solid var(--border)",
            borderRadius: 8, padding: "10px 16px",
            opacity: f.is_active ? 1 : 0.5 }}>
            <div style={{ flex: 1 }}>
              <span style={{ fontWeight: 600, fontSize: 13, color: "var(--text-primary)" }}>{f.label}</span>
              <span style={{ color: "var(--text-muted)", fontSize: 11, marginLeft: 10 }}>
                key: {f.field_key} • {f.field_type}{f.is_required ? " • required" : ""}
              </span>
            </div>
            <button onClick={() => toggleActive(f)}
              style={{ background: f.is_active ? "rgba(0,191,166,0.15)" : "rgba(255,255,255,0.05)",
                border: "none", borderRadius: 6, padding: "4px 12px", cursor: "pointer",
                color: f.is_active ? "#00BFA6" : "var(--text-muted)", fontSize: 12 }}>
              {f.is_active ? "Active" : "Inactive"}
            </button>
            <button onClick={() => deleteField(f.id)}
              style={{ background: "none", border: "none", color: "#F44336", cursor: "pointer", fontSize: 16 }}>
              ✕
            </button>
          </div>
        ))}
      </div>

      {/* Add new field */}
      <form onSubmit={addField} style={{ display: "flex", gap: 10, flexWrap: "wrap",
        background: "var(--bg-card)", border: "1px solid var(--border)", borderRadius: 10, padding: 16 }}>
        <input className="field-input" style={{ flex: "0 0 140px" }} placeholder="Field key (snake_case)"
          value={newField.fieldKey} required
          onChange={e => setNewField(f => ({ ...f, fieldKey: e.target.value }))} />
        <input className="field-input" style={{ flex: 1, minWidth: 140 }} placeholder="Display label"
          value={newField.label} required
          onChange={e => setNewField(f => ({ ...f, label: e.target.value }))} />
        <select className="field-input" style={{ flex: "0 0 110px" }} value={newField.fieldType}
          onChange={e => setNewField(f => ({ ...f, fieldType: e.target.value }))}>
          <option value="text">Text</option>
          <option value="number">Number</option>
          <option value="date">Date</option>
          <option value="select">Select</option>
        </select>
        <label style={{ display: "flex", alignItems: "center", gap: 8, color: "var(--text-muted)", fontSize: 13 }}>
          <input type="checkbox" checked={newField.isRequired}
            onChange={e => setNewField(f => ({ ...f, isRequired: e.target.checked }))} />
          Required
        </label>
        <button type="submit" className="btn-primary" style={{ background: "#7C6FFF" }}>+ Add</button>
      </form>
    </div>
  );
}

function UsersManager({ apiFetch }) {
  const [users, setUsers] = useState([]);
  const load = async () => {
    const d = await apiFetch("/admin/users");
    if (d) setUsers(d.users || []);
  };
  useEffect(() => { load(); }, []);

  const toggle = async (u) => {
    await apiFetch(`/admin/users/${u.id}`, { method: "PUT", body: { isActive: !u.is_active } });
    load();
  };

  return (
    <div>
      <h3 className="section-title">User Management</h3>
      <div style={{ display: "flex", flexDirection: "column", gap: 10 }}>
        {users.map(u => (
          <div key={u.id} className="record-row">
            <div>
              <div style={{ fontWeight: 600, fontSize: 13, color: "var(--text-primary)" }}>{u.name}</div>
              <div style={{ fontSize: 12, color: "var(--text-muted)" }}>{u.email} • {u.role}</div>
            </div>
            <button onClick={() => toggle(u)}
              style={{ background: u.is_active ? "rgba(0,191,166,0.15)" : "rgba(244,67,54,0.15)",
                border: "none", borderRadius: 6, padding: "4px 14px", cursor: "pointer",
                color: u.is_active ? "#00BFA6" : "#F44336", fontSize: 12, fontWeight: 600 }}>
              {u.is_active ? "Active" : "Disabled"}
            </button>
          </div>
        ))}
      </div>
    </div>
  );
}

// ═══════════════════════════════════════════════════════════════
// DHO PORTAL — Regional Multi-Pharmacy View
// ═══════════════════════════════════════════════════════════════
function DHOPortal() {
  const { user, logout, apiFetch } = useAuth();
  const [tab, setTab] = useState("map");

  return (
    <Shell role="dho" user={user} onLogout={logout}>
      <div style={{ display: "flex", borderBottom: "1px solid var(--border)" }}>
        {[["map","Regional View"], ["risk","Stockout Risk"], ["pharmacies","Pharmacies"],
          ["register","Register Pharmacy"]].map(([key, label]) => (
          <button key={key} onClick={() => setTab(key)}
            style={{ padding: "14px 22px", background: "none", border: "none",
              borderBottom: tab === key ? "2px solid #FFB547" : "2px solid transparent",
              color: tab === key ? "#FFB547" : "var(--text-muted)", fontFamily: "var(--font-display)",
              fontSize: 13, cursor: "pointer", fontWeight: tab === key ? 600 : 400, letterSpacing: 0.5 }}>
            {label}
          </button>
        ))}
      </div>
      <div style={{ padding: 24 }}>
        {tab === "map"        && <RegionalMapView apiFetch={apiFetch} />}
        {tab === "risk"       && <StockoutRiskView apiFetch={apiFetch} />}
        {tab === "pharmacies" && <PharmacyList apiFetch={apiFetch} />}
        {tab === "register"   && <RegisterPharmacy apiFetch={apiFetch} />}
      </div>
    </Shell>
  );
}

function RegionalMapView({ apiFetch }) {
  const [mapData, setMapData] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(async () => {
    const d = await apiFetch("/dho/regional-map");
    if (d) setMapData(d.mapData || []);
    setLoading(false);
  }, []);

  if (loading) return <Spinner />;

  const riskColor = { ok: "#00BFA6", caution: "#FFB547", warning: "#FF7043", critical: "#F44336" };

  return (
    <div>
      <h3 className="section-title">Regional Pharmacy Map</h3>
      {mapData.length === 0 ? (
        <div style={{ color: "var(--text-muted)", fontSize: 13, textAlign: "center", padding: 40 }}>
          No pharmacies with GPS coordinates registered yet.
          Add coordinates when registering a pharmacy.
        </div>
      ) : (
        <>
          {/* Legend */}
          <div style={{ display: "flex", gap: 16, marginBottom: 20, flexWrap: "wrap" }}>
            {Object.entries(riskColor).map(([k, c]) => (
              <div key={k} style={{ display: "flex", alignItems: "center", gap: 6 }}>
                <div style={{ width: 10, height: 10, borderRadius: "50%", background: c }} />
                <span style={{ fontSize: 12, color: "var(--text-muted)", textTransform: "capitalize" }}>{k}</span>
              </div>
            ))}
          </div>
          {/* Map placeholder — replace with Leaflet or Mapbox in production */}
          <div style={{ background: "var(--bg-card)", border: "1px solid var(--border)",
            borderRadius: 12, padding: 20, marginBottom: 20 }}>
            <div style={{ color: "var(--text-muted)", fontSize: 12, marginBottom: 16, textAlign: "center" }}>
              Map view — integrate Leaflet.js with coordinates for full GIS display
            </div>
            {/* Visual grid substitute */}
            <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))", gap: 12 }}>
              {mapData.map(p => (
                <div key={p.id} style={{ background: "var(--bg-deep)", borderRadius: 10,
                  border: `1px solid ${riskColor[p.risk_level]}`,
                  padding: "14px 16px" }}>
                  <div style={{ display: "flex", alignItems: "center", gap: 8, marginBottom: 8 }}>
                    <div style={{ width: 8, height: 8, borderRadius: "50%",
                      background: riskColor[p.risk_level], flexShrink: 0 }} />
                    <span style={{ fontWeight: 600, fontSize: 13, color: "var(--text-primary)" }}>{p.name}</span>
                  </div>
                  <div style={{ fontSize: 11, color: "var(--text-muted)", lineHeight: 1.6 }}>
                    {p.district}, {p.region}<br />
                    {p.expiring_count || 0} expiring • {p.low_stock_count || 0} low stock<br />
                    {p.visits_7d || 0} visits (7d)
                  </div>
                </div>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
}

function StockoutRiskView({ apiFetch }) {
  const [risks, setRisks] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(async () => {
    const d = await apiFetch("/analytics/stockout-risk");
    if (d) setRisks(d.risks || []);
    setLoading(false);
  }, []);

  if (loading) return <Spinner />;

  const riskColors = { critical: "#F44336", warning: "#FF7043", expiring: "#FFB547",
    low_stock: "#FF9800", ok: "#00BFA6" };

  const grouped = risks.reduce((acc, r) => {
    const key = r.risk_level;
    if (!acc[key]) acc[key] = [];
    acc[key].push(r); return acc;
  }, {});

  return (
    <div>
      <h3 className="section-title">AI Stockout Risk Predictions</h3>
      <p style={{ color: "var(--text-muted)", fontSize: 13, marginBottom: 24 }}>
        Based on average daily consumption over the last 30 days vs remaining stock.
      </p>
      {["critical", "warning", "expiring", "low_stock"].filter(k => grouped[k]?.length).map(level => (
        <div key={level} style={{ marginBottom: 24 }}>
          <h4 style={{ color: riskColors[level], fontSize: 12, letterSpacing: 1,
            textTransform: "uppercase", marginBottom: 12 }}>
            {level.replace("_", " ")} ({grouped[level].length})
          </h4>
          <div style={{ display: "flex", flexDirection: "column", gap: 8 }}>
            {grouped[level].map(r => (
              <div key={r.id} className="record-row">
                <div>
                  <div style={{ fontWeight: 600, fontSize: 13, color: "var(--text-primary)" }}>{r.name}</div>
                  <div style={{ fontSize: 12, color: "var(--text-muted)" }}>
                    {r.pharmacy_name} • {r.quantity_remaining} remaining •{" "}
                    avg {r.avg_daily_use}/day
                  </div>
                </div>
                <div style={{ textAlign: "right" }}>
                  <div style={{ color: riskColors[level], fontWeight: 700, fontSize: 18,
                    fontFamily: "var(--font-display)" }}>
                    {r.estimated_days_remaining != null ? `${r.estimated_days_remaining}d` : "—"}
                  </div>
                  <div style={{ color: "var(--text-muted)", fontSize: 11 }}>estimated</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      ))}
      {risks.every(r => r.risk_level === "ok") && (
        <div style={{ color: "#00BFA6", textAlign: "center", padding: 40, fontSize: 14 }}>
          ✓ All pharmacies show adequate stock levels
        </div>
      )}
    </div>
  );
}

function PharmacyList({ apiFetch }) {
  const [pharmacies, setPharmacies] = useState([]);
  const [loading, setLoading] = useState(true);
  useEffect(async () => {
    const d = await apiFetch("/dho/pharmacies");
    if (d) setPharmacies(d.pharmacies || []);
    setLoading(false);
  }, []);

  if (loading) return <Spinner />;

  return (
    <div>
      <h3 className="section-title">Registered Pharmacies ({pharmacies.length})</h3>
      <div style={{ display: "flex", flexDirection: "column", gap: 12 }}>
        {pharmacies.map(p => (
          <div key={p.id} style={{ background: "var(--bg-card)", border: "1px solid var(--border)",
            borderRadius: 12, padding: 18, display: "flex", gap: 16, alignItems: "flex-start" }}>
            <div style={{ flex: 1 }}>
              <div style={{ fontWeight: 700, fontSize: 14, color: "var(--text-primary)",
                marginBottom: 6, fontFamily: "var(--font-display)" }}>{p.name}</div>
              <div style={{ fontSize: 12, color: "var(--text-muted)", lineHeight: 1.7 }}>
                {p.district}, {p.region} •{" "}
                {p.admin_count} admin(s) • {p.de_count} DE(s)
              </div>
            </div>
            <div style={{ textAlign: "right", fontSize: 12, color: "var(--text-muted)" }}>
              <div style={{ color: "#FFB547", fontWeight: 600, fontSize: 18,
                fontFamily: "var(--font-display)" }}>{p.medicine_types}</div>
              medicine types
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}

function RegisterPharmacy({ apiFetch }) {
  const [form, setForm] = useState({
    name: "", region: "", district: "", lat: "", lng: "",
    contactPhone: "", whatsappNumber: "",
    adminName: "", adminEmail: "", adminPassword: ""
  });
  const [status, setStatus] = useState(null);
  const set = (k, v) => setForm(f => ({ ...f, [k]: v }));

  const submit = async (e) => {
    e.preventDefault(); setStatus("loading");
    try {
      await apiFetch("/dho/pharmacies", { method: "POST", body: form });
      setForm({ name: "", region: "", district: "", lat: "", lng: "", contactPhone: "",
        whatsappNumber: "", adminName: "", adminEmail: "", adminPassword: "" });
      setStatus("success");
      setTimeout(() => setStatus(null), 4000);
    } catch (e) { setStatus("error:" + e.message); }
  };

  return (
    <div style={{ maxWidth: 580 }}>
      <h3 className="section-title">Register New Pharmacy</h3>
      <form onSubmit={submit}>
        <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16, marginBottom: 20 }}>
          <div className="form-group" style={{ gridColumn: "1/-1" }}>
            <label className="field-label">Pharmacy Name *</label>
            <input className="field-input" value={form.name} onChange={e => set("name", e.target.value)} required />
          </div>
          {[["region","Region"], ["district","District"]].map(([k, l]) => (
            <div key={k} className="form-group">
              <label className="field-label">{l}</label>
              <input className="field-input" value={form[k]} onChange={e => set(k, e.target.value)} />
            </div>
          ))}
          {[["lat","GPS Latitude"], ["lng","GPS Longitude"]].map(([k, l]) => (
            <div key={k} className="form-group">
              <label className="field-label">{l}</label>
              <input className="field-input" type="number" step="any" value={form[k]}
                onChange={e => set(k, e.target.value)} />
            </div>
          ))}
          <div className="form-group">
            <label className="field-label">Contact Phone</label>
            <input className="field-input" value={form.contactPhone} onChange={e => set("contactPhone", e.target.value)} />
          </div>
          <div className="form-group">
            <label className="field-label">WhatsApp Number</label>
            <input className="field-input" value={form.whatsappNumber}
              onChange={e => set("whatsappNumber", e.target.value)} placeholder="+256700000000" />
          </div>
        </div>

        <div style={{ borderTop: "1px solid var(--border)", paddingTop: 20, marginBottom: 20 }}>
          <h4 style={{ color: "var(--text-muted)", fontSize: 12, letterSpacing: 1, marginBottom: 16 }}>
            INITIAL ADMIN ACCOUNT
          </h4>
          <div style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: 16 }}>
            {[["adminName","Admin Full Name *", false], ["adminEmail","Admin Email *", false],
              ["adminPassword","Temporary Password *", true]].map(([k, l, isPwd]) => (
              <div key={k} className="form-group" style={{ gridColumn: isPwd ? "1/-1" : undefined }}>
                <label className="field-label">{l}</label>
                <input className="field-input" type={isPwd ? "password" : "text"}
                  value={form[k]} onChange={e => set(k, e.target.value)} required={l.includes("*")} />
              </div>
            ))}
          </div>
        </div>

        {status === "success" && <div className="alert-success" style={{ marginBottom: 12 }}>Pharmacy registered successfully ✓</div>}
        {status?.startsWith("error") && <div className="alert-danger" style={{ marginBottom: 12 }}>{status.replace("error:", "")}</div>}

        <button type="submit" className="btn-primary" disabled={status === "loading"}
          style={{ background: "#FFB547", color: "#1A1535" }}>
          {status === "loading" ? "Registering…" : "Register Pharmacy"}
        </button>
      </form>
    </div>
  );
}

// ═══════════════════════════════════════════════════════════════
// SHARED COMPONENTS
// ═══════════════════════════════════════════════════════════════
function Shell({ role, user, onLogout, children }) {
  const { fullReset } = useAuth();
  const meta = ROLE_META[role];

  return (
    <div style={{ display: "flex", flexDirection: "column", minHeight: "100vh",
      background: "var(--bg-deep)" }}>
      {/* Top bar */}
      <header style={{ display: "flex", alignItems: "center", padding: "0 24px",
        height: 56, borderBottom: "1px solid var(--border)",
        background: "var(--bg-card)", flexShrink: 0 }}>
        <div style={{ display: "flex", alignItems: "center", gap: 10, flex: 1 }}>
          <span style={{ fontSize: 18, color: meta.color }}>✚</span>
          <span style={{ fontFamily: "var(--font-display)", fontWeight: 700, fontSize: 14,
            letterSpacing: 1, color: "var(--text-primary)" }}>MED PREDICT</span>
          <span style={{ background: meta.bg, color: meta.color, fontSize: 10, fontWeight: 700,
            padding: "2px 8px", borderRadius: 100, letterSpacing: 1, marginLeft: 4 }}>
            {meta.short}
          </span>
        </div>
        <div style={{ display: "flex", alignItems: "center", gap: 16 }}>
          <span style={{ color: "var(--text-muted)", fontSize: 12 }}>{user?.name}</span>
          <button onClick={onLogout}
            style={{ background: "var(--bg-deep)", border: "1px solid var(--border)",
              borderRadius: 6, padding: "5px 14px", color: "var(--text-muted)",
              cursor: "pointer", fontSize: 12 }}>
            Sign out
          </button>
          <button onClick={fullReset}
            style={{ background: "none", border: "none", color: "var(--text-muted)",
              cursor: "pointer", fontSize: 11 }}>
            Switch portal
          </button>
        </div>
      </header>
      <main style={{ flex: 1, overflowY: "auto" }}>{children}</main>
    </div>
  );
}

function Spinner() {
  return (
    <div style={{ display: "flex", justifyContent: "center", padding: 60 }}>
      <div style={{ width: 32, height: 32, border: "2px solid var(--border)",
        borderTopColor: "#7C6FFF", borderRadius: "50%",
        animation: "spin 0.7s linear infinite" }} />
    </div>
  );
}

// ─── Global CSS ───────────────────────────────────────────────
const GLOBAL_CSS = `
  @import url('https://fonts.googleapis.com/css2?family=Syne:wght@400;600;700;800&family=DM+Sans:wght@300;400;500&display=swap');

  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }

  :root {
    --font-display: 'Syne', sans-serif;
    --font-body: 'DM Sans', sans-serif;
    --bg-deep: #0A0E14;
    --bg-card: #111722;
    --bg-surface: #161D2A;
    --text-primary: #E8EAF0;
    --text-muted: #6B7280;
    --border: #1E2840;
  }

  body { font-family: var(--font-body); background: var(--bg-deep); color: var(--text-primary);
    line-height: 1.5; -webkit-font-smoothing: antialiased; }

  .screen-center { display: flex; align-items: center; justify-content: center; padding: 40px 0; }

  .role-card { display: flex; align-items: center; padding: 18px 22px; background: var(--bg-card);
    border: 1px solid var(--border); border-radius: 12px; cursor: pointer;
    transition: all 0.18s ease; width: 100%; }

  .field-input { width: 100%; background: var(--bg-surface); border: 1px solid var(--border);
    border-radius: 8px; padding: 10px 14px; color: var(--text-primary);
    font-family: var(--font-body); font-size: 14px; outline: none; transition: border-color 0.15s; }
  .field-input:focus { border-color: rgba(124,111,255,0.5); }
  .field-label { display: block; font-size: 11px; font-weight: 600; letter-spacing: 0.8px;
    color: var(--text-muted); margin-bottom: 6px; text-transform: uppercase; }
  .form-group { display: flex; flex-direction: column; }

  .btn-primary { background: #7C6FFF; color: white; border: none; border-radius: 8px;
    padding: 10px 20px; font-family: var(--font-display); font-size: 14px; font-weight: 600;
    cursor: pointer; transition: opacity 0.15s; letter-spacing: 0.3px; }
  .btn-primary:hover:not(:disabled) { opacity: 0.87; }
  .btn-primary:disabled { opacity: 0.5; cursor: not-allowed; }

  .btn-sm-success { background: rgba(0,191,166,0.15); color: #00BFA6; border: 1px solid rgba(0,191,166,0.3);
    border-radius: 6px; padding: 5px 12px; cursor: pointer; font-size: 12px; font-weight: 600; }
  .btn-sm-danger { background: rgba(244,67,54,0.12); color: #F44336; border: 1px solid rgba(244,67,54,0.3);
    border-radius: 6px; padding: 5px 12px; cursor: pointer; font-size: 12px; font-weight: 600; }

  .alert-success { background: rgba(0,191,166,0.1); border: 1px solid rgba(0,191,166,0.3);
    color: #00BFA6; border-radius: 8px; padding: 12px 16px; font-size: 13px; }
  .alert-danger { background: rgba(244,67,54,0.1); border: 1px solid rgba(244,67,54,0.3);
    color: #F44336; border-radius: 8px; padding: 12px 16px; font-size: 13px; }

  .section-title { font-family: var(--font-display); font-size: 16px; font-weight: 700;
    color: var(--text-primary); margin-bottom: 20px; letter-spacing: 0.5px; }

  .record-row { display: flex; align-items: center; justify-content: space-between;
    background: var(--bg-card); border: 1px solid var(--border); border-radius: 8px;
    padding: 12px 16px; }

  select.field-input { appearance: none; cursor: pointer; }

  @keyframes spin { to { transform: rotate(360deg); } }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(8px); } to { opacity: 1; transform: none; } }
`;
