package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yungkhann/echo-server/internal/database"
	custommiddleware "github.com/yungkhann/echo-server/internal/middleware"
)







func main() {
	godotenv.Load()

	db := database.InitDB()
	defer db.Close()
    e := echo.New()

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.POST("/api/auth/register", database.RegisterHandler)
	e.POST("/api/auth/login", database.LoginHandler)
	e.GET("/api/users/me", database.GetMeHandler, custommiddleware.AuthMiddleware)
	e.GET("/student/:id", database.GetStudentHandler)
    e.GET("/all_class_schedule", database.GetAllScheduleHandler)
    e.GET("/schedule/group/:id", database.GetScheduleByGroupHandler)
    e.POST("/attendance/subject", database.PostAttendanceHandler)
    e.GET("/attendanceByStudentId/:id", database.GetAttendanceByStudentIdHandler)
    e.GET("/attendanceBySubjectId/:id", database.GetAttendanceBySubjectIdHandler)
	  e.Start(":8080")
}