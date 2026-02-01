import { useState } from "react";
import apiService, { type Student } from "../services/apiService";

export default function StudentView() {
  const [studentId, setStudentId] = useState("");
  const [student, setStudent] = useState<Student | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const loadStudent = async () => {
    if (!studentId) return;
    setLoading(true);
    setError("");
    try {
      const data = await apiService.getStudent(Number(studentId));
      setStudent(data);
    } catch (err: any) {
      setError(err.response?.data?.error || "Student not found");
      setStudent(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-sm p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-6">
        Student Information
      </h2>

      <div className="flex gap-2 mb-6">
        <input
          type="number"
          placeholder="Enter Student ID"
          value={studentId}
          onChange={(e) => setStudentId(e.target.value)}
          className="px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-gray-900 focus:border-transparent"
        />
        <button
          onClick={loadStudent}
          className="px-6 py-2 bg-gray-900 text-white rounded-md hover:bg-gray-800 transition-colors font-medium"
        >
          Search
        </button>
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md mb-4">
          {error}
        </div>
      )}

      {loading ? (
        <div className="text-center py-8 text-gray-600">Loading...</div>
      ) : student ? (
        <div className="bg-gray-50 rounded-lg p-6 space-y-4">
          <div className="flex justify-between py-3 border-b border-gray-200">
            <span className="font-semibold text-gray-700">Full Name:</span>
            <span className="text-gray-900">{student.full_name}</span>
          </div>
          <div className="flex justify-between py-3 border-b border-gray-200">
            <span className="font-semibold text-gray-700">Gender:</span>
            <span className="text-gray-900">{student.gender}</span>
          </div>
          <div className="flex justify-between py-3 border-b border-gray-200">
            <span className="font-semibold text-gray-700">Birth Date:</span>
            <span className="text-gray-900">{student.birth_date}</span>
          </div>
          <div className="flex justify-between py-3 border-b border-gray-200">
            <span className="font-semibold text-gray-700">Group ID:</span>
            <span className="text-gray-900">{student.group_id}</span>
          </div>
          {student.group_name && (
            <div className="flex justify-between py-3">
              <span className="font-semibold text-gray-700">Group Name:</span>
              <span className="text-gray-900">{student.group_name}</span>
            </div>
          )}
        </div>
      ) : null}
    </div>
  );
}
