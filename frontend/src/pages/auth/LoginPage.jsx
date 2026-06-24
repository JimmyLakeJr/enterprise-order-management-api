import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { API_BASE_URL, getMessage } from "../../api/apiClient";
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
    if (!form.email.trim()) return "Vui lòng nhập email.";
    if (!form.email.includes("@")) return "Email không hợp lệ.";
    if (!form.password) return "Vui lòng nhập mật khẩu.";
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

  function handleGoogleLogin() {
    window.location.assign(`${API_BASE_URL}/auth/google/login`);
  }

  return (
    <div className="auth-shell">
      <GlassCard strong className="auth-card">
        <h1>Đăng nhập</h1>
        <p className="muted">Truy cập hệ thống để tạo đơn hàng, theo dõi tồn kho và quản lý tài khoản.</p>
        <ErrorMessage message={error} />
        <form className="form-stack" onSubmit={handleSubmit}>
          <Input
            label="Email"
            type="email"
            autoComplete="email"
            value={form.email}
            onChange={(event) => setForm({ ...form, email: event.target.value })}
          />
          <Input
            label="Mật khẩu"
            type="password"
            autoComplete="current-password"
            value={form.password}
            onChange={(event) => setForm({ ...form, password: event.target.value })}
          />
          <Button type="submit" disabled={loading}>
            {loading ? "Đang đăng nhập..." : "Đăng nhập"}
          </Button>
          <Button type="button" variant="secondary" onClick={handleGoogleLogin} disabled={loading}>
            Đăng nhập bằng Google
          </Button>
        </form>
        <p className="muted">
          Chưa có tài khoản? <Link to="/register">Đăng ký</Link>
        </p>
      </GlassCard>
    </div>
  );
}
