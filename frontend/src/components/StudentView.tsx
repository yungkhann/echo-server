import { useState, useEffect } from "react";
import axios from "axios";
import authService from "../services/authService";
import apiService, { type Student } from "../services/apiService";

const API_URL = import.meta.env.VITE_API_URL;

export default function StudentView() {
  const [students, setStudents] = useState<Student[]>([]);
  const [studentId, setStudentId] = useState("");
  const [selectedStudent, setSelectedStudent] = useState<Student | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [showSearch, setShowSearch] = useState(false);

  useEffect(() => {
    loadAllStudents();
  }, []);

  const loadAllStudents = async () => {
    setLoading(true);
    setError("");
    try {
      const response = await axios.get(`${API_URL}/students`, {
        headers: {
          Authorization: `Bearer ${authService.getToken()}`,
        },
      });
      setStudents(response.data);
    } catch (err: any) {
      setError(err.response?.data?.error || "Failed to load students");
    } finally {
      setLoading(false);
    }
  };

  const loadStudent = async () => {
    if (!studentId) return;
    setLoading(true);
    setError("");
    try {
      const data = await apiService.getStudent(Number(studentId));
      setSelectedStudent(data);
      setShowSearch(false);
    } catch (err: any) {
      setError(err.response?.data?.error || "Student not found");
      setSelectedStudent(null);
    } finally {
      setLoading(false);
    }
  };

  const viewStudent = (student: Student) => {
    setSelectedStudent(student);
  };

  const closeModal = () => {
    setSelectedStudent(null);
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Students</h2>
        <div className="flex gap-2">
          <button
            onClick={() => setShowSearch(!showSearch)}
            className="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors font-medium"
          >
            {showSearch ? "Hide Search" : "Search by ID"}
          </button>
          <button
            onClick={loadAllStudents}
            className="px-4 py-2 bg-gray-900 text-white rounded-md hover:bg-gray-800 transition-colors font-medium"
          >
            Refresh
          </button>
        </div>
      </div>

      {showSearch && (
        <div className="flex gap-2 mb-6">
          <input
            type="number"
            placeholder="Enter Student ID"
            value={studentId}
            onChange={(e) => setStudentId(e.target.value)}
            className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          />
          <button
            onClick={loadStudent}
            className="px-6 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 transition-colors font-medium"
          >
            Search
          </button>
        </div>
      )}

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4">
          {error}
        </div>
      )}

      {loading ? (
        <div className="text-center py-8 text-gray-600">Loading...</div>
      ) : (
        <div className="overflow-x-auto">
          {students.length === 0 ? (
            <p className="text-center py-8 text-gray-500">No students found</p>
          ) : (
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    ID
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    Full Name
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    Gender
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    Birth Date
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    Group
                  </th>
                  <th className="px-6 py-3 text-left text-xs font-medium text-gray-700 uppercase tracking-wider">
                    Actions
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {students.map((student) => (
                  <tr key={student.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      <span
                        className={`${
                          student.id < 0
                            ? "text-orange-600 font-semibold"
                            : "text-gray-900"
                        }`}
                      >
                        {student.id}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {student.full_name}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {student.gender}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {student.birth_date || "-"}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {student.group_name || "No Group"}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm">
                      <button
                        onClick={() => viewStudent(student)}
                        className="text-indigo-600 hover:text-indigo-900 font-medium"
                      >
                        View Details
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          )}
        </div>
      )}

      <div className="mt-4 text-sm text-gray-600">
        Total students: {students.length}
        {students.some((s) => s.id < 0) && (
          <span className="ml-2 text-orange-600">
            (IDs in orange need profile completion)
          </span>
        )}
      </div>

      {/* Student Details Modal */}
      {selectedStudent && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg p-6 w-full max-w-md">
            <h3 className="text-xl font-bold text-gray-900 mb-4">
              Student Details
            </h3>

            <div className="bg-gray-50 rounded-lg p-4 space-y-3">
              <div className="flex justify-between py-2 border-b border-gray-200">
                <span className="font-semibold text-gray-700">ID:</span>
                <span className="text-gray-900">{selectedStudent.id}</span>
              </div>
              <div className="flex justify-between py-2 border-b border-gray-200">
                <span className="font-semibold text-gray-700">Full Name:</span>
                <span className="text-gray-900">
                  {selectedStudent.full_name}
                </span>
              </div>
              <div className="flex justify-between py-2 border-b border-gray-200">
                <span className="font-semibold text-gray-700">Gender:</span>
                <span className="text-gray-900">{selectedStudent.gender}</span>
              </div>
              <div className="flex justify-between py-2 border-b border-gray-200">
                <span className="font-semibold text-gray-700">Birth Date:</span>
                <span className="text-gray-900">
                  {selectedStudent.birth_date || "-"}
                </span>
              </div>
              <div className="flex justify-between py-2 border-b border-gray-200">
                <span className="font-semibold text-gray-700">Group ID:</span>
                <span className="text-gray-900">
                  {selectedStudent.group_id || "-"}
                </span>
              </div>
              {selectedStudent.group_name && (
                <div className="flex justify-between py-2">
                  <span className="font-semibold text-gray-700">
                    Group Name:
                  </span>
                  <span className="text-gray-900">
                    {selectedStudent.group_name}
                  </span>
                </div>
              )}
            </div>

            <div className="mt-6">
              <button
                onClick={closeModal}
                className="w-full px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors font-medium"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
