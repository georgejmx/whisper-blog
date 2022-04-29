package routes

import (
	"encoding/json"
	d "whisper-blog/controller"
	x "whisper-blog/security"
	tp "whisper-blog/types"
	u "whisper-blog/utils"
	w "whisper-blog/words"

	"github.com/gin-gonic/gin"
)

const MAX_ANON_REACTIONS int = 6

var dbo tp.ControllerTemplate

/* Establishes database connection and controller object, else panics */
func SetupDatabase() {
	dbo = &d.DbController{}
	if err := dbo.Init(); err != nil {
		panic("unable to initialise database")
	}
}

/* Gets the chain stored in backend. This inlcudes all posts and the top 3
reactions for each post */
func GetChain(c *gin.Context) {
	posts, err := dbo.SelectPosts()
	if err != nil {
		sendFailure(c, "database operation failed")
		return
	}
	c.JSON(200, posts)
}

/* Adds a Post contained in the request body to database, subject to
validation */
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
	post.Descriptors, err = w.GenerateDescriptors()
	if err != nil {
		sendFailure(c, "unable to generate descriptors for post")
		return
	} else if err = dbo.InsertPost(post); err != nil {
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

/* Adds a Reaction contained in the request body to databse, subject to
validation Input reaction should be of the format:
{postId, descriptor, gravitasHash}  */
func AddReaction(c *gin.Context) {
	// Parsing request body
	var reaction tp.Reaction
	body, err := c.GetRawData()
	err2 := json.Unmarshal(body, &reaction)
	if err != nil || err2 != nil {
		sendFailure(c, "invalid request body")
		return
	}

	// Checking that we have a correct descriptor and gravitas hash
	descriptors, err := dbo.SelectDescriptors(reaction.PostId)
	if err != nil {
		sendFailure(c, "db error when selecting descriptors")
		return
	} else if !u.CheckDescriptor(reaction.Descriptor, descriptors) {
		sendFailure(c, "invalid reaction descriptor provided")
		return
	}

	// Determining the gravitas of reaction and its validity, handling errors.
	// Also setting the correct gravitas value
	isValidHash, gravitas, err := x.ValidateReactionHash(
		dbo, reaction.GravitasHash, reaction.PostId)
	if err != nil {
		sendFailure(c, err.Error())
		return
	}
	reaction.Gravitas = gravitas

	// If valid hash provided proceed
	if !isValidHash {
		// Proceed to adding an anonymous hash if following conditions skip
		count, err := dbo.SelectAnonReactionCount(reaction.PostId)
		if err != nil {
			sendFailure(c, "error selecting number of anonymous reactions")
			return
		} else if count >= MAX_ANON_REACTIONS {
			sendFailure(c, "no more anonymous reactions can be made")
			return
		}

	}

	if err := dbo.InsertReaction(reaction); err != nil {
		sendFailure(c, "error when performing db insert")
		return
	}

	// Sending success response
	c.JSON(201, gin.H{
		"message": "reaction successful",
		"reply":   w.GenerateDescriptor(1349),
		"marker":  1,
	})
}

// Sends a HTTP failure response
func sendFailure(context *gin.Context, msg string) {
	context.JSON(400, gin.H{
		"message": msg,
		"marker":  0,
	})
}
