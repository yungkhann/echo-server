import { useState, useEffect } from "react";
import axios from "axios";
import authService from "../services/authService";
import apiService, { type Attendance } from "../services/apiService";

const API_URL = "http://localhost:8080";

interface Subject {
  id: number;
  subject_name: string;
  subject_code: string;
  credits: number;
}

export default function AttendanceView() {
  const [studentId, setStudentId] = useState("");
  const [subjectId, setSubjectId] = useState("");
  const [attendances, setAttendances] = useState<Attendance[]>([]);
  const [subjects, setSubjects] = useState<Subject[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [viewMode, setViewMode] = useState<"student" | "subject" | null>(null);

  const [newAttendance, setNewAttendance] = useState({
    student_id: "",
    subject_id: "",
    visit_day: "",
    visited: true,
  });
  const [postSuccess, setPostSuccess] = useState("");

  useEffect(() => {
    loadSubjects();
  }, []);

  const loadSubjects = async () => {
    try {
      const response = await axios.get(`${API_URL}/subjects`, {
        headers: {
          Authorization: `Bearer ${authService.getToken()}`,
        },
      });
      setSubjects(response.data);
    } catch (err: any) {
      console.error("Failed to load subjects:", err);
    }
  };

  const loadAttendanceByStudent = async () => {
    if (!studentId) return;
    setLoading(true);
    setError("");
    setPostSuccess("");
    try {
      const data = await apiService.getAttendanceByStudent(Number(studentId));
      setAttendances(data);
      setViewMode("student");
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to load attendance");
    } finally {
      setLoading(false);
    }
  };

  const loadAttendanceBySubject = async () => {
    if (!subjectId) return;
    setLoading(true);
    setError("");
    setPostSuccess("");
    try {
      const data = await apiService.getAttendanceBySubject(Number(subjectId));
      setAttendances(data);
      setViewMode("subject");
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to load attendance");
    } finally {
      setLoading(false);
    }
  };

  const handlePostAttendance = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setPostSuccess("");
    setLoading(true);

    try {
      await apiService.postAttendance({
        student_id: Number(newAttendance.student_id),
        subject_id: Number(newAttendance.subject_id),
        visit_day: newAttendance.visit_day,
        visited: newAttendance.visited,
      });
      setPostSuccess("Attendance recorded successfully!");
      setNewAttendance({
        student_id: "",
        subject_id: "",
        visit_day: "",
        visited: true,
      });
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to post attendance");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">
        Attendance Management
      </h2>

      {/* Post Attendance Form */}
      <div className="bg-gray-50 rounded-lg p-6 mb-6">
        <h3 className="text-lg font-semibold text-gray-900 mb-4">
          Record Attendance
        </h3>
        <form onSubmit={handlePostAttendance} className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <input
              type="number"
              placeholder="Student ID"
              value={newAttendance.student_id}
              onChange={(e) =>
                setNewAttendance({
                  ...newAttendance,
                  student_id: e.target.value,
                })
              }
              required
              className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
            />
            <select
              value={newAttendance.subject_id}
              onChange={(e) =>
                setNewAttendance({
                  ...newAttendance,
                  subject_id: e.target.value,
                })
              }
              required
              className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
            >
              <option value="">Select Subject</option>
              {subjects.map((subject) => (
                <option key={subject.id} value={subject.id}>
                  {subject.subject_name} ({subject.subject_code})
                </option>
              ))}
            </select>
          </div>
          <div className="flex gap-4 items-center">
            <input
              type="date"
              value={newAttendance.visit_day}
              onChange={(e) =>
                setNewAttendance({
                  ...newAttendance,
                  visit_day: e.target.value,
                })
              }
              required
              className="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
            />
            <label className="flex items-center gap-2 text-sm font-medium text-gray-700 cursor-pointer">
              <input
                type="checkbox"
                checked={newAttendance.visited}
                onChange={(e) =>
                  setNewAttendance({
                    ...newAttendance,
                    visited: e.target.checked,
                  })
                }
                className="w-4 h-4 text-gray-900 border-gray-300 rounded focus:ring-gray-900"
              />
              Attended
            </label>
          </div>
          <button
            type="submit"
            disabled={loading}
            className="w-full px-6 py-2 bg-gray-900 text-white rounded-md hover:bg-gray-800 transition-colors disabled:opacity-50 disabled:cursor-not-allowed font-medium"
          >
            {loading ? "Recording..." : "Record Attendance"}
          </button>
        </form>
      </div>

      {postSuccess && (
        <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-md mb-4">
          {postSuccess}
        </div>
      )}

      {/* View Attendance */}
      <div className="flex gap-2 flex-wrap mb-6">
        <div className="flex gap-2">
          <input
            type="number"
            placeholder="Student ID"
            value={studentId}
            onChange={(e) => setStudentId(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
          />
          <button
            onClick={loadAttendanceByStudent}
            disabled={!studentId}
            className="px-6 py-2 bg-gray-900 text-white rounded-md hover:bg-gray-800 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          >
            View by Student
          </button>
        </div>
        <div className="flex gap-2">
          <select
            value={subjectId}
            onChange={(e) => setSubjectId(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
          >
            <option value="">Select Subject</option>
            {subjects.map((subject) => (
              <option key={subject.id} value={subject.id}>
                {subject.subject_name} ({subject.subject_code})
              </option>
            ))}
          </select>
          <button
            onClick={loadAttendanceBySubject}
            disabled={!subjectId}
            className="px-6 py-2 bg-gray-900 text-white rounded-md hover:bg-gray-800 transition-colors font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          >
            View by Subject
          </button>
        </div>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4">
          {error}
        </div>
      )}

      {loading && viewMode ? (
        <div className="text-center py-8 text-gray-600">Loading...</div>
      ) : (
        <div className="overflow-x-auto">
          {!attendances || attendances.length === 0 ? (
            viewMode && (
              <p className="text-center py-8 text-gray-500">
                No attendance records found
              </p>
            )
          ) : (
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    Student ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    Subject ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    Visit Day
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    Status
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {attendances.map((attendance) => (
                  <tr key={attendance.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {attendance.id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {attendance.student_id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {attendance.subject_id}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {attendance.visit_day}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap">
                      <span
                        className={`px-2 py-1 text-xs font-semibold rounded-full ${
                          attendance.visited
                            ? "bg-green-100 text-green-800"
                            : "bg-red-100 text-red-800"
                        }`}
                      >
                        {attendance.visited ? "Present" : "Absent"}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}
    </div>
  );
}
