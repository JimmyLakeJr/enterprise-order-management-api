import { Link } from "react-router-dom";
import Button from "../../components/common/Button";
import Card from "../../components/common/Card";

export default function ForbiddenPage() {
  return (
    <Card>
      <h1>Forbidden</h1>
      <p className="muted">Your account does not have permission to access this page.</p>
      <Link to="/">
        <Button>Back to products</Button>
      </Link>
    </Card>
  );
}
