import { Navigate, Route, Routes } from "react-router-dom";
import { useState } from "react";
import Cookies from "js-cookie";
import MainLayout from "./components/layouts/MainLayout";
import Analytics from "./components/pages/Analytics";
import Login from "./components/pages/Auth/Login";
import Register from "./components/pages/Auth/Register";
import { decodeJwt } from "jose";
import { useProfile } from "./store/useProfile";
import { useShallow } from "zustand/shallow";
import Wallets from "./components/pages/Wallet/Wallets";
import Transactions from "./components/pages/Transaction/Transactions";
import Investments from "./components/pages/Investments";
import CreateWallet from "./components/pages/Wallet/CreateWallet";
import AddTransaction from "./components/pages/Transaction/AddTransaction";

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

  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route
          path={"/"}
          element={
            <ProtectedRoute isAuthenticated={isAuthenticated}>
              <Analytics />
            </ProtectedRoute>
          }
        />
        <Route
          path={"/wallets"}
          element={
            <ProtectedRoute isAuthenticated={isAuthenticated}>
              <Wallets />
            </ProtectedRoute>
          }
        />
        <Route
          path={"/wallets/create"}
          element={
            <ProtectedRoute isAuthenticated={isAuthenticated}>
              <CreateWallet />
            </ProtectedRoute>
          }
        />
        <Route
          path={"/transactions"}
          element={
            <ProtectedRoute isAuthenticated={isAuthenticated}>
              <Transactions />
            </ProtectedRoute>
          }
        />
        <Route
          path={"/transactions/add/:type"}
          element={
            <ProtectedRoute isAuthenticated={isAuthenticated}>
              <AddTransaction />
            </ProtectedRoute>
          }
        />
        <Route
          path={"/investments"}
          element={
            <ProtectedRoute isAuthenticated={isAuthenticated}>
              <Investments />
            </ProtectedRoute>
          }
        />
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
