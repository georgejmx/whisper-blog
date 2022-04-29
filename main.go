/* Copyright 2022 George Miller <georgejmx@pm.me>.
Usage of this code is subject to a GNU license as detailed in the
LICENSE file. */

package main

import (
	"net/http"
	"os"

	r "github.com/georgejmx/whisper-blog/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Setting up database connection and routes
	r.SetupDatabase()
	router := gin.Default()
	router.GET("/data", r.GetChain)
	router.POST("/data/post", r.AddPost)
	router.POST("/data/reaction", r.AddReaction)

	// Serving frontend at root path, then running
	router.GET("/", gin.WrapH(http.FileServer(http.FS(os.DirFS("client/dist")))))
	router.Run("localhost:8080")
}
