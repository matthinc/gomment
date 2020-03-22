package api

import (
    "github.com/gin-gonic/gin"
    "github.com/matthinc/gomment/logic"
)

type routeHandlerType func(*gin.Context, *logic.BusinessLogic)

func injectLogic(routeHandler routeHandlerType, logic *logic.BusinessLogic) gin.HandlerFunc {
    return func(c *gin.Context) {
        routeHandler(c, logic)
    }
}

func StartApi(logic *logic.BusinessLogic) {
    router := gin.Default()
    router.GET("/status", injectLogic(routeStatus, logic))
    router.GET("/comments", injectLogic(routeGetComments, logic))
    router.POST("/comment", injectLogic(routePostComment, logic))
    router.GET("/comments/preview", injectLogic(routePreviewComments, logic))

    if len(logic.PwHash) > 0 {
        // enable admin routes
        router.POST("/admin/login", injectLogic(routeAdminLogin, logic))
    }
    
    router.Run(":8000")
}
