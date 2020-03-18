package api

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func routeStatus(c *gin.Context) {
    c.String(http.StatusOK, "running")
}

func StartApi() {
    router := gin.Default()
    router.GET("/status", routeStatus)
    router.Run(":8000")
}
