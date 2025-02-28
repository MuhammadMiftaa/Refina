import { Navigate, Route, Routes } from "react-router-dom";
import { useState } from "react";
import Cookies from "js-cookie";
import MainLayout from "./components/layouts/MainLayout";
import Home from "./components/pages/Home";
import Login from "./components/pages/Login";
import Register from "./components/pages/Register";
import { decodeJwt } from "jose";
import { useProfile } from "./store/useProfile";
import { useShallow } from "zustand/shallow";

function App() {
  const { setProfile } = useProfile(
    useShallow((state) => ({ setProfile: state.setProfile }))
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
          path="/"
          element={
            <ProtectedRoute isAuthenticated={isAuthenticated}>
              <Home />
            </ProtectedRoute>
          }
        />
        {/* <Route path="dashboard" element={<Home />} /> */}
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
