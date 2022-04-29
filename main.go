/* Copyright 2022 George Miller <georgejmx@pm.me>.
Usage of this code is subject to a GNU public license as detailed in the
LICENSE file. */

package main

import (
	"net/http"
	"os"

	config "github.com/georgejmx/whisper-blog/config"
	r "github.com/georgejmx/whisper-blog/routes"

	"github.com/gin-gonic/gin"
)

/* Program entry point when used in production */
func main() {
	setup(true).Run("localhost:8080")
}

/* Read production configuration and setup production server */
func setup(isProduction bool) *gin.Engine {
	// Setting config
	config.SetupEnv(isProduction)

	// Setting up database connection and routes
	r.SetupDatabase()
	router := gin.Default()
	router.GET("/data", r.GetChain)
	router.POST("/data/post", r.AddPost)
	router.POST("/data/reaction", r.AddReaction)

	// Serving frontend at root path, then running
	router.GET("/", gin.WrapH(
		http.FileServer(http.FS(os.DirFS("client/dist")))))
	return router
}
