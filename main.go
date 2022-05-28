/* Copyright 2022 George Miller <georgejmx@pm.me>.
Usage of this code is subject to a GNU public license as detailed in the
LICENSE file. *This code* is defined to be everything statically linked
to this file by golang */

package main

import (
	"embed"
	"io/fs"
	"net/http"

	config "github.com/georgejmx/whisper-blog/config"
	r "github.com/georgejmx/whisper-blog/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

//go:embed client/public/*
var client embed.FS

/* Program entry point when used in production */
func main() {
	rl := ratelimit.New(150)
	setup(true, rl).RunTLS(":8000", "private.crt", "public.key")
}

/* Read configuration and setup production or test server */
func setup(isProduction bool, rl ratelimit.Limiter) *gin.Engine {
	// Setting config
	config.SetupEnv(isProduction)

	// Setting up database connection, rate limiting, router and cors
	r.SetupDatabase()
	r.Rl = rl
	router := gin.Default()
	router.Use(cors.Default())

	// Defining routes
	router.GET("/data/chain", r.GetRawChain)
	router.POST("/data/post", r.AddPost)
	router.POST("/data/react", r.AddReaction)
	router.GET("/html/chain", r.GetHtmlChain)
	router.GET("/html/reaction/:id", r.GetHtmlReactions)

	// Serving client at root directory
	stripped, err := fs.Sub(client, "client/public")
	if err != nil {
		panic("error when bundling client files")
	}
	router.StaticFS("/w", http.FS(stripped))
	router.GET("/", func(c *gin.Context) {
		c.Redirect(301, "/w")
	})

	return router
}
