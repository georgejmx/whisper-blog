package routes

import (
	"bytes"
	"html/template"
	"strconv"
	"strings"

	tp "github.com/georgejmx/whisper-blog/types"
	u "github.com/georgejmx/whisper-blog/utils"
	"github.com/gin-gonic/gin"
)

/* Gets HTML markup for the frontend chain, dependent on current backup data */
func GetHtmlChain(c *gin.Context) {
	var htmlPosts []tp.PostHtmlContent
	_, stampedPosts := getChain(c)

	// Converting stamped posts to html suitable types
	for _, stamped := range stampedPosts {
		// Need to know this to not append an arrow to bottom of genesis post
		var isSuccessor bool
		if stamped.Tag == 0 {
			isSuccessor = false
		} else {
			isSuccessor = true
		}
		htmlPost := tp.PostHtmlContent{
			Colour:      u.GetTagColour(stamped.Tag),
			Timestring:  u.GetTimestring(stamped.Time),
			IsSuccessor: isSuccessor,
			Id:          stamped.Id,
			Title:       stamped.Title,
			Contents:    stamped.Contents,
			Author:      stamped.Author,
			Reactions:   stamped.Reactions,
		}

		htmlPosts = append(htmlPosts, htmlPost)
	}

	// Getting our template, and its structure
	t, err := template.ParseFiles("templates/chain.gohtml")
	if err != nil {
		sendFailure(c, "error parsing html template")
		return
	}
	htmlStructure := tp.HtmlPostContainer{HtmlPosts: htmlPosts}

	// Executing template, to return byte array. Sending this to client
	var buf bytes.Buffer
	t.Execute(&buf, htmlStructure)
	c.Data(200, "text/html; charset=utf-8", buf.Bytes())
}

/* Gets html reactions that will be passed to frontend */
func GetHtmlReactions(c *gin.Context) {
	attachHeaders(c)

	postId, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		sendFailure(c, "error parsing url parameter")
		return
	}

	descriptorsStr, err := dbo.SelectDescriptors(int(postId))
	if err != nil {
		sendFailure(c, "error getting post descriptors")
		return
	}

	descriptors := strings.Split(descriptorsStr, ";")

	// Getting our template, and its structure
	t, err := template.ParseFiles("templates/descriptors.gohtml")
	if err != nil {
		sendFailure(c, "error parsing html template")
		return
	}
	htmlStructure := tp.HtmlReactionContainer{Descriptors: descriptors}

	// Executing template, to return byte array. Sending this to client
	var buf bytes.Buffer
	t.Execute(&buf, htmlStructure)
	c.Data(200, "text/html; charset=utf-8", buf.Bytes())

}
