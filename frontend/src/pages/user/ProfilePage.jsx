import Card from "../../components/common/Card";
import Badge from "../../components/common/Badge";
import { useAuth } from "../../contexts/AuthContext";

export default function ProfilePage() {
  const { user } = useAuth();

  return (
    <Card>
      <h1>Tài khoản</h1>
      <p><strong>Họ tên:</strong> {user?.name}</p>
      <p><strong>Email:</strong> {user?.email}</p>
      <p><strong>Vai trò:</strong> <Badge>{user?.role}</Badge></p>
    </Card>
  );
}
