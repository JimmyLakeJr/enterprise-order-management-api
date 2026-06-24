import { useCallback, useEffect, useState } from "react";
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
import { useConfirm } from "../../hooks/useConfirm";
import { formatCurrency } from "../../utils/formatCurrency";
import { resolveAssetUrl } from "../../utils/resolveAssetUrl";

const IMAGE_TYPES = new Set(["image/jpeg", "image/png", "image/webp", "image/gif", "image/avif"]);
const MAX_IMAGE_SIZE = 5 * 1024 * 1024;

const INITIAL_FORM = {
  category_id: "",
  name: "",
  description: "",
  price: 0,
  stock: 0,
  image_url: "",
  is_active: true,
};

const INITIAL_FILTERS = {
  page: 1,
  limit: 10,
  keyword: "",
  category_id: "",
  min_price: "",
  max_price: "",
  status: "all",
};

export default function AdminProductsPage() {
  const { confirm } = useConfirm();
  const [products, setProducts] = useState([]);
  const [categories, setCategories] = useState([]);
  const [meta, setMeta] = useState(null);
  const [draftFilters, setDraftFilters] = useState(INITIAL_FILTERS);
  const [appliedFilters, setAppliedFilters] = useState(INITIAL_FILTERS);
  const [form, setForm] = useState(INITIAL_FORM);
  const [editingId, setEditingId] = useState(null);
  const [selectedImageFile, setSelectedImageFile] = useState(null);
  const [localPreviewUrl, setLocalPreviewUrl] = useState("");
  const [failedPreviewUrl, setFailedPreviewUrl] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [uploadingImage, setUploadingImage] = useState(false);
  const [deletingId, setDeletingId] = useState(null);
  const [restoringId, setRestoringId] = useState(null);
  const [error, setError] = useState("");
  const [categoryError, setCategoryError] = useState("");
  const [filterError, setFilterError] = useState("");
  const [formError, setFormError] = useState("");

  useEffect(() => {
    return () => {
      if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl);
    };
  }, [localPreviewUrl]);

  const loadProducts = useCallback(async (filters) => {
    setLoading(true);
    setError("");
    try {
      const result = await productApi.listAdmin(filters);
      setProducts(result.data);
      setMeta(result.meta);
    } catch (err) {
      setProducts([]);
      setMeta(null);
      setError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    let active = true;
    categoryApi
      .list()
      .then((data) => {
        if (active) setCategories(data);
      })
      .catch((err) => {
        if (active) setCategoryError(getMessage(err));
      });
    return () => {
      active = false;
    };
  }, []);

  useEffect(() => {
    const timer = window.setTimeout(() => void loadProducts(appliedFilters), 0);
    return () => window.clearTimeout(timer);
  }, [appliedFilters, loadProducts]);

  function handleFilterChange(event) {
    const { name, value } = event.target;
    setDraftFilters((current) => ({ ...current, [name]: value }));
    setFilterError("");
  }

  function handleFilterSubmit(event) {
    event.preventDefault();
    const minPrice = draftFilters.min_price === "" ? null : Number(draftFilters.min_price);
    const maxPrice = draftFilters.max_price === "" ? null : Number(draftFilters.max_price);

    if (minPrice !== null && maxPrice !== null && minPrice > maxPrice) {
      setFilterError("Giá tối thiểu không được lớn hơn giá tối đa.");
      return;
    }

    setAppliedFilters({ ...draftFilters, keyword: draftFilters.keyword.trim(), page: 1 });
  }

  function resetFilters() {
    setDraftFilters(INITIAL_FILTERS);
    setAppliedFilters(INITIAL_FILTERS);
    setFilterError("");
  }

  function handleFormChange(event) {
    const { name, value, checked, type } = event.target;
    setForm((current) => ({ ...current, [name]: type === "checkbox" ? checked : value }));
    setFormError("");
    if (name === "image_url") {
      setFailedPreviewUrl("");
    }
  }

  function clearSelectedImage() {
    if (localPreviewUrl) URL.revokeObjectURL(localPreviewUrl);
    setLocalPreviewUrl("");
    setSelectedImageFile(null);
  }

  function resetForm() {
    setEditingId(null);
    setForm(INITIAL_FORM);
    setFormError("");
    setFailedPreviewUrl("");
    clearSelectedImage();
  }

  function isValidImageURL(value) {
    const url = value.trim();
    if (!url) return true;
    if (url.startsWith("/uploads/")) return true;
    try {
      const parsed = new URL(url);
      return Boolean(parsed.protocol && parsed.host);
    } catch {
      return false;
    }
  }

  function validateForm() {
    if (!form.name.trim()) return "Tên sản phẩm là bắt buộc.";
    if (!Number(form.category_id)) return "Danh mục là bắt buộc.";

    const price = Number(form.price);
    const stock = Number(form.stock);
    if (!Number.isFinite(price) || price < 0) return "Giá sản phẩm phải là số lớn hơn hoặc bằng 0.";
    if (!Number.isInteger(stock) || stock < 0) return "Tồn kho phải là số nguyên lớn hơn hoặc bằng 0.";
    if (!isValidImageURL(form.image_url)) return "Đường dẫn ảnh không hợp lệ.";
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
      await loadProducts(appliedFilters);
    } catch (err) {
      setFormError(getMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  function handleImageFileChange(event) {
    const file = event.target.files?.[0];
    setFormError("");

    if (!file) {
      clearSelectedImage();
      return;
    }
    if (!IMAGE_TYPES.has(file.type)) {
      setFormError("Chỉ chấp nhận JPG/JPEG, PNG, WebP, GIF hoặc AVIF.");
      event.target.value = "";
      return;
    }
    if (file.size > MAX_IMAGE_SIZE) {
      setFormError("Ảnh sản phẩm không được vượt quá 5 MB.");
      event.target.value = "";
      return;
    }

    clearSelectedImage();
    setSelectedImageFile(file);
    setLocalPreviewUrl(URL.createObjectURL(file));
  }

  async function handleUploadImage() {
    if (!selectedImageFile) return;
    setUploadingImage(true);
    setFormError("");
    try {
      const result = await productApi.uploadImage(selectedImageFile);
      setForm((current) => ({ ...current, image_url: result.url }));
      setFailedPreviewUrl("");
      clearSelectedImage();
    } catch (err) {
      setFormError(getMessage(err));
    } finally {
      setUploadingImage(false);
    }
  }

  function startEdit(product) {
    clearSelectedImage();
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
    setFailedPreviewUrl("");
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  async function handleDelete(product) {
    const accepted = await confirm({
      title: "Ẩn sản phẩm?",
      message: `Sản phẩm “${product.name}” sẽ được ẩn khỏi danh sách đang hoạt động.`,
      confirmLabel: "Ẩn sản phẩm",
      danger: true,
    });
    if (!accepted) return;

    setDeletingId(product.id);
    setError("");
    try {
      await productApi.remove(product.id);
      if (editingId === product.id) resetForm();
      await loadProducts(appliedFilters);
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setDeletingId(null);
    }
  }

  async function handleRestore(product) {
    setRestoringId(product.id);
    setError("");
    try {
      await productApi.restore(product.id);
      await loadProducts(appliedFilters);
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setRestoringId(null);
    }
  }

  const currentPage = Number(meta?.page || appliedFilters.page);
  const totalPages = Number(meta?.total_pages || 0);
  const totalProducts = Number(meta?.total || products.length);
  const remotePreviewUrl = resolveAssetUrl(form.image_url);
  const previewUrl = localPreviewUrl || remotePreviewUrl;

  return (
    <div className="grid admin-page-grid">
      <Card className="admin-form-card">
        <div className="page-header compact-header">
          <div>
            <span className="eyebrow">Chỉnh sửa sản phẩm</span>
            <h1>{editingId ? `Sửa sản phẩm #${editingId}` : "Tạo sản phẩm"}</h1>
            <p className="muted">Cập nhật thông tin, giá bán, tồn kho và hình ảnh sản phẩm.</p>
          </div>
          {editingId && (
            <Button type="button" variant="secondary" onClick={resetForm}>
              Tạo mới
            </Button>
          )}
        </div>

        <ErrorMessage message={categoryError ? `Không tải được danh mục: ${categoryError}` : ""} />
        <ErrorMessage message={formError} />
        <form className="form-stack" onSubmit={handleSubmit}>
          <div className="grid admin-product-form-grid">
            <Input label="Tên sản phẩm" name="name" required maxLength="150" value={form.name} onChange={handleFormChange} />
            <Select label="Danh mục" name="category_id" required value={form.category_id} onChange={handleFormChange}>
              <option value="">Chọn danh mục đang hoạt động</option>
              {categories.map((category) => (
                <option key={category.id} value={category.id}>
                  {category.name}
                </option>
              ))}
            </Select>
            <Input label="Giá" name="price" type="number" min="0" step="1" value={form.price} onChange={handleFormChange} />
            <Input label="Tồn kho" name="stock" type="number" min="0" step="1" value={form.stock} onChange={handleFormChange} />
            <Input
              label="Đường dẫn ảnh"
              name="image_url"
              type="text"
              placeholder="https://... hoặc /uploads/products/images/..."
              value={form.image_url}
              onChange={handleFormChange}
            />
            <label className="checkbox-field product-active-field">
              <input type="checkbox" name="is_active" checked={form.is_active} onChange={handleFormChange} />
              <span>Đang hoạt động</span>
            </label>
          </div>

          <label className="field">
            <span>Hoặc tải ảnh sản phẩm lên server</span>
            <input
              type="file"
              accept=".jpg,.jpeg,.png,.webp,.gif,.avif,image/jpeg,image/png,image/webp,image/gif,image/avif"
              onChange={handleImageFileChange}
            />
          </label>

          {selectedImageFile && (
            <div className="actions">
              <Button type="button" onClick={handleUploadImage} disabled={uploadingImage}>
                {uploadingImage ? "Đang tải ảnh..." : "Tải ảnh lên server"}
              </Button>
              <Button type="button" variant="secondary" onClick={clearSelectedImage} disabled={uploadingImage}>
                Hủy ảnh đã chọn
              </Button>
            </div>
          )}

          {previewUrl && (
            <div className="image-preview admin-image-preview">
              {failedPreviewUrl === previewUrl ? (
                <span>Không tải được ảnh xem trước</span>
              ) : (
                <img src={previewUrl} alt="Xem trước sản phẩm" onError={() => setFailedPreviewUrl(previewUrl)} />
              )}
            </div>
          )}

          <p className="preview-only-banner product-media-note">
            Bạn có thể dán URL ảnh công khai hoặc tải ảnh trực tiếp lên server. Sau khi upload thành công, URL nội bộ sẽ được lưu vào
            `image_url`.
          </p>

          <Textarea label="Mô tả" name="description" maxLength="1000" value={form.description} onChange={handleFormChange} />
          {!form.is_active && <p className="inline-note">Sản phẩm ngừng hoạt động sẽ được ẩn khỏi danh sách.</p>}

          <div className="actions">
            <Button type="submit" disabled={submitting || uploadingImage}>
              {submitting ? "Đang lưu..." : editingId ? "Lưu thay đổi" : "Tạo sản phẩm"}
            </Button>
            {editingId && (
              <Button type="button" variant="secondary" onClick={resetForm}>
                Hủy sửa
              </Button>
            )}
          </div>
        </form>
      </Card>

      <Card className="admin-list-card">
        <div className="page-header admin-list-header">
          <div>
            <span className="eyebrow">Quản lý sản phẩm</span>
            <h2>Danh sách sản phẩm</h2>
            <p className="muted">{loading ? "Đang tải..." : `${totalProducts} sản phẩm phù hợp bộ lọc.`}</p>
          </div>
        </div>

        <form className="admin-filter admin-product-filter" onSubmit={handleFilterSubmit}>
          <Input label="Từ khóa" name="keyword" placeholder="Tên sản phẩm..." value={draftFilters.keyword} onChange={handleFilterChange} />
          <Select label="Danh mục" name="category_id" value={draftFilters.category_id} onChange={handleFilterChange}>
            <option value="">Tất cả danh mục</option>
            {categories.map((category) => (
              <option key={category.id} value={category.id}>
                {category.name}
              </option>
            ))}
          </Select>
          <Select label="Trạng thái" name="status" value={draftFilters.status} onChange={handleFilterChange}>
            <option value="all">Tất cả</option>
            <option value="active">Đang hoạt động</option>
            <option value="inactive">Đã ẩn</option>
          </Select>
          <Input label="Giá từ" name="min_price" type="number" min="0" value={draftFilters.min_price} onChange={handleFilterChange} />
          <Input label="Giá đến" name="max_price" type="number" min="0" value={draftFilters.max_price} onChange={handleFilterChange} />
          <div className="filter-actions admin-filter-actions">
            <Button type="submit" disabled={loading}>
              Lọc
            </Button>
            <Button type="button" variant="secondary" disabled={loading} onClick={resetFilters}>
              Xóa lọc
            </Button>
          </div>
        </form>

        <ErrorMessage message={filterError} />
        <ErrorMessage message={error} />

        {loading ? (
          <Loading label="Đang tải sản phẩm..." variant="table" count={6} />
        ) : products.length === 0 ? (
          <EmptyState title="Không có sản phẩm phù hợp" description="Thử đổi bộ lọc hoặc tạo sản phẩm mới." />
        ) : (
          <Table
            rows={products}
            columns={[
              { key: "id", title: "ID" },
              { key: "name", title: "Tên" },
              { key: "category", title: "Danh mục", render: (product) => product.category?.name || `#${product.category_id}` },
              { key: "price", title: "Giá", render: (product) => formatCurrency(product.price) },
              { key: "stock", title: "Tồn kho" },
              {
                key: "status",
                title: "Trạng thái",
                render: (product) => {
                  if (!product.category?.is_active) return <Badge tone="warning">Danh mục đã ẩn</Badge>;
                  return <Badge tone={product.is_active ? "success" : "default"}>{product.is_active ? "Đang hoạt động" : "Đã ẩn"}</Badge>;
                },
              },
              {
                key: "actions",
                title: "Thao tác",
                render: (product) => (
                  <div className="actions table-actions">
                    {product.is_active ? (
                      <>
                        <Button type="button" variant="secondary" disabled={!product.category?.is_active} onClick={() => startEdit(product)}>
                          Sửa
                        </Button>
                        <Button type="button" variant="danger" disabled={deletingId === product.id} onClick={() => handleDelete(product)}>
                          {deletingId === product.id ? "Đang ẩn..." : "Ẩn"}
                        </Button>
                      </>
                    ) : (
                      <Button type="button" disabled={restoringId === product.id} onClick={() => handleRestore(product)}>
                        {restoringId === product.id ? "Đang khôi phục..." : "Khôi phục"}
                      </Button>
                    )}
                  </div>
                ),
              },
            ]}
          />
        )}

        {totalPages > 1 && (
          <div className="pagination">
            <Button
              type="button"
              variant="secondary"
              disabled={currentPage <= 1 || loading}
              onClick={() => setAppliedFilters((current) => ({ ...current, page: currentPage - 1 }))}
            >
              Trang trước
            </Button>
            <span>
              Trang {currentPage} / {totalPages}
            </span>
            <Button
              type="button"
              variant="secondary"
              disabled={currentPage >= totalPages || loading}
              onClick={() => setAppliedFilters((current) => ({ ...current, page: currentPage + 1 }))}
            >
              Trang sau
            </Button>
          </div>
        )}
      </Card>
    </div>
  );
}
