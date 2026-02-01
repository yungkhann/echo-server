import { useState, useEffect } from "react";
import Login from "./components/Login";
import Register from "./components/Register";
import StudentView from "./components/StudentView";
import ScheduleView from "./components/ScheduleView";
import AttendanceView from "./components/AttendanceView";
import UsersView from "./components/UsersView";
import authService, { type User } from "./services/authService";

type View = "dashboard" | "students" | "schedule" | "attendance" | "users";

function App() {
  const [isLogin, setIsLogin] = useState(true);
  const [user, setUser] = useState<User | null>(null);
  const [currentView, setCurrentView] = useState<View>("dashboard");

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
    setCurrentView("dashboard");
  };

  if (user) {
    return (
      <div className="min-h-screen bg-gray-50">
        <nav className="bg-white border-b border-gray-200">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center h-16">
              <div className="text-xl font-bold text-gray-900">
                University System
              </div>

              <div className="flex items-center space-x-1">
                <button
                  className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                    currentView === "dashboard"
                      ? "bg-gray-900 text-white"
                      : "text-gray-600 hover:bg-gray-100"
                  }`}
                  onClick={() => setCurrentView("dashboard")}
                >
                  Dashboard
                </button>
                {(user.role === "admin" ||
                  user.role === "teacher" ||
                  user.role === "student") && (
                  <button
                    className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                      currentView === "students"
                        ? "bg-gray-900 text-white"
                        : "text-gray-600 hover:bg-gray-100"
                    }`}
                    onClick={() => setCurrentView("students")}
                  >
                    Students
                  </button>
                )}
                {(user.role === "admin" ||
                  user.role === "teacher" ||
                  user.role === "student") && (
                  <button
                    className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                      currentView === "schedule"
                        ? "bg-gray-900 text-white"
                        : "text-gray-600 hover:bg-gray-100"
                    }`}
                    onClick={() => setCurrentView("schedule")}
                  >
                    Schedule
                  </button>
                )}
                {(user.role === "admin" || user.role === "teacher") && (
                  <button
                    className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                      currentView === "attendance"
                        ? "bg-gray-900 text-white"
                        : "text-gray-600 hover:bg-gray-100"
                    }`}
                    onClick={() => setCurrentView("attendance")}
                  >
                    Attendance
                  </button>
                )}
                {user.role === "admin" && (
                  <button
                    className={`px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                      currentView === "users"
                        ? "bg-gray-900 text-white"
                        : "text-gray-600 hover:bg-gray-100"
                    }`}
                    onClick={() => setCurrentView("users")}
                  >
                    Users
                  </button>
                )}
              </div>

              <div className="flex items-center space-x-3">
                <span className="px-3 py-1 bg-gray-900 text-white text-xs font-semibold rounded-full uppercase">
                  {user.role}
                </span>
                <button
                  onClick={handleLogout}
                  className="px-4 py-2 text-sm font-medium text-gray-700 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
                >
                  Logout
                </button>
              </div>
            </div>
          </div>
        </nav>

        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
          {currentView === "dashboard" && (
            <div className="bg-white rounded-lg shadow-sm p-8">
              <h1 className="text-3xl font-bold text-gray-900 mb-8">
                Welcome, {user.full_name || user.email}
              </h1>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
                <div className="bg-gray-900 text-white rounded-lg p-6">
                  <h3 className="text-sm font-medium opacity-80 mb-2">Role</h3>
                  <p className="text-2xl font-bold capitalize">{user.role}</p>
                </div>
                <div className="bg-gray-900 text-white rounded-lg p-6">
                  <h3 className="text-sm font-medium opacity-80 mb-2">
                    User ID
                  </h3>
                  <p className="text-2xl font-bold">{user.id}</p>
                </div>
                <div className="bg-gray-900 text-white rounded-lg p-6">
                  <h3 className="text-sm font-medium opacity-80 mb-2">
                    Account Created
                  </h3>
                  <p className="text-2xl font-bold">
                    {new Date(user.created_at).toLocaleDateString()}
                  </p>
                </div>
              </div>
              <p className="text-center text-gray-600">
                {user.role === "admin" &&
                  "You have full access to all features"}
                {user.role === "teacher" &&
                  "You can view schedules and manage attendance"}
                {user.role === "student" &&
                  "You can view your schedules and student information"}
              </p>
            </div>
          )}

          {currentView === "students" && <StudentView />}
          {currentView === "schedule" && <ScheduleView />}
          {currentView === "attendance" && <AttendanceView />}
          {currentView === "users" && <UsersView />}
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-900">
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
