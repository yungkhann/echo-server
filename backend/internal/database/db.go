package database

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

type Student struct{
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
        SELECT s.full_name, s.gender, s.birth_date::text, s.group_id, sg.group_name 
        FROM students s
        LEFT JOIN student_groups sg ON s.group_id = sg.id
        WHERE s.id = $1
    `
    row := pool.QueryRow(context.Background(), query, id)

    err := row.Scan(&student.FullName, &student.Gender, &student.BirthDate, &student.GroupID, &student.GroupName)
    if err != nil {
        return nil, err
    }

    return student, nil
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
        return c.JSON(http.StatusNotFound, map[string]string{"error": "Attendace not found"})
    }
    
    if !(attendance.StudentId > 0) || !(attendance.SubjectID > 0) || (attendance.VisitDay == "") {
                return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bad request"})

    }

    createdAttendance, err := createAttendance(pool, &attendance)
     if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
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



// func getStudentById(pool *pgxpool.Pool, id int) (*Student, error) {
//     student := &Student{}

//     query := `
//         SELECT s.full_name, s.gender, s.birth_date::text, s.group_id, sg.group_name 
//         FROM students s
//         LEFT JOIN student_groups sg ON s.group_id = sg.id
//         WHERE s.id = $1
//     `
//     row := pool.QueryRow(context.Background(), query, id)

//     err := row.Scan(&student.FullName, &student.Gender, &student.BirthDate, &student.GroupID, &student.GroupName)
//     if err != nil {
//         return nil, err
//     }

//     return student, nil
// }


func GetAttendanceByStudentIdHandler(c echo.Context) error {
    studentId, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid student ID"})
    }

    attendances, err := getAttendanceByStudentId(pool, studentId)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }

    return c.JSON(http.StatusOK, attendances)
}

func getAttendanceByStudentId(pool *pgxpool.Pool, id int) ([]Attendance, error) {
    query := `
        SELECT a.id, a.subject_id, a.visit_day, a.visited, a.student_id
        FROM attendance a
        WHERE a.student_id = $1
        ORDER BY a.visit_day DESC
        LIMIT 5
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
        SELECT a.id, a.subject_id, a.visit_day, a.visited, a.student_id
        FROM attendance a
        WHERE a.subject_id = $1
        ORDER BY a.visit_day DESC
        LIMIT 5
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
            return nil, err
        }
        attendances = append(attendances, attendance)
    }

    return attendances, nil
}

