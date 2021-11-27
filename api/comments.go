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

func (api *Api) routePostComment(c *gin.Context) {
	var commentCreation model.CommentCreation
	err := c.BindJSON(&commentCreation)
	if err != nil {
		zap.L().Sugar().Info("routePostComment ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
		})
		return
	}

	var result struct {
		Id int64 `json:"id"`
	}

	result.Id, err = api.logic.CreateComment(&commentCreation)
	if err != nil {
		if valErr, ok := err.(logic.ValidationError); ok {
			handleValidationError(c, valErr)
			return
		}
		zap.L().Sugar().Error("routePostComment", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	resultJson, _ := json.Marshal(result)
	c.String(http.StatusOK, string(resultJson))
}

func (api *Api) routeGetCommentsNbf(c *gin.Context) {
	api.routeGetComments(orderNbf, c)
}

func (api *Api) routeGetCommentsNsf(c *gin.Context) {
	api.routeGetComments(orderNsf, c)
}

func (api *Api) routeGetCommentsOsf(c *gin.Context) {
	api.routeGetComments(orderOsf, c)
}

func (api *Api) routeGetComments(order orderType, c *gin.Context) {
	// Required thread parameter
	threadPath, err := getStringQueryParameter(c, "threadPath")
	if err != nil {
		zap.L().Sugar().Info("routeGetComments ", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
		})
		return
	}

	// Optional query parameters
	depth := getIntQueryParameter(c, "depth", 0)
	max := getIntQueryParameter(c, "max", 0)
	preview := getIntQueryParameter(c, "preview", 0)
	parent := getInt64QueryParameter(c, "parent", 0)

	// Query commentResult tree
	var commentResult logic.CommentResult
	switch order {
	case orderNbf:
		commentResult, err = api.logic.GetCommentsNbf(threadPath, parent, depth, max)
	case orderNsf:
		commentResult, err = api.logic.GetCommentsNsf(threadPath, parent, depth, max)
	case orderOsf:
		commentResult, err = api.logic.GetCommentsOsf(threadPath, parent, depth, max)
	}
	if err != nil {
		zap.L().Sugar().Error("routeGetComments", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	validation := api.logic.GetValidation()
	administration := api.logic.GetAdministration()

	commentResponse := ThreadCommentsResponse{
		Config: ServerConfig{
			RequireEmail:    validation.RequireEmail,
			RequireAuthor:   validation.RequireAuthor,
			RequireApproval: administration.RequireApproval,
		},
		Thread: commentResult,
	}

	// Return JSON or generate preview
	if preview == 0 {
		commentsJson, _ := json.Marshal(commentResponse)
		c.String(http.StatusOK, string(commentsJson))
	} else {
		preview := api.logic.GenerateHTMLThreadPreview(commentResult)
		c.Header("Content-Type", "text/html")
		c.String(http.StatusOK, preview)
	}
}

func (api *Api) routeGetMoreCommentsNbf(c *gin.Context) {
	api.routeGetMoreComments(orderNbf, c)
}

func (api *Api) routeGetMoreCommentsNsf(c *gin.Context) {
	api.routeGetMoreComments(orderNsf, c)
}

func (api *Api) routeGetMoreCommentsOsf(c *gin.Context) {
	api.routeGetMoreComments(orderOsf, c)
}

func (api *Api) routeGetMoreComments(order orderType, c *gin.Context) {
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
		comments, err = api.logic.GetMoreCommentsNbf(threadId, parentId, newestCreatedAt, excludeIds, limit)
	case orderNsf:
		comments, err = api.logic.GetMoreCommentsNsf(threadId, parentId, newestCreatedAt, excludeIds, limit)
	case orderOsf:
		comments, err = api.logic.GetMoreCommentsOsf(threadId, parentId, newestCreatedAt, excludeIds, limit)
	}
	if err != nil {
		zap.L().Sugar().Error("routeGetMoreComments", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "error",
		})
		return
	}

	commentsJson, _ := json.Marshal(comments)
	c.String(http.StatusOK, string(commentsJson))
}
