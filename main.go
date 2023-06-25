package main

import (
	"github.com/gin-gonic/gin"
	"github.com/yudhifadilah/back/controllers/notecontroller"
	"github.com/yudhifadilah/back/models"
)

func main() {
	r := gin.Default()
	models.ConnectDatabase()

	r.GET("/api/notes", notecontroller.Index)
	r.GET("/api/notes/:id", notecontroller.Show)
	r.POST("/api/notes", notecontroller.Create)
	r.PUT("/api/notes/:id", notecontroller.Update)
	r.DELETE("/api/notes", notecontroller.Delete)

	r.Run()
}
