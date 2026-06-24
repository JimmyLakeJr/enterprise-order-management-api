import { useEffect, useRef, useState } from "react";
import { getMessage } from "../../api/apiClient";
import { getCategories } from "../../api/categoryApi";
import { getProducts } from "../../api/productApi";
import Button from "../../components/common/Button";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Loading from "../../components/common/Loading";
import Select from "../../components/common/Select";
import ProductCard from "../../components/products/ProductCard";
import Toast from "../../components/common/Toast";
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
  const requestIdRef = useRef(0);
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [meta, setMeta] = useState(null);
  const [draftFilters, setDraftFilters] = useState(DEFAULT_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(DEFAULT_FILTERS);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [filterError, setFilterError] = useState("");
  const [categoryError, setCategoryError] = useState("");
  const [successMessage, setSuccessMessage] = useState("");
  const [retryToken, setRetryToken] = useState(0);

  useEffect(() => {
    let active = true;

    getCategories()
      .then((data) => {
        if (active) setCategories(data);
      })
      .catch((err) => {
        if (active) setCategoryError(`Không tải được danh mục: ${getMessage(err)}`);
      });

    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    const requestId = ++requestIdRef.current;

    async function loadProducts() {
      setLoading(true);
      setError("");

      try {
        const result = await getProducts(appliedFilters);
        if (requestId !== requestIdRef.current) return;
        setProducts(result.data);
        setMeta(result.meta);
      } catch (err) {
        if (requestId !== requestIdRef.current) return;
        setProducts([]);
        setMeta(null);
        setError(`Không tải được danh sách sản phẩm: ${getMessage(err)}`);
      } finally {
        if (requestId === requestIdRef.current) setLoading(false);
      }
    }

    loadProducts();
  }, [appliedFilters, retryToken]);

  function handleChange(event) {
    const { name, value } = event.target;
    setDraftFilters((current) => ({ ...current, [name]: value }));
    setFilterError("");
  }

  function handleSubmit(event) {
    event.preventDefault();
    const minPrice = draftFilters.min_price === "" ? null : Number(draftFilters.min_price);
    const maxPrice = draftFilters.max_price === "" ? null : Number(draftFilters.max_price);

    if (minPrice !== null && maxPrice !== null && minPrice > maxPrice) {
      setFilterError("Giá tối thiểu không được lớn hơn giá tối đa.");
      return;
    }

    setFilterError("");
    setAppliedFilters({ ...draftFilters, keyword: draftFilters.keyword.trim(), page: 1 });
  }

  function handleReset() {
    setDraftFilters(DEFAULT_FILTERS);
    setAppliedFilters(DEFAULT_FILTERS);
    setFilterError("");
  }

  function handleAddToCart(product) {
    addToCart(product, 1);
    setSuccessMessage(`Đã thêm “${product.name}” vào giỏ hàng.`);
  }

  function goToPage(page) {
    setAppliedFilters((current) => ({ ...current, page }));
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  const currentPage = Number(meta?.page || appliedFilters.page);
  const totalPages = Number(meta?.total_pages || 0);
  const totalProducts = Number(meta?.total || products.length);

  return (
    <>
      <header className="page-header product-list-header">
        <div>
          <span className="eyebrow">Danh mục sản phẩm</span>
          <h1>Chọn sản phẩm phù hợp với bạn</h1>
          <p className="muted">Tìm kiếm theo tên, danh mục và khoảng giá từ dữ liệu sản phẩm hiện có.</p>
        </div>
        {!loading && !error && <span className="result-count">{totalProducts} sản phẩm</span>}
      </header>

      <form className="filters" onSubmit={handleSubmit} aria-label="Bộ lọc sản phẩm">
        <Input
          label="Từ khóa"
          name="keyword"
          placeholder="Nhập tên sản phẩm..."
          value={draftFilters.keyword}
          onChange={handleChange}
        />
        <Select label="Danh mục" name="category_id" value={draftFilters.category_id} onChange={handleChange}>
          <option value="">Tất cả danh mục</option>
          {categories.map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </Select>
        <Input
          label="Giá từ"
          name="min_price"
          type="number"
          min="0"
          step="1"
          placeholder="0"
          value={draftFilters.min_price}
          onChange={handleChange}
        />
        <Input
          label="Giá đến"
          name="max_price"
          type="number"
          min="0"
          step="1"
          placeholder="Không giới hạn"
          value={draftFilters.max_price}
          onChange={handleChange}
        />
        <div className="filter-actions">
          <Button type="submit" disabled={loading}>Áp dụng</Button>
          <Button type="button" variant="secondary" onClick={handleReset} disabled={loading}>
            Xóa lọc
          </Button>
        </div>
      </form>

      <ErrorMessage message={filterError} />
      <ErrorMessage message={categoryError} />
      <ErrorMessage message={error} onRetry={() => setRetryToken((current) => current + 1)} />
      <Toast message={successMessage} tone="success" onDismiss={() => setSuccessMessage("")} />

      {loading ? (
        <Loading label="Đang tải sản phẩm..." variant="cards" count={8} />
      ) : error ? null : products.length === 0 ? (
        <EmptyState title="Không tìm thấy sản phẩm" description="Hãy thử từ khóa, danh mục hoặc khoảng giá khác." />
      ) : (
        <div className="grid product-grid">
          {products.map((product) => (
            <ProductCard key={product.id} product={product} onAdd={handleAddToCart} />
          ))}
        </div>
      )}

      {!error && totalPages > 1 && (
        <nav className="pagination" aria-label="Phân trang sản phẩm">
          <Button type="button" variant="secondary" disabled={currentPage <= 1 || loading} onClick={() => goToPage(currentPage - 1)}>
            Trang trước
          </Button>
          <span>Trang {currentPage} / {totalPages}</span>
          <Button type="button" variant="secondary" disabled={currentPage >= totalPages || loading} onClick={() => goToPage(currentPage + 1)}>
            Trang sau
          </Button>
        </nav>
      )}
    </>
  );
}
