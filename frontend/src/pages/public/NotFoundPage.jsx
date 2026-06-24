import { Link } from "react-router-dom";
import Button from "../../components/common/Button";
import GlassCard from "../../components/common/GlassCard";

export default function NotFoundPage() {
  return (
    <main className="container page">
      <GlassCard strong>
        <h1>404</h1>
        <p className="muted">Trang bạn tìm không tồn tại.</p>
        <Link to="/products">
          <Button>Về trang sản phẩm</Button>
        </Link>
      </GlassCard>
    </main>
  );
}
