import { useState, useEffect } from "react";
import Login from "./components/Login";
import Register from "./components/Register";
import authService, { type User } from "./services/authService";
import "./App.css";

function App() {
  const [isLogin, setIsLogin] = useState(true);
  const [user, setUser] = useState<User | null>(null);

  useEffect(() => {
    const currentUser = authService.getCurrentUser();
    if (currentUser) {
      setUser(currentUser);
    }
  }, []);

  const handleLoginSuccess = () => {
    const currentUser = authService.getCurrentUser();
    setUser(currentUser);
  };

  const handleRegisterSuccess = () => {
    const currentUser = authService.getCurrentUser();
    setUser(currentUser);
  };

  const handleLogout = () => {
    authService.logout();
    setUser(null);
  };

  if (user) {
    return (
      <div className="app">
        <div className="dashboard">
          <h1>Welcome to Echo Server</h1>
          <div className="user-info">
            <p>
              <strong>Email:</strong> {user.email}
            </p>
            <p>
              <strong>User ID:</strong> {user.id}
            </p>
            <p>
              <strong>Account created:</strong>{" "}
              {new Date(user.created_at).toLocaleDateString()}
            </p>
          </div>
          <button onClick={handleLogout} className="logout-btn">
            Logout
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="app">
      {isLogin ? (
        <Login
          onSwitchToRegister={() => setIsLogin(false)}
          onLoginSuccess={handleLoginSuccess}
        />
      ) : (
        <Register
          onSwitchToLogin={() => setIsLogin(true)}
          onRegisterSuccess={handleRegisterSuccess}
        />
      )}
    </div>
  );
}

export default App;
