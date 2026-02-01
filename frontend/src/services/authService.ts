import axios from "axios";

const API_URL = `${import.meta.env.VITE_API_URL}/api/auth`;

export interface User {
  id: number;
  email: string;
  role: string;
  full_name?: string;
  created_at: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface RegisterData {
  email: string;
  password: string;
  role: string;
  full_name: string;
}

export interface LoginData {
  email: string;
  password: string;
}

class AuthService {
  async register(data: RegisterData): Promise<LoginResponse> {
    const response = await axios.post<LoginResponse>(
      `${API_URL}/register`,
      data,
    );
    if (response.data.token) {
      localStorage.setItem("token", response.data.token);
      localStorage.setItem("user", JSON.stringify(response.data.user));
    }
    return response.data;
  }

  async login(data: LoginData): Promise<LoginResponse> {
    const response = await axios.post<LoginResponse>(`${API_URL}/login`, data);
    if (response.data.token) {
      localStorage.setItem("token", response.data.token);
      localStorage.setItem("user", JSON.stringify(response.data.user));
    }
    return response.data;
  }

  logout(): void {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
  }

  getCurrentUser(): User | null {
    const userStr = localStorage.getItem("user");
    if (userStr) {
      return JSON.parse(userStr);
    }
    return null;
  }

  getToken(): string | null {
    return localStorage.getItem("token");
  }

  isAuthenticated(): boolean {
    return !!this.getToken();
  }
}

export default new AuthService();
