import { Link } from "react-router-dom";
import Button from "../../components/common/Button";
import GlassCard from "../../components/common/GlassCard";

export default function ForbiddenPage() {
  return (
    <GlassCard strong>
      <h1>Không có quyền truy cập</h1>
      <p className="muted">Tài khoản của bạn không có quyền mở trang này.</p>
      <Link to="/">
        <Button>Về trang sản phẩm</Button>
      </Link>
    </GlassCard>
  );
}
