package main

import (
	r "whisper-blog/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r.Setup()
	router := gin.Default()
	router.GET("/init", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "A place to mysteriously connect with friends :-)",
			"version": 0.2,
		})
	})
	router.GET("/data", r.GetChain)
	router.POST("/data/post", r.AddPost)
	router.Run("localhost:8080")
}
