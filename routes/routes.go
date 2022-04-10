package routes

import (
	"encoding/json"
	d "whisper-blog/controller"
	x "whisper-blog/security"
	tp "whisper-blog/types"
	u "whisper-blog/utils"

	"github.com/gin-gonic/gin"
)

var dbo tp.ControllerTemplate

/* Establishes database connection and controller object, else panics */
func SetupDatabase() {
	dbo = &d.DbController{}
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

	// Performing hash validation
	isValidated, err := x.ValidateHash(dbo, post.Hash)
	if err != nil {
		sendFailure(c, "unable to perform passcode validation")
		return
	} else if !isValidated {
		sendFailure(c, "passcode validation failed")
		return
	}

	// Generating post descriptors then performing db insert of post
	post.Descriptors, err = u.GenerateDescriptors()
	if err != nil {
		sendFailure(c, "unable to generate descriptors for post")
		return
	} else if err = dbo.AddPost(post); err != nil {
		sendFailure(c, "database operation failed")
		return
	}

	// Inserting new passcode and getting cipher
	cipher, err := x.SetHashAndRetrieveCipher(dbo)
	if err != nil {
		sendFailure(c, "error when setting new passcode and/or getting cipher")
		return
	}

	// Sending success response
	c.JSON(201, gin.H{
		"message": "post successful",
		"data":    cipher,
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
