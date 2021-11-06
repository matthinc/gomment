package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/matthinc/gomment/logic"
	"github.com/matthinc/gomment/model"
	"net/http"
)

func routePostComment(c *gin.Context, logic *logic.BusinessLogic) {
	var commentCreation model.CommentCreation
	err := c.BindJSON(&commentCreation)
	if err != nil {
		fmt.Println(err)
	}

	var result struct {
		Id int64 `json:"id"`
	}

	result.Id = logic.CreateComment(&commentCreation)

	resultJson, _ := json.Marshal(result)
	c.String(http.StatusOK, string(resultJson))
}

func routeGetComments(c *gin.Context, logic *logic.BusinessLogic) {
	// Required thread parameter
	threadPath, err := getStringQueryParameter(c, "threadPath")
	if err != nil {
		c.String(http.StatusBadRequest, "'threadPath'-parameter is not optional")
		return
	}

	// Optional query parameters
	depth := getIntQueryParameter(c, "depth", 0)
	max := getIntQueryParameter(c, "max", 0)
	offset := getIntQueryParameter(c, "offset", 0)
	preview := getIntQueryParameter(c, "preview", 0)
	parent := getIntQueryParameter(c, "parent", 0)

	// Query comments tree
	comments, err := logic.GetCommentsTree(threadPath, parent, depth, max, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	// Return JSON or generate preview
	if preview == 0 {
		commentsJson, _ := json.Marshal(comments)
		c.String(http.StatusOK, string(commentsJson))
	} else {
		preview := logic.GenerateHTMLThreadPreview(comments)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, preview)
	}
}
