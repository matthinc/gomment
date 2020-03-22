package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/matthinc/gomment/auth"
	"github.com/matthinc/gomment/logic"
)

type loginRequest struct {
	Password string `json:"password" binding:"required"`
}

func routeAdminLogin(c *gin.Context, logic *logic.BusinessLogic) {
	var loginRequest loginRequest
	err := c.BindJSON(&loginRequest)
	if err != nil {
		fmt.Println(err)
	}

	isValid := auth.ValidatePw(loginRequest.Password, logic.PwHash)

	if isValid {
		c.String(http.StatusOK, `{"status": "success"}`)
	} else {
		c.String(http.StatusBadRequest, `{"status": "error"}`)
	}
}
