package routes

import (
	"fmt"
	"log"

	d "github.com/georgejmx/whisper-blog/controller"
	tp "github.com/georgejmx/whisper-blog/types"
	u "github.com/georgejmx/whisper-blog/utils"
	"go.uber.org/ratelimit"

	"github.com/gin-gonic/gin"
)

var (
	Rl  ratelimit.Limiter
	dbo tp.ControllerTemplate
)

/* Establishes database connection and controller object, else panics */
func SetupDatabase() {
	dbo = &d.DbController{}
	if err := dbo.Init(); err != nil {
		log.Fatal("unable to initialise database")
	}
}

/* Gets chain from backend, returning it as a type. This means output can be
parsed both as JSON and HTML */
func getChain(c *gin.Context) (int, []tp.Post) {
	attachHeaders(c)

	// Selecting posts data
	posts, err := dbo.SelectPosts()
	if err != nil {
		sendFailure(c, "selecting posts database operation failed")
		return -1, []tp.Post{}
	}

	// Attaching top reactions to each post, in a modified slice
	var stampedPosts []tp.Post
	for _, val := range posts {
		postReactions, err := dbo.SelectPostReactions(val.Id)
		if err != nil {
			sendFailure(c, fmt.Sprintf("error getting reactions of %v", val.Id))
			return -1, []tp.Post{}
		}
		val.Reactions = postReactions
		stampedPosts = append(stampedPosts, val)
	}

	// Calculating days since previous post
	var daysSince int
	if len(stampedPosts) == 0 {
		daysSince = 0
		stampedPosts = []tp.Post{}
	} else {
		daysSince = u.TimeSincePost(
			true, stampedPosts[len(stampedPosts)-1].Time)
	}

	return daysSince, stampedPosts
}

/* Allowing test database to be cleared by integration tests */
func Clear() bool { return dbo.Clear() }

/* Attaches CORS headers to the current context */
func attachHeaders(c *gin.Context) *gin.Context {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Headers",
		`Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, 
		Authorization, accept, origin, Cache-Control, X-Requested-With`)
	c.Header("Access-Control-Allow-Methods", "GET,POST,HEAD,OPTIONS")
	return c
}

/* Sends a HTTP failure response */
func sendFailure(context *gin.Context, msg string) {
	context.JSON(400, gin.H{
		"message": msg,
		"marker":  0,
	})
}
