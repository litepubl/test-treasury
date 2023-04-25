package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitRoutes(router *gin.Engine, updateRoutes *UpdateRoutes, findnameRoutes *FindnameRoutes) {
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	router.GET("/update", updateRoutes.Update)
	router.GET("/state", updateRoutes.State)

	router.GET("/get_names", findnameRoutes.GetNames)
}
