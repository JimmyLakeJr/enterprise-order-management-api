import { Link } from "react-router-dom";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";

export default function NotFoundPage() {
  return (
    <Card>
      <h1>404</h1>
      <p className="muted">Trang bạn tìm không tồn tại.</p>
      <Link to="/">
        <Button>Về trang sản phẩm</Button>
      </Link>
    </Card>
  );
}
