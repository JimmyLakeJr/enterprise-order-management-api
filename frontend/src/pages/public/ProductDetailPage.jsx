import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { getProductById } from "../../api/productApi";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Loading from "../../components/common/Loading";
import { useCart } from "../../contexts/CartContext";
import { useAsync } from "../../hooks/useAsync";
import { formatCurrency } from "../../utils/format";

export default function ProductDetailPage() {
  const { id } = useParams();
  const { addToCart } = useCart();
  const { data: product, loading, error } = useAsync(() => getProductById(id), [id]);
  const [quantity, setQuantity] = useState(1);
  const [message, setMessage] = useState("");

  if (loading) return <Loading />;
  if (error) return <ErrorMessage message={error} />;
  if (!product) return <ErrorMessage message="Không tìm thấy sản phẩm." />;

  const stock = Number(product.stock || 0);
  const categoryName = product.category?.name || product.category_name || "Sản phẩm";
  const quantityError =
    Number(quantity) <= 0 ? "Số lượng phải lớn hơn 0." : Number(quantity) > stock ? "Số lượng không được vượt tồn kho." : "";

  function handleQuantityChange(event) {
    const nextQuantity = Number(event.target.value);
    if (nextQuantity > stock) {
      setQuantity(stock);
      return;
    }
    setQuantity(nextQuantity);
  }

  function handleAddToCart() {
    if (quantityError || stock <= 0) return;
    addToCart(product, Number(quantity));
    setMessage(`Đã thêm ${quantity} sản phẩm vào giỏ hàng.`);
  }

  return (
    <Card>
      <div className="product-detail">
        <div className="product-detail-image">
          {product.image_url ? <img src={product.image_url} alt={product.name} /> : <span>Không có ảnh</span>}
        </div>

        <div className="form-stack">
          <div>
            <Link to="/" className="muted">
              Quay lại danh sách
            </Link>
            <h1>{product.name}</h1>
            <Badge tone="primary">{categoryName}</Badge>
          </div>

          <p>{product.description || "Chưa có mô tả cho sản phẩm này."}</p>
          <h2>{formatCurrency(product.price)}</h2>
          <p className={stock > 0 ? "muted" : "stock-danger"}>Tồn kho: {stock}</p>

          <div className="quantity-box">
            <Input
              label="Số lượng"
              type="number"
              min="1"
              max={stock}
              value={quantity}
              onChange={handleQuantityChange}
              disabled={stock <= 0}
              error={quantityError}
            />
          </div>

          {message && <div className="alert alert-success">{message}</div>}

          <Button onClick={handleAddToCart} disabled={stock <= 0 || Boolean(quantityError)}>
            Thêm vào giỏ
          </Button>
        </div>
      </div>
    </Card>
  );
}
