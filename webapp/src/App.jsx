import { useEffect, useState } from "react";

import LoginPage from "./components/LoginPage";
import HomePage from "./components/home/Page";

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  const isLoggedIn = () => {
    return !!localStorage.getItem("adminSecret");
  };

  const login = (input) => {
    if (!input || !input.adminSecret) {
      return;
    }
    localStorage.setItem("adminSecret", input.adminSecret);
    setIsAuthenticated(true);
  };

  const logout = () => {
    localStorage.removeItem("adminSecret");
    setIsAuthenticated(false);
  };

  useEffect(() => {
    if (isLoggedIn()) {
      setIsAuthenticated(true);
    }
  }, []);

  if (isAuthenticated) {
    return <HomePage logout={logout} />;
  }

  return <LoginPage login={login} />;
}

export default App;
