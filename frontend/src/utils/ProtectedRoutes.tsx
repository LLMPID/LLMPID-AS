import { useAtomValue } from "jotai";
import { authTokenAtom } from "@/atoms/authAtom";
import { Navigate, Outlet } from "react-router-dom";
export default function ProtectedRoutes() {
  const auth = useAtomValue(authTokenAtom);

  if (!auth) {
    return <Navigate to="/login" replace />;
  }
  return <Outlet />;
}
