import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { getMessage } from "../../api/apiClient";
import Button from "../../components/common/Button";
import ErrorMessage from "../../components/common/ErrorMessage";
import GlassCard from "../../components/common/GlassCard";
import Input from "../../components/common/Input";
import { ROLES } from "../../constants/domain";
import { useAuth } from "../../contexts/AuthContext";

export default function LoginPage() {
  const { login, loading } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [form, setForm] = useState({ email: "", password: "" });
  const [error, setError] = useState("");

  function validate() {
    if (!form.email.trim()) return "Email is required";
    if (!form.email.includes("@")) return "Email is invalid";
    if (!form.password) return "Password is required";
    return "";
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setError("");

    const validationError = validate();
    if (validationError) {
      setError(validationError);
      return;
    }

    try {
      const data = await login(form);
      const redirectTo = location.state?.from || (data.user?.role === ROLES.ADMIN ? "/admin" : "/");
      navigate(redirectTo, { replace: true });
    } catch (err) {
      setError(getMessage(err));
    }
  }

  return (
    <div className="auth-shell">
      <GlassCard strong className="auth-card">
        <h1>Login</h1>
        <p className="muted">Sign in to create orders and manage your account.</p>
        <ErrorMessage message={error} />
        <form className="form-stack" onSubmit={handleSubmit}>
          <Input label="Email" type="email" autoComplete="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
          <Input label="Password" type="password" autoComplete="current-password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} />
          <Button type="submit" disabled={loading}>{loading ? "Logging in..." : "Login"}</Button>
        </form>
        <p className="muted">
          Do not have an account? <Link to="/register">Register</Link>
        </p>
      </GlassCard>
    </div>
  );
}
