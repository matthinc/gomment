package api

import (
	"github.com/gin-gonic/gin"
	"github.com/matthinc/gomment/logic"
	"net/http"
)

func routeStatus(c *gin.Context, logic *logic.BusinessLogic) {
	c.String(http.StatusOK, `{"status": "ok"}`)
}
