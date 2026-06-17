import { useEffect, useState } from "react";
import { getCategories } from "../../api/categoryApi";
import { getProducts } from "../../api/productApi";
import Button from "../../components/common/Button";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Loading from "../../components/common/Loading";
import Select from "../../components/common/Select";
import ProductCard from "../../components/products/ProductCard";
import { useCart } from "../../contexts/CartContext";

const DEFAULT_FILTERS = {
  page: 1,
  limit: 12,
  keyword: "",
  category_id: "",
  min_price: "",
  max_price: "",
};

export default function ProductListPage() {
  const { addToCart } = useCart();
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [meta, setMeta] = useState(null);
  const [filters, setFilters] = useState(DEFAULT_FILTERS);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");

  useEffect(() => {
    getCategories()
      .then(setCategories)
      .catch(() => setCategories([]));
  }, []);

  useEffect(() => {
    loadProducts(filters);
  }, [filters.page]);

  async function loadProducts(nextFilters) {
    setLoading(true);
    setError("");
    try {
      const result = await getProducts(nextFilters);
      setProducts(result.data);
      setMeta(result.meta);
    } catch {
      setError("Không tải được danh sách sản phẩm. Vui lòng kiểm tra backend API.");
    } finally {
      setLoading(false);
    }
  }

  function handleChange(event) {
    const { name, value } = event.target;
    setFilters((current) => ({ ...current, [name]: value }));
  }

  function handleSubmit(event) {
    event.preventDefault();
    const nextFilters = { ...filters, page: 1 };
    setFilters(nextFilters);
    loadProducts(nextFilters);
  }

  function handleReset() {
    setFilters(DEFAULT_FILTERS);
    loadProducts(DEFAULT_FILTERS);
  }

  function handleAddToCart(product) {
    addToCart(product, 1);
    setSuccessMessage(`Đã thêm "${product.name}" vào giỏ hàng.`);
    window.setTimeout(() => setSuccessMessage(""), 1800);
  }

  function goToPage(page) {
    setFilters((current) => ({ ...current, page }));
  }

  const currentPage = meta?.page || filters.page;
  const totalPages = meta?.total_pages || 1;

  return (
    <>
      <div className="page-header">
        <div>
          <h1>Sản phẩm</h1>
          <p className="muted">Tìm kiếm sản phẩm, xem tồn kho và thêm vào giỏ hàng để tạo đơn.</p>
        </div>
      </div>

      <form className="filters" onSubmit={handleSubmit}>
        <Input name="keyword" placeholder="Tìm theo tên sản phẩm" value={filters.keyword} onChange={handleChange} />
        <Select name="category_id" value={filters.category_id} onChange={handleChange}>
          <option value="">Tất cả danh mục</option>
          {categories.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </Select>
        <Input name="min_price" type="number" min="0" placeholder="Giá từ" value={filters.min_price} onChange={handleChange} />
        <Input name="max_price" type="number" min="0" placeholder="Giá đến" value={filters.max_price} onChange={handleChange} />
        <div className="filter-actions">
          <Button type="submit">Tìm kiếm</Button>
          <Button type="button" variant="secondary" onClick={handleReset}>
            Xóa lọc
          </Button>
        </div>
      </form>

      {successMessage && <div className="alert alert-success">{successMessage}</div>}
      <ErrorMessage message={error} />

      {loading ? (
        <Loading />
      ) : products.length === 0 ? (
        <EmptyState title="Không có sản phẩm" description="Thử thay đổi từ khóa, danh mục hoặc khoảng giá." />
      ) : (
        <div className="grid product-grid">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} onAdd={handleAddToCart} />
          ))}
        </div>
      )}

      <div className="pagination">
        <Button type="button" variant="secondary" disabled={currentPage <= 1 || loading} onClick={() => goToPage(currentPage - 1)}>
          Previous
        </Button>
        <span>
          Page {currentPage} / {totalPages}
        </span>
        <Button type="button" variant="secondary" disabled={currentPage >= totalPages || loading} onClick={() => goToPage(currentPage + 1)}>
          Next
        </Button>
      </div>
    </>
  );
}
