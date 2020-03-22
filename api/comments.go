package api

import (
    "github.com/gin-gonic/gin"
    "github.com/matthinc/gomment/model"
    "github.com/matthinc/gomment/logic"
    "net/http"
    "fmt"
    "encoding/json"
)

func routePostComment(c *gin.Context, logic *logic.BusinessLogic) {
    var comment model.Comment
    err := c.BindJSON(&comment)
    if err != nil {
        fmt.Println(err)
    }
    logic.AddComment(&comment)
    c.String(http.StatusOK, "ok")
}

func routeGetComments(c *gin.Context, logic *logic.BusinessLogic) {
    comments := logic.GetCommentsTree(0)
    commentsJson, _ := json.Marshal(comments)
    c.String(http.StatusOK, string(commentsJson))
}

func routePreviewComments(c *gin.Context, logic *logic.BusinessLogic) {
    preview := logic.GenerateHTMLThreadPreview(0)
    c.Header("Content-Type", "text/html")
    c.String(http.StatusOK, preview)
}

