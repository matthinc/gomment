package api

import (
    "github.com/gin-gonic/gin"
    "github.com/matthinc/gomment/model"
    "net/http"
	"fmt"
)

func routeStatus(c *gin.Context) {
    c.String(http.StatusOK, "running")
}

func routePostComment(c *gin.Context) {
    var comment model.Comment
    c.BindJSON(&comment)
    fmt.Println(comment.Text)
    c.String(http.StatusOK, "ok")
}

func StartApi() {
    router := gin.Default()
    router.GET("/status", routeStatus)
    router.POST("/comment", routePostComment)
    router.Run(":8000")
}
