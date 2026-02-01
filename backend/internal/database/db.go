package database

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Student struct{
    ID        int    `json:"id,omitempty"`
    FullName  string `json:"full_name"`
    Gender    string `json:"gender"`
    BirthDate string `json:"birth_date"`
    GroupID   int    `json:"group_id"`
    GroupName string `json:"group_name"`
}

type Group struct{
    ID        int    `json:"id"`
    GroupName string `json:"group_name"`
    FacultyID int    `json:"faculty_id"`
}
type Schedule struct{
    ID          int    `json:"id"`
    SubjectName string `json:"subject_name"`
    TimeSlot    string `json:"time_slot"`
    GroupID     int    `json:"group_id"`
    GroupName   string `json:"group_name"`
}

type Attendance struct{
    ID          int    `json:"id"`
    SubjectID   int    `json:"subject_id"`
    VisitDay    string `json:"visit_day"`
    Visited     bool   `json:"visited"`
    StudentId   int    `json:"student_id"`
}

var pool *pgxpool.Pool

func InitDB() *pgxpool.Pool {
    connStr := os.Getenv("DATABASE_URL")
    if connStr == "" {
        log.Fatal("DATABASE_URL is not set in .env file")
    }

    config, err := pgxpool.ParseConfig(connStr)
    if err != nil {
        log.Fatalf("Unable to parse connection string: %v", err)
    }

    pool, err = pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        log.Fatalf("Unable to create connection pool: %v", err)
    }

    err = pool.Ping(context.Background())
    if err != nil {
        log.Fatalf("Could not ping database: %v", err)
    }

    fmt.Println("Connected to PostgreSQL")
    return pool
}

func GetStudentHandler(c echo.Context) error{
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
    }

    userRole := c.Get("user_role").(string)
    userID := c.Get("user_id").(int)
    
    if userRole == "student" {
        var studentID int
        err := pool.QueryRow(context.Background(), "SELECT student_id FROM users WHERE id = $1", userID).Scan(&studentID)
        if err != nil || studentID != id {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied: you can only view your own profile"})
        }
    }

    student, err := getStudentById(pool, id)
    if err != nil {
        if err == pgx.ErrNoRows {
            return c.JSON(http.StatusNotFound, map[string]string{"error": "Student not found"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }
    return c.JSON(http.StatusOK, student)
}

func getStudentById(pool *pgxpool.Pool, id int) (*Student, error) {
    student := &Student{}

    query := `
        SELECT s.id, s.full_name, s.gender, s.birth_date::text, s.group_id, sg.group_name 
        FROM students s
        LEFT JOIN student_groups sg ON s.group_id = sg.id
        WHERE s.id = $1
    `
    row := pool.QueryRow(context.Background(), query, id)

    err := row.Scan(&student.ID, &student.FullName, &student.Gender, &student.BirthDate, &student.GroupID, &student.GroupName)
    if err != nil {
        return nil, err
    }

    return student, nil
}

func GetAllStudentsHandler(c echo.Context) error {
    query := `
        SELECT 
            COALESCE(s.id, -u.id) as id,
            COALESCE(s.full_name, u.full_name, u.email) as full_name,
            COALESCE(s.gender, 'Unknown') as gender,
            COALESCE(s.birth_date::text, '') as birth_date,
            COALESCE(s.group_id, 0) as group_id,
            COALESCE(sg.group_name, 'No Group') as group_name
        FROM users u
        LEFT JOIN students s ON u.student_id = s.id
        LEFT JOIN student_groups sg ON s.group_id = sg.id
        WHERE u.role = 'student'
        
        UNION
        
        SELECT 
            s.id,
            s.full_name,
            s.gender,
            s.birth_date::text,
            s.group_id,
            COALESCE(sg.group_name, 'No Group') as group_name
        FROM students s
        LEFT JOIN student_groups sg ON s.group_id = sg.id
        WHERE s.id NOT IN (SELECT student_id FROM users WHERE student_id IS NOT NULL)
        
        ORDER BY id DESC
    `
    rows, err := pool.Query(context.Background(), query)
    if err != nil {
        fmt.Printf("GetAllStudentsHandler error: %v\n", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }
    defer rows.Close()

    var students []Student
    for rows.Next() {
        var student Student
        err := rows.Scan(&student.ID, &student.FullName, &student.Gender, &student.BirthDate, &student.GroupID, &student.GroupName)
        if err != nil {
            fmt.Printf("Row scan error: %v\n", err)
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
        }
        students = append(students, student)
    }

    return c.JSON(http.StatusOK, students)
}

type CreateStudentRequest struct {
    FullName  string `json:"full_name"`
    Gender    string `json:"gender"`
    BirthDate string `json:"birth_date"`
    GroupID   int    `json:"group_id"`
    UserID    int    `json:"user_id,omitempty"`
}

func CreateStudentHandler(c echo.Context) error {
    var req CreateStudentRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
    }

    if req.FullName == "" || req.Gender == "" || req.BirthDate == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields"})
    }

    var studentID int
    query := `INSERT INTO students (full_name, gender, birth_date, group_id) VALUES ($1, $2, $3, $4) RETURNING id`
    err := pool.QueryRow(context.Background(), query, req.FullName, req.Gender, req.BirthDate, req.GroupID).Scan(&studentID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create student"})
    }

    if req.UserID > 0 {
        _, err = pool.Exec(context.Background(), "UPDATE users SET student_id = $1 WHERE id = $2", studentID, req.UserID)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Student created but failed to link to user"})
        }
    }

    return c.JSON(http.StatusCreated, map[string]interface{}{
        "id":         studentID,
        "full_name":  req.FullName,
        "gender":     req.Gender,
        "birth_date": req.BirthDate,
        "group_id":   req.GroupID,
    })
}

