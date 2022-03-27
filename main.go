package main

import (
	"net/http"
	"os"
	r "whisper-blog/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Setting up database connection and routes
	r.SetupDatabase()
	router := gin.Default()
	router.GET("/data", r.GetChain)
	router.POST("/data/post", r.AddPost)

	// Serving frontend at root path, then running
	router.GET("/", gin.WrapH(http.FileServer(http.FS(os.DirFS("client/dist")))))
	router.Run("localhost:8080")
}
