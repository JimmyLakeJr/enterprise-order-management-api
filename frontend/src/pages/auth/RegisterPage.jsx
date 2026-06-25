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
    phone: "",
    password: "",
    confirm_password: "",
  });
  const [error, setError] = useState("");

  function validate() {
    if (form.full_name.trim().length < 2) return "Họ tên phải có ít nhất 2 ký tự.";
    if (!form.email.trim() && !form.phone.trim()) return "Cần nhập email hoặc số điện thoại.";
    if (form.email.trim() && !form.email.includes("@")) return "Email không hợp lệ.";

    const normalizedPhone = form.phone.replace(/[^\d+]/g, "");
    if (form.phone.trim() && normalizedPhone.length < 9) return "Số điện thoại không hợp lệ.";

    if (form.password.length < 6) return "Mật khẩu phải có ít nhất 6 ký tự.";
    if (form.password !== form.confirm_password) return "Mật khẩu xác nhận không khớp.";
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
        phone: form.phone,
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
        <h1>Đăng ký</h1>
        <p className="muted">Tạo tài khoản bằng email hoặc số điện thoại để bắt đầu đặt hàng và theo dõi đơn hàng.</p>
        <ErrorMessage message={error} />
        <form className="form-stack" onSubmit={handleSubmit}>
          <Input
            label="Họ tên"
            autoComplete="name"
            value={form.full_name}
            onChange={(event) => setForm({ ...form, full_name: event.target.value })}
          />
          <Input
            label="Email"
            type="email"
            autoComplete="email"
            value={form.email}
            onChange={(event) => setForm({ ...form, email: event.target.value })}
          />
          <Input
            label="Số điện thoại"
            type="tel"
            autoComplete="tel"
            value={form.phone}
            onChange={(event) => setForm({ ...form, phone: event.target.value })}
          />
          <Input
            label="Mật khẩu"
            type="password"
            autoComplete="new-password"
            value={form.password}
            onChange={(event) => setForm({ ...form, password: event.target.value })}
          />
          <Input
            label="Xác nhận mật khẩu"
            type="password"
            autoComplete="new-password"
            value={form.confirm_password}
            onChange={(event) => setForm({ ...form, confirm_password: event.target.value })}
          />
          <Button type="submit" disabled={loading}>
            {loading ? "Đang tạo tài khoản..." : "Đăng ký"}
          </Button>
        </form>
        <p className="muted">
          Đã có tài khoản? <Link to="/login">Đăng nhập</Link>
        </p>
      </GlassCard>
    </div>
  );
}
