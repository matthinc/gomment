package api

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"

	"github.com/matthinc/gomment/logic"
)

type Api struct {
	logic *logic.BusinessLogic
}

func NewApi(businessLogic *logic.BusinessLogic) Api {
	return Api{
		logic: businessLogic,
	}
}

type routeHandlerType func(*gin.Context)

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

func (api *Api) adminJsonMiddleware(routeHandler routeHandlerType) gin.HandlerFunc {
	return func(c *gin.Context) {
		if isAdmin(c, api.logic) {
			routeHandler(c)
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

func (api *Api) StartApi() {
	router := gin.Default()

	// Static files
	router.Static("/static", "./frontend")

	v1 := router.Group("/api/v1")
	v1.GET("/status", api.routeStatus)

	// create comment
	v1.POST("/comment", api.routePostComment)

	// newest branch first
	v1.GET("/comments/nbf", api.routeGetCommentsNbf)
	v1.GET("/morecomments/nbf", api.routeGetMoreCommentsNbf)
	// newest sibling first
	v1.GET("/comments/nsf", api.routeGetCommentsNsf)
	v1.GET("/morecomments/nsf", api.routeGetMoreCommentsNsf)
	// oldest sibling first
	v1.GET("/comments/osf", api.routeGetCommentsOsf)
	v1.GET("/morecomments/osf", api.routeGetMoreCommentsOsf)

	if len(api.logic.Administration.PasswordHash) > 0 {
		// enable admin API routes
		v1.POST("/admin/login", api.routeAdminLogin)
		v1.GET("/admin/threads", api.adminJsonMiddleware(api.routeAdminThreads))

		// unauthenticated login route
		router.StaticFile("/admin/login", "./frontend/admin/login.html")
		router.StaticFile("/admin/style.css", "./frontend/admin/style.css")

		// enable admin static routes
		adminArea := router.Group("/admin")
		adminArea.Use(adminRedirectMiddleware(api.logic))
		adminArea.StaticFile("/", "./frontend/admin/index.html")
		adminArea.StaticFile("/gomment-admin.js", "./frontend/admin/gomment-admin.js")
	}

	router.Run(":8000")
}
