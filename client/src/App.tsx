import { Routes, Route, Navigate } from "react-router-dom";
import LoginPage from "./pages/login";
import RegisterPage from "./pages/register";
import { useState } from "react";
import HomePage from "./pages/home";
import TransactionsPage from "./pages/transactions/page";
import FormAddPage from "./pages/transactions/form-add";

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState<boolean>(() => {
    const cookieString = document.cookie;
    const cookies = cookieString.split("; ").reduce((acc, cookie) => {
      const [key, value] = cookie.split("=");
      acc[key] = value;
      return acc;
    }, {} as Record<string, string>);
    return cookies.token !== undefined;
  });

  const handleLogin = () => {
    setIsAuthenticated(true);
  };

  return (
    <Routes>
      <Route path="/" element={<ProtectedRoute isAuthenticated={isAuthenticated}><HomePage /></ProtectedRoute>}/>
      <Route path="transactions">
        <Route index element={<ProtectedRoute isAuthenticated={isAuthenticated}><TransactionsPage /></ProtectedRoute>} />
        <Route path="add" element={<ProtectedRoute isAuthenticated={isAuthenticated}><FormAddPage /></ProtectedRoute>} />
      </Route>
      <Route path="login" element={<LoginPage handleLogin={handleLogin} isAuthenticated={isAuthenticated}/>}/>
      <Route path="register" element={<RegisterPage handleLogin={handleLogin} isAuthenticated={isAuthenticated}/>}/>
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
