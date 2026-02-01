package database

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	FullName  string    `json:"full_name,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	FullName string `json:"full_name"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type JWTClaims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

func init() {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "yungkhann-keep-ur-secret-bruh"
	}
	jwtSecret = []byte(secret)
}

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateJWT(userID int, email string, role string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func RegisterHandler(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if !ValidateEmail(req.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email format"})
	}

	if len(req.Password) < 6 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be 6 words length"})
	}

	if req.Role == "" {
		req.Role = "student"
	}
	if req.Role != "student" && req.Role != "teacher" && req.Role != "admin" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid role. Must be student, teacher, or admin"})
	}

	var existingUser User
	err := pool.QueryRow(context.Background(), "SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingUser.ID)
	if err == nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": "This email is taken"})
	}

	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
	}

	var userID int
	var role string
	query := `INSERT INTO users (email, password, role, full_name, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id, role`
	err = pool.QueryRow(context.Background(), query, req.Email, hashedPassword, req.Role, req.FullName, time.Now()).Scan(&userID, &role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	token, err := GenerateJWT(userID, req.Email, role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusCreated, LoginResponse{
		Token: token,
		User: User{
			ID:        userID,
			Email:     req.Email,
			Role:      role,
			FullName:  req.FullName,
			CreatedAt: time.Now(),
		},
	})
}

func LoginHandler(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	var user User
	query := `SELECT id, email, password, role, full_name, created_at FROM users WHERE email = $1`
	err := pool.QueryRow(context.Background(), query, req.Email).Scan(&user.ID, &user.Email, &user.Password, &user.Role, &user.FullName, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if !CheckPasswordHash(req.Password, user.Password) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
	}

	token, err := GenerateJWT(user.ID, user.Email, user.Role)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: User{
			ID:        user.ID,
			Email:     user.Email,
			Role:      user.Role,
			FullName:  user.FullName,
			CreatedAt: user.CreatedAt,
		},
	})
}

func GetMeHandler(c echo.Context) error {
	userID := c.Get("user_id").(int)

	var user User
	query := `SELECT id, email, COALESCE(role, 'student') as role, COALESCE(full_name, '') as full_name, created_at FROM users WHERE id = $1`
	err := pool.QueryRow(context.Background(), query, userID).Scan(&user.ID, &user.Email, &user.Role, &user.FullName, &user.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		fmt.Printf("GetMeHandler error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusOK, user)
}

func GetAllUsersHandler(c echo.Context) error {
	fmt.Println("GetAllUsersHandler called")
	query := `SELECT id, email, COALESCE(role, 'student') as role, COALESCE(full_name, '') as full_name, created_at FROM users ORDER BY created_at DESC`
	rows, err := pool.Query(context.Background(), query)
	if err != nil {
		fmt.Printf("Database query error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Email, &user.Role, &user.FullName, &user.CreatedAt)
		if err != nil {
			fmt.Printf("Row scan error: %v\n", err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		fmt.Printf("Rows iteration error: %v\n", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	return c.JSON(http.StatusOK, users)
}
