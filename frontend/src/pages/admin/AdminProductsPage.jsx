import { useEffect, useState } from "react";
import { getMessage } from "../../api/apiClient";
import { categoryApi } from "../../api/categoryApi";
import { productApi } from "../../api/productApi";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Loading from "../../components/common/Loading";
import Select from "../../components/common/Select";
import Table from "../../components/common/Table";
import Textarea from "../../components/common/Textarea";
import { formatCurrency } from "../../utils/format";

const initialForm = {
  category_id: "",
  name: "",
  description: "",
  price: 0,
  stock: 0,
  image_url: "",
  is_active: true,
};

const initialFilters = {
  page: 1,
  limit: 10,
  keyword: "",
  category_id: "",
};

export default function AdminProductsPage() {
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [meta, setMeta] = useState(null);
  const [filters, setFilters] = useState(initialFilters);
  const [form, setForm] = useState(initialForm);
  const [editingId, setEditingId] = useState(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [formError, setFormError] = useState("");

  useEffect(() => {
    loadCategories();
  }, []);

  useEffect(() => {
    loadProducts(filters);
  }, [filters.page]);

  async function loadCategories() {
    try {
      setCategories(await categoryApi.list());
    } catch {
      setCategories([]);
    }
  }

  async function loadProducts(nextFilters) {
    setLoading(true);
    setError("");
    try {
      const result = await productApi.list({
        page: nextFilters.page,
        limit: nextFilters.limit,
        keyword: nextFilters.keyword,
        category_id: nextFilters.category_id,
      });
      setProducts(result.data);
      setMeta(result.meta);
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }

  function handleFilterChange(event) {
    const { name, value } = event.target;
    setFilters((current) => ({ ...current, [name]: value }));
  }

  function handleFilterSubmit(event) {
    event.preventDefault();
    const nextFilters = { ...filters, page: 1 };
    setFilters(nextFilters);
    loadProducts(nextFilters);
  }

  function handleFormChange(event) {
    const { name, value, checked, type } = event.target;
    setForm((current) => ({ ...current, [name]: type === "checkbox" ? checked : value }));
  }

  function resetForm() {
    setEditingId(null);
    setForm(initialForm);
    setFormError("");
  }

  function validateForm() {
    if (!form.name.trim()) return "Tên sản phẩm là bắt buộc.";
    if (!Number(form.category_id)) return "Category là bắt buộc.";
    if (Number(form.price) < 0) return "Giá sản phẩm phải >= 0.";
    if (Number(form.stock) < 0) return "Tồn kho phải >= 0.";
    return "";
  }

  async function handleSubmit(event) {
    event.preventDefault();
    const validationMessage = validateForm();
    if (validationMessage) {
      setFormError(validationMessage);
      return;
    }

    setSubmitting(true);
    setFormError("");
    try {
      const payload = {
        category_id: Number(form.category_id),
        name: form.name.trim(),
        description: form.description.trim(),
        price: Number(form.price),
        stock: Number(form.stock),
        image_url: form.image_url.trim(),
        is_active: form.is_active,
      };

      if (editingId) await productApi.update(editingId, payload);
      else await productApi.create(payload);

      resetForm();
      await loadProducts(filters);
    } catch (err) {
      setFormError(getMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  function startEdit(product) {
    setEditingId(product.id);
    setForm({
      category_id: product.category_id || "",
      name: product.name || "",
      description: product.description || "",
      price: product.price ?? 0,
      stock: product.stock ?? 0,
      image_url: product.image_url || "",
      is_active: product.is_active ?? true,
    });
    setFormError("");
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  async function handleDelete(product) {
    if (!window.confirm(`Xóa sản phẩm "${product.name}"?`)) return;

    setError("");
    try {
      await productApi.remove(product.id);
      await loadProducts(filters);
    } catch (err) {
      setError(getMessage(err));
    }
  }

  function goToPage(page) {
    const nextFilters = { ...filters, page };
    setFilters(nextFilters);
  }

  const currentPage = meta?.page || filters.page;
  const totalPages = meta?.total_pages || 1;

  return (
    <div className="grid">
      <Card>
        <div className="page-header">
          <div>
            <h1>Products</h1>
            <p className="muted">Quản lý sản phẩm, tồn kho, giá bán và trạng thái active.</p>
          </div>
          <Button type="button" variant="secondary" onClick={resetForm}>
            Add Product
          </Button>
        </div>

        <ErrorMessage message={formError} />
        <form className="form-stack" onSubmit={handleSubmit}>
          <div className="grid grid-2">
            <Input label="Name" name="name" value={form.name} onChange={handleFormChange} />
            <Select label="Category" name="category_id" value={form.category_id} onChange={handleFormChange}>
              <option value="">Chọn category</option>
              {categories.map((category) => (
                <option key={category.id} value={category.id}>
                  {category.name}
                </option>
              ))}
            </Select>
            <Input label="Price" name="price" type="number" min="0" value={form.price} onChange={handleFormChange} />
            <Input label="Stock" name="stock" type="number" min="0" value={form.stock} onChange={handleFormChange} />
            <Input label="Image URL" name="image_url" value={form.image_url} onChange={handleFormChange} />
            <label className="checkbox-field product-active-field">
              <input type="checkbox" name="is_active" checked={form.is_active} onChange={handleFormChange} />
              <span>Active</span>
            </label>
          </div>

          {form.image_url && (
            <div className="image-preview">
              <img src={form.image_url} alt="Product preview" />
            </div>
          )}

          <Textarea label="Description" name="description" value={form.description} onChange={handleFormChange} />

          <div className="actions">
            <Button type="submit" disabled={submitting}>
              {submitting ? "Saving..." : editingId ? "Save Changes" : "Create Product"}
            </Button>
            {editingId && (
              <Button type="button" variant="secondary" onClick={resetForm}>
                Cancel
              </Button>
            )}
          </div>
        </form>
      </Card>

      <Card>
        <form className="admin-filter" onSubmit={handleFilterSubmit}>
          <Input name="keyword" placeholder="Keyword" value={filters.keyword} onChange={handleFilterChange} />
          <Select name="category_id" value={filters.category_id} onChange={handleFilterChange}>
            <option value="">Tất cả category</option>
            {categories.map((category) => (
              <option key={category.id} value={category.id}>
                {category.name}
              </option>
            ))}
          </Select>
          <Button type="submit">Filter</Button>
        </form>

        <ErrorMessage message={error} />
        {loading ? (
          <Loading />
        ) : products.length === 0 ? (
          <EmptyState title="Không có sản phẩm" description="Thử đổi filter hoặc tạo sản phẩm mới." />
        ) : (
          <Table
            rows={products}
            columns={[
              { key: "id", title: "ID" },
              { key: "name", title: "Name" },
              { key: "category", title: "Category", render: (product) => product.category?.name || product.category_id },
              { key: "price", title: "Price", render: (product) => formatCurrency(product.price) },
              { key: "stock", title: "Stock" },
              {
                key: "status",
                title: "Status",
                render: (product) => (
                  <Badge tone={product.is_active ? "success" : "danger"}>
                    {product.is_active ? "active" : "inactive"}
                  </Badge>
                ),
              },
              {
                key: "actions",
                title: "Action",
                render: (product) => (
                  <div className="actions">
                    <Button type="button" onClick={() => startEdit(product)}>
                      Edit
                    </Button>
                    <Button type="button" variant="danger" onClick={() => handleDelete(product)}>
                      Delete
                    </Button>
                  </div>
                ),
              },
            ]}
          />
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
      </Card>
    </div>
  );
}
