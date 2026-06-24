import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import ErrorMessage from "../../components/common/ErrorMessage";
import GlassCard from "../../components/common/GlassCard";
import { ROLES } from "../../constants/domain";
import { useAuth } from "../../contexts/AuthContext";

function parseHashParams() {
  const raw = window.location.hash.startsWith("#") ? window.location.hash.slice(1) : "";
  return new URLSearchParams(raw);
}

export default function GoogleAuthCallbackPage() {
  const navigate = useNavigate();
  const { completeOAuthLogin } = useAuth();
  const [error, setError] = useState("");

  useEffect(() => {
    let cancelled = false;

    async function completeLogin() {
      const params = parseHashParams();
      const status = params.get("status");
      const accessToken = params.get("access_token");
      const refreshToken = params.get("refresh_token");
      const message = params.get("error");

      if (status !== "success" || !accessToken || !refreshToken) {
        setError(message || "Đăng nhập Google không thành công.");
        window.history.replaceState(null, "", "/auth/google/callback");
        return;
      }

      try {
        const user = await completeOAuthLogin({ accessToken, refreshToken });
        if (cancelled) return;
        window.history.replaceState(null, "", "/auth/google/callback");
        navigate(user?.role === ROLES.ADMIN ? "/admin" : "/", { replace: true });
      } catch (err) {
        if (cancelled) return;
        setError(err?.response?.data?.message || err?.message || "Không thể hoàn tất đăng nhập Google.");
      }
    }

    void completeLogin();
    return () => {
      cancelled = true;
    };
  }, [completeOAuthLogin, navigate]);

  return (
    <div className="auth-shell">
      <GlassCard strong className="auth-card">
        <h1>Đăng nhập bằng Google</h1>
        {error ? (
          <>
            <p className="muted">Không thể hoàn tất đăng nhập.</p>
            <ErrorMessage message={error} />
          </>
        ) : (
          <p className="muted">Đang xác thực tài khoản Google và nạp phiên đăng nhập...</p>
        )}
      </GlassCard>
    </div>
  );
}

