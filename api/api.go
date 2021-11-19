package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

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

func getInt64QueryParameter(c *gin.Context, parameter string, defaultValue int64) int64 {
	queryParameters := c.Request.URL.Query()

	value := defaultValue
	valueStr, hasValue := queryParameters[parameter]
	if hasValue {
		value, _ = strconv.ParseInt(valueStr[0], 10, 64)
	}

	return value
}

func getInt64ListQueryParameter(c *gin.Context, parameter string) []int64 {
	queryParameters := c.Request.URL.Query()

	ret := make([]int64, 0)

	valueStr, hasValue := queryParameters[parameter]
	if hasValue {
		stringList := strings.Split(valueStr[0], ",")
		for _, s := range stringList {
			value, err := strconv.ParseInt(s, 10, 64)
			if err == nil {
				ret = append(ret, value)
			}
		}
	}

	return ret
}

func getStringQueryParameter(c *gin.Context, parameter string) (string, error) {
	queryParameters := c.Request.URL.Query()

	valueStr, hasValue := queryParameters[parameter]
	if hasValue {
		return url.QueryUnescape(valueStr[0])
	}

	return "", errors.New("required parameter '" + parameter + "' not provided in query string")
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

	// create comment
	v1.POST("/comment", injectLogic(routePostComment, logic))

	// newest branch first
	v1.GET("/comments/nbf", injectLogic(routeGetCommentsNbf, logic))
	v1.GET("/morecomments/nbf", injectLogic(routeGetMoreCommentsNbf, logic))
	// newest sibling first
	v1.GET("/comments/nsf", injectLogic(routeGetCommentsNsf, logic))
	v1.GET("/morecomments/nsf", injectLogic(routeGetMoreCommentsNsf, logic))
	// oldest sibling first
	v1.GET("/comments/osf", injectLogic(routeGetCommentsOsf, logic))
	v1.GET("/morecomments/osf", injectLogic(routeGetMoreCommentsOsf, logic))

	if len(logic.Administration.PasswordHash) > 0 {
		// enable admin API routes
		v1.POST("/admin/login", injectLogic(routeAdminLogin, logic))
		v1.GET("/admin/threads", adminJsonMiddleware(routeAdminThreads, logic))

		// unauthenticated login route
		router.StaticFile("/admin/login", "./frontend/admin/login.html")
		router.StaticFile("/admin/style.css", "./frontend/admin/style.css")

		// enable admin static routes
		adminArea := router.Group("/admin")
		adminArea.Use(adminRedirectMiddleware(logic))
		adminArea.StaticFile("/", "./frontend/admin/index.html")
		adminArea.StaticFile("/gomment-admin.js", "./frontend/admin/gomment-admin.js")
	}

	router.Run(":8000")
}
