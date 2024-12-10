package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Comment  string `json:"comment"`
	Avatar   string `json:"avatar"`
}

var DB *sql.DB
var ErrorLogger *logrus.Logger

func initDatabase() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		ErrorLogger.Error("Database connection error:", err)
		log.Fatal(err)
	}

	for i := 0; i < 5; i++ {
		if err = DB.Ping(); err == nil {
			break
		}
		time.Sleep(5 * time.Second)
	}
}

func uploadAvatar(c *gin.Context) {
	userID := c.Param("id")
	file, err := c.FormFile("file")
	if err != nil {
		ErrorLogger.Error("File upload error:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	filename := fmt.Sprintf("avatars/%s%s", userID, filepath.Ext(filepath.Base(file.Filename)))

	if err := c.SaveUploadedFile(file, filename); err != nil {
		ErrorLogger.Error("File save error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	_, err = DB.Exec("UPDATE users SET avatar = $1 WHERE id = $2", filename, userID)
	if err != nil {
		ErrorLogger.Error("Database update error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update database"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Avatar uploaded successfully!", "avatar_url": filename})
}

func addUser(c *gin.Context) {
	var avatarPath string
	file, _ := c.FormFile("avatar")
	if file != nil {
		avatarPath = fmt.Sprintf("uploads/%s", filepath.Base(file.Filename))
		_ = c.SaveUploadedFile(file, avatarPath)
	}

	var lastID int
	err := DB.QueryRow(
		"INSERT INTO users(username, email, comment, avatar) VALUES ($1, $2, $3, $4) RETURNING id",
		c.PostForm("username"), c.PostForm("email"), c.PostForm("comment"), avatarPath,
	).Scan(&lastID)

	if err != nil {
		ErrorLogger.Error("Insert error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added", "id": lastID, "avatar": avatarPath})
}

func getUsers(c *gin.Context) {
	rows, err := DB.Query("SELECT id, username, email, comment, avatar FROM users")
	if err != nil {
		ErrorLogger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		var avatar sql.NullString // 使用 sql.NullString 处理可能为 NULL 的字段
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Comment, &avatar)
		if err != nil {
			ErrorLogger.Error(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 如果 avatar 是 NULL，则设置为空字符串
		if avatar.Valid {
			user.Avatar = avatar.String
		} else {
			user.Avatar = "alex.jpg" // 或者可以设置一个默认头像的路径
		}

		users = append(users, user)
	}

	c.JSON(http.StatusOK, users)
}

func deleteUser(c *gin.Context) {
	username := c.Param("username")
	result, err := DB.Exec("DELETE FROM users WHERE username=$1", username)
	if err != nil {
		ErrorLogger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "User does not exist", "username": username})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

func main() {
	InitLogger()
	initDatabase()

	r := gin.Default()
	r.Static("/uploads", "./uploads")
	r.Use(cors.Default())

	r.POST("/users/:id/avatar", uploadAvatar)
	r.GET("/users", getUsers)
	r.POST("/users", addUser)
	r.DELETE("/users/:username", deleteUser)

	if err := r.Run(":8080"); err != nil {
		ErrorLogger.Error("Server error:", err)
	}
}

func InitLogger() {
	ErrorLogger = logrus.New()
	ErrorLogger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: "2006-01-02 15:04:05.000"})
	ErrorLogger.SetOutput(os.Stdout)
	ErrorLogger.SetLevel(logrus.ErrorLevel)
}
