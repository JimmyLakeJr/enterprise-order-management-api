import { useEffect, useState } from "react";
import { userApi } from "../../api/userApi";
import { getMessage } from "../../api/apiClient";
import Badge from "../../components/common/Badge";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";
import EmptyState from "../../components/common/EmptyState";
import ErrorMessage from "../../components/common/ErrorMessage";
import Input from "../../components/common/Input";
import Loading from "../../components/common/Loading";
import Select from "../../components/common/Select";
import Table from "../../components/common/Table";
import { ROLES } from "../../constants/domain";
import { useAuth } from "../../contexts/AuthContext";
import { formatDate } from "../../utils/format";

const initialForm = {
  name: "",
  email: "",
  role: ROLES.USER,
};

export default function AdminUsersPage() {
  const { user: currentUser } = useAuth();
  const [users, setUsers] = useState([]);
  const [meta, setMeta] = useState(null);
  const [search, setSearch] = useState("");
  const [page, setPage] = useState(1);
  const [form, setForm] = useState(initialForm);
  const [editingId, setEditingId] = useState(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [formError, setFormError] = useState("");

  useEffect(() => {
    loadUsers({ page, search });
  }, [page]);

  async function loadUsers(params = {}) {
    setLoading(true);
    setError("");
    try {
      const result = await userApi.list({
        page: params.page || 1,
        limit: 10,
        search: params.search || "",
      });
      setUsers(result.data);
      setMeta(result.meta);
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }

  function handleSearch(event) {
    event.preventDefault();
    setPage(1);
    loadUsers({ page: 1, search });
  }

  function handleFormChange(event) {
    const { name, value } = event.target;
    setForm((current) => ({ ...current, [name]: value }));
  }

  function startEdit(user) {
    setEditingId(user.id);
    setForm({
      name: user.name || "",
      email: user.email || "",
      role: user.role || ROLES.USER,
    });
    setFormError("");
  }

  function cancelEdit() {
    setEditingId(null);
    setForm(initialForm);
    setFormError("");
  }

  async function handleSubmit(event) {
    event.preventDefault();
    setFormError("");

    if (!form.name.trim()) {
      setFormError("Tên user là bắt buộc.");
      return;
    }
    if (!form.email.trim()) {
      setFormError("Email là bắt buộc.");
      return;
    }

    setSubmitting(true);
    try {
      await userApi.update(editingId, {
        name: form.name.trim(),
        email: form.email.trim(),
        role: form.role,
      });
      cancelEdit();
      await loadUsers({ page, search });
    } catch (err) {
      setFormError(getMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(user) {
    if (!window.confirm(`Deactivate/delete user "${user.email}"?`)) return;

    setError("");
    try {
      await userApi.remove(user.id);
      await loadUsers({ page, search });
    } catch (err) {
      setError(getMessage(err));
    }
  }

  const currentPage = meta?.page || page;
  const totalPages = meta?.total_pages || 1;

  return (
    <div className="grid">
      {editingId && (
        <Card>
          <h2>Edit User</h2>
          <ErrorMessage message={formError} />
          <form className="form-stack" onSubmit={handleSubmit}>
            <div className="grid grid-2">
              <Input label="Name" name="name" value={form.name} onChange={handleFormChange} />
              <Input label="Email" name="email" type="email" value={form.email} onChange={handleFormChange} />
              <Select label="Role" name="role" value={form.role} onChange={handleFormChange}>
                <option value={ROLES.USER}>{ROLES.USER}</option>
                <option value={ROLES.ADMIN}>{ROLES.ADMIN}</option>
              </Select>
            </div>
            <div className="actions">
              <Button type="submit" disabled={submitting}>
                {submitting ? "Saving..." : "Save"}
              </Button>
              <Button type="button" variant="secondary" onClick={cancelEdit}>
                Cancel
              </Button>
            </div>
          </form>
        </Card>
      )}

      <Card>
        <div className="page-header">
          <div>
            <h1>Users</h1>
            <p className="muted">Quản lý tài khoản user/admin nếu backend user API được bật.</p>
          </div>
          <form className="actions" onSubmit={handleSearch}>
            <Input placeholder="Search email/name" value={search} onChange={(event) => setSearch(event.target.value)} />
            <Button type="submit">Search</Button>
          </form>
        </div>

        <ErrorMessage message={error} />

        {loading ? (
          <Loading />
        ) : users.length === 0 ? (
          <EmptyState title="Không có user" description="Thử đổi từ khóa tìm kiếm." />
        ) : (
          <Table
            rows={users}
            columns={[
              { key: "id", title: "ID" },
              { key: "name", title: "Name" },
              { key: "email", title: "Email" },
              { key: "role", title: "Role", render: (user) => <Badge tone={user.role === ROLES.ADMIN ? "primary" : "default"}>{user.role}</Badge> },
              {
                key: "is_active",
                title: "Active",
                render: (user) => (
                  <Badge tone={user.is_active ? "success" : "danger"}>{user.is_active ? "active" : "inactive"}</Badge>
                ),
              },
              { key: "created_at", title: "Created", render: (user) => formatDate(user.created_at) || "N/A" },
              {
                key: "actions",
                title: "Action",
                render: (user) => (
                  <div className="actions">
                    <Button type="button" onClick={() => startEdit(user)}>
                      Edit
                    </Button>
                    <Button
                      type="button"
                      variant="danger"
                      disabled={currentUser?.id === user.id}
                      onClick={() => handleDelete(user)}
                    >
                      Delete
                    </Button>
                  </div>
                ),
              },
            ]}
          />
        )}

        <div className="pagination">
          <Button type="button" variant="secondary" disabled={currentPage <= 1 || loading} onClick={() => setPage(currentPage - 1)}>
            Previous
          </Button>
          <span>
            Page {currentPage} / {totalPages}
          </span>
          <Button type="button" variant="secondary" disabled={currentPage >= totalPages || loading} onClick={() => setPage(currentPage + 1)}>
            Next
          </Button>
        </div>
      </Card>
    </div>
  );
}
