import axios from "axios";
import authService from "./authService";

const API_URL = import.meta.env.VITE_API_URL;

const getAuthHeaders = () => ({
  Authorization: `Bearer ${authService.getToken()}`,
});

axios.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      authService.logout();
      window.location.href = "/login";
    }
    return Promise.reject(error);
  },
);

export interface Student {
  id?: number;
  full_name: string;
  gender: string;
  birth_date: string;
  group_id: number;
  group_name?: string;
}

export interface Schedule {
  id: number;
  subject_name: string;
  time_slot: string;
  group_id: number;
  group_name: string;
}

export interface Attendance {
  id?: number;
  subject_id: number;
  visit_day: string;
  visited: boolean;
  student_id: number;
}

class ApiService {
  async getStudent(id: number): Promise<Student> {
    const response = await axios.get(`${API_URL}/student/${id}`, {
      headers: getAuthHeaders(),
    });
    return response.data;
  }

  async getAllSchedule(): Promise<Schedule[]> {
    const response = await axios.get(`${API_URL}/all_class_schedule`, {
      headers: getAuthHeaders(),
    });
    return response.data;
  }

  async getScheduleByGroup(groupId: number): Promise<Schedule[]> {
    const response = await axios.get(`${API_URL}/schedule/group/${groupId}`, {
      headers: getAuthHeaders(),
    });
    return response.data;
  }

  async postAttendance(data: Attendance): Promise<any> {
    const response = await axios.post(`${API_URL}/attendance/subject`, data, {
      headers: getAuthHeaders(),
    });
    return response.data;
  }

  async getAttendanceByStudent(studentId: number): Promise<Attendance[]> {
    const response = await axios.get(
      `${API_URL}/attendanceByStudentId/${studentId}`,
      {
        headers: getAuthHeaders(),
      },
    );
    return response.data;
  }

  async getAttendanceBySubject(subjectId: number): Promise<Attendance[]> {
    const response = await axios.get(
      `${API_URL}/attendanceBySubjectId/${subjectId}`,
      {
        headers: getAuthHeaders(),
      },
    );
    return response.data;
  }
}

export default new ApiService();
