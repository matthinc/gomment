package api

import (
    "github.com/gin-gonic/gin"
    "github.com/matthinc/gomment/model"
    "github.com/matthinc/gomment/logic"
    "net/http"
    "fmt"
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
