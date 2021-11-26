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

const AdminSid = "GOMMENT_SID"

func (api *Api) routeAdminLogin(c *gin.Context) {
	var loginRequest loginRequest
	err := c.BindJSON(&loginRequest)
	if err != nil {
		fmt.Println(err)
	}

	isValid := auth.ValidatePw(loginRequest.Password, api.logic.Administration.PasswordHash)
	var sessionId string
	var sessionData logic.SessionData
	if isValid {
		sessionId, sessionData, err = api.logic.CreateSession()
	}
	if isValid && err == nil {
		c.SetSameSite(http.SameSiteStrictMode)
		c.SetCookie(AdminSid, sessionId, int(logic.SessionDuration.Seconds()), "", "", false, true)

		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			// "session_id": sessionId,
			"valid_until": sessionData.ValidUntil,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
		})
	}
}

func (api *Api) routeAdminThreads(c *gin.Context) {
	threads, err := api.logic.DB.GetThreads()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
		})
		return
	}

	c.JSON(http.StatusOK, threads)
}
