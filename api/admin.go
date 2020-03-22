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



func routeAdminLogin(c *gin.Context, l *logic.BusinessLogic) {
	var loginRequest loginRequest
	err := c.BindJSON(&loginRequest)
	if err != nil {
		fmt.Println(err)
	}

	isValid := auth.ValidatePw(loginRequest.Password, l.PwHash)
    var sessionId string
    var sessionData logic.SessionData
    if isValid {
        sessionId, sessionData, err = l.CreateSession()
    }
	if isValid && err == nil {
		c.JSON(http.StatusOK, gin.H{
            "status": "success",
            "session_id": sessionId,
            "valid_until": sessionData.ValidUntil,
        })
	} else {
        c.JSON(http.StatusBadRequest, gin.H{
            "status": "error",
        })
	}
}

func routeAdminThreads(c *gin.Context, l *logic.BusinessLogic) {
    threads := l.DB.QueryThreads()
    c.JSON(http.StatusOK, threads)
}
