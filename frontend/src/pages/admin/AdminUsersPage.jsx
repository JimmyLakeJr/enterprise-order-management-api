import { useCallback, useEffect, useState } from "react";
import { getMessage } from "../../api/apiClient";
import { userApi } from "../../api/userApi";
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
import { useConfirm } from "../../hooks/useConfirm";

const INITIAL_FORM = { name: "", email: "", phone: "", role: ROLES.USER };
const INITIAL_QUERY = { page: 1, limit: 10, search: "" };

export default function AdminUsersPage() {
  const { user: currentUser } = useAuth();
  const { confirm } = useConfirm();
  const [users, setUsers] = useState([]);
  const [meta, setMeta] = useState(null);
  const [searchDraft, setSearchDraft] = useState("");
  const [query, setQuery] = useState(INITIAL_QUERY);
  const [form, setForm] = useState(INITIAL_FORM);
  const [editingId, setEditingId] = useState(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [deletingId, setDeletingId] = useState(null);
  const [error, setError] = useState("");
  const [formError, setFormError] = useState("");

  const loadUsers = useCallback(async (params) => {
    setLoading(true);
    setError("");
    try {
      const result = await userApi.list(params);
      setUsers(result.data);
      setMeta(result.meta);
    } catch (err) {
      setUsers([]);
      setMeta(null);
      setError(getMessage(err));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    const timer = window.setTimeout(() => void loadUsers(query), 0);
    return () => window.clearTimeout(timer);
  }, [loadUsers, query]);

  function handleSearch(event) {
    event.preventDefault();
    setQuery((current) => ({ ...current, page: 1, search: searchDraft.trim() }));
  }

  function clearSearch() {
    setSearchDraft("");
    setQuery(INITIAL_QUERY);
  }

  function handleFormChange(event) {
    const { name, value } = event.target;
    setForm((current) => ({ ...current, [name]: value }));
    setFormError("");
  }

  function startEdit(user) {
    setEditingId(user.id);
    setForm({
      name: user.name || "",
      email: user.email || "",
      phone: user.phone || "",
      role: user.role || ROLES.USER,
    });
    setFormError("");
    window.scrollTo({ top: 0, behavior: "smooth" });
  }

  function cancelEdit() {
    setEditingId(null);
    setForm(INITIAL_FORM);
    setFormError("");
  }

  function validateForm() {
    if (!form.name.trim()) return "Tên người dùng là bắt buộc.";
    if (form.name.trim().length < 2) return "Tên người dùng phải có ít nhất 2 ký tự.";
    if (!form.email.trim() && !form.phone.trim()) return "Cần email hoặc số điện thoại.";
    if (form.email.trim() && !/^\S+@\S+\.\S+$/.test(form.email.trim())) return "Email không hợp lệ.";
    if (form.phone.trim() && form.phone.replace(/[^\d+]/g, "").length < 9) return "Số điện thoại không hợp lệ.";
    if (![ROLES.USER, ROLES.ADMIN].includes(form.role)) return "Role không hợp lệ.";
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
      await userApi.update(editingId, {
        name: form.name.trim(),
        email: form.email.trim(),
        phone: form.phone.trim(),
        role: form.role,
      });
      cancelEdit();
      await loadUsers(query);
    } catch (err) {
      setFormError(getMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(user) {
    if (currentUser?.id === user.id) return;
    const accepted = await confirm({
      title: "Vô hiệu hóa tài khoản?",
      message: `Tài khoản “${user.email || user.phone || user.name}” sẽ được chuyển sang trạng thái ngừng hoạt động.`,
      confirmLabel: "Vô hiệu hóa",
      danger: true,
    });
    if (!accepted) return;

    setDeletingId(user.id);
    setError("");
    try {
      await userApi.remove(user.id);
      if (editingId === user.id) cancelEdit();
      await loadUsers(query);
    } catch (err) {
      setError(getMessage(err));
    } finally {
      setDeletingId(null);
    }
  }

  const currentPage = Number(meta?.page || query.page);
  const totalPages = Number(meta?.total_pages || 0);
  const totalUsers = Number(meta?.total || users.length);

  return (
    <div className="grid admin-page-grid">
      {editingId && (
        <Card className="admin-form-card">
          <div className="page-header compact-header">
            <div>
              <span className="eyebrow">Chỉnh sửa người dùng</span>
              <h2>Sửa người dùng #{editingId}</h2>
            </div>
            <Button type="button" variant="secondary" onClick={cancelEdit}>Đóng form</Button>
          </div>
          <ErrorMessage message={formError} />
          <form className="form-stack" onSubmit={handleSubmit}>
            <div className="grid admin-user-form-grid">
              <Input label="Họ tên" name="name" required maxLength="100" value={form.name} onChange={handleFormChange} />
              <Input label="Email" name="email" type="email" maxLength="255" value={form.email} onChange={handleFormChange} />
              <Input label="Số điện thoại" name="phone" type="tel" maxLength="20" value={form.phone} onChange={handleFormChange} />
              <Select label="Vai trò" name="role" value={form.role} onChange={handleFormChange}>
                <option value={ROLES.USER}>Người dùng</option>
                <option value={ROLES.ADMIN}>Quản trị viên</option>
              </Select>
            </div>
            {currentUser?.id === editingId && <p className="inline-note">Bạn đang sửa chính tài khoản admin đang đăng nhập.</p>}
            <div className="actions">
              <Button type="submit" disabled={submitting}>{submitting ? "Đang lưu..." : "Lưu người dùng"}</Button>
              <Button type="button" variant="secondary" onClick={cancelEdit}>Hủy</Button>
            </div>
          </form>
        </Card>
      )}

      <Card className="admin-list-card">
        <div className="page-header admin-list-header">
          <div>
            <span className="eyebrow">Tài khoản hoạt động</span>
            <h1>Quản lý người dùng</h1>
            <p className="muted">{loading ? "Đang tải..." : `${totalUsers} tài khoản đang hoạt động.`}</p>
          </div>
          <form className="admin-user-search" onSubmit={handleSearch}>
            <Input aria-label="Tìm người dùng" placeholder="Tìm tên, email hoặc số điện thoại..." value={searchDraft} onChange={(event) => setSearchDraft(event.target.value)} />
            <Button type="submit" disabled={loading}>Tìm</Button>
            <Button type="button" variant="secondary" disabled={loading} onClick={clearSearch}>Xóa lọc</Button>
          </form>
        </div>

        <ErrorMessage message={error} />

        {loading ? (
          <Loading label="Đang tải người dùng..." variant="table" count={6} />
        ) : users.length === 0 ? (
          <EmptyState title="Không có người dùng phù hợp" description="Thử tìm bằng tên, email hoặc số điện thoại khác." />
        ) : (
          <Table
            rows={users}
            columns={[
              { key: "id", title: "ID" },
              { key: "name", title: "Họ tên" },
              { key: "email", title: "Email", render: (user) => user.email || "—" },
              { key: "phone", title: "SĐT", render: (user) => user.phone || "—" },
              { key: "role", title: "Vai trò", render: (user) => <Badge tone={user.role === ROLES.ADMIN ? "primary" : "default"}>{user.role === ROLES.ADMIN ? "Quản trị viên" : "Người dùng"}</Badge> },
              { key: "status", title: "Trạng thái", render: () => <Badge tone="success">Đang hoạt động</Badge> },
              { key: "created_at", title: "Ngày tạo", render: (user) => formatDate(user.created_at) || "—" },
              {
                key: "actions",
                title: "Thao tác",
                render: (user) => {
                  const isCurrentUser = currentUser?.id === user.id;
                  return (
                    <div className="actions table-actions">
                      <Button type="button" variant="secondary" onClick={() => startEdit(user)}>Sửa</Button>
                      <Button
                        type="button"
                        variant="danger"
                        disabled={isCurrentUser || deletingId === user.id}
                        title={isCurrentUser ? "Không thể vô hiệu hóa tài khoản đang đăng nhập" : "Vô hiệu hóa tài khoản"}
                        onClick={() => handleDelete(user)}
                      >
                        {isCurrentUser ? "Tài khoản hiện tại" : deletingId === user.id ? "Đang xử lý..." : "Vô hiệu hóa"}
                      </Button>
                    </div>
                  );
                },
              },
            ]}
          />
        )}

        {totalPages > 1 && (
          <div className="pagination">
            <Button type="button" variant="secondary" disabled={currentPage <= 1 || loading} onClick={() => setQuery((current) => ({ ...current, page: currentPage - 1 }))}>Trang trước</Button>
            <span>Trang {currentPage} / {totalPages}</span>
            <Button type="button" variant="secondary" disabled={currentPage >= totalPages || loading} onClick={() => setQuery((current) => ({ ...current, page: currentPage + 1 }))}>Trang sau</Button>
          </div>
        )}
      </Card>
    </div>
  );
}
