package api

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/matthinc/gomment/logic"
)

type routeHandlerType func(*gin.Context, *logic.BusinessLogic)

func injectLogic(routeHandler routeHandlerType, logic *logic.BusinessLogic) gin.HandlerFunc {
	return func(c *gin.Context) {
		routeHandler(c, logic)
	}
}

func getIntQueryParameter(c *gin.Context, parameter string, defaultValue int) int {
	queryParameters := c.Request.URL.Query()

	value := defaultValue
	valueStr, hasValue := queryParameters[parameter]
	if hasValue {
		value, _ = strconv.Atoi(valueStr[0])
	}

	return value
}

func isAdmin(c *gin.Context, logic *logic.BusinessLogic) bool {
	sidCookie, err := c.Request.Cookie(AdminSid)

	if err != nil {
		return false
	}

	// https://github.com/gin-gonic/gin/issues/1717
	sidCookieValue, err := url.QueryUnescape(sidCookie.Value)

	if err != nil {
		return false
	}

	// everyone that has a valid session is a admin
	fmt.Println(sidCookie.Raw)
	_, err = logic.GetSession(sidCookieValue)
	if err == nil {
		return true
	} else {
		return false
	}
}

func adminJsonMiddleware(routeHandler routeHandlerType, logic *logic.BusinessLogic) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isAdmin(c, logic) {
			routeHandler(c, logic)
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{})
		}
	}
}

func redirectWithPrefix(c *gin.Context, destination string) {
	if prefix := c.Request.Header.Get("X-Forwarded-Prefix"); len(prefix) > 0 {
		destination = prefix + "/" + destination
	}
	c.Redirect(http.StatusTemporaryRedirect, destination)
}

func adminRedirectMiddleware(logic *logic.BusinessLogic) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !isAdmin(c, logic) {
			redirectWithPrefix(c, "/admin/login")
			c.Abort()
		}
	}
}

func StartApi(logic *logic.BusinessLogic) {
	router := gin.Default()

	// Static files
	router.Static("/static", "./frontend")

	v1 := router.Group("/api/v1")
	v1.GET("/status", injectLogic(routeStatus, logic))
	v1.GET("/comments", injectLogic(routeGetComments, logic))
	v1.POST("/comment", injectLogic(routePostComment, logic))

	if len(logic.PwHash) > 0 {
		// enable admin API routes
		v1.POST("/admin/login", injectLogic(routeAdminLogin, logic))
		v1.GET("/admin/threads", adminJsonMiddleware(routeAdminThreads, logic))

		// unauthenticated login route
		router.StaticFile("/admin/login", "./frontend/admin/login.html")

		// enable admin static routes
		adminArea := router.Group("/admin")
		adminArea.Use(adminRedirectMiddleware(logic))
		adminArea.StaticFile("/", "./frontend/admin/index.html")
	}

	router.Run(":8000")
}
