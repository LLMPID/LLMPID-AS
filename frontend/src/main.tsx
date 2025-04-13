import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { Provider } from "@/components/ui/provider";
//import App from "@/App.tsx";
import Login from "@/pages/Login.tsx";
//import Dashboard from "@/pages/Dashboard.tsx"
//import ChangePassword from "./pages/ChangePassword";
createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <Provider forcedTheme="light">
        <Login />
    </Provider>
  </StrictMode>
);
