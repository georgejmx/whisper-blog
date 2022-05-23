package routes

import (
	"encoding/json"

	x "github.com/georgejmx/whisper-blog/security"
	tp "github.com/georgejmx/whisper-blog/types"
	u "github.com/georgejmx/whisper-blog/utils"
	w "github.com/georgejmx/whisper-blog/words"
	"github.com/gin-gonic/gin"
)

const MAX_ANON_REACTIONS int = 6

/* Gets the chain stored in backend as JSON. This inlcudes all posts and the
top 3 reactions for each post */
func GetRawChain(c *gin.Context) {
	// Sending success json response with chain data
	daysSince, stampedPosts := getChain(c)
	if daysSince != -1 {
		c.JSON(200, gin.H{
			"marker":     1,
			"days_since": daysSince,
			"chain":      stampedPosts,
		})
	}
}

/* Adds a Post contained in the request body to database, subject to
validation */
func AddPost(c *gin.Context) {
	var (
		post   tp.Post
		marker int
	)

	// Parsing request body
	body, err := c.GetRawData()
	err2 := json.Unmarshal(body, &post)
	if err != nil || err2 != nil || len(post.Author) > 10 {
		sendFailure(c, "invalid request body")
		return
	}

	// If this is the genesis post, can skip hash validation
	isGenesis, err := checkForGenesis()
	if err != nil {
		sendFailure(c, "error when determing if genesis post")
		return
	}

	// Need to perform hash validation if not genesis post
	if isGenesis {
		marker = 2
		post.Tag = 0
	} else {
		marker = 1
		isValidated, err := x.ValidateHash(dbo, post.Hash)
		if err != nil || post.Tag == 0 {
			sendFailure(c, "unable to perform passcode validation")
			return
		} else if !isValidated {
			sendFailure(c, "passcode validation failed")
			return
		}
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
	cipher, err := x.SetHashAndRetrieveCipher(dbo, isGenesis, post.Hash)
	if err != nil {
		sendFailure(c, "error when setting new passcode and/or getting cipher")
		return
	}

	// Sending success response
	c.JSON(201, gin.H{
		"message": "post successful",
		"data":    cipher,
		"marker":  marker,
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

	// If valid hash provided proceed without calling this block
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

	// Finally adding reaction with the correct gravitas
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

/* Determining if this is the genesis post */
func checkForGenesis() (bool, error) {
	// Selecting the existing chain
	posts, err := dbo.SelectPosts()
	if len(posts) == 0 {
		return true, err
	}
	return false, err
}
