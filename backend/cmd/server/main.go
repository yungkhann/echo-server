package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/yungkhann/echo-server/internal/database"
	custommiddleware "github.com/yungkhann/echo-server/internal/middleware"
)

func main() {
	godotenv.Load()

	db := database.InitDB()
	defer db.Close()
    e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.PATCH},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	e.POST("/api/auth/register", database.RegisterHandler)
	e.POST("/api/auth/login", database.LoginHandler)
	e.GET("/api/users/me", database.GetMeHandler, custommiddleware.AuthMiddleware)
	e.GET("/api/users", database.GetAllUsersHandler, custommiddleware.AuthMiddleware, custommiddleware.RequireRole("admin"))
	e.GET("/students", database.GetAllStudentsHandler, custommiddleware.AuthMiddleware, custommiddleware.RequireRole("teacher", "admin"))
	e.POST("/students", database.CreateStudentHandler, custommiddleware.AuthMiddleware, custommiddleware.RequireRole("admin"))
	e.POST("/students/from-user", database.CreateStudentFromUserHandler, custommiddleware.AuthMiddleware, custommiddleware.RequireRole("admin"))
	e.GET("/student/:id", database.GetStudentHandler, custommiddleware.AuthMiddleware)
	e.GET("/groups", database.GetAllGroupsHandler, custommiddleware.AuthMiddleware)
	e.GET("/subjects", database.GetAllSubjectsHandler, custommiddleware.AuthMiddleware)
    e.GET("/all_class_schedule", database.GetAllScheduleHandler, custommiddleware.AuthMiddleware)
    e.GET("/schedule/group/:id", database.GetScheduleByGroupHandler, custommiddleware.AuthMiddleware)
    e.POST("/attendance/subject", database.PostAttendanceHandler, custommiddleware.AuthMiddleware, custommiddleware.RequireRole("teacher", "admin"))
    e.GET("/attendanceByStudentId/:id", database.GetAttendanceByStudentIdHandler, custommiddleware.AuthMiddleware)
    e.GET("/attendanceBySubjectId/:id", database.GetAttendanceBySubjectIdHandler, custommiddleware.AuthMiddleware, custommiddleware.RequireRole("teacher", "admin"))
	

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	
	e.Start(":8080")
}