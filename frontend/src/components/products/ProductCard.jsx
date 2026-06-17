import { Link } from "react-router-dom";
import Button from "../common/Button";
import Card from "../common/Card";
import { formatCurrency } from "../../utils/format";

export default function ProductCard({ product, onAdd }) {
  const categoryName = product.category?.name || product.category_name;
  const isOutOfStock = Number(product.stock) <= 0;

  return (
    <Card className="product-card">
      <div className="product-image">
        {product.image_url ? <img src={product.image_url} alt={product.name} /> : <span>Không có ảnh</span>}
      </div>

      <div className="product-body">
        <Link to={`/products/${product.id}`} className="product-title">
          {product.name}
        </Link>
        {categoryName && <span className="muted">{categoryName}</span>}
        <strong>{formatCurrency(product.price)}</strong>
        <span className={isOutOfStock ? "stock-danger" : "muted"}>Tồn kho: {product.stock}</span>
      </div>

      <div className="product-actions">
        <Link className="btn btn-secondary" to={`/products/${product.id}`}>
          Xem chi tiết
        </Link>
        <Button onClick={() => onAdd(product, 1)} disabled={isOutOfStock}>
          Thêm vào giỏ
        </Button>
      </div>
    </Card>
  );
}
