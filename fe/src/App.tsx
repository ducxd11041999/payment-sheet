import React, { useState, useEffect } from "react";
import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import LoginForm from "./components/login_form";
import Dashboard from "./components/dashboard";
import BlockDetail from "./components/block_detail";
import { CssBaseline } from "@mui/material";

const App = () => {
  const [token, setToken] = useState<string | null>(localStorage.getItem("token"));

  useEffect(() => {
    if (token) {
      localStorage.setItem("token", token);
    } else {
      localStorage.removeItem("token");
    }
  }, [token]);

  return (
    <>
      <CssBaseline />
      <Router>
        <Routes>
          {!token && (
            <Route
              path="/login"
              element={<LoginForm onLogin={(t) => setToken(t)} />}
            />
          )}
          {token ? (
            <>
              <Route path="/home" element={<Dashboard />} />
              <Route path="/blocks/:month" element={<BlockDetail />} />
              <Route path="*" element={<Navigate to="/home" replace />} />
            </>
          ) : (
            <Route path="*" element={<Navigate to="/login" replace />} />
          )}
        </Routes>
      </Router>
    </>
  );
};

export default App;
