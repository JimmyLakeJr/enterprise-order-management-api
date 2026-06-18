import { useEffect, useMemo, useState } from "react";
import { categoryApi } from "../../api/categoryApi";
import { getMessage } from "../../api/apiClient";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Loading from "../../components/common/Loading";
import Table from "../../components/common/Table";
import Textarea from "../../components/common/Textarea";

const initialForm = {
  name: "",
  description: "",
  is_active: true,
};

export default function AdminCategoriesPage() {
  const [categories, setCategories] = useState([]);
  const [form, setForm] = useState(initialForm);
  const [editingId, setEditingId] = useState(null);
  const [search, setSearch] = useState("");
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [formError, setFormError] = useState("");

  useEffect(() => {
    loadCategories();
  }, []);

  async function loadCategories() {
    setLoading(true);
    setError("");
    try {
      setCategories(await categoryApi.list());
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }

  const filteredCategories = useMemo(() => {
    const keyword = search.trim().toLowerCase();
    if (!keyword) return categories;
    return categories.filter((category) => {
      return (
        category.name?.toLowerCase().includes(keyword) ||
        category.description?.toLowerCase().includes(keyword)
      );
    });
  }, [categories, search]);

  function handleChange(event) {
    const { name, value, checked, type } = event.target;
    setForm((current) => ({ ...current, [name]: type === "checkbox" ? checked : value }));
  }

  function startCreate() {
    setEditingId(null);
    setForm(initialForm);
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
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setFormError("");

    if (!form.name.trim()) {
      setFormError("Tên danh mục là bắt buộc.");
      return;
    }

    setSubmitting(true);
    try {
      const payload = {
        name: form.name.trim(),
        description: form.description.trim(),
        is_active: form.is_active,
      };

      if (editingId) await categoryApi.update(editingId, payload);
      else await categoryApi.create(payload);

      startCreate();
      await loadCategories();
    } catch (err) {
      setFormError(getMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(category) {
    if (!window.confirm(`Xóa danh mục "${category.name}"?`)) return;

    setError("");
    try {
      await categoryApi.remove(category.id);
      await loadCategories();
    } catch (err) {
      setError(getMessage(err));
    }
  }

  return (
    <div className="admin-two-column">
      <Card>
        <div className="page-header">
          <div>
            <h1>Categories</h1>
            <p className="muted">Tạo, cập nhật và soft delete danh mục sản phẩm.</p>
          </div>
          <Button type="button" variant="secondary" onClick={startCreate}>
            Add Category
          </Button>
        </div>

        <ErrorMessage message={formError} />
        <form className="form-stack" onSubmit={handleSubmit}>
          <Input label="Name" name="name" value={form.name} onChange={handleChange} />
          <Textarea label="Description" name="description" value={form.description} onChange={handleChange} />
          <label className="checkbox-field">
            <input type="checkbox" name="is_active" checked={form.is_active} onChange={handleChange} />
            <span>Active</span>
          </label>
          <div className="actions">
            <Button type="submit" disabled={submitting}>
              {submitting ? "Saving..." : editingId ? "Save Changes" : "Create Category"}
            </Button>
            {editingId && (
              <Button type="button" variant="secondary" onClick={startCreate}>
                Cancel
              </Button>
            )}
          </div>
        </form>
      </Card>

      <Card>
        <div className="page-header">
          <div>
            <h2>Danh sách danh mục</h2>
            <p className="muted">Search đang lọc client-side trên dữ liệu đã tải.</p>
          </div>
          <Input placeholder="Search category" value={search} onChange={(event) => setSearch(event.target.value)} />
        </div>

        <ErrorMessage message={error} />
        {loading ? (
          <Loading />
        ) : filteredCategories.length === 0 ? (
          <EmptyState title="Không có danh mục" description="Thử thay đổi keyword hoặc tạo danh mục mới." />
        ) : (
          <Table
            rows={filteredCategories}
            columns={[
              { key: "id", title: "ID" },
              { key: "name", title: "Name" },
              { key: "description", title: "Description" },
              {
                key: "is_active",
                title: "Status",
                render: (category) => (
                  <Badge tone={category.is_active ? "success" : "danger"}>
                    {category.is_active ? "active" : "inactive"}
                  </Badge>
                ),
              },
              {
                key: "actions",
                title: "Action",
                render: (category) => (
                  <div className="actions">
                    <Button type="button" onClick={() => startEdit(category)}>
                      Edit
                    </Button>
                    <Button type="button" variant="danger" onClick={() => handleDelete(category)}>
                      Delete
                    </Button>
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
