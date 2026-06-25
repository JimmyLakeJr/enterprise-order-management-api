import { useState } from "react";
import { Link } from "react-router-dom";
import Button from "../common/Button";
import Card from "../common/Card";
import { formatCurrency } from "../../utils/formatCurrency";
import { resolveAssetUrl } from "../../utils/resolveAssetUrl";

export default function ProductCard({ product, onAdd, addLabel = "Thêm vào giỏ" }) {
  const [failedImageUrl, setFailedImageUrl] = useState("");
  const categoryName = product.category?.name || product.category_name || "Chưa phân loại";
  const stock = Number(product.stock);
  const hasStock = Number.isFinite(stock);
  const isOutOfStock = hasStock && stock <= 0;
  const imageSrc = resolveAssetUrl(product.image_url);

  return (
    <Card className="product-card" as="article">
      <Link to={`/products/${product.id}`} className="product-image" aria-label={`Xem ${product.name}`}>
        {imageSrc && failedImageUrl !== imageSrc ? (
          <img src={imageSrc} alt={product.name} onError={() => setFailedImageUrl(imageSrc)} />
        ) : (
          <span className="image-placeholder" aria-label="Sản phẩm chưa có ảnh">
            <span aria-hidden="true">▧</span>
            Chưa có ảnh
          </span>
        )}
      </Link>

      <div className="product-body">
        <span className="product-category">{categoryName}</span>
        <Link to={`/products/${product.id}`} className="product-title">
          {product.name}
        </Link>
        <strong className="product-price">{formatCurrency(product.price)}</strong>
        <span className={isOutOfStock ? "stock-danger" : "product-stock"}>
          {isOutOfStock ? "Hết hàng" : hasStock ? `Còn ${stock} sản phẩm` : "Còn hàng"}
        </span>
      </div>

      <div className="product-actions">
        <Link className="btn btn-secondary" to={`/products/${product.id}`}>
          Xem chi tiết
        </Link>
        <Button onClick={() => onAdd(product, 1)} disabled={isOutOfStock} aria-label={`Thêm ${product.name} vào giỏ`}>
          {addLabel}
        </Button>
      </div>
    </Card>
  );
}
