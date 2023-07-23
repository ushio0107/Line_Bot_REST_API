package api

import "github.com/gin-gonic/gin"

func NewRouter(a *API) *gin.Engine {
	router := gin.Default()

	r := router.Group("/linebot")
	r.POST("", a.receiveHandler)

	return router
}
