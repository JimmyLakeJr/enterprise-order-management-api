import { Outlet } from "react-router-dom";
import AppHeader from "../components/common/AppHeader";

export default function PublicLayout() {
  return (
    <>
      <AppHeader />
      <main className="container page">
        <Outlet />
      </main>
      <footer className="footer">
        <div className="container">Enterprise product and order management demo.</div>
      </footer>
    </>
  );
}
