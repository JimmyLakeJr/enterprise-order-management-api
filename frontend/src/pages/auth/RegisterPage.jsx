import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { getMessage } from "../../api/apiClient";
import Button from "../../components/common/Button";
import ErrorMessage from "../../components/common/ErrorMessage";
import GlassCard from "../../components/common/GlassCard";
import Input from "../../components/common/Input";
import { useAuth } from "../../contexts/AuthContext";

export default function RegisterPage() {
  const { register, loading } = useAuth();
  const navigate = useNavigate();
  const [form, setForm] = useState({
    full_name: "",
    email: "",
    password: "",
    confirm_password: "",
  });
  const [error, setError] = useState("");

  function validate() {
    if (form.full_name.trim().length < 2) return "Full name must have at least 2 characters";
    if (!form.email.includes("@")) return "Email is invalid";
    if (form.password.length < 6) return "Password must have at least 6 characters";
    if (form.password !== form.confirm_password) return "Confirm password does not match";
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
      await register({
        name: form.full_name,
        email: form.email,
        password: form.password,
      });
      navigate("/", { replace: true });
    } catch (err) {
      setError(getMessage(err));
    }
  }

  return (
    <div className="auth-shell">
      <GlassCard strong className="auth-card">
        <h1>Register</h1>
        <p className="muted">Create a user account for the demo client.</p>
        <ErrorMessage message={error} />
        <form className="form-stack" onSubmit={handleSubmit}>
          <Input label="Full name" autoComplete="name" value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} />
          <Input label="Email" type="email" autoComplete="email" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
          <Input label="Password" type="password" autoComplete="new-password" value={form.password} onChange={(e) => setForm({ ...form, password: e.target.value })} />
          <Input label="Confirm password" type="password" autoComplete="new-password" value={form.confirm_password} onChange={(e) => setForm({ ...form, confirm_password: e.target.value })} />
          <Button type="submit" disabled={loading}>{loading ? "Creating account..." : "Register"}</Button>
        </form>
        <p className="muted">
          Already have an account? <Link to="/login">Login</Link>
        </p>
      </GlassCard>
    </div>
  );
}
