package api

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matthinc/gomment/logic"
	"github.com/matthinc/gomment/model"
	"go.uber.org/zap"
)

type orderType int

const (
	orderNbf orderType = iota
	orderNsf orderType = iota
	orderOsf orderType = iota
)

func routePostComment(c *gin.Context, logic *logic.BusinessLogic) {
	var commentCreation model.CommentCreation
	err := c.BindJSON(&commentCreation)
	if err != nil {
		zap.L().Sugar().Infow("routePostComment", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
		})
		return
	}

	var result struct {
		Id int64 `json:"id"`
	}

	result.Id, err = logic.CreateComment(&commentCreation)
	if err != nil {
		zap.L().Sugar().Errorw("routePostComment", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	resultJson, _ := json.Marshal(result)
	c.String(http.StatusOK, string(resultJson))
}

func routeGetCommentsNbf(c *gin.Context, logic *logic.BusinessLogic) {
	routeGetComments(orderNbf, c, logic)
}

func routeGetCommentsNsf(c *gin.Context, logic *logic.BusinessLogic) {
	routeGetComments(orderNsf, c, logic)
}

func routeGetCommentsOsf(c *gin.Context, logic *logic.BusinessLogic) {
	routeGetComments(orderOsf, c, logic)
}

func routeGetComments(order orderType, c *gin.Context, logic *logic.BusinessLogic) {
	// Required thread parameter
	threadPath, err := getStringQueryParameter(c, "threadPath")
	if err != nil {
		zap.L().Sugar().Infow("routeGetComments", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	// Optional query parameters
	depth := getIntQueryParameter(c, "depth", 0)
	max := getIntQueryParameter(c, "max", 0)
	preview := getIntQueryParameter(c, "preview", 0)
	parent := getInt64QueryParameter(c, "parent", 0)

	// Query comments tree
	var comments model.CommentsResponse
	switch order {
	case orderNbf:
		comments, err = logic.GetCommentsNbf(threadPath, parent, depth, max)
	case orderNsf:
		comments, err = logic.GetCommentsNsf(threadPath, parent, depth, max)
	case orderOsf:
		comments, err = logic.GetCommentsOsf(threadPath, parent, depth, max)
	}
	if err != nil {
		zap.L().Sugar().Errorw("routeGetComments", err)
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

func routeGetMoreCommentsNbf(c *gin.Context, logic *logic.BusinessLogic) {
	routeGetMoreComments(orderNbf, c, logic)
}

func routeGetMoreCommentsNsf(c *gin.Context, logic *logic.BusinessLogic) {
	routeGetMoreComments(orderNsf, c, logic)
}

func routeGetMoreCommentsOsf(c *gin.Context, logic *logic.BusinessLogic) {
	routeGetMoreComments(orderOsf, c, logic)
}

func routeGetMoreComments(order orderType, c *gin.Context, logic *logic.BusinessLogic) {
	// Optional query parameters
	threadId := getInt64QueryParameter(c, "threadId", 0)
	parentId := getInt64QueryParameter(c, "parentId", 0)
	newestCreatedAt := getInt64QueryParameter(c, "newestCreatedAt", 0)
	excludeIds := getInt64ListQueryParameter(c, "excludeIds")
	limit := getIntQueryParameter(c, "limit", 0)

	// Query comments tree
	var comments []model.Comment
	var err error
	switch order {
	case orderNbf:
		comments, err = logic.GetMoreCommentsNbf(threadId, parentId, newestCreatedAt, excludeIds, limit)
	case orderNsf:
		comments, err = logic.GetMoreCommentsNsf(threadId, parentId, newestCreatedAt, excludeIds, limit)
	case orderOsf:
		comments, err = logic.GetMoreCommentsOsf(threadId, parentId, newestCreatedAt, excludeIds, limit)
	}
	if err != nil {
		zap.L().Sugar().Errorw("routeGetMoreComments", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	commentsJson, _ := json.Marshal(comments)
	c.String(http.StatusOK, string(commentsJson))
}