type CreateStudentFromUserRequest struct {
    UserID    int    `json:"user_id"`
    Gender    string `json:"gender"`
    BirthDate string `json:"birth_date"`
    GroupID   int    `json:"group_id"`
}

func CreateStudentFromUserHandler(c echo.Context) error {
    var req CreateStudentFromUserRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
    }

    if req.UserID == 0 || req.Gender == "" || req.BirthDate == "" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing required fields: user_id, gender, birth_date"})
    }

    var fullName, email, role string
    var existingStudentID *int
    err := pool.QueryRow(context.Background(), 
        "SELECT COALESCE(full_name, email), email, role, student_id FROM users WHERE id = $1", 
        req.UserID).Scan(&fullName, &email, &role, &existingStudentID)
    if err != nil {
        if err == pgx.ErrNoRows {
            return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }

    if role != "student" {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "User is not a student"})
    }

    if existingStudentID != nil {
        return c.JSON(http.StatusConflict, map[string]string{"error": "User already has a student profile"})
    }

    var studentID int
    query := `INSERT INTO students (full_name, email, gender, birth_date, group_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`
    err = pool.QueryRow(context.Background(), query, fullName, email, req.Gender, req.BirthDate, req.GroupID).Scan(&studentID)
    if err != nil {
        fmt.Printf("Failed to create student: %v\n", err)
        if strings.Contains(err.Error(), "foreign key constraint") || strings.Contains(err.Error(), "group_id_fkey") {
            return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group_id. Please select a valid group"})
        }
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create student"})
    }

    _, err = pool.Exec(context.Background(), "UPDATE users SET student_id = $1 WHERE id = $2", studentID, req.UserID)
    if err != nil {
        fmt.Printf("Failed to link user to student: %v\n", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Student created but failed to link to user"})
    }

    return c.JSON(http.StatusCreated, map[string]interface{}{
        "id":         studentID,
        "user_id":    req.UserID,
        "full_name":  fullName,
        "gender":     req.Gender,
        "birth_date": req.BirthDate,
        "group_id":   req.GroupID,
    })
}

