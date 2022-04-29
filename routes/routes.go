package routes

import (
	"encoding/json"
	"fmt"
	"log"

	d "github.com/georgejmx/whisper-blog/controller"
	x "github.com/georgejmx/whisper-blog/security"
	tp "github.com/georgejmx/whisper-blog/types"
	u "github.com/georgejmx/whisper-blog/utils"
	w "github.com/georgejmx/whisper-blog/words"

	"github.com/gin-gonic/gin"
)

const MAX_ANON_REACTIONS int = 6

var dbo tp.ControllerTemplate

/* Establishes database connection and controller object, else panics */
func SetupDatabase() {
	dbo = &d.DbController{}
	if err := dbo.Init(); err != nil {
		log.Fatal("unable to initialise database")
	}
}

/* Gets the chain stored in backend. This inlcudes all posts and the top 3
reactions for each post */
func GetChain(c *gin.Context) {
	// Selecting posts data
	posts, err := dbo.SelectPosts()
	if err != nil {
		sendFailure(c, "selecting posts database operation failed")
		return
	}

	// Attaching top reactions to each post, in a modified slice
	var stampedPosts []tp.Post
	for _, val := range posts {
		postReactions, err := dbo.SelectPostReactions(val.Id)
		if err != nil {
			sendFailure(c, fmt.Sprintf("error getting reactions of %v", val.Id))
		}
		val.Reactions = postReactions
		stampedPosts = append(stampedPosts, val)
	}

	c.JSON(200, stampedPosts)
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
		"data":    w.GenerateDescriptor(1349),
		"marker":  1,
	})
}

/* Allowing test database to be cleared by integration tests */
func Clear() bool { return dbo.Clear() }

/* Sends a HTTP failure response */
func sendFailure(context *gin.Context, msg string) {
	context.JSON(400, gin.H{
		"message": msg,
		"marker":  0,
	})
}
