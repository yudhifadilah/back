package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/yudhifadilah/back/controllers/notecontroller"
	"github.com/yudhifadilah/back/models"

	//"github.com/gin-contrib/recovery"

	jwt_middleware "github.com/yudhifadilah/back/middleware"
)

func main() {
	r := gin.Default()
	models.ConnectDatabase()
	//r.Use(cors.Default())

	// Middleware Logger
	r.Use(logger.SetLogger())
	r.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"*"},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"GET,POST,PUT,DELETE"},
	}))
	// Middleware Recovery
	//r.Use(recovery.Default())
	//login
	r.POST("/api/user/register", jwt_middleware.RegisterUser)
	r.POST("/api/user/login", jwt_middleware.LoginUser)

	// Menggunakan middleware untuk otentikasi JWT
	r.Use(jwt_middleware.Authenticate())
	r.GET("/api/user/getme", jwt_middleware.GetMe)

	r.POST("/api/user/logout", jwt_middleware.LogoutUser)

	r.GET("/api/notes", notecontroller.Index)
	r.GET("/api/notes/:id", notecontroller.Show)
	r.POST("/api/notes", notecontroller.Create)
	r.PUT("/api/notes/:id", notecontroller.Update)
	r.DELETE("/api/dell/:id", notecontroller.Delete)
	r.Run()
}
