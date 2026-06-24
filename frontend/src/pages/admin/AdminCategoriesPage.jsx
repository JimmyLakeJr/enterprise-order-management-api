import { useEffect, useMemo, useState } from "react";
import { getMessage } from "../../api/apiClient";
import { categoryApi } from "../../api/categoryApi";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Loading from "../../components/common/Loading";
import Table from "../../components/common/Table";
import Textarea from "../../components/common/Textarea";
import { useConfirm } from "../../hooks/useConfirm";

const INITIAL_FORM = { name: "", description: "", is_active: true };

export default function AdminCategoriesPage() {
  const { confirm } = useConfirm();
  const [categories, setCategories] = useState([]);
  const [form, setForm] = useState(INITIAL_FORM);
  const [editingId, setEditingId] = useState(null);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [deletingId, setDeletingId] = useState(null);
  const [restoringId, setRestoringId] = useState(null);
  const [error, setError] = useState("");
  const [formError, setFormError] = useState("");

  async function loadCategories() {
    setLoading(true);
    setError("");
    try {
      setCategories(await categoryApi.listAdmin("all"));
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    const timer = window.setTimeout(() => void loadCategories(), 0);
    return () => window.clearTimeout(timer);
  }, []);

  const filteredCategories = useMemo(() => {
    const keyword = search.trim().toLowerCase();
    if (!keyword) return categories;
    return categories.filter((category) =>
      [category.name, category.description].some((value) => value?.toLowerCase().includes(keyword))
    );
  }, [categories, search]);

  function handleChange(event) {
    const { name, value, checked, type } = event.target;
    setForm((current) => ({ ...current, [name]: type === "checkbox" ? checked : value }));
    setFormError("");
  }

  function resetForm() {
    setEditingId(null);
    setForm(INITIAL_FORM);
    setFormError("");
  }

  function startEdit(category) {
    setEditingId(category.id);
    setForm({
      name: category.name || "",
      description: category.description || "",
      is_active: category.is_active ?? true,
    });
    setFormError("");
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  async function handleSubmit(event) {
    event.preventDefault();
    const name = form.name.trim();

    if (!name) {
      setFormError("Tên danh mục là bắt buộc.");
      return;
    }
    if (name.length < 2) {
      setFormError("Tên danh mục phải có ít nhất 2 ký tự.");
      return;
    }

    setSubmitting(true);
    setFormError("");
    try {
      const payload = {
        name,
        description: form.description.trim(),
        is_active: form.is_active,
      };

      if (editingId) await categoryApi.update(editingId, payload);
      else await categoryApi.create(payload);

      resetForm();
      await loadCategories();
    } catch (err) {
      setFormError(getMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(category) {
    const accepted = await confirm({
      title: "Ẩn danh mục?",
      message: `Danh mục “${category.name}” sẽ được ẩn khỏi danh sách đang hoạt động.`,
      confirmLabel: "Ẩn danh mục",
      danger: true,
    });
    if (!accepted) return;

    setDeletingId(category.id);
    setError("");
    try {
      await categoryApi.remove(category.id);
      if (editingId === category.id) resetForm();
      await loadCategories();
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setDeletingId(null);
    }
  }

  async function handleRestore(category) {
    setRestoringId(category.id);
    setError("");
    try {
      await categoryApi.restore(category.id);
      await loadCategories();
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setRestoringId(null);
    }
  }

  return (
    <div className="admin-two-column">
      <Card className="admin-form-card">
        <div className="page-header compact-header">
          <div>
            <span className="eyebrow">Chỉnh sửa danh mục</span>
            <h1>{editingId ? `Sửa danh mục #${editingId}` : "Tạo danh mục"}</h1>
          </div>
          {editingId && <Button type="button" variant="secondary" onClick={resetForm}>Tạo mới</Button>}
        </div>

        <ErrorMessage message={formError} />
        <form className="form-stack" onSubmit={handleSubmit}>
          <Input label="Tên danh mục" name="name" required maxLength="100" value={form.name} onChange={handleChange} />
          <Textarea label="Mô tả" name="description" maxLength="1000" value={form.description} onChange={handleChange} />
          <label className="checkbox-field">
            <input type="checkbox" name="is_active" checked={form.is_active} onChange={handleChange} />
            <span>Đang hoạt động</span>
          </label>
          {!form.is_active && <p className="inline-note">Danh mục ngừng hoạt động sẽ được ẩn khỏi danh sách.</p>}
          <div className="actions">
            <Button type="submit" disabled={submitting}>{submitting ? "Đang lưu..." : editingId ? "Lưu thay đổi" : "Tạo danh mục"}</Button>
            {editingId && <Button type="button" variant="secondary" onClick={resetForm}>Hủy sửa</Button>}
          </div>
        </form>
      </Card>

      <Card className="admin-list-card">
        <div className="page-header admin-list-header">
          <div>
            <span className="eyebrow">Tất cả danh mục</span>
            <h2>Danh sách danh mục</h2>
            <p className="muted">Danh mục đã ẩn có thể được khôi phục tại đây.</p>
            <p className="inline-note">Ẩn danh mục có thể làm sản phẩm public trong danh mục đó không còn hiển thị; hệ thống cũng chặn ẩn khi vẫn còn product active.</p>
          </div>
          <Input aria-label="Tìm danh mục" placeholder="Tìm tên hoặc mô tả..." value={search} onChange={(event) => setSearch(event.target.value)} />
        </div>

        <ErrorMessage message={error} />

        {loading ? (
          <Loading label="Đang tải danh mục..." variant="table" count={5} />
        ) : filteredCategories.length === 0 ? (
          <EmptyState title="Không có danh mục phù hợp" description="Thử từ khóa khác hoặc tạo danh mục mới." />
        ) : (
          <Table
            rows={filteredCategories}
            columns={[
              { key: "id", title: "ID" },
              { key: "name", title: "Tên" },
              { key: "description", title: "Mô tả", render: (category) => category.description || "—" },
              { key: "status", title: "Trạng thái", render: (category) => <Badge tone={category.is_active ? "success" : "default"}>{category.is_active ? "Đang hoạt động" : "Đã ẩn"}</Badge> },
              {
                key: "actions",
                title: "Thao tác",
                render: (category) => (
                  <div className="actions table-actions">
                    {category.is_active ? (
                      <>
                        <Button type="button" variant="secondary" onClick={() => startEdit(category)}>Sửa</Button>
                        <Button type="button" variant="danger" disabled={deletingId === category.id} onClick={() => handleDelete(category)}>
                          {deletingId === category.id ? "Đang ẩn..." : "Ẩn"}
                        </Button>
                      </>
                    ) : (
                      <Button type="button" disabled={restoringId === category.id} onClick={() => handleRestore(category)}>
                        {restoringId === category.id ? "Đang khôi phục..." : "Khôi phục"}
                      </Button>
                    )}
                  </div>
                ),
              },
            ]}
          />
        )}
      </Card>
    </div>
  );
}
