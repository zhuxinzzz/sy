package main

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	//"github.com/lib/pq"
	"log"
)

var db *sql.DB

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Comment  string `json:"comment"`
}

func init() {
	// 数据库连接代码保留但注释
	/*
	   var err error
	   connStr := "user=yourusername dbname=yourdbname sslmode=disable"
	   db, err = sql.Open("postgres", connStr)
	   if err != nil {
	       log.Fatal(err)
	   }
	*/
}

func getUsers(c *gin.Context) {
	// 模拟数据
	mockUsers := []User{
		{
			Username: "john_doe",
			Email:    "john@example.com",
			Comment:  "Hello World test",
		},
		{
			Username: "jane_doe",
			Email:    "jane@example.com",
			Comment:  "Welcome",
		},
	}

	c.JSON(200, mockUsers)

	// 数据库查询代码保留但注释
	/*
	   rows, err := db.Query("SELECT username, email, comment FROM users")
	   if err != nil {
	       c.JSON(500, gin.H{"error": err.Error()})
	       return
	   }
	   defer rows.Close()

	   var users []User
	   for rows.Next() {
	       var user User
	       if err := rows.Scan(&user.Username, &user.Email, &user.Comment); err != nil {
	           c.JSON(500, gin.H{"error": err.Error()})
	           return
	       }
	       users = append(users, user)
	   }
	   if err := rows.Err(); err != nil {
	       c.JSON(500, gin.H{"error": err.Error()})
	       return
	   }
	   c.JSON(200, users)
	*/
}

func main() {
	router := gin.Default()

	// 配置 CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200"} // Angular 开发服务器地址
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}

	router.Use(cors.New(config))

	// 路由设置
	//router.GET("/api/users", getUsers)
	router.GET("/users", getUsers)

	// 启动服务器
	log.Println("Server is running on port 8080")
	router.Run(":8080")
}