func GetAllGroupsHandler(c echo.Context) error {
    query := `SELECT id, group_name, faculty_id, course_year FROM student_groups ORDER BY group_name`
    rows, err := pool.Query(context.Background(), query)
    if err != nil {
        fmt.Printf("Failed to get groups: %v\n", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch groups"})
    }
    defer rows.Close()

    var groups []map[string]interface{}
    for rows.Next() {
        var id, facultyID, courseYear int
        var groupName string
        err := rows.Scan(&id, &groupName, &facultyID, &courseYear)
        if err != nil {
            fmt.Printf("Failed to scan group: %v\n", err)
            continue
        }
        groups = append(groups, map[string]interface{}{
            "id":          id,
            "group_name":  groupName,
            "faculty_id":  facultyID,
            "course_year": courseYear,
        })
    }

    return c.JSON(http.StatusOK, groups)
}

func GetAllSubjectsHandler(c echo.Context) error {
    query := `SELECT id, subject_name, subject_code, credits FROM subjects ORDER BY subject_name`
    rows, err := pool.Query(context.Background(), query)
    if err != nil {
        fmt.Printf("Failed to get subjects: %v\n", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch subjects"})
    }
    defer rows.Close()

    var subjects []map[string]interface{}
    for rows.Next() {
        var id, credits int
        var subjectName, subjectCode string
        err := rows.Scan(&id, &subjectName, &subjectCode, &credits)
        if err != nil {
            fmt.Printf("Failed to scan subject: %v\n", err)
            continue
        }
        subjects = append(subjects, map[string]interface{}{
            "id":           id,
            "subject_name": subjectName,
            "subject_code": subjectCode,
            "credits":      credits,
        })
    }

    return c.JSON(http.StatusOK, subjects)
}

func GetAllScheduleHandler(c echo.Context) error{
    schedules, err := getAllSchedules(pool)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }
    return c.JSON(http.StatusOK, schedules)
}

func getAllSchedules(pool *pgxpool.Pool) ([]Schedule, error) {
    query := `
        SELECT s.id, s.subject_name, s.time_slot, s.group_id, sg.group_name
        FROM schedule s
        LEFT JOIN student_groups sg ON s.group_id = sg.id
        ORDER BY s.id
    `
    rows, err := pool.Query(context.Background(), query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var schedules []Schedule
    for rows.Next() {
        var schedule Schedule
        err := rows.Scan(&schedule.ID, &schedule.SubjectName, &schedule.TimeSlot, &schedule.GroupID, &schedule.GroupName)
        if err != nil {
            return nil, err
        }
        schedules = append(schedules, schedule)
    }

    return schedules, nil
}

func GetScheduleByGroupHandler(c echo.Context) error{
    idStr := c.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid group ID"})
    }

    schedules, err := getScheduleByGroupId(pool, id)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }
    return c.JSON(http.StatusOK, schedules)
}

func getScheduleByGroupId(pool *pgxpool.Pool, groupId int) ([]Schedule, error) {
    query := `
        SELECT s.id, s.subject_name, s.time_slot, s.group_id, sg.group_name
        FROM schedule s
        LEFT JOIN student_groups sg ON s.group_id = sg.id
        WHERE s.group_id = $1
    `
    rows, err := pool.Query(context.Background(), query, groupId)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var schedules []Schedule
    for rows.Next() {
        var schedule Schedule
        err := rows.Scan(&schedule.ID, &schedule.SubjectName, &schedule.TimeSlot, &schedule.GroupID, &schedule.GroupName)
        if err != nil {
            return nil, err
        }
        schedules = append(schedules, schedule)
    }

    return schedules, nil
}

