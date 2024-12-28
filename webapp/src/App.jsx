import { useEffect, useState} from 'react';

import LoginPage from './components/LoginPage'
import HomePage from './components/HomePage'

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
      const adminSecret = localStorage.getItem('adminSecret');
      if (adminSecret) {
          setIsAuthenticated(true);
      }
  }, []);

  if (isAuthenticated) {
    return (
      <HomePage />
    )
  }

  return (
    <LoginPage setIsAuthenticated={setIsAuthenticated} />
  )
}

export default App
