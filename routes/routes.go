package routes

import (
	"encoding/json"
	d "whisper-blog/controller"
	tp "whisper-blog/types"

	"github.com/gin-gonic/gin"
)

var dbo d.DatabaseObject

// Establishes database connection, else panics
func SetupDatabase() {
	if err := dbo.Init(); err != nil {
		panic("unable to initialise database")
	}
}

func AddPost(c *gin.Context) {
	// Parsing request body
	var post tp.Post
	body, err := c.GetRawData()
	err2 := json.Unmarshal(body, &post)
	if err != nil || err2 != nil {
		sendFailure(c, "invalid request body")
		return
	}

	// Generating post descriptors then attempting db insert
	post.Descriptors = `fun;proverbial;feindish;flat;dummylike;
		serendipitous;granular;bespoke;frosted;zyrgony`
	err3 := dbo.AddPost(post)
	if err3 != nil {
		sendFailure(c, "database operation failed")
		return
	}

	// Sending success response
	c.JSON(201, gin.H{
		"message": "post successful",
		"marker":  1,
	})
}

func GetChain(c *gin.Context) {
	posts, err := dbo.GrabPosts()
	if err != nil {
		sendFailure(c, "database operation failed")
		return
	}
	c.JSON(200, posts)
}

// Sends a HTTP failure response
func sendFailure(context *gin.Context, msg string) {
	context.JSON(400, gin.H{
		"message": msg,
		"marker":  0,
	})
}