func PostAttendanceHandler(c echo.Context) error{
    var attendance Attendance
    if err := c.Bind(&attendance); err!=nil{
        fmt.Printf("Bind error: %v\n", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
    }
    
    fmt.Printf("Received attendance: StudentId=%d, SubjectID=%d, VisitDay=%s, Visited=%v\n", 
        attendance.StudentId, attendance.SubjectID, attendance.VisitDay, attendance.Visited)
    
    if !(attendance.StudentId > 0) || !(attendance.SubjectID > 0) || (attendance.VisitDay == "") {
        fmt.Printf("Validation failed: StudentId=%d, SubjectID=%d, VisitDay=%s\n", 
            attendance.StudentId, attendance.SubjectID, attendance.VisitDay)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request: missing required fields"})
    }

    createdAttendance, err := createAttendance(pool, &attendance)
     if err != nil {
        fmt.Printf("Database error creating attendance: %v\n", err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error", "details": err.Error()})
    }

    return c.JSON(http.StatusCreated, createdAttendance)
}


func createAttendance(pool *pgxpool.Pool, attendance *Attendance) (*Attendance, error) {
    query := `
        INSERT INTO attendance (subject_id, visit_day, visited, student_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id
    `
    
    err := pool.QueryRow(
        context.Background(),
        query,
        attendance.SubjectID,
        attendance.VisitDay,
        attendance.Visited,
        attendance.StudentId,
    ).Scan(&attendance.ID)
    
    if err != nil {
        return nil, err
    }

    return attendance, nil
}

func GetAttendanceByStudentIdHandler(c echo.Context) error {
    studentId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        fmt.Printf("Invalid student ID: %v\n", err)
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
    }

    userRole := c.Get("user_role").(string)
    userID := c.Get("user_id").(int)
    
    if userRole == "student" {
        var studentIDFromUser *int
        err := pool.QueryRow(context.Background(), "SELECT student_id FROM users WHERE id = $1", userID).Scan(&studentIDFromUser)
        if err != nil {
            fmt.Printf("Error getting student_id from users: %v\n", err)
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied: unable to verify student"})
        }
        if studentIDFromUser == nil || *studentIDFromUser != studentId {
            return c.JSON(http.StatusForbidden, map[string]string{"error": "Access denied: you can only view your own attendance"})
        }
    }

    attendances, err := getAttendanceByStudentId(pool, studentId)
    if err != nil {
        fmt.Printf("Error getting attendance by student ID %d: %v\n", studentId, err)
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }

    return c.JSON(http.StatusOK, attendances)
}

func getAttendanceByStudentId(pool *pgxpool.Pool, id int) ([]Attendance, error) {
    query := `
        SELECT a.id, a.subject_id, a.visit_day::text, a.visited, a.student_id
        FROM attendance a
        WHERE a.student_id = $1
        ORDER BY a.visit_day DESC
        LIMIT 50
    `
    
    rows, err := pool.Query(context.Background(), query, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var attendances []Attendance
    for rows.Next() {
        var attendance Attendance
        err := rows.Scan(&attendance.ID, &attendance.SubjectID, &attendance.VisitDay, &attendance.Visited, &attendance.StudentId)
        if err != nil {
            fmt.Printf("Error scanning attendance row: %v\n", err)
            return nil, err
        }
        attendances = append(attendances, attendance)
    }

    return attendances, nil
}

func GetAttendanceBySubjectIdHandler(c echo.Context) error {
    subjectId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid subject ID"})
    }

    attendances, err := getAttendanceBySubjectId(pool, subjectId)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }

    return c.JSON(http.StatusOK, attendances)
}

func getAttendanceBySubjectId(pool *pgxpool.Pool, id int) ([]Attendance, error) {
    query := `
        SELECT a.id, a.subject_id, a.visit_day::text, a.visited, a.student_id
        FROM attendance a
        WHERE a.subject_id = $1
        ORDER BY a.visit_day DESC
        LIMIT 50
    `

    rows, err := pool.Query(context.Background(), query, id)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var attendances []Attendance
    for rows.Next() {
        var attendance Attendance
        err := rows.Scan(&attendance.ID, &attendance.SubjectID, &attendance.VisitDay, &attendance.Visited, &attendance.StudentId)
        if err != nil {
            fmt.Printf("Error scanning attendance row: %v\n", err)
            return nil, err
        }
        attendances = append(attendances, attendance)
    }

    return attendances, nil
}

