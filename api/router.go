package api

import "github.com/gin-gonic/gin"

func NewRouter(a *API) *gin.Engine {
	router := gin.Default()

	r := router.Group("/linebot")
	r.POST("", a.filterMessage /* Middleware */, a.receiveHandler)
	r.POST("/broadcast", a.broadcastMessage)
	r.GET("/get", a.getAllMessages)

	return router
}
