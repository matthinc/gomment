package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (api *Api) routeStatus(c *gin.Context) {
	c.String(http.StatusOK, `{"status": "ok"}`)
}
