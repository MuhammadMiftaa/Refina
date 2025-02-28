import { Route, Routes } from "react-router";
import "./App.css";
import { useState } from "react";
import Cookies from "js-cookie";
import MainLayout from "./components/layouts/MainLayout";
import Home from "./components/pages/Home";

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(() => {
    return Cookies.get("token") ? true : false;
  });

  const handleLogin = () => {
    setIsAuthenticated(true);
  };
  return (
    <Routes>
      <Route path="/" element={<MainLayout />}>
        <Route path="dashboard" element={<Home />} />
      </Route>

      <Route path="login" element />
    </Routes>
  );
}

export default App;
