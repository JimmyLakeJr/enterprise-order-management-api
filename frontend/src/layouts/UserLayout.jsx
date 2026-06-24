import { Outlet } from "react-router-dom";
import AppHeader from "../components/common/AppHeader";

export default function UserLayout() {
  return (
    <>
      <AppHeader />
      <main className="container page">
        <Outlet />
      </main>
    </>
  );
}
