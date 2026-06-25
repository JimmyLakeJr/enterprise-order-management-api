import { useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getProductById } from "../../api/productApi";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Loading from "../../components/common/Loading";
import { useAuth } from "../../contexts/AuthContext";
import { useCart } from "../../contexts/CartContext";
import { useAsync } from "../../hooks/useAsync";
import { formatCurrency } from "../../utils/formatCurrency";
import { resolveAssetUrl } from "../../utils/resolveAssetUrl";

export default function ProductDetailPage() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated } = useAuth();
  const { addToCart } = useCart();
  const { data: product, loading, error, reload } = useAsync(() => getProductById(id), [id]);
  const [quantityState, setQuantityState] = useState({ productId: id, value: 1 });
  const [messageState, setMessageState] = useState({ productId: id, value: "" });
  const [failedImageUrl, setFailedImageUrl] = useState("");
  const quantity = quantityState.productId === id ? quantityState.value : 1;
  const message = messageState.productId === id ? messageState.value : "";

  if (loading) return <Loading label="Đang tải chi tiết sản phẩm..." variant="detail" count={2} />;

  if (error) {
    return (
      <div className="detail-error">
        <ErrorMessage message={`Không tải được sản phẩm: ${error}`} />
        <div className="actions">
          <Button onClick={reload}>Thử lại</Button>
          <Link className="btn btn-secondary" to="/">Về danh sách</Link>
        </div>
      </div>
    );
  }

  if (!product) return <ErrorMessage message="Không tìm thấy sản phẩm." />;

  const parsedStock = Number(product.stock);
  const hasStock = Number.isFinite(parsedStock);
  const stock = hasStock ? Math.max(0, parsedStock) : null;
  const isOutOfStock = hasStock && stock <= 0;
  const numericQuantity = Number(quantity);
  const quantityError =
    quantity === "" || !Number.isInteger(numericQuantity) || numericQuantity <= 0
      ? "Số lượng phải là số nguyên lớn hơn 0."
      : hasStock && numericQuantity > stock
        ? `Số lượng không được vượt quá tồn kho (${stock}).`
        : "";
  const categoryName = product.category?.name || product.category_name || "Chưa phân loại";
  const imageSrc = resolveAssetUrl(product.image_url);

  function handleQuantityChange(event) {
    setQuantityState({ productId: id, value: event.target.value });
    setMessageState({ productId: id, value: "" });
  }

  function handleAddToCart() {
    if (quantityError || isOutOfStock) return;
    if (!isAuthenticated) {
      navigate("/login", { state: { from: `/products/${id}` } });
      return;
    }

    addToCart(product, numericQuantity);
    setMessageState({ productId: id, value: `Đã thêm ${numericQuantity} x "${product.name}" vào giỏ hàng.` });
  }

  return (
    <>
      <Link to="/" className="back-link">← Quay lại danh sách sản phẩm</Link>

      <Card className="product-detail-card">
        <div className="product-detail">
          <div className="product-detail-image">
            {imageSrc && failedImageUrl !== imageSrc ? (
              <img src={imageSrc} alt={product.name} onError={() => setFailedImageUrl(imageSrc)} />
            ) : (
              <span className="image-placeholder image-placeholder-large">
                <span aria-hidden="true">▧</span>
                Sản phẩm chưa có ảnh
              </span>
            )}
          </div>

          <div className="product-detail-content">
            <Badge tone="primary">{categoryName}</Badge>
            <div>
              <h1>{product.name}</h1>
              <p className="product-description">{product.description || "Sản phẩm này chưa có mô tả."}</p>
            </div>

            <strong className="product-detail-price">{formatCurrency(product.price)}</strong>
            <div className={isOutOfStock ? "stock-danger stock-panel" : "stock-panel stock-available"}>
              <span>{isOutOfStock ? "Hết hàng" : "Sẵn sàng đặt hàng"}</span>
              <strong>{hasStock ? `${stock} sản phẩm trong kho` : "Còn hàng"}</strong>
            </div>

            <div className="purchase-panel">
              <Input
                label="Số lượng"
                type="number"
                inputMode="numeric"
                min="1"
                max={hasStock ? stock : undefined}
                step="1"
                value={quantity}
                onChange={handleQuantityChange}
                disabled={isOutOfStock}
                error={quantityError}
                fieldClassName="quantity-box"
              />
              <Button className="add-to-cart-button" onClick={handleAddToCart} disabled={isOutOfStock || Boolean(quantityError)}>
                {isOutOfStock ? "Tạm hết hàng" : isAuthenticated ? "Thêm vào giỏ" : "Đăng nhập để mua"}
              </Button>
            </div>

            {message && <div className="alert alert-success" role="status">{message}</div>}
          </div>
        </div>
      </Card>
    </>
  );
}
