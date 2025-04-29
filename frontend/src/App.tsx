import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import ChangePassword from "@/pages/ChangePassword";
import Dashboard from "@/pages/Dashboard";
import Login from "@/pages/Login";
import ProtectedRoutes from "@/utils/ProtectedRoutes";
// import Header from "@/components/ui/Header";
export default function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Navigate to="/dashboard" replace />} />
        <Route path="/login" element={<Login />} />
        <Route element={<ProtectedRoutes />}>
          <Route path="/change" element={<ChangePassword />} />
          <Route path="/dashboard" element={<Dashboard />} />
        </Route>
      </Routes>
    </Router>
  );
}
