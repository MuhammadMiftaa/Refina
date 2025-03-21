import { Navigate, Route, Routes } from "react-router-dom";
import { useState } from "react";
import Cookies from "js-cookie";
import MainLayout from "./components/layouts/MainLayout";
import Analytics from "./components/pages/Analytics";
import Login from "./components/pages/Login";
import Register from "./components/pages/Register";
import { decodeJwt } from "jose";
import { useProfile } from "./store/useProfile";
import { useShallow } from "zustand/shallow";
import { data } from "./helper/Data";
import Wallets from "./components/pages/Wallets";
import Transactions from "./components/pages/Transactions";
import Investments from "./components/pages/Investments";
import { createElement } from "react";

function App() {
  const { setProfile } = useProfile(
    useShallow((state) => ({ setProfile: state.setProfile })),
  );

  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    return Cookies.get("token") ? true : false;
  });

  const handleLogin = () => {
    const token = Cookies.get("token");
    if (token) {
      const decoded = decodeJwt(token);
      setProfile({
        username: decoded.username as string,
        email: decoded.email as string,
      });
      setIsAuthenticated(true);
    }
  };

  const components = {
    Analytics,
    Wallets,
    Transactions,
    Investments,
  } as const;

  return (
    <Routes>
      <Route element={<MainLayout />}>
        {data.navMain.map((item) => (
          <Route
            key={item.title}
            path={item.url}
            element={
              <ProtectedRoute isAuthenticated={isAuthenticated}>
                {createElement(
                  components[item.title as keyof typeof components],
                )}
              </ProtectedRoute>
            }
          />
        ))}
      </Route>

      <Route
        path="login"
        element={
          <Login isAuthenticated={isAuthenticated} handleLogin={handleLogin} />
        }
      />
      <Route
        path="register"
        element={
          <Register
            isAuthenticated={isAuthenticated}
            handleLogin={handleLogin}
          />
        }
      />
    </Routes>
  );
}

export default App;

function ProtectedRoute(props: {
  isAuthenticated: boolean;
  children: React.ReactNode;
}) {
  return props.isAuthenticated ? props.children : <Navigate to="/login" />;
}
