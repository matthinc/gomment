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
    // Required thread parameter
    thread := getIntQueryParameter(c, "thread", -1)
    if thread == -1 {
        c.String(http.StatusBadRequest, "'thread'-parameter is not optional")
    }

    // Optional query parameters
    depth := getIntQueryParameter(c, "depth", 0)
    max := getIntQueryParameter(c, "max", 0)
    offset := getIntQueryParameter(c, "offset", 0)
    tree := getIntQueryParameter(c, "tree", 0)
    preview := getIntQueryParameter(c, "preview", 0)

    // Query comments tree
    comments := logic.GetCommentsTree(thread, depth, max, offset, tree)

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


