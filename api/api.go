package api

import (
    "net/http"
    "strings"
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

const AUTH_HEADER_PREFIX = "Bearer "

func isAdmin(c *gin.Context, logic *logic.BusinessLogic) bool {
    authHeader := c.Request.Header.Get("Authorization")
    if strings.HasPrefix(authHeader, AUTH_HEADER_PREFIX) {
        authHeader = authHeader[len(AUTH_HEADER_PREFIX):]
    } else {
        return false
    }
    
    // everyone that has a valid session is a admin
    _, err := logic.GetSession(authHeader)
    if err == nil {
        return true
    } else {
        return false
    }
}

func adminArea(routeHandler routeHandlerType, logic *logic.BusinessLogic) gin.HandlerFunc {
    return func(c *gin.Context) {
        if isAdmin(c, logic) {
            routeHandler(c, logic)
        } else {
            c.JSON(http.StatusUnauthorized, gin.H{})
        }
    }
}

func StartApi(logic *logic.BusinessLogic) {
    router := gin.Default()
    router.GET("/status", injectLogic(routeStatus, logic))
    router.GET("/comments", injectLogic(routeGetComments, logic))
    router.POST("/comment", injectLogic(routePostComment, logic))
  
    if len(logic.PwHash) > 0 {
        // enable admin routes
        router.POST("/admin/login", injectLogic(routeAdminLogin, logic))
        router.GET("/admin/threads", adminArea(routeAdminThreads, logic))
    }
    
    router.Run(":8000")
}
